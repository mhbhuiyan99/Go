# CSV File I/O in Go — `encoding/csv`, `os.OpenFile`, and Common Patterns

> **Who this is for:** Go developers who want to read from and write to CSV
> files using the standard library. No third-party packages required.

---

## What `encoding/csv` Provides

Go's standard library includes `encoding/csv` for reading and writing
comma-separated values. It handles quoting, escaping, and multi-line fields
automatically — you work with `[]string` slices, not raw text.

```go
import (
    "encoding/csv"
    "os"
)
```

There are two main types:
- `csv.Reader` — reads CSV data row by row or all at once
- `csv.Writer` — writes rows to a CSV file

Both wrap an `io.Reader` or `io.Writer` — meaning they work with files,
network connections, buffers, or anything that implements those interfaces.

---

## Understanding CSV Structure in Go

Every CSV file is a sequence of **rows**. Every row is a slice of strings.

```
id,name,email                       ← header row  →  []string{"id","name","email"}
1,Alice,alice@example.com           ← data row    →  []string{"1","Alice","alice@example.com"}
2,Bob,bob@example.com               ← data row    →  []string{"2","Bob","bob@example.com"}
```

`encoding/csv` never parses values beyond strings. Converting `"1"` to `int`
or `"350.50"` to `float64` is always your responsibility using `strconv`.

---

## `os.OpenFile` — Opening Files with the Right Flags

Before reading or writing, you need a file handle. `os.OpenFile` gives you
precise control over how the file is opened.

```go
f, err := os.OpenFile("data/users.csv", flags, permissions)
```

### The flags

Flags are combined with the bitwise OR operator `|`:

| Flag | Meaning |
|---|---|
| `os.O_RDONLY` | Open for reading only |
| `os.O_WRONLY` | Open for writing only |
| `os.O_RDWR` | Open for reading and writing |
| `os.O_CREATE` | Create the file if it does not exist |
| `os.O_APPEND` | Append to the file instead of overwriting |
| `os.O_TRUNC` | Truncate (clear) the file on open |

### The permissions argument

The third argument is the Unix file permission, used only when creating
a new file. `0644` is the standard for readable data files:
- Owner: read + write
- Group: read only
- Others: read only

### Common flag combinations

```go
// Read an existing file
f, err := os.OpenFile("data.csv", os.O_RDONLY, 0644)

// Create if not exists, append new rows (for adding records)
f, err := os.OpenFile("data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// Create if not exists, overwrite everything (for full rewrites)
f, err := os.OpenFile("data.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
```

Always close the file when done. Use `defer` immediately after a successful
open so you never accidentally leave a file handle open:

```go
f, err := os.OpenFile("data.csv", os.O_RDONLY, 0644)
if err != nil {
    return err
}
defer f.Close()   // ← always defer right here, before any other logic
```

---

## Reading CSV Files

### `csv.NewReader` + `ReadAll` — read everything at once

```go
func ReadAll(filename string) ([][]string, error) {
    f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    reader := csv.NewReader(f)
    rows, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    return rows, nil
}
```

`ReadAll` returns `[][]string` — a slice of rows, each row being a slice
of string fields. It reads the entire file into memory.

```go
rows, err := ReadAll("data/users.csv")
// rows[0] → []string{"id", "name", "email"}         (header)
// rows[1] → []string{"1", "Alice", "alice@example.com"}
// rows[2] → []string{"2", "Bob",   "bob@example.com"}
```

### Skipping the header row

The header is just `rows[0]`. Skip it by slicing:

```go
rows, err := reader.ReadAll()
if err != nil {
    return nil, err
}

if len(rows) <= 1 {
    return [][]string{}, nil  // file is empty or header-only
}

dataRows := rows[1:]  // skip the header
```

### `Read` — read one row at a time

Use `Read` when the file is large and you do not want to load it all into
memory at once:

```go
reader := csv.NewReader(f)

// Skip header
// If the file is completely empty, Read returns io.EOF immediately.
// Treat that as "no data" rather than a real error.
if _, err := reader.Read(); err != nil {
    if err == io.EOF {
        return nil  // empty file — nothing to process, not a failure
    }
    return err      // real read error
}

for {
    row, err := reader.Read()
    if err == io.EOF {
        break        // end of file — normal exit
    }
    if err != nil {
        return err   // actual error
    }
    // process row
    fmt.Println(row[0], row[1])
}
```

`Read` returns `io.EOF` when there are no more rows. This is not an error —
it is the normal signal that the file has been fully read. The same `io.EOF`
is also returned immediately on the header skip if the file is completely
empty (0 bytes), so that case must be handled separately before the loop.

---

## Writing CSV Files

### `csv.NewWriter` + `Write` + `Flush`

