package main

import "fmt"

type User struct{
	// member variable or property
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

	user1 = User{ // Instance or Object
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

/* 
Output:
Name =  Mojammel Age =  24
Name:  Saimon Age:  17
Name:  Mojammel ID:  12345
*/
