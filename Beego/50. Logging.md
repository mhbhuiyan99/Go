# Beego Logging — `logs.Info`, `logs.Warn`, `logs.Error` and When to Use Each

> **Who this is for:** Go developers who want to replace `fmt.Println`
> with Beego's structured logging system, and understand when to use
> each log level.

---

## A Note on the Import Path

Beego's logging package is **not** accessed through `beego.Info()`.
It lives in a separate package:

```go
import "github.com/beego/beego/v2/core/logs"

logs.Info("Server started")
logs.Warn("Slow response")
logs.Error("Failed to connect")
```

If you see `beego.Info()` referenced anywhere, that is the old Beego v1
import path (`github.com/astaxie/beego`) and does not apply to v2 projects.

---

## Why Not `fmt.Println`?

`fmt.Println` writes plain text with no structure:

```go
fmt.Println("User registered:", email)
```

Output:
```
User registered: alice@example.com
```

Beego's logger adds structure automatically — timestamp, log level,
and source file/line:

```go
logs.Info("User registered: %s", email)
```

Output:
```
2025/06/01 10:30:00.123 [I] [auth.go:42] User registered: alice@example.com
```

The extra structure matters once your application has real traffic:
- **Timestamp** — when did this happen?
- **Level** — how serious is this?
- **File/line** — which code produced this log?

`fmt.Println` gives you none of that. You would have to add it manually
to every single print statement, which nobody does consistently.

---

## The Four Log Levels

```go
logs.Debug("Parsing row: %v", row)
logs.Info("Server started on port %d", 8080)
logs.Warn("Slow response: %v", duration)
logs.Error("Failed to write CSV: %v", err)
```

### `logs.Debug` — fine-grained development detail

Use for information that is only useful while actively developing or
debugging a specific issue. Too noisy for normal operation.

```go
logs.Debug("Parsed CSV row: %+v", row)
logs.Debug("Cache hit for key: %s", key)
logs.Debug("SQL query: %s", query)
```

### `logs.Info` — normal operations

Use for events that are expected and routine. Confirms the application
is working as intended.

```go
logs.Info("Server started on port %d", 8080)
logs.Info("User registered: %s", email)
logs.Info("Expense created: id=%d amount=%.2f", expense.ID, expense.Amount)
```

### `logs.Warn` — something unexpected but not breaking

Use when something unusual happened, the application handled it, and
execution continues normally. Worth investigating but not urgent.

```go
logs.Warn("Skipping invalid CSV row: %v", err)
logs.Warn("Partner API slow: %v", elapsed)
logs.Warn("Retrying failed request, attempt %d", attempt)
```

### `logs.Error` — something failed

Use when an operation could not complete and the failure has real
consequences — a request cannot be served, data cannot be saved.

```go
logs.Error("Failed to write CSV: %v", err)
logs.Error("Database connection lost: %v", err)
logs.Error("Failed to parse JSON: %v", err)
```

---

## The Decision Tree

```
Did the operation succeed as expected?
│
├── Yes, and it's routine           → logs.Info
├── Yes, but only useful in dev     → logs.Debug
│
Did something go wrong?
│
├── Handled gracefully, continuing   → logs.Warn
└── Could not recover, operation failed → logs.Error
```

---

## Practical Examples from a Typical API

```go
func (c *ExpenseController) Create() {
    userID, ok := c.getUserID()
    if !ok {
        return // AuthFilter already logged this
    }

    var req models.CreateExpenseRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
        // client sent bad data — not a server problem, no log needed
        // or, if you want visibility into malformed requests:
        logs.Warn("[Create] invalid JSON from user %d: %v", userID, err)
        c.RespondError(400, "Invalid JSON")
        return
    }

    expense, err := services.CreateExpense(userID, req)
    if err != nil {
        // could be validation (client's fault) or a real failure
        logs.Warn("[Create] failed for user %d: %v", userID, err)
        c.RespondError(400, err.Error())
        return
    }

    // success — routine confirmation
    logs.Info("[Create] expense %d created for user %d", expense.ID, userID)
    c.RespondCreated(expense)
}
```

```go
func GetAllExpenses() ([]Expense, error) {
    f, err := os.OpenFile(path, os.O_RDONLY, 0644)
    if err != nil {
        if os.IsNotExist(err) {
            return []Expense{}, nil // expected on first run — no log needed
        }
        // unexpected failure opening the file — this is serious
        logs.Error("[GetAllExpenses] failed to open file: %v", err)
        return nil, err
    }
    defer f.Close()

    rows, err := csv.NewReader(f).ReadAll()
    if err != nil {
        logs.Error("[GetAllExpenses] failed to read CSV: %v", err)
        return nil, err
    }

    var expenses []Expense
    for _, row := range rows[1:] {
        expense, err := rowToExpense(row)
        if err != nil {
            // one bad row shouldn't stop the whole read — skip and warn
            logs.Warn("[GetAllExpenses] skipping invalid row: %v", err)
            continue
        }
        expenses = append(expenses, expense)
    }
    return expenses, nil
}
```

---

## Validation Errors — Usually Not Worth Logging

A common mistake is logging every validation failure as a `Warn` or `Error`.
Client validation failures (missing fields, bad formats) are normal and
expected — they happen constantly and are not actionable for you as the
developer. Logging every one creates noise that buries real problems.

