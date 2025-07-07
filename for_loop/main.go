package main

import "fmt"

func For(x int) {
	fmt.Print("for loop : ")

	for i := 0; i <= x; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
}

func For_as_While(x int) {
	fmt.Print("for as while : ")

	for i := x; i >= 0; i-- {
		fmt.Print(i, " ")
	}
	fmt.Println()
}

func For_Range(x int) {
	fmt.Print("for...range Loop : ")

	for i := range x {
		fmt.Print(i, " ")
	}
	fmt.Println()
}

func main() {

	n := 7
	For(n)
	For_as_While(n)
	For_Range(n)
}
