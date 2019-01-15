package main

import (
	"fmt"
	"sync"
)

func GetMapString(m sync.Map, k string) (res string) {
	fmt.Println("GetMapString")
	if v, ok := m.Load(k); ok {
		if str, ok := v.(string); ok {
			res = str
			fmt.Println("res:", res)
		} else {
			fmt.Println("str:", str)

		}
	} else {
		fmt.Println("v:", v)
	}
	return
}
func main() {
	var m sync.Map
	m.Store("a", 1)
	m.Store("v", "aa")
	fmt.Println(GetMapString(m, "v"))
	// if v, ok := m.LoadOrStore("b", 2); !ok {
	// 	fmt.Println(" LoadOrStore not exist:", v)
	// }
	// if v, ok := m.Load("c"); !ok {
	// 	fmt.Println("Load not exist:", v)
	// }

	// m.Delete("d")
	// f := func(k, v interface{}) bool {
	// 	fmt.Println(k, v)
	// 	return true
	// }
	// m.Range(f)
}
