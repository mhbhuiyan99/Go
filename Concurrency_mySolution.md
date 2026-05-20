
## Exercise 1.1: Where's My Output?

### Option A: `sync.WaitGroup`
You wait for completion, not time.
```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    wg.Add(1) // tell WaitGroup: "we're waiting for 1 thing"

    go func() {
        defer wg.Done() // decrement counter when done
        fmt.Println("done")
    }()

    wg.Wait() // block until counter hits 0
}
```

### Option B: Channel
Channels are Go's primary concurrency primitive. `struct{}` is used when you only care about the signal, not data.
```go
package main

import "fmt"

func main() {
    done := make(chan struct{}) // empty struct = zero-byte signal

    go func() {
        fmt.Println("done")
        close(done) // signal completion
    }()

    <-done // block until channel is closed
}
```

### Option C: For a single value, just use the channel directly
```go
func main() {
    result := make(chan string)

    go func() {
        result <- "done" // send data back
    }()

    fmt.Println(<-result) // receive and print
}
```

🧠 **Key Insight**

| Approach         | Use When                                                  |
| ---------------- | --------------------------------------------------------- |
| `time.Sleep`     | **Never** for synchronization (only for simulating delay) |
| `sync.WaitGroup` | Waiting for multiple goroutines to finish                 |
| `chan struct{}`  | Signaling completion / cancellation                       |
| `chan T`         | Passing data between goroutines                           |

-------

## Exercise 1.2: The Channel Handshake
```go
func main() {
	ch := make(chan string) // unbuffered channel created
    //  send in goroutine
	go func() {
		ch <- "hello from goroutine"
	}()
	fmt.Println(<-ch) // receive in main
}
```
#### Q1: What happens if you swap roles (goroutine receives, main sends)?
```go
func main() {
    ch := make(chan string)

    go func() {
        msg := <-ch      // goroutine tries to RECEIVE
        fmt.Println(msg)
    }()

    ch <- "hello from main"   // main tries to SEND
}
```
It works perfectly fine. An unbuffered channel requires both a sender and receiver to be ready at the same time — but it doesn't matter which goroutine is which. The channel blocks until the handshake happens. <br> <br>

#### Q2: What happens if you remove the goroutine and do both in main?
```go
func main() {
    ch := make(chan string)

    ch <- "hello"    // ❌ BLOCKS FOREVER (deadlock!)
    msg := <-ch      // this line never runs
    fmt.Println(msg)
}
```
**Result:** `fatal error: all goroutines are asleep - deadlock!` <br>
Why? On an unbuffered channel, the send blocks until someone receives. But there's no other goroutine to receive — main is stuck on the send and can never reach the receive. <br>

🧠 **The Rule**
| Scenario                                                  | Result                          |
| --------------------------------------------------------- | ------------------------------- |
| Send & receive in **different goroutines**                | ✅ Works (synchronous handshake) |
| Send & receive in **same goroutine**, unbuffered          | ❌ Deadlock                      |
| Send & receive in **same goroutine**, buffered (size ≥ 1) | ✅ Works (send doesn't block)    |
<br>

**Does it deadlock?**
```go
func main() {
    ch := make(chan string, 1)  // buffered with capacity 1

    ch <- "hello"
    msg := <-ch
    fmt.Println(msg)
}
```
**No deadlock.** The buffer holds the value, so the send doesn't block. The receive happens later in the same goroutine — fine.

-------
## Exercise 1.3: Parallel Squares
