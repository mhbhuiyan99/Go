Go tutorial: 
1. [A Tour of Go](https://go.dev/tour/list)
2. [Go Basics](https://github.com/gophertuts/go-basics/tree/master)
   
------
## imports
    "fmt"       // For formatted I/O : printing, scanning
    "math/rand" // For random number generation
    "math"      // Provides mathematical functions and constants - no randomness 
------
## Exported Names in Go

In Go, exported names (also called "public" identifiers) are variables, functions, types, or constants that are accessible from outside their package. The rule is simple:

If a name starts with a capital letter, it is exported (public).

If it starts with a lowercase letter, it is unexported (private).
```
package mypkg
const Answer = 42 // Exported (accessible outside mypkg)
const secret = 0  // Unexported (private to mypkg)
```
### Why Go Uses Capitalization
**Simple and explicit:** No need for keywords like public or private.

**Readability:** Easily spot exported names in code (e.g., strings.ToUpper vs strings.toLower).

