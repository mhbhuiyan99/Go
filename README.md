Resources: 
1. Learn: [A Tour of Go](https://go.dev/tour/list) + [Golang Backend Development](https://youtube.com/playlist?list=PLpCqPSEm2Xe8sEY2haMDUVgwbkIs5NCJI&si=Krk7dnfCoNG7ulpY)
2. Practice: [Exercism Go Track](https://exercism.org/tracks/go)
3. Find something: [Go Basics](https://github.com/gophertuts/go-basics/tree/master)
4. [My Notes on Medium](https://medium.com/@mhbhuiyan)
   
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
---
Output:
30
10
```
### Type conversions
The expression T(v) converts the value v to the type T.<br>

Some numeric conversions:
```
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)
```
Or, put more simply:
```
i := 42
f := float64(i)
u := uint(f)
```
Unlike in C, in Go assignment between items of different type requires an explicit conversion.
```
var x int = 10
var y float64 = float64(x)  // ✅ 
var y float64 = x  // ❌
```
--------
### Packages and Scopes

[Understanding Packages and Scopes in Golang](https://medium.com/@mhbhuiyan10023/understanding-packages-and-scopes-in-golang-fc0b11d65001)

------
## Function Types 
### First Order Function
A first-order function is a regular function that:
1. Does not take another function as an argument
2. Does not return another function

#### 🛡️ Standard Function or Named Function
```
func functionName(parameters) returnType {
    // simply a regular function that has a name
}
```
#### 🛡️ Init Function

The init() function in Go is a special built-in function that is automatically called before the main() function or when a package is imported (you can not call this).<br>

It is commonly used for:
1. Setting up initial values
2. Initializing state
3. Connecting to databases or reading config files
4. Running setup logic before any other function is called
```
package main
import "fmt"

var x = 10
func main(){
	fmt.Println(x)
}
func init(){
	fmt.Println(x)
	x = 20
}
---
Output:
10
20
```
##### Syntax
1. It takes no arguments
2. It returns nothing
3. You can have multiple init() functions in a package (in different files or the same file)

#### 🛡️ Anonymous Function

An anonymous function is a function without a name. 
```
package main

import "fmt"

func main() {
    add := func(a int, b int) int { // anonymous function stored in the variable add
        return a + b
    }

    result := add(5, 7)
    fmt.Println("Sum:", result) // Output: Sum: 12
}
```
#### 🛡️ IIFE (Immediately Invoked Function Expression)

In Go, an IIFE is an ***anonymous function*** that is defined and called immediately.
```
package main

import "fmt"

func main() {
    result := func(x, y int) int {
        return x + y
    }(10, 15) // 👈 Call immediately

    fmt.Println("IIFE Result:", result) // Output: 25
}
```
#### Parameter and Argument
Parameters are placeholders, arguments are real values.
```
func greet(name string) { // "name" is a parameter
    fmt.Println("Hello", name)
}

func main() {
    greet("Mojammel")     // "Mojammel" is an argument
}
```
## Higher-order functions OR First-Class Functions
A function which does at least one of the following
1. takes one or more functions as arguments
2. returns a function as its result

🧪 Example 1: Function as Parameter
```
func calculate(a int, b int, op func(int, int) int) int {
                  // 🔍syntax for assigning a function to a variable >> anonymous function
    return op(a, b)
}

func multiply(x, y int) int {
    return x * y
}

func main() {
    result := calculate(3, 4, multiply)
    fmt.Println("Result:", result) // Output: 12
}
```
🧪 Example 2: Returning a Function
```
func makeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

func main() {
    add5 := makeAdder(5)
    fmt.Println(add5(10)) // Output: 15
}
```
### First-Class Citizen (or First-Class Object)
In programming, a first-class citizen (or first-class object/value) is any ✨entity✨ (Variables, Functions, Structs, Slices, Maps, etc.) that can be used like any other value — i.e., assigned to variables, passed as arguments, returned from functions, and stored in data structures.

### First-Class Functions
Functions are treated like values.<br>
A first-class function is a function ( = ✨entity✨ ) that is treated as a first-class citizen in a programming language.

1. Assign functions to variables
2. Pass them as arguments
3. Return them from other functions
4. Store them in data structures (like slices, maps)

### Callback Function
A callback function is a function passed as an argument to another function, which is then called ("called back") inside that outer function.
```
package main

import "fmt"

// This is the "caller" function that receives a callback
func calculate(a int, b int, callback func(int, int) int) int {
    return callback(a, b)
}

// This is the "callback" function
func add(x int, y int) int {
    return x + y
}

func main() {
    result := calculate(3, 4, add) // Passing 'add' as a callback
    fmt.Println("Result:", result) // Output: 7
}
```
## [Closures](https://gobyexample.com/closures)
some example: [Function values](https://go.dev/tour/moretypes/24), [Function closures](https://go.dev/tour/moretypes/25) <br>
Exercise: [Fibonacci closure](https://go.dev/tour/moretypes/26)

------
## Flow control statements
### For

Go has only one looping construct, the ```for``` loop.
**Syntax:** 
```
for initialization; condition; update {
  // Statements
}
```
#### 🛡️ for Loop as a While Loop
```
for condition {
  // Statements
}
```
#### 🛡️ Infinite Loops
```
for {
  fmt.Println("Running forever...")
}
```
#### 🛡️ for...range Loop
```
for index, value = range nums {
  // Statements
}
```
You can also skip the index or value if not needed:
```
for _, value := range nums { // skip index
    // Statements
}
```
[learn for loop with code](https://github.com/mhbhuiyan99/Go/tree/main/for_loop)

### if, else if, else
```
if condition1 {
    // code1
} else if condition2 {
    // code2
} else {
    // code3
}
```
#### Short Statement in if

You can declare and initialize a variable inside the if statement.
```
if x := 5; x > 3 {
    fmt.Println("x is greater than 3")
}
```
❌ No parentheses required
```
// ✅ correct
if x > 5 { ... }

// ❌ wrong
if (x > 5) { ... } // Compiler allows but not idiomatic Go
```
#### Exercise: [Loops and Functions](https://go.dev/tour/flowcontrol/8) > [solution](https://github.com/mhbhuiyan99/Go/blob/main/Exercise/Loops_and_Functions.go) <br>
🧠 [Newton-Raphson method for finding square roots](https://mhbhuiyan.medium.com/newton-raphson-method-for-finding-square-roots-30d0f9021869)

### Switch

A ```switch``` statement is a shorter way to write a sequence of ```if - else``` statements. It runs the first case whose value is equal to the condition expression.<br>
Switch cases evaluate cases from top to bottom, stopping when a case succeeds.<br>
Switch without a condition is the same as switch ```true```.
```
package main

import "fmt"

func main() {
	mood := "hungry"

	switch mood {
	case "happy":
		fmt.Println("Let's write some awesome Go! 😄")
	case "sleepy":
		fmt.Println("Need... more... coffee... ☕😴")
	case "hungry":
		fmt.Println("Feed me bytes and burgers! 🍔💻")
	default:
		fmt.Println("Unknown mood. Rebooting... 🤖🔄")
	}
}

```
[more examples](https://go.dev/wiki/Switch)

### Defer

```defer``` is a keyword in Go used to delay the execution of a function until the surrounding function returns.
```
package main

import "fmt"

func main() {
	defer fmt.Println("World")
	fmt.Println("Hello")
}
/* Output: 
Hello
World
*/
```
#### Multiple defer statements

➡️ LIFO (Last In, First Out) — like a stack.
```
func main() {
	defer fmt.Println("One")
	defer fmt.Println("Two")
	defer fmt.Println("Three")
}
/* Output :
Three
Two
One
*/
```
---------
## [Go’s memory model](https://mhbhuiyan.medium.com/gos-memory-model-092546edd714)


## Struct
#### Defining a struct:
```
type User struct{
	// member variable or property
	Name string
	Age int
}
```
#### Create Object / Instance :: Instantiate
```
	var user1 User
	user1 = User{ // Instance or Object
		Name: "Mojammel",
		Age: 24,
	}
	user2 := User{
		Name: "Saimon",
		Age: 17,
	}
```

### Receiver Function:
Syntax:
```
func (r ReceiverType) MethodName(params) ReturnType {
    // method body
}
```
A Receiver Function is bound to a specific type (usually a struct), using a receiver.
```
package main
import "fmt"

type User struct{
	Name string
	Age int
}

func printDetails(user User){
	fmt.Println("Name = ", user.Name, "Age = ", user.Age)
}

// receiver function
func (info User) printUsingReceiverFunction(){
	fmt.Println("Name: ", info.Name, "Age: ", info.Age)
}
func (info User) printUsingReceiverFunction2(id string){
	fmt.Println("Name: ", info.Name, "ID: ", id)
}

func main(){
	var user1 User
	user1 = User{
		Name: "Mojammel",
		Age: 24,
	}
	user2 := User{
		Name: "Saimon",
		Age: 17,
	}

	printDetails(user1)
	user2.printUsingReceiverFunction()
	user1.printUsingReceiverFunction2("12345")
}
/* Output:
Name =  Mojammel Age =  24
Name:  Saimon Age:  17
Name:  Mojammel ID:  12345
*/
```
### Pointers to structs
To access the field X of a struct when we have the struct pointer p we could write (*p).X. However, that notation is cumbersome, so the language permits us instead to write just p.X, without the explicit dereference.
```
type Vertex struct {
	X int
	Y int
}
func main() {
	v := Vertex{1, 2}
	p := &v
	p.X = 1e9
	fmt.Println(p.X)
}
```
### Struct Literals

A [struct literal](https://go.dev/tour/moretypes/5) denotes a newly allocated struct value by listing the values of its fields.<br>

You can list just a subset of fields by using the Name: syntax. (And the order of named fields is irrelevant.)<br>

The special prefix & returns a pointer to the struct value.
```
type Vertex struct {
	X, Y int
}

var (
	v1 = Vertex{1, 2}  // has type Vertex
	v2 = Vertex{X: 1}  // Y:0 is implicit
	v3 = Vertex{}      // X:0 and Y:0
	p  = &Vertex{1, 2} // has type *Vertex
)

func main() {
	fmt.Println(v1, p, v2, v3)
}
```

-------
## Array
```
package main
import "fmt"

func main(){
	var arr [2] int
	fmt.Println(arr)
	// Output : [0 0], 0 is the default array value in Go

	/* 
	arr[2] = 7
	fmt.Println(arr)

	Output: Error
	0 index array : index start from 0 [0 1 ...] */

	arr[1] = 7
	fmt.Println(arr)
	// Output: [0 7]

	
	arr2 := [2] int{3,9}
	fmt.Println(arr2)
	// Output: [3 9]

	// print specific index value
	fmt.Println(arr[1], arr2[1])
	// Output: 7 9
}
```
-----------
## Pointers to Slices
read it: [From Pointers to Slices: Exploring Go’s Underlying Memory](https://mhbhuiyan.medium.com/from-pointers-to-slices-exploring-gos-underlying-memory-d611b408870c)

----------
## Maps
map = Key value pairs. <br>
Array and Slice use only integer type index. But in ```map``` index can be any type.
```
func main() {
	menu := map[string]float64{
		"rice": 11.23,
		"milk": 5.67,
		"tea" : 3.43,
		"water": 1.87, // last comma must. otherwise syntax error
	}
	fmt.Println(menu)
}
```
**printng using ```loop```:** <br>
The	for...range	loop processes map keys and values in a	random order because a map is an unordered collection of keys and values.
```
	for key, value := range menu{
		fmt.Println(key, "\t:: ", value) // the \t for tab
	}
```
Try another way:
```
func main() {
	var menu map[string]int
	menu["rice"] = 123
	fmt.Println(menu) 
}
// result: ❌ error
panic: assignment to entry in nil map
```
The zero value of a map is nil. A nil map has no keys, nor can keys be added.<br>
How can we fixed it?<br>
The ```make``` function returns a map of the given type, initialized and ready for use.
```
func main() {
	menu := make(map[string]int)
	menu["rice"] = 123
	fmt.Println(menu) // Output: map[rice:123]
}
```
```make``` allocate and initiate - non zeroed storage<br>
There is an another keyword ```new```, it only allocate, no initiate of memory. 

### Map literals
Map literals are like struct literals, but the keys are required.
```
type Vertex struct {
	Lat, Long float64
}

var m = map[string]Vertex{
	"Bell Labs": Vertex{
		40.68433, -74.39967,
	},
	"Google": Vertex{
		37.42202, -122.08408,
	},
}

func main() {
	fmt.Println(m)
}
// Output: map[Bell Labs:{40.68433 -74.39967} Google:{37.42202 -122.08408}]
```
If the top-level type is just a type name, you can omit it from the elements of the literal.
```
var m = map[string]Vertex{
	"Bell Labs": {40.68433, -74.39967},
	"Google":    {37.42202, -122.08408},
}
```
#### Insert or update an element in map m:
```m[key] = elem```
#### Retrieve an element:
```elem = m[key]```
#### Delete an element:
```delete(m, key)```
#### Test that a key is present with a two-value assignment:
```elem, ok = m[key]```<br>
If key is in m, ok is true. If not, ok is false.<br>
If key is not in the map, then elem is the zero value for the map's element type.
```
func main() {
	m := make(map[string]int)

	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])

	m["Answer"] = 48
	fmt.Println("The value:", m["Answer"])

	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])

	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)
} /* Output:
The value: 42
The value: 48
The value: 0
The value: 0 Present? false */
```
Exercise: [Maps + BONUS(strings.Fields)](https://go.dev/tour/moretypes/23) : [Solution](https://github.com/mhbhuiyan99/Go/blob/main/Exercise/map__withStringsFields.go) <br>

**Removing key/value pairs with the ```delete``` function:**
```
func main() {
	menu := map[string]float64{
		"rice": 11.23,
		"milk": 5.67,
		"tea" : 3.43,
	}
	delete(menu, "rice")
	fmt.Println(menu)
} // Output: map[milk:5.67 tea:3.43]
```
--------


