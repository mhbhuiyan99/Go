# Testing HTTP Handlers in Go — Table-Driven Pattern

> **Who this is for:** Go developers who understand the basics of
> table-driven tests and want to apply that pattern to testing HTTP
> handlers — both raw `net/http` handlers and framework-based handlers
> (like Beego).

---

## The Core Tool — `net/http/httptest`

Go's standard library includes `net/http/httptest` specifically for
testing HTTP handlers without starting a real server.

Two types you use in every HTTP test:

```go
import "net/http/httptest"

// httptest.NewRecorder() — a fake http.ResponseWriter
// captures the status code, headers, and body your handler writes
w := httptest.NewRecorder()

// http.NewRequest() — builds a real *http.Request
// without needing a real network connection
r := http.NewRequest("POST", "/api/v1/users", body)
```

Your handler writes to `w` exactly as it would in production.
After the handler returns, you inspect `w` to verify the output.

---

## A Minimal Example — Raw `net/http`

```go
package handlers_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}

func TestHealthHandler(t *testing.T) {
    r := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()

    HealthHandler(w, r)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }
    if w.Body.String() != `{"status":"ok"}` {
        t.Errorf("unexpected body: %s", w.Body.String())
    }
}
```

`w.Code` holds the status code.
`w.Body` is a `*bytes.Buffer` — call `.String()` to read it.

---

## The Full Table-Driven Pattern for HTTP Handlers

The test case struct for HTTP handlers typically has these fields:

```go
tests := []struct {
    name       string            // describes what this case tests
    method     string            // "GET", "POST", "PUT", "DELETE"
    path       string            // URL path, may include query params
    body       string            // JSON request body (empty string if none)
    headers    map[string]string // request headers
    wantStatus int               // expected HTTP status code
    wantMsg    string            // expected message in JSON response
}{
    // test cases here
}
```

Not every field is needed in every test. Add fields as your handler
requires them — query params, auth tokens, path parameters.

---

## Complete Example — Testing a Create User Handler

```go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestCreateUserHandler(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
        wantMsg    string
    }{
        {
            name:       "valid user returns 201",
            body:       `{"name":"Alice","email":"alice@example.com","password":"secret123"}`,
            wantStatus: 201,
            wantMsg:    "User created",
        },
        {
            name:       "missing name returns 400",
            body:       `{"email":"alice@example.com","password":"secret123"}`,
            wantStatus: 400,
            wantMsg:    "Name is required",
        },
        {
            name:       "invalid email returns 400",
            body:       `{"name":"Alice","email":"notanemail","password":"secret123"}`,
            wantStatus: 400,
            wantMsg:    "Invalid email",
        },
        {
            name:       "password too short returns 400",
            body:       `{"name":"Alice","email":"alice@example.com","password":"abc"}`,
            wantStatus: 400,
            wantMsg:    "Password must be at least 6 characters",
        },
        {
            name:       "malformed JSON returns 400",
            body:       `{invalid}`,
            wantStatus: 400,
            wantMsg:    "Invalid JSON",
        },
        {
            name:       "empty body returns 400",
            body:       "",
            wantStatus: 400,
            wantMsg:    "Invalid JSON",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // build the request
            r := httptest.NewRequest("POST", "/api/v1/users",
                strings.NewReader(tt.body))
            r.Header.Set("Content-Type", "application/json")

            // build the recorder
            w := httptest.NewRecorder()

            // call the handler
            CreateUserHandler(w, r)

            // assert status code
            if w.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d — body: %s",
                    w.Code, tt.wantStatus, w.Body.String())
            }

            // assert message in response body
            if tt.wantMsg != "" && !strings.Contains(w.Body.String(), tt.wantMsg) {
                t.Errorf("body %q does not contain %q", w.Body.String(), tt.wantMsg)
            }
        })
    }
}
```

---

## Parsing the JSON Response Body

