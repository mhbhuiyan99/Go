
## Exercise 1.1:

### Option A: `sync.WaitGroup`
You wait for completion, not time.
```
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
```
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
```
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
