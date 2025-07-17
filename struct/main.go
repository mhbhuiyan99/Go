package main

import (
	"fmt"
	"learn/others"
)

func printDetails(user others.User) {
	fmt.Println("Name = ", user.Name, "Age = ", user.Age)
}

func main() {
	var user1 others.User

	user1 = others.User{ // Instance or Object
		Name: "Mojammel",
		Age:  24,
	}

	user2 := others.User{
		Name: "Saimon",
		Age:  17,
	}

	printDetails(user1)
	user2.PrintUsingReceiverFunction()
	user1.PrintUsingReceiverFunction2("12345")

	others.Literals()
}
