package main

import "fmt"

const (
	a = iota
	b
	c
	d
)
const (
	e = iota
	f
	g
)

func main() {
	defer func() {
		fmt.Println("a:", a)
	}()
	defer func() {
		fmt.Println("d:", d)
	}()
	defer func() {
		fmt.Println("e:", e)
	}()
	defer func() {
		fmt.Println("f:", f)
	}()
	defer func() {
		fmt.Println("g:", g)
	}()
}
