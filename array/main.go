package main
import "fmt"

func main(){
	var arr [2] int
	fmt.Println(arr)

	// Print : [0 0], 0 is the default array value in Go

	/* 
	arr[2] = 7
	fmt.Println(arr)

	Output: Error
	0 index array : index start from 0 [0 1 ...] */

	arr[1] = 7
	fmt.Println(arr)

	
	arr2 := [2] int{3,9}
	fmt.Println(arr2)

	// print specific index value
	fmt.Println(arr[1], arr2[1])
}
/* Output:
[0 0]
[0 7]
[3 9]
7 9
*/