`strings.Contains` is simple but fragile — if the message wording changes,
the test breaks. For structured assertions, decode the response:

```go
// helper — decode response body into a map
func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
    t.Helper()
    var result map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
        t.Fatalf("failed to parse response body %q: %v", w.Body.String(), err)
    }
    return result
}
```

```go
// use in tests
result := parseBody(t, w)

// check the status field
if result["status"] != "success" {
    t.Errorf("expected status=success, got %v", result["status"])
}

// check the message field
if result["message"] != "User created" {
    t.Errorf("expected message %q, got %v", "User created", result["message"])
}

// check a nested data field
data, ok := result["data"].(map[string]interface{})
if !ok {
    t.Fatalf("expected data object, got %T", result["data"])
}
if data["email"] != "alice@example.com" {
    t.Errorf("expected email in data, got %v", data["email"])
}
```

---

## Testing Handlers That Require Headers

Add a `headers` field to the test struct:

```go
func TestGetExpenseHandler(t *testing.T) {
    tests := []struct {
        name       string
        path       string
        headers    map[string]string
        wantStatus int
    }{
        {
            name:       "valid user ID returns 200",
            path:       "/api/v1/expenses",
            headers:    map[string]string{"X-User-ID": "1"},
            wantStatus: 200,
        },
        {
            name:       "missing header returns 401",
            path:       "/api/v1/expenses",
            headers:    map[string]string{},
            wantStatus: 401,
        },
        {
            name:       "invalid header value returns 401",
            path:       "/api/v1/expenses",
            headers:    map[string]string{"X-User-ID": "notanumber"},
            wantStatus: 401,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            r := httptest.NewRequest("GET", tt.path, nil)

            // set all headers from the test case
            for key, val := range tt.headers {
                r.Header.Set(key, val)
            }

            w := httptest.NewRecorder()
            GetExpenseHandler(w, r)

            if w.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
            }
        })
    }
}
```

---

## Testing Handlers That Use Path Parameters

Raw `net/http` does not parse path parameters — that is a router concern.
Two approaches:

**Option 1 — pass the value directly to the handler (no router involved):**

```go
// If your handler reads the ID from the URL itself using a helper:
r := httptest.NewRequest("GET", "/api/v1/expenses/42", nil)
```

**Option 2 — use your router to dispatch the request (tests the full stack):**

```go
// register your routes into a mux, then send requests through it
mux := http.NewServeMux()
mux.HandleFunc("/api/v1/expenses/", GetExpenseHandler)

r := httptest.NewRequest("GET", "/api/v1/expenses/42", nil)
w := httptest.NewRecorder()
mux.ServeHTTP(w, r) // routes the request through the mux
```

The second approach tests routing + handler together, which is closer
to production behaviour.

---

## Testing Beego Handlers

Beego uses its own router. To test Beego handlers, send requests through
`beego.BeeApp.Handlers` — Beego's internal dispatcher:

```go
package controllers_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    beego "github.com/beego/beego/v2/server/web"
    _ "myapp/routers" // blank import registers all routes via init()
)

// doRequest sends a request through Beego's router and returns the recorder.
func doRequest(method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
    var buf *bytes.Buffer
    if body != "" {
        buf = bytes.NewBufferString(body)
    } else {
        buf = bytes.NewBufferString("")
    }

    r, _ := http.NewRequest(method, path, buf)
    r.Header.Set("Content-Type", "application/json")
    for k, v := range headers {
        r.Header.Set(k, v)
    }

    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r) // send through Beego's router
    return w
}
```

Now you test exactly as if a real client sent the request:

