
**Basic:** [Go Concurrency: How Goroutines and Channels Work](https://medium.com/gitconnected/go-concurrency-how-goroutines-and-channels-work-4ee2c5a2e045)

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

## Exercise 1.3: "Parallel Squares"
Write a program that:
1. Creates []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
2. Spawns one goroutine per number to compute its square
3. Sends each square on a single results channel
4. main collects exactly 9 results and prints them <br>

Constraints:
1. No sync.WaitGroup
2. Must close the results channel so range works
3. Must not leak goroutines

**Hint:** You need a way to know when all 9 goroutines are done so you can close the channel. Since you can't use WaitGroup, think about counting.
