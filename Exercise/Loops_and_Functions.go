Problem Link : https://go.dev/tour/flowcontrol/8

package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := float64(1)
	
	for {
		new_result := z - (z*z - x)/(2*z)
		
		if math.Abs(new_result - z) < 1e-6{
			return new_result
		}
		z = new_result
	}	
}

func main() {
	fmt.Println(Sqrt(144))
}