```go
func AppendRow(filename string, row []string) error {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    writer := csv.NewWriter(f)

    if err := writer.Write(row); err != nil {
        return err
    }

    writer.Flush()  // ← required — writes buffered data to the file

    return writer.Error()  // check if Flush itself encountered an error
}
```

### Why `Flush` is required

`csv.Writer` buffers writes internally for efficiency. `Write` puts data
into the buffer — it does not write to the file immediately. `Flush` pushes
all buffered data to the underlying file.

**If you forget `Flush`, your data is silently lost** — no error, no panic,
just an empty or incomplete file.

```go
writer := csv.NewWriter(f)
writer.Write([]string{"1", "Alice", "alice@example.com"})
writer.Write([]string{"2", "Bob",   "bob@example.com"})
// Without Flush — nothing is written to disk
writer.Flush()
// Now both rows are written
```

### Writing multiple rows at once — `WriteAll`

```go
rows := [][]string{
    {"id", "name", "email"},
    {"1", "Alice", "alice@example.com"},
    {"2", "Bob",   "bob@example.com"},
}

writer := csv.NewWriter(f)
if err := writer.WriteAll(rows); err != nil {
    return err
}
// WriteAll calls Flush internally — no separate Flush needed
```

`WriteAll` is convenient for writing an entire dataset in one call.
It calls `Flush` for you.

---

## Mapping CSV Rows to Structs

`encoding/csv` works only with `[]string`. You must convert between
string slices and your own structs manually.

### Define the struct

```go
type User struct {
    ID        int
    Name      string
    Email     string
    CreatedAt string
}
```

### CSV row → struct (deserialize)

```go
import (
    "fmt"
    "strconv"
)

func rowToUser(row []string) (User, error) {
    if len(row) < 4 {
        return User{}, fmt.Errorf("invalid row: expected 4 fields, got %d", len(row))
    }

    id, err := strconv.Atoi(row[0])
    if err != nil {
        return User{}, fmt.Errorf("invalid id %q: %w", row[0], err)
    }

    return User{
        ID:        id,
        Name:      row[1],
        Email:     row[2],
        CreatedAt: row[3],
    }, nil
}
```

### Struct → CSV row (serialize)

```go
func userToRow(u User) []string {
    return []string{
        strconv.Itoa(u.ID),
        u.Name,
        u.Email,
        u.CreatedAt,
    }
}
```

### Reading all users from a file

```go
func GetAllUsers(filename string) ([]User, error) {
    f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
    if err != nil {
        if os.IsNotExist(err) {
            return []User{}, nil  // file not yet created — return empty slice
        }
        return nil, err
    }
    defer f.Close()

    rows, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return nil, err
    }

    if len(rows) <= 1 {
        return []User{}, nil  // empty or header-only
    }

    var users []User
    for _, row := range rows[1:] {  // skip header
        user, err := rowToUser(row)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}
```

### Appending a new record

```go
func CreateUser(filename string, u User) error {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    writer := csv.NewWriter(f)
    if err := writer.Write(userToRow(u)); err != nil {
        return err
    }
    writer.Flush()
    return writer.Error()
}
```

---

## The Header Problem — Writing It Only Once

When using `O_APPEND|O_CREATE`, the file is created empty on first run.
You need to write the header row once, but never again on subsequent runs.

The simplest approach: check the file size. If it is zero, write the header.

```go
func CreateUser(filename string, u User) error {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    writer := csv.NewWriter(f)

    // Write header only when file is new (empty)
    info, err := f.Stat()
    if err != nil {
        return err
    }
    if info.Size() == 0 {
        if err := writer.Write([]string{"id", "name", "email", "created_at"}); err != nil {
            return err
        }
    }

    if err := writer.Write(userToRow(u)); err != nil {
        return err
    }
    writer.Flush()
    return writer.Error()
}
```

---

## The Rewrite Pattern — Update and Delete

CSV has no concept of editing a row in place. To update or delete a record,
you must:

1. Read all rows into memory
2. Modify the slice (change or remove the target row)
3. Write all rows back to the file (overwrite from scratch)

```
Read all → Modify in memory → Write all back
```

### Update a record

```go
func UpdateUser(filename string, updated User) error {
    // Step 1 — read everything
    users, err := GetAllUsers(filename)
    if err != nil {
        return err
    }

    // Step 2 — find and replace
    found := false
    for i, u := range users {
        if u.ID == updated.ID {
            users[i] = updated
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("user with ID %d not found", updated.ID)
    }

    // Step 3 — write everything back
    return writeAllUsers(filename, users)
}
```

### Delete a record