```go
func TestAuthRegister(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
        wantField  string // "success" or "error"
    }{
        {
            name:       "valid registration returns 201",
            body:       `{"name":"Alice","email":"alice@example.com","password":"secret123"}`,
            wantStatus: 201,
            wantField:  "success",
        },
        {
            name:       "duplicate email returns 409",
            body:       `{"name":"Alice","email":"alice@example.com","password":"secret123"}`,
            wantStatus: 409,
            wantField:  "error",
        },
        {
            name:       "missing name returns 400",
            body:       `{"email":"alice@example.com","password":"secret123"}`,
            wantStatus: 400,
            wantField:  "error",
        },
        {
            name:       "password too short returns 400",
            body:       `{"name":"Alice","email":"new@example.com","password":"abc"}`,
            wantStatus: 400,
            wantField:  "error",
        },
    }

    cleanData() // reset test CSV files once before this table

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := doRequest("POST", "/api/v1/auth/register", tt.body, nil)

            if w.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d — body: %s",
                    w.Code, tt.wantStatus, w.Body.String())
            }

            result := parseBody(t, w)
            if result["status"] != tt.wantField {
                t.Errorf("status field = %v, want %q", result["status"], tt.wantField)
            }
        })
    }
}
```

---

## `TestMain` for Beego Controller Tests

Beego reads `conf/app.conf` relative to the working directory.
When `go test` runs a package inside `controllers/`, the working directory
is `controllers/` — not the project root. Beego cannot find the config file
and falls back to defaults, or fails silently.

Fix this in `TestMain`:

```go
// controllers/auth_test.go
package controllers_test

import (
    "os"
    "testing"

    beego "github.com/beego/beego/v2/server/web"
    _ "myapp/routers"
)

func TestMain(m *testing.M) {
    // move to project root so beego finds conf/app.conf
    os.Chdir("../")

    // redirect CSV reads/writes to test-specific files
    // so tests never touch real data
    beego.AppConfig.Set("datadir", "testdata")
    beego.AppConfig.Set("userfile", "users_test.csv")
    beego.AppConfig.Set("expensefile", "expenses_test.csv")

    code := m.Run()

    // clean up test data after all tests finish
    os.RemoveAll("testdata")

    os.Exit(code)
}
```

**Why `package controllers_test` not `package controllers`:**

Using the external test package (`_test` suffix) means your test file
cannot access unexported identifiers from `controllers`. This is
intentional — it forces you to test the public HTTP interface, not
internal implementation details. It also avoids circular imports when
`routers` imports `controllers`.

---

## The `registerAndLogin` Helper Pattern

Controller tests that require authentication need a user to exist first.
Extract this into a helper rather than repeating it in every test function:

```go
// registerAndLogin creates a user, logs in, and returns the user ID as string.
func registerAndLogin(t *testing.T, email string) string {
    t.Helper()

    // register
    w := doRequest("POST", "/api/v1/auth/register",
        fmt.Sprintf(`{"name":"Test","email":%q,"password":"secret123"}`, email),
        nil,
    )
    if w.Code != 201 {
        t.Fatalf("register failed: status=%d body=%s", w.Code, w.Body.String())
    }

    // login
    w = doRequest("POST", "/api/v1/auth/login",
        fmt.Sprintf(`{"email":%q,"password":"secret123"}`, email),
        nil,
    )
    if w.Code != 200 {
        t.Fatalf("login failed: status=%d body=%s", w.Code, w.Body.String())
    }

    result := parseBody(t, w)
    data := result["data"].(map[string]interface{})
    // JSON numbers decode as float64 — convert to int string
    return fmt.Sprintf("%d", int(data["user_id"].(float64)))
}
```

Usage in a test:

