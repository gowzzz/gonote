package main

import (
	"fmt"
)

type Slice []int

func NewSlice() Slice {
	return make(Slice, 0)
}
func (s *Slice) Add(elem int) *Slice {
	*s = append(*s, elem)
	fmt.Println(elem)
	return s
}

// ---------------------
type student struct {
	Name string
	Age  int
}

func pase_student() {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for k, stu := range stus {
		m[stu.Name] = &stus[k]
	}
	for k, v := range m {
		println(k, "=>", v.Name)
	}
}
func main() {
	l := new([]int)
	*l = append(*l, 1)
	fmt.Println(*l)
	// ----------
	// pase_student()
	// ----------------
	// s := NewSlice()
	// defer s.Add(1).Add(2)
	// fmt.Println("aaaaaaaaaaaa")
	// s.Add(3)
	// fmt.Println("bbbbbbbb")

	// -------------
	// x := []string{"a", "b", "c"}
	// for v := range x {
	// 	fmt.Print(v)
	// }
	// -------------------------
	// for i := 0; i < 5; i++ {
	// 	go func() {
	// 		fmt.Printf("iii:%d ", i)
	// 	}()
	// }
	// for i := 0; i < 5; i++ {
	// 	go func(num int) {
	// 		fmt.Printf("num:%d ", num)
	// 	}(i)
	// }
	// time.Sleep(1 * time.Second)
	// for i := 0; i < 5; i++ {
	// 	defer fmt.Printf("%d ", i)
	// }
	// defer fmt.Printf("ccc:%d ", c())
	// defer func(){
	// 	fmt.Println("111")
	// }()
	// defer func(){
	// 	fmt.Println("222")
	// }()
	// defer func(){
	// 	fmt.Println("333")
	// }()
	// panic(444)
}
func c() (i int) {
	defer func() {
		fmt.Printf("%d ", i)
		i++
		fmt.Printf("%d ", i)
	}()
	return 1
}
