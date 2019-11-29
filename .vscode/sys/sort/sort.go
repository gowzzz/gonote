package main

import (
	"fmt"
	"sort"
)

func main() {
	mapInfo := map[string]int32{
		"roy":     18,
		"1jason2": 23,
		"kitty":   16,
		"2jason3": 23,
		"hugo":    21,
		"3jason4": 23,
		"tina":    35,
		"4jason5": 23,
	}

	type peroson struct {
		Name string
		Age  int32
	}

	var lstPerson []peroson
	for k, v := range mapInfo {
		lstPerson = append(lstPerson, peroson{k, v})
	}

	sort.SliceStable(lstPerson, func(i, j int) bool {
		return lstPerson[i].Name > lstPerson[j].Name // 降序
		// return lstPerson[i].Age < lstPerson[j].Age  // 升序
	})
	fmt.Println(lstPerson)
	sort.SliceStable(lstPerson, func(i, j int) bool {
		return lstPerson[i].Age > lstPerson[j].Age // 降序
		// return lstPerson[i].Age < lstPerson[j].Age  // 升序
	})
	fmt.Println(lstPerson)
}
