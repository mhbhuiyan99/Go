package others

import "fmt"

func (info User) PrintUsingReceiverFunction() {
	fmt.Println("Name: ", info.Name, "Age: ", info.Age)
}
func (info User) PrintUsingReceiverFunction2(id string) {
	fmt.Println("Name: ", info.Name, "ID: ", id)
}