```go
// Usually unnecessary — this is normal client behavior, not a system issue
if req.Amount <= 0 {
    logs.Warn("invalid amount from client")  // noisy, not useful
    return fmt.Errorf("amount must be greater than zero")
}

// Better — just return the error, let the controller respond
if req.Amount <= 0 {
    return fmt.Errorf("amount must be greater than zero")
}
```

**Reserve logging for things a developer needs to know about** — file
I/O failures, unexpected panics, external service failures, and other
issues that indicate something is actually broken in the system itself.

---

## Prefixing Logs by Context

Use a bracketed prefix to identify where a log came from. This becomes
essential once you have many endpoints logging simultaneously:

```go
logs.Error("[CreateExpense] %v", err)
logs.Info("[Login] user authenticated: %s", email)
logs.Warn("[CSV] skipping malformed row: %v", err)
```

Output:
```
2025/06/01 10:30:00 [E] [CreateExpense] failed to save expense: disk full
2025/06/01 10:30:01 [I] [Login] user authenticated: alice@example.com
2025/06/01 10:30:02 [W] [CSV] skipping malformed row: invalid amount
```

Without prefixes, when ten different features log errors, it becomes
impossible to tell them apart at a glance.

---

## Configuring Log Level in `app.conf`

You can control which levels actually get printed. This is useful in
production — you often want `Info` and above, but not `Debug`.

```ini
# conf/app.conf
loglevel = 6
```

| Value | Level | Meaning |
|---|---|---|
| 7 | Debug | show everything, including debug details |
| 6 | Info | show info, warnings, and errors (typical default) |
| 4 | Warning | show only warnings and errors |
| 3 | Error | show only errors |

Setting `loglevel = 6` means `logs.Debug()` calls are silently suppressed
without needing to remove them from your code. This is why it is safe to
leave `logs.Debug` calls in place — they simply do not print unless you
lower the threshold during active debugging.

---

## Logging to a File Instead of the Console

By default, Beego logs to stdout (the console). For production, you
often want logs written to a file:

```ini
# conf/app.conf
LogAdapter = file
```

Or configure it directly in code at startup:

```go
// main.go
import "github.com/beego/beego/v2/core/logs"

func main() {
    logs.SetLogger(logs.AdapterFile, `{"filename":"logs/app.log"}`)
    // ... rest of main
}
```

---

## `logs.Emergency`, `logs.Alert`, `logs.Critical` — Rarely Needed

Beego actually supports more levels than the four commonly used:

```go
logs.Emergency("...")  // system is unusable
logs.Alert("...")      // action must be taken immediately
logs.Critical("...")   // critical conditions
logs.Error("...")      // error conditions
logs.Warning("...")    // warning conditions (alias: logs.Warn)
logs.Notice("...")     // normal but significant condition
logs.Informational("...") // informational messages (alias: logs.Info)
logs.Debug("...")      // debug-level messages
```

For almost all application code, `Debug`, `Info`, `Warn`, and `Error`
cover every real scenario. The others (`Emergency`, `Alert`, `Critical`)
are typically reserved for infrastructure-level monitoring systems, not
day-to-day application logic.

---

## Common Mistakes

### Using `fmt.Println` in production code

```go
// WRONG — no timestamp, no level, cannot be filtered or redirected
fmt.Println("Error creating expense:", err)

// CORRECT
logs.Error("[CreateExpense] %v", err)
```

---

### Logging sensitive data

```go
// WRONG — password ends up in log files
logs.Info("Login attempt: email=%s password=%s", email, password)

// CORRECT — never log passwords, tokens, or secrets
logs.Info("Login attempt: email=%s", email)
```

---

### Using `Error` for expected client mistakes

```go
// WRONG — this happens constantly and is not a system problem
logs.Error("validation failed: title is required")

// CORRECT — just return the error to the controller, no log needed
return fmt.Errorf("title is required")
```

---

### Forgetting the import

```go
// WRONG import — this is the old Beego v1 path
import "github.com/astaxie/beego"
beego.Info("...")

// CORRECT import for Beego v2
import "github.com/beego/beego/v2/core/logs"
logs.Info("...")
```

---

## Quick Reference

```go
import "github.com/beego/beego/v2/core/logs"

logs.Debug("Detail for developers: %v", data)     // dev-only detail
logs.Info("Normal operation: %s", "user logged in") // expected events
logs.Warn("Recovered from issue: %v", err)          // handled, continuing
logs.Error("Operation failed: %v", err)             // real failure
```

| Level | Use for | Example |
|---|---|---|
| `Debug` | Fine-grained dev detail | Parsing steps, cache hits |
| `Info` | Normal, expected events | Server started, user registered |
| `Warn` | Handled but unusual | Skipped bad row, slow response |
| `Error` | Failed operations | File write failed, DB down |

```ini
# conf/app.conf — control what gets printed
loglevel = 6   # 7=Debug 6=Info 4=Warning 3=Error
```

---

## Further Reading

- [Beego Logging Docs](https://beego.wiki/docs/module/logs/)
- [Beego `logs` package source](https://github.com/beego/beego/tree/develop/core/logs)
