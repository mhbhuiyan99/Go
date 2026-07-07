# Testing Beego Controllers — Request Setup and Response Checking

> **Who this is for:** Go developers who want to understand the underlying
> mechanics of testing a Beego controller — building a request, sending it
> through Beego's router, and inspecting what came back — before applying
> the table-driven pattern on top.
>
> **See also:** `12-testing-controllers.md` covers the full table-driven
> pattern for controller tests. This guide focuses on the request/response
> mechanics in isolation.

---

## The Problem Beego Testing Solves

A controller method never runs on its own — it needs:
- A `*http.Request` (the incoming request)
- An `http.ResponseWriter` (something to write the response to)
- Beego's router to dispatch the request to the right controller method

`httptest` gives you all three without starting a real server on a real port.

```
Real request flow:
  Browser → network → Beego server → router → controller → response

Test request flow:
  Your test → httptest.Request → Beego's router (in-process) → controller → httptest.Recorder
```

No network, no port, no external process. Everything happens in memory,
inside the same test binary.

---

## Step 1 — Build the Request

```go
import "net/http"

r, err := http.NewRequest("GET", "/api/v1/health", nil)
```

Three arguments:
- **Method** — `"GET"`, `"POST"`, `"PUT"`, `"DELETE"`
- **URL path** — can include query params: `"/api/v1/expenses?category=Food"`
- **Body** — an `io.Reader`, or `nil` for requests with no body

### Building a request with a JSON body

```go
import (
    "net/http"
    "strings"
)

body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`

r, err := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
if err != nil {
    t.Fatalf("failed to build request: %v", err)
}
```

`strings.NewReader` converts a plain string into an `io.Reader`, which is
what `http.NewRequest` expects for the body argument.

### Setting headers

```go
r.Header.Set("Content-Type", "application/json")
r.Header.Set("X-User-ID", "1")
```

Headers must be set **after** the request is created — `http.NewRequest`
does not accept headers as an argument.

---

## Step 2 — Build the Response Recorder

```go
import "net/http/httptest"

w := httptest.NewRecorder()
```

`httptest.NewRecorder()` returns a `*httptest.ResponseRecorder` — a fake
`http.ResponseWriter` that captures everything written to it instead of
sending it over a network.

After the handler runs, the recorder holds:

```go
w.Code          // int — the HTTP status code, e.g. 200, 404
w.Body          // *bytes.Buffer — the raw response body
w.Header()      // http.Header — response headers
```

---

## Step 3 — Send the Request Through Beego's Router

Beego's global router dispatcher is `beego.BeeApp.Handlers`. It implements
`http.Handler`, so you call `.ServeHTTP(w, r)` exactly as a real server would:

```go
import beego "github.com/beego/beego/v2/server/web"

beego.BeeApp.Handlers.ServeHTTP(w, r)
```

This one line:
1. Looks at `r.Method` and `r.URL.Path`
2. Matches it against every route registered in `routers/router.go`
3. Creates an instance of the matching controller
4. Calls the matching method (`Get`, `Post`, or an explicitly mapped name)
5. The controller writes its response into `w`

After this line runs, `w` contains everything the controller produced.

---

## Step 4 — Inspect the Response

### Checking the status code

```go
if w.Code != 200 {
    t.Errorf("expected status 200, got %d", w.Code)
}
```

### Checking the raw body as a string

```go
body := w.Body.String()
if !strings.Contains(body, "Server is running") {
    t.Errorf("body does not contain expected text: %s", body)
}
```

### Parsing the body as JSON (preferred)

```go
import "encoding/json"

var result map[string]interface{}
if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
    t.Fatalf("failed to parse JSON response: %v", err)
}

if result["status"] != "success" {
    t.Errorf("expected status=success, got %v", result["status"])
}
```

`w.Body` is a `*bytes.Buffer`. Use `.Bytes()` for `json.Unmarshal`,
or `.String()` for plain text comparisons.

---

## A Complete Minimal Test

```go
package controllers_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"

    beego "github.com/beego/beego/v2/server/web"
    _ "myapp/routers" // registers routes via init()
)

func TestMain(m *testing.M) {
    os.Chdir("../") // so Beego finds conf/app.conf
    os.Exit(m.Run())
}