```go
func TestExpenseCreate(t *testing.T) {
    userID := registerAndLogin(t, "create@example.com")

    tests := []struct {
        name       string
        body       string
        wantStatus int
    }{
        {
            name:       "valid expense returns 201",
            body:       `{"title":"Lunch","amount":350.50,"category":"Food","expense_date":"2025-06-10"}`,
            wantStatus: 201,
        },
        {
            name:       "missing title returns 400",
            body:       `{"amount":350.50,"category":"Food","expense_date":"2025-06-10"}`,
            wantStatus: 400,
        },
        {
            name:       "no auth header returns 401",
            body:       `{"title":"Lunch","amount":100,"category":"Food","expense_date":"2025-06-10"}`,
            wantStatus: 401,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            headers := map[string]string{}
            // only set auth header for non-401 cases
            if tt.wantStatus != 401 {
                headers["X-User-ID"] = userID
            }
            w := doRequest("POST", "/api/v1/expenses", tt.body, headers)
            if w.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d — body: %s",
                    w.Code, tt.wantStatus, w.Body.String())
            }
        })
    }
}
```

---

## Common Mistakes

### Not importing the routers package

```go
// WRONG — routes are never registered, every request returns 404
package controllers_test

import beego "github.com/beego/beego/v2/server/web"

// CORRECT — blank import triggers init() which registers routes
import _ "myapp/routers"
```

---

### Using the wrong working directory

```go
// WRONG — conf/app.conf not found, config falls back to defaults
func TestMain(m *testing.M) {
    code := m.Run()
    os.Exit(code)
}

// CORRECT — change to project root first
func TestMain(m *testing.M) {
    os.Chdir("../")
    code := m.Run()
    os.Exit(code)
}
```

---

### Forgetting `os.Exit(code)` in `TestMain`

```go
// WRONG — always exits 0 even when tests fail
func TestMain(m *testing.M) {
    os.Chdir("../")
    m.Run()    // result ignored
    os.RemoveAll("testdata")
    // no os.Exit — Go uses default exit 0
}

// CORRECT
func TestMain(m *testing.M) {
    os.Chdir("../")
    code := m.Run()
    os.RemoveAll("testdata")
    os.Exit(code)   // preserve the actual pass/fail result
}
```

---

### Not resetting state between tests

```go
// WRONG — test 2 sees data created by test 1
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // no cleanup — state leaks between cases
        doRequest("POST", "/api/v1/users", tt.body, nil)
    })
}

// CORRECT — decide explicitly: shared state or clean state
// For independent cases:
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        cleanData()  // reset before each case
        doRequest("POST", "/api/v1/users", tt.body, nil)
    })
}

// For sequential cases (e.g. register then duplicate):
cleanData() // reset once before the whole table
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        doRequest("POST", "/api/v1/users", tt.body, nil)
    })
}
```

---

## Running Controller Tests

```bash
# run all tests across all packages
go test ./... -v

# run only controller tests
go test ./controllers/... -v

# run a specific test function
go test ./controllers/... -run TestAuthRegister

# run a specific subtest
go test ./controllers/... -run TestAuthRegister/valid_registration

# with coverage
go test ./... -cover
```

---

## Quick Reference

```go
// Setup in TestMain
func TestMain(m *testing.M) {
    os.Chdir("../")                          // find conf/app.conf
    beego.AppConfig.Set("datadir","testdata") // isolate test data
    code := m.Run()
    os.RemoveAll("testdata")
    os.Exit(code)
}

// Helper — send a request through Beego
func doRequest(method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
    r, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
    r.Header.Set("Content-Type", "application/json")
    for k, v := range headers { r.Header.Set(k, v) }
    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r)
    return w
}

// Helper — parse JSON response
func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
    t.Helper()
    var result map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &result)
    return result
}

// Test structure
tests := []struct {
    name       string
    body       string
    headers    map[string]string
    wantStatus int
    wantMsg    string
}{ ... }

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        w := doRequest("POST", "/path", tt.body, tt.headers)
        if w.Code != tt.wantStatus {
            t.Errorf("status = %d, want %d — body: %s",
                w.Code, tt.wantStatus, w.Body.String())
        }
    })
}
```

---

## Further Reading

- [Go `net/http/httptest` package](https://pkg.go.dev/net/http/httptest)
- [Go `testing` package](https://pkg.go.dev/testing)
- [Go blog — Using subtests](https://go.dev/blog/subtests)
