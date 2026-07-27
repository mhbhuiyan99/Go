# Go Profiling — A Practical Guide

> A structured, from-scratch guide to Go's `pprof`-based profiling tools — what each profile measures, how to collect it, and how to actually read the output. Written to be followed independently, no prior profiling experience assumed.
>
> Further reading: [go.dev/doc/diagnostics](https://go.dev/doc/diagnostics) (official overview) · [JetBrains profiling guide](https://blog.jetbrains.com/go/2026/05/20/golang-profiling-guide/) (deep dive) · [dev.to: Mastering Go profiling](https://dev.to/vplatonov/mastering-go-profiling-3fi3) (quick reference) · [Oodle: Go Profiling in Production](https://blog.oodle.ai/go-profiling-in-production/) (real-world case study)

---

## 1. What is profiling, and why does it exist?

Profiling answers one question: **where does my program spend its resources (CPU time, memory, waiting time)?**

Go's runtime samples your program's call stack — either on a timer (e.g. 100 times/sec for CPU) or on specific events (e.g. every time a mutex is contended) — and writes that data out in a format the `pprof` tool understands. You then use `pprof` to turn raw samples into something readable: a table, a call graph, or a flame graph.

Profiling sits alongside three other diagnostic categories in Go:
- **Profiling** — where resources go (this doc)
- **Tracing** — the timeline of events (scheduling, GC, blocking) over a window of execution — `go tool trace`
- **Debugging** — pausing execution to inspect state — Delve
- **Runtime stats** — high-level health metrics (`runtime.ReadMemStats`, `GODEBUG=gctrace=1`, etc.)

⚠️ Rule of thumb: **profile first to find *what* is slow, trace to understand *why* it's scheduled that way.**

---

## 2. The profile types

| Profile | What it measures | Enabled by default? | Use it when... |
|---|---|---|---|
| **cpu** | Where CPU cycles are spent (active execution only) | No — must be started/stopped explicitly | App is CPU-bound / slow processing |
| **heap** | Currently-live memory allocations (as of last GC) | Yes (always collecting samples) | Suspected memory leak, high RAM usage |
| **allocs** | *All* allocations since process start (including freed) | Yes (same underlying data as heap) | High GC pressure, allocation-heavy hot loops |
| **goroutine** | Stack traces of all goroutines *right now* | Yes | App hangs, goroutine count keeps growing, deadlock |
| **block** | Where goroutines wait on sync primitives (mutex, channel, waitgroup) | No — `runtime.SetBlockProfileRate` | Low CPU usage but high latency |
| **mutex** | Which goroutines *cause* other goroutines to wait (lock contention) | No — `runtime.SetMutexProfileFraction` | Throughput doesn't scale with concurrency |
| **threadcreate** | Where new OS threads get created | Yes | Rarely used directly |

**Important distinction (this trips people up):**
- `heap` and `allocs` are *the same underlying data* — just presented with a different default sample type.
  - `inuse_space` / `inuse_objects` → what's live right now (heap's default)
  - `alloc_space` / `alloc_objects` → cumulative since start, including freed memory (allocs' default)
  - Growth in `allocs` ≠ a leak. It just means the app is allocating a lot — check `inuse_*` to know what's actually still held.

- `block` tells you **what** is waiting. `mutex` tells you **what's causing** the wait. You often read them together: goroutine profile shows the pile-up → block profile shows who's stuck → mutex profile shows who's holding the lock too long.

---

## 3. Collecting profiles — two main ways

### a) `runtime/pprof` — explicit, code-controlled

Good for CLI tools, one-off scripts, tests, or anywhere you want to trigger a profile programmatically.

```go
import (
    "os"
    "runtime"
    "runtime/pprof"
)

// CPU has a start/stop API because it streams over a time window
func captureCPU() error {
    f, _ := os.Create("cpu.pprof")
    defer f.Close()
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    doWork() // the code you want to measure
    return nil
}

// Heap should be captured right after a forced GC to see current live state
func captureHeap() error {
    runtime.GC()
    f, _ := os.Create("heap.pprof")
    defer f.Close()
    return pprof.Lookup("heap").WriteTo(f, 0)
}
```

### b) `net/http/pprof` — the standard for long-running services

```go
import (
    "log"
    "net/http"
    _ "net/http/pprof" // side-effect import registers handlers
)

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    runServer()
}
```

Then pull profiles live over HTTP:

```bash
# CPU profile, sampled for 30s
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Current live heap (force GC first for an accurate snapshot)
go tool pprof http://localhost:6060/debug/pprof/heap?gc=1

# All allocations since start
go tool pprof http://localhost:6060/debug/pprof/allocs

# Block / mutex — only meaningful if enabled via SetBlockProfileRate / SetMutexProfileFraction
go tool pprof http://localhost:6060/debug/pprof/block
go tool pprof http://localhost:6060/debug/pprof/mutex

# Human-readable goroutine dump (no pprof needed)
curl "http://localhost:6060/debug/pprof/goroutine?debug=2"
```

Enabling block/mutex profiling in code (do this deliberately, it adds overhead):
```go
runtime.SetBlockProfileRate(1)        // sample every blocking event
runtime.SetMutexProfileFraction(1)    // sample every mutex contention
```

> Visiting `http://localhost:6060/debug/pprof/` directly in a browser gives you an index of everything available.

---

## 4. Reading the output

Once you have a profile file, `go tool pprof` reads it the same way regardless of how it was collected.

```bash
go tool pprof ./my-binary cpu.pprof     # interactive shell
go tool pprof -http=:8081 ./my-binary cpu.pprof   # web UI with graphs + flame graph
```

### Key interactive commands
| Command | What it shows |
|---|---|
| `top` | Table of the most expensive functions |
| `list <func>` | Source code annotated with per-line cost |
| `tree` | Same data as top, organized by caller |
| `web` | Opens call graph visualization |
| `peek <func>` | Statistical view of a specific function's callers/callees |

### The two numbers that matter: flat vs cumulative
- **flat** — resource consumed *inside that function itself*, excluding what it calls. High flat = the function itself is expensive.
- **cumulative (cum)** — resource consumed by the function *and everything it calls*. High cum, low flat = the problem is in a callee, not this function.

### Flame graphs
- Width = cost (time or memory) — **this is the number that matters**.
- Height = call stack depth — a tall skinny tower is *not* automatically your bottleneck.
- Read bottom-up: widest blocks at the bottom of the stack are where to look first.

### Red flags to watch for
- `runtime.mallocgc` dominating → too many/too-small allocations
- `sync.(*Mutex).Lock` high in a CPU profile → contention, not real work
- Many narrow repeated blocks in a flame graph → inefficient work inside a loop
- `reflect.*` showing up hot → often reflection-based (de)serialization — a common one is `encoding/json` under load; consider a codegen-based encoder if it matters

---

## 5. Practical rules of thumb

- **It's safe to profile in production**, but not free — expect some CPU overhead while a CPU profile is running. Test the overhead before turning it on live.
- Collect **one profile type at a time** — some profiles interfere with each other (e.g. heap profiling can skew CPU profiles).
- CPU profiles should run for at least **10–30 seconds** to get a meaningful sample.
- **heap** profiles are cheap/safe to collect continuously; **block** and **mutex** should only be enabled briefly when actively investigating something, since they add real overhead.
- Compare before/after with diffing:
  ```bash
  go tool pprof -diff_base old_heap.pprof new_heap.pprof
  ```
- In production, restrict access to `/debug/pprof/` (internal network only / auth) — it can leak information about your binary and, in theory, aid a DoS if left wide open.

---

## 6. Suggested exercises

> Reading about profiling only gets you so far — the concepts click once you've actually looked at real output. Try these on any small Go program (even a toy one):

1. **CPU hotspot:** Write a function that's deliberately inefficient (e.g. string concatenation with `+` in a loop, or repeated `fmt.Sprintf`). Profile it, find it in `top`, then fix it and profile again to confirm the fix worked.
2. **Heap vs allocs:** Write a function that allocates a large slice and holds onto it, and another that allocates the same size but lets it go out of scope. Compare `heap` and `allocs` profiles between the two.
3. **Goroutine leak:** Start a goroutine that reads from a channel that's never written to or closed. Watch it show up in the goroutine profile as the count grows.
4. **Mutex contention:** Have several goroutines fight over one `sync.Mutex` guarding a slow operation. Compare the block and mutex profiles.
5. **Diffing:** Capture a heap profile, cause an allocation-heavy operation, capture another, and diff them with `-diff_base`.

### Template for recording what you find

Copy this block per exercise — writing down what you actually saw (not just what you expected to see) is what makes it stick.

```markdown
### Exercise: <name>
**Command(s) run:**
```bash```


**What `top` / flame graph showed:**
-

**What the fix/change was:**
-

**Gotchas hit:**
-

```

---

## 7. Open questions / things to look up later
- [ ] `go tool trace` — how it differs in practice from CPU profiling (JetBrains article mentions it but doesn't cover it)
- [ ] Profile-guided optimization (`go doc pgo`) — CPU profiles can feed the compiler itself
- [ ] Continuous/production profiling tools (Parca, Pyroscope) — only relevant once basics are solid
