
# Phase 1: Feel the Goroutine
## Exercise 1.1: "Where's My Output?"
Run this code. Predict the output. Then fix it so it always prints "done" before the program exits.
```
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("done")
    }()
    // your main goroutine exits immediately
}
```
**Questions to answer before fixing:**
1. Why doesn't it print sometimes?
2. If you add time.Sleep(2 * time.Second) at the end, does that "fix" it? Is that a good fix?
3. What's the idiomatic way to wait?

## Exercise 1.2: "The Channel Handshake"
Write a program where:
1. `main` creates an unbuffered string channel.
2. `main` spawns a goroutine.
3. The goroutine sends "hello from goroutine" on the channel.
4. `main` receives from the channel and prints it.

```
package main

import "fmt"

func main() {
    // 1. create an unbuffered string channel
    // 2. spawn a goroutine that sends "hello from goroutine" on it
    // 3. main receives and prints it
}
```

Then answer:
1. What happens if you swap the send and receive order (goroutine receives, main sends)?
2. What happens if you remove the goroutine and do both send and receive in main?

