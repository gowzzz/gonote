package main

import (
	"fmt"
	"regexp"
	"unicode/utf8"
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
func test(a []int) {
	fmt.Println(a)
	a[1] = 99
	a = []int{4, 5, 6, 7}
	fmt.Println(a)
}

//字符串验证：中英文数字
func ValidateString(str string, len int) (code int) {
	nl := regexp.MustCompile("[^\u4e00-\u9fa50-9A-Za-z]+") //是否包含之外的字符，true包含 false不包含
	if nl.MatchString(str) {
		return -1
	}
	sLen := utf8.RuneCountInString(str)
	if sLen > len || sLen < 1 {
		return -1
	}

	return
}

//参数：手机号
func IsPhoneNum(str string) bool {
	// reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	// reg := `^\d{11}$` //11位数字
	reg := `^1\d{10}$` //以1开头，后面跟10位数字
	zl := regexp.MustCompile(reg)
	//是否符合规则
	if zl.MatchString(str) {
		return true
	}
	return false
}
func main() {
	fmt.Println(IsPhoneNum("13521146683"))
	return
	var aaa = make(map[string]interface{})
	// var bbb = make(map[string]string)
	aaa["aaa"] = 1
	aaa["bbb"] = "abc"
	fmt.Println(aaa["a"])
	tmp, ok := aaa["a"].(string)
	if ok {
		fmt.Println(tmp)
	} else {
		fmt.Println("no a")
	}

	if aaa["bbb"] == "abc1" {
		fmt.Println(1)
	} else {
		fmt.Println(2)
	}

	// fmt.Println(ValidateString(aaa["bbb"], 2))

	return

	a := []int{1, 2, 3}
	test(a)
	fmt.Println("main:", a)

	return
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