func TestHealthCheck(t *testing.T) {
    // 1. build the request
    r, err := http.NewRequest("GET", "/api/v1/health", nil)
    if err != nil {
        t.Fatalf("failed to build request: %v", err)
    }

    // 2. build the recorder
    w := httptest.NewRecorder()

    // 3. send through Beego's router
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    // 4. check the status code
    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }

    // 5. check the body
    if !strings.Contains(w.Body.String(), "running") {
        t.Errorf("unexpected body: %s", w.Body.String())
    }
}
```

Run it:

```bash
go test ./controllers/... -v -run TestHealthCheck
```

Expected output:

```
=== RUN   TestHealthCheck
--- PASS: TestHealthCheck (0.00s)
PASS
```

---

## Testing a POST Request with a Body

```go
func TestRegisterUser(t *testing.T) {
    body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`

    r, err := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
    if err != nil {
        t.Fatalf("failed to build request: %v", err)
    }
    r.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    if w.Code != http.StatusCreated {
        t.Errorf("expected status 201, got %d — body: %s", w.Code, w.Body.String())
    }

    var result map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
        t.Fatalf("failed to parse response: %v", err)
    }
    if result["status"] != "success" {
        t.Errorf("expected status=success, got %v", result["status"])
    }
}
```

---

## Testing a Request with Headers (Auth Pattern)

```go
func TestGetExpenses_Unauthorized(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/v1/expenses", nil)
    // no X-User-ID header set — should be rejected

    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    if w.Code != http.StatusUnauthorized {
        t.Errorf("expected status 401, got %d", w.Code)
    }
}

func TestGetExpenses_Authorized(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/v1/expenses", nil)
    r.Header.Set("X-User-ID", "1")

    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }
}
```

---

## Testing Query Parameters

Query parameters are just part of the URL string — no special handling needed:

```go
func TestListExpenses_FilterByCategory(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/v1/expenses?category=Food", nil)
    r.Header.Set("X-User-ID", "1")

    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }
}
```

---

## Testing Path Parameters

Path parameters are also just part of the URL — Beego's router extracts
them internally when it matches the route:

```go
func TestGetExpenseByID(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/v1/expenses/1", nil)
    r.Header.Set("X-User-ID", "1")

    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }
}
```

No manual extraction needed in the test — Beego's router does this the
same way it does for real requests.

---

## Why the Working Directory Matters

If you skip `os.Chdir("../")` in `TestMain`, you may see:

```
init global config instance failed. open conf/app.conf: no such file or directory
```

This happens because `go test` runs with the working directory set to the
**package folder** (e.g. `controllers/`), not the project root. Beego looks
for `conf/app.conf` relative to the working directory and cannot find it.

```go
func TestMain(m *testing.M) {
    os.Chdir("../") // move up one level to the project root
    os.Exit(m.Run())
}
```

If your test file is nested deeper, adjust accordingly:
```go
os.Chdir("../../") // two levels up
```

---

## Why the Blank Import Is Required

```go
import _ "myapp/routers"
```

Without this import, `routers/router.go`'s `init()` function never runs,
which means no routes are ever registered with `beego.BeeApp.Handlers`.
Every request in your test would return `404`, regardless of what your
controller code actually does.

```go
// Without the blank import — this always fails, controller code is irrelevant
beego.BeeApp.Handlers.ServeHTTP(w, r) // r matches nothing, returns 404
```

---

## Common Mistakes

### Forgetting to set `Content-Type` on POST/PUT requests

```go
// May cause Beego to fail parsing the body correctly
r, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(body))

// Always set this for JSON bodies
r.Header.Set("Content-Type", "application/json")
```

---

### Reading `w.Body` before the handler has run

```go
// WRONG — reading before ServeHTTP is called, body is empty
fmt.Println(w.Body.String())
beego.BeeApp.Handlers.ServeHTTP(w, r)

// CORRECT — read after
beego.BeeApp.Handlers.ServeHTTP(w, r)
fmt.Println(w.Body.String())
```

---

### Reusing the same recorder across multiple requests

```go
// WRONG — second request's output mixes with the first
w := httptest.NewRecorder()
beego.BeeApp.Handlers.ServeHTTP(w, r1)
beego.BeeApp.Handlers.ServeHTTP(w, r2) // w still has r1's data

// CORRECT — new recorder per request
w1 := httptest.NewRecorder()
beego.BeeApp.Handlers.ServeHTTP(w1, r1)

w2 := httptest.NewRecorder()
beego.BeeApp.Handlers.ServeHTTP(w2, r2)
```

---

### Using `httptest.NewRequest` vs `http.NewRequest`

Both work, but they differ slightly:

```go
// http.NewRequest — returns (*http.Request, error), you must check the error
r, err := http.NewRequest("GET", "/path", nil)
if err != nil {
    t.Fatalf("failed to build request: %v", err)
}

// httptest.NewRequest — returns *http.Request only, panics on invalid input
r := httptest.NewRequest("GET", "/path", nil)
```

`httptest.NewRequest` is more convenient for tests since there is rarely
a real error case with a hardcoded valid path. `http.NewRequest` is the
safer choice if the path or method is built dynamically and could be invalid.

---

## Quick Reference

```go
import (
    "net/http"
    "net/http/httptest"
    "encoding/json"
    beego "github.com/beego/beego/v2/server/web"
    _ "myapp/routers"
)

func TestMain(m *testing.M) {
    os.Chdir("../")
    os.Exit(m.Run())
}

func TestSomething(t *testing.T) {
    // 1. build request
    r, _ := http.NewRequest("POST", "/api/v1/path", strings.NewReader(`{"key":"value"}`))
    r.Header.Set("Content-Type", "application/json")
    r.Header.Set("X-User-ID", "1")

    // 2. build recorder
    w := httptest.NewRecorder()

    // 3. dispatch through Beego
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    // 4. check status
    if w.Code != http.StatusOK {
        t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
    }

    // 5. check body
    var result map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &result)
    if result["status"] != "success" {
        t.Errorf("unexpected status field: %v", result["status"])
    }
}
```

| Step | Code |
|---|---|
| Build request | `http.NewRequest(method, path, body)` |
| Set headers | `r.Header.Set(key, value)` |
| Build recorder | `httptest.NewRecorder()` |
| Dispatch | `beego.BeeApp.Handlers.ServeHTTP(w, r)` |
| Check status | `w.Code` |
| Check raw body | `w.Body.String()` |
| Parse JSON body | `json.Unmarshal(w.Body.Bytes(), &result)` |

---

## Further Reading

- [Go `net/http/httptest` package](https://pkg.go.dev/net/http/httptest)
- `12-testing-controllers.md` — table-driven pattern applied to these mechanics
