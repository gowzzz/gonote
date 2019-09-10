package main
import(
	"fmt"
)
const ( 
    a = 10
    b 
    c 
    d =iota * 2 
    e 
    f 
    g 
)
func main(){
	fmt.Println("a:",a)
	fmt.Println("b:",b)
	fmt.Println("c:",c)
	fmt.Println("d:",d)
	fmt.Println("e:",e)
	fmt.Println("f:",f)
	fmt.Println("g:",g)
}
/*
a: 10
b: 10
c: 10
d: 6
e: 8
f: 10
g: 12
*/ 