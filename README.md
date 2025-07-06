Resources: 
1. [A Tour of Go](https://go.dev/tour/list) + [Golang Backend Development](https://youtube.com/playlist?list=PLpCqPSEm2Xe8sEY2haMDUVgwbkIs5NCJI&si=Krk7dnfCoNG7ulpY)
2. [Exercism Go Track](https://exercism.org/tracks/go)
3. [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests)
4. [Go Basics](https://github.com/gophertuts/go-basics/tree/master)
   
------
## imports
```
    "fmt"       	// For formatted I/O : printing, scanning
    "math/rand" 	// For random number generation
    "math"      	// Provides mathematical functions and constants - no randomness 
    "math/cmplx" 	// Provides mathematical functions for complex numbers
```
Go refuses to compile programs with unused variables or imports, trading short-term convenience for long-term build speed and program clarity.

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
--------
## Variables

The var statement declares a list of variables; as in function argument lists, the type is last.
```
var c, python, java bool
```
### Variables with initializers ###

A var declaration can include initializers, one per variable. ```var i, j int = 1, 2```
If an initializer is present, the type can be omitted; the variable will take the type of the initializer. ```var c, python, java = true, false, "no!"```

Inside a function, the ```:=``` short assignment statement can be used in place of a var declaration with implicit type.
```
func main() {
    x := 42          // x is declared as an int
    name := "Alice"   // name is declared as a string
    y := 3.14        // y is declared as a float64
}
```
The ```:=``` syntax is not allowed at package scope (outside any function).

### Basic types

```
bool
string
int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr
byte // alias for uint8
rune // alias for int32
     // represents a Unicode code point
float32 float64
complex64 complex128
```
#### Zero Values
Variables declared without an explicit initial value are given their zero value. <br>
The zero value is:
```
0 for numeric types,
false for the boolean type, and
"" (the empty string) for strings.
```
### Variable Shadowing
Variable shadowing happens when a variable declared within a narrower scope (like inside a function or block) has the same name as a variable in an outer scope. The inner variable "shadows" or hides the outer one within its scope.
```
package main
import "fmt"
var x = 10
func main(){
	done := false
	if !done {
		x := 30
		fmt.Println(x)
	}
	fmt.Println(x)
}
```

--------
### Packages and Scopes

[Understanding Packages and Scopes in Golang](https://medium.com/@mhbhuiyan10023/understanding-packages-and-scopes-in-golang-fc0b11d65001)