```go
func DeleteUser(filename string, id int) error {
    // Step 1 — read everything
    users, err := GetAllUsers(filename)
    if err != nil {
        return err
    }

    // Step 2 — rebuild slice without the target row
    filtered := make([]User, 0, len(users))
    for _, u := range users {
        if u.ID != id {
            filtered = append(filtered, u)
        }
    }

    // Step 3 — write everything back
    return writeAllUsers(filename, filtered)
}
```

### The shared write-back helper

Both update and delete use the same write-back function:

```go
func writeAllUsers(filename string, users []User) error {
    // O_TRUNC clears the file before writing
    f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    writer := csv.NewWriter(f)

    // Write header
    if err := writer.Write([]string{"id", "name", "email", "created_at"}); err != nil {
        return err
    }

    // Write all rows
    for _, u := range users {
        if err := writer.Write(userToRow(u)); err != nil {
            return err
        }
    }

    writer.Flush()
    return writer.Error()
}
```

The key difference here is `O_TRUNC` — it clears the file on open so you
write fresh from the beginning. Without it, the new content would be appended
on top of the old content.

---

## Generating the Next ID

Since there is no auto-increment in CSV, you calculate the next ID by
finding the current maximum:

```go
func GetNextID(filename string) (int, error) {
    users, err := GetAllUsers(filename)
    if err != nil {
        return 0, err
    }

    maxID := 0
    for _, u := range users {
        if u.ID > maxID {
            maxID = u.ID
        }
    }
    return maxID + 1, nil
}
```

Using `maxID + 1` instead of `len(users) + 1` is important. If a record is
deleted, `len(users)` could produce a duplicate ID. Taking the max always
produces a new unique ID.

---

## Common Mistakes

### Forgetting `Flush`

```go
writer := csv.NewWriter(f)
writer.Write(row)
// No Flush — data is in the buffer and never reaches the file
// File appears empty or unchanged
```

Always call `writer.Flush()` after your last `Write` call.

---

### Not checking `writer.Error()` after `Flush`

`Flush` does not return an error directly. It stores it internally.
Check it with `writer.Error()`:

```go
writer.Flush()
if err := writer.Error(); err != nil {
    return err
}
```

Or check both in one step:

```go
writer.Flush()
return writer.Error()  // nil if everything succeeded
```

---

### Using `O_APPEND` for rewrites

```go
// WRONG for rewrite — appends new content on top of old
f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// CORRECT for rewrite — clears first, then writes fresh
f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
```

---

### Not handling `os.IsNotExist`

If the file has never been created yet, `os.OpenFile` with `O_RDONLY` returns
an error. Handle it explicitly instead of treating it as a real failure:

```go
f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
if err != nil {
    if os.IsNotExist(err) {
        return []User{}, nil   // no file yet — return empty, not an error
    }
    return nil, err            // real error
}
```

---

### Off-by-one when skipping the header

```go
rows, _ := reader.ReadAll()

// WRONG — rows[0] is the header, this processes it as data
for _, row := range rows { ... }

// CORRECT — skip the first row
for _, row := range rows[1:] { ... }
```

---

## Quick Reference

```go
// Open for reading
f, err := os.OpenFile("data.csv", os.O_RDONLY, 0644)

// Open for appending (adding records)
f, err := os.OpenFile("data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// Open for rewriting (update / delete)
f, err := os.OpenFile("data.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

// Always
defer f.Close()

// Read all rows
rows, err := csv.NewReader(f).ReadAll()
dataRows := rows[1:]  // skip header

// Read one row at a time
reader := csv.NewReader(f)
for {
    row, err := reader.Read()
    if err == io.EOF { break }
    if err != nil { return err }
}

// Write rows
writer := csv.NewWriter(f)
writer.Write([]string{"field1", "field2"})
writer.Flush()
err = writer.Error()

// Write all at once (Flush included)
writer.WriteAll(rows)
```

| Goal | Flags | Notes |
|---|---|---|
| Read file | `O_RDONLY` | File must exist |
| Append new record | `O_APPEND\|O_CREATE\|O_WRONLY` | Creates file if missing |
| Rewrite entire file | `O_TRUNC\|O_CREATE\|O_WRONLY` | Clears before writing |
| Convert field to int | `strconv.Atoi(row[n])` | Always returns `(int, error)` |
| Convert field to float | `strconv.ParseFloat(row[n], 64)` | Always returns `(float64, error)` |
| Convert int to field | `strconv.Itoa(n)` | |
| Convert float to field | `strconv.FormatFloat(f, 'f', 2, 64)` | `'f'` = decimal, `2` = decimal places |

---

## Further Reading

- [Go `encoding/csv` package](https://pkg.go.dev/encoding/csv)
- [Go `os` package — OpenFile](https://pkg.go.dev/os#OpenFile)
- [Go `strconv` package](https://pkg.go.dev/strconv)
