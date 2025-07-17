package main

import "fmt"

func print(num *[3]int) {
	fmt.Println(num)
}

func main() {
	x := 21
	fmt.Println("Value of x: ", x) // Output: Value of x: 21

	p := &x // p holds the address of x
	*p = 27 // change the value of x through its address

	fmt.Println("Address of x: ", p) // Output: Address of x: 0xc00000a0d8
	fmt.Println("Value at address p(&x):", *p)
	// Output: Value at address p(&x): 27
	fmt.Println("Value of x: ", x) // Output: Value of x: 27

	/* why pointer?
	- Avoid copying
	- Modify the original array
	- Reduced memory usage */

	arr := [3]int{3, 5, 7}
	print(&arr)
}
