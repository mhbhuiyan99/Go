package others

import "fmt"

type what struct{
	x,y int
}

var (
	v1 = what{1,2}
	v2 = what{x:3}
	v3 = what{}
	p = &what{4,5}
)

func Literals(){
	fmt.Println(v1, v2, v3, p)
}