# Table-Driven Tests in Go

> **Who this is for:** Go developers who want to write clean, maintainable
> tests using the table-driven pattern — the standard testing style in Go.

---

## What Is a Table-Driven Test?

A table-driven test is a pattern where you define all your test cases as a
slice of structs, then loop over them with a single test function. Instead
of writing one function per scenario, you write one function and one row per
scenario.

**Without table-driven tests:**

```go
func TestAdd_TwoPositives(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}

func TestAdd_NegativeNumbers(t *testing.T) {
    result := Add(-1, -2)
    if result != -3 {
        t.Errorf("expected -3, got %d", result)
    }
}

func TestAdd_Zero(t *testing.T) {
    result := Add(0, 5)
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

**With table-driven tests:**

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {name: "two positives",    a: 2,  b: 3,  want: 5},
        {name: "negative numbers", a: -1, b: -2, want: -3},
        {name: "zero",             a: 0,  b: 5,  want: 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

Same coverage, far less repetition. Adding a new case is one line.

---

## The Anatomy of a Table-Driven Test

```go
func TestFunctionName(t *testing.T) {
    // 1. Define the test table
    tests := []struct {
        name    string   // descriptive name for this case
        input   string   // inputs to the function under test
        want    string   // expected output
        wantErr bool     // whether an error is expected
    }{
        // 2. Define each test case as a row
        {
            name:    "valid input",
            input:   "hello",
            want:    "HELLO",
            wantErr: false,
        },
        {
            name:    "empty input",
            input:   "",
            want:    "",
            wantErr: false,
        },
    }

    // 3. Loop and run each case as a subtest
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 4. Call the function under test
            got, err := ToUpper(tt.input)

            // 5. Assert
            if tt.wantErr && err == nil {
                t.Errorf("expected error, got nil")
                return
            }
            if !tt.wantErr && err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            if got != tt.want {
                t.Errorf("ToUpper(%q) = %q, want %q", tt.input, got, tt.want)
            }
        })
    }
}
```

---

## `t.Run` — Subtests

`t.Run(name, func)` runs a named subtest. Each subtest:

- Has its own pass/fail status
- Appears separately in test output
- Can be run individually with `-run`

```
=== RUN   TestAdd
=== RUN   TestAdd/two_positives
--- PASS: TestAdd/two_positives (0.00s)
=== RUN   TestAdd/negative_numbers
--- PASS: TestAdd/negative_numbers (0.00s)
=== RUN   TestAdd/zero
--- PASS: TestAdd/zero (0.00s)
--- PASS: TestAdd (0.00s)
```

---

## Naming Test Cases

Test case names become part of the subtest name in output.
Use descriptive names that read like sentences:

```go
// Good — explains what the input is and what should happen
{name: "empty string returns error"}
{name: "valid email passes validation"}
{name: "negative amount is rejected"}
{name: "duplicate email returns conflict"}

// Bad — vague, not useful when a test fails
{name: "test1"}
{name: "case2"}
{name: "invalid"}
```

Spaces in names are converted to underscores in output:
`"empty string returns error"` → `TestValidate/empty_string_returns_error`

---

## Testing Functions That Return Errors

This is the most common pattern. Define `wantErr bool` and optionally
`errMsg string` to check the exact error message:

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid division",
            a: 10, b: 2,
            want:    5.0,
            wantErr: false,
        },
        {
            name:    "divide by zero returns error",
            a: 10, b: 0,
            wantErr: true,
            errMsg:  "division by zero",
        },
        {
            name: "negative dividend",
            a: -6, b: 2,
            want:    -3.0,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error %q, got nil", tt.errMsg)
                    return
                }
                if tt.errMsg != "" && err.Error() != tt.errMsg {
                    t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
                }
                return
            }

            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            if got != tt.want {
                t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

---

## `t.Error` vs `t.Fatal` vs `t.Fatalf`

| Function | What it does |
|---|---|
| `t.Error(msg)` | Marks test as failed, continues running |
| `t.Errorf(fmt, ...)` | Like `t.Error` but with formatting |
| `t.Fatal(msg)` | Marks test as failed, **stops the current test immediately** |
| `t.Fatalf(fmt, ...)` | Like `t.Fatal` but with formatting |

```go
// Use t.Errorf when you want to collect multiple failures
if got.Name != tt.wantName {
    t.Errorf("Name = %q, want %q", got.Name, tt.wantName)
}
if got.Age != tt.wantAge {
    t.Errorf("Age = %d, want %d", got.Age, tt.wantAge)  // still runs even if Name failed
}

// Use t.Fatalf when continuing makes no sense
result, err := SomeSetup()
if err != nil {
    t.Fatalf("setup failed: %v", err)  // no point continuing if setup broke
}
// use result here safely
```

**Rule of thumb:**
- Setup failures → `t.Fatal` (rest of the test would be meaningless)
- Assertion failures → `t.Error` (show all failures at once)

---

## `t.Helper()` — Cleaner Error Output

When a failure occurs inside a helper function, Go normally points to the
line inside the helper, not the line in the test that called it.
`t.Helper()` fixes that — it marks the function as a helper so error
output points to the caller instead.

```go
// Without t.Helper() — error points to line inside assertEqual, not useful
func assertEqual(t *testing.T, got, want int) {
    if got != want {
        t.Errorf("got %d, want %d", got, want)  // line reported here
    }
}

