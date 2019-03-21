package main

import (
	"fmt"

	"github.com/casbin/casbin"
)

func main() {
	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.
	e := casbin.NewEnforcer("./model.conf", "./policy.csv")

	// EnforceSafe
	// res, err := e.EnforceSafe(sub, obj, act)
	// if err != nil {
	// 	fmt.Println("err:", err)
	// 	return
	// }
	// if res {
	// 	fmt.Println("res ok")

	// }
	if e.Enforce(sub, obj, act) == true {
		// permit alice to read data1
		fmt.Println("ok")
	} else {
		// deny the request, show an error
		fmt.Println("deny")
	}
}
