Link : https://go.dev/tour/basics/1

package main

import (
    "fmt"       // For formatted I/O (printing, scanning)
    "math/rand" // For random number generation
)

func main() {
	fmt.Println("My favorite number is", rand.Intn(10))
}
