package main

import "fmt"

func main() {
	s := GetSingleton()
	fmt.Println(s)
	s = nil
	s = GetSingleton()
}
