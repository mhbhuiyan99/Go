Resources: 
1. [A Tour of Go](https://go.dev/tour/list)
2. [Exercism Go Track](https://exercism.org/tracks/go)
3. [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests)
   
Others:
1. [Learn X in Y minutes](https://learnxinyminutes.com/go/)
2. [Go by Example](https://gobyexample.com/)
3. [Go Basics](https://github.com/gophertuts/go-basics/tree/master)
   
------
## imports
    "fmt"       // For formatted I/O : printing, scanning
    "math/rand" // For random number generation
    "math"      // Provides mathematical functions and constants - no randomness 
------
## Exported Names in Go

In Go, exported names (also called "public" identifiers) are variables, functions, types, or constants that are accessible from outside their package. The rule is simple:

If a name starts with a capital letter, it is exported (public).<br>
If it starts with a lowercase letter, it is unexported (private).
```
package mypkg
const Answer = 42 // Exported (accessible outside mypkg)
const secret = 0  // Unexported (private to mypkg)
```
### Why Go Uses Capitalization
**Simple and explicit:** No need for keywords like public or private.<br>
**Readability:** Easily spot exported names in code (e.g., strings.ToUpper vs strings.toLower).
------
## Function

Type comes after the variable name.
```
func add(x int, y int) int { // func add(x, y int) int {
	return x + y;
}
```
A function can return any number of results.
```
func swap(x, y string) (string, string) {
	return y, x
}
```
### Named return values
Go's return values may be named. If so, they are treated as variables defined at the top of the function.<br>
These names should be used to document the meaning of the return values.<br>
A return statement without arguments returns the named return values. This is known as a "naked" return.<br>
Naked return statements should be used only in short functions, as with the example shown here. They can harm readability in longer functions.
```
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}
```

