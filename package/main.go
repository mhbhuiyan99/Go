package main

import (
	"fmt"
	"Bhuiyan/package/mathematical_operations"
)

func main(){
	x, y := 3, 5

	fmt.Println("Sum = ", mathematical_operations.Sum(x,y))
	fmt.Println("Multiply = ", mathematical_operations.Multiply(x,y))
}