// With t.Helper() — error points to the test case that called assertEqual
func assertEqual(t *testing.T, got, want int) {
    t.Helper()  // ← add this
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}
```

Always add `t.Helper()` to any helper function that calls `t.Error` or `t.Fatal`.

---

## The `mustX` Pattern — Setup Helpers That Fail Fast

When setting up test data, use helpers that call `t.Fatalf` on failure.
This stops the test immediately with a clear message if setup breaks,
rather than letting a nil pointer or empty result cause a confusing failure
later in the actual assertion.

```go
// mustCreateUser creates a user and stops the test immediately if it fails.
func mustCreateUser(t *testing.T, name, email string) *User {
    t.Helper()
    user, err := CreateUser(name, email)
    if err != nil {
        t.Fatalf("mustCreateUser: %v", err)
    }
    return user
}

// Use in tests
func TestGetUser(t *testing.T) {
    created := mustCreateUser(t, "Alice", "alice@example.com")
    // if mustCreateUser fails, we stop here with a clear message
    // instead of crashing later when accessing created.ID on a nil pointer

    user, err := GetUserByID(created.ID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Email != "alice@example.com" {
        t.Errorf("email = %q, want %q", user.Email, "alice@example.com")
    }
}
```

---

## Shared State Between Test Cases

By default, all cases in a table test share the same state. Sometimes
that is intentional (each case builds on the previous), sometimes it is not.

**Intentional — cases run in sequence, state persists:**

```go
func TestRegisterDuplicateEmail(t *testing.T) {
    tests := []struct {
        name       string
        email      string
        wantErr    bool
    }{
        {name: "first registration succeeds",   email: "alice@test.com", wantErr: false},
        {name: "duplicate email fails",         email: "alice@test.com", wantErr: true},
        {name: "different email succeeds",      email: "bob@test.com",   wantErr: false},
    }

    cleanTestData()  // clean once before all cases

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Register(tt.email)
            if tt.wantErr && err == nil {
                t.Errorf("expected error, got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

**Independent — each case starts clean:**

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        cleanTestData()  // reset before each case
        // now each case is independent
    })
}
```

Be explicit about which approach you are using. Mixing them accidentally
is a common source of flaky tests.

---

## `TestMain` — Package-Level Setup and Teardown

`TestMain` runs once before all tests in the package. Use it for:
- Changing working directory
- Setting configuration values
- Creating shared fixtures
- Cleaning up after all tests finish

```go
func TestMain(m *testing.M) {
    // setup — runs once before all tests
    os.Chdir("../")                      // change working directory
    os.MkdirAll("testdata", 0755)        // create test data directory

    // run all tests
    code := m.Run()

    // teardown — runs once after all tests
    os.RemoveAll("testdata")             // clean up

    os.Exit(code)
}
```

`m.Run()` returns an exit code — `0` for pass, non-zero for fail.
You must call `os.Exit(code)` at the end, not `return`.
Without it, the teardown runs but the exit code is always 0 (always passes).

---

## Running Tests

```bash
# run all tests in the current package
go test .

# run all tests in all packages
go test ./...

# run with verbose output (shows each subtest)
go test ./... -v

# run with coverage
go test ./... -cover

# run a specific test function
go test -run TestAdd

# run a specific subtest
go test -run TestAdd/two_positives

# run all subtests matching a pattern
go test -run TestAdd/negative
```

---

## Coverage

```bash
# show coverage percentage
go test ./... -cover

# generate a coverage report file
go test ./... -coverprofile=coverage.out

# view coverage in browser (highlights covered lines in green/red)
go tool cover -html=coverage.out
```

Coverage shows what percentage of your code is exercised by tests.
The `-html` report lets you see exactly which lines are not covered,
which makes it easy to identify what tests to add next.

---

## Complete Example

```go
package math_test

import (
    "testing"
    "errors"
)

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
        errMsg  string
    }{
        {
            name: "positive numbers",
            a: 10, b: 2, want: 5,
        },
        {
            name: "negative dividend",
            a: -6, b: 2, want: -3,
        },
        {
            name: "decimal result",
            a: 1, b: 3, want: 0.3333333333333333,
        },
        {
            name:    "divide by zero returns error",
            a: 10, b: 0,
            wantErr: true,
            errMsg:  "division by zero",
        },
        {
            name: "zero dividend",
            a: 0, b: 5, want: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)

            if tt.wantErr {
                if err == nil {
                    t.Fatalf("expected error %q, got nil", tt.errMsg)
                }
                if err.Error() != tt.errMsg {
                    t.Errorf("error = %q, want %q", err.Error(), tt.errMsg)
                }
                return
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tt.want {
                t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

---

## Quick Reference

```go
// Basic structure
func TestXxx(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {name: "case description", input: ..., want: ...},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionUnderTest(tt.input)
            if tt.wantErr {
                if err == nil { t.Errorf("expected error, got nil") }
                return
            }
            if err != nil { t.Fatalf("unexpected error: %v", err) }
            if got != tt.want { t.Errorf("got %v, want %v", got, tt.want) }
        })
    }
}
```

| Concept | Usage |
|---|---|
| `t.Run(name, fn)` | Run a named subtest |
| `t.Error / t.Errorf` | Fail but continue |
| `t.Fatal / t.Fatalf` | Fail and stop immediately |
| `t.Helper()` | Mark function as helper for cleaner error output |
| `TestMain(m)` | Package-level setup and teardown |
| `m.Run()` | Run all tests, returns exit code |
| `go test -run Name` | Run specific test or subtest |
| `go test -cover` | Show coverage percentage |

---

## Further Reading

- [Go testing package](https://pkg.go.dev/testing)
- [Go blog — Using subtests and sub-benchmarks](https://go.dev/blog/subtests)
- [Go test flags reference](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
