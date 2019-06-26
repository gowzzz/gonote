package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"
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
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main() {

	//获取当前时间
	t := time.Now() //2018-07-11 15:07:51.8858085 +0800 CST m=+0.004000001
	fmt.Println(t.Year())
	fmt.Println(t.Month())
	fmt.Println(t.Day())
	fmt.Println(t.Hour())
	fmt.Println(t.Minute())
	fmt.Println(t.Unix() % (60 * 60 * 24)) //1531293019
	// fmt.Println(t.Second())    //1531293019
	// // fmt.Println(t.Format("2006-01-02 15:04:05"))
	t1 := time.Unix(t.Unix()-t.Unix()%(60*60*24), 0)
	t2 := time.Unix(t.Unix()-int64(t.Second()), 0)
	fmt.Println(t1.Format("2006-01-02 15:04:05"))
	fmt.Println(t2.Format("2006-01-02 15:04:05"))
	// fmt.Println(GetRandomString(5))
	return
	path := `G:\工作目录\homey2.2\aaa`
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)
			if err != nil {
				log.Println("PathExists file err2:", err)
			}
		} else {
			log.Println("PathExists file err:", err)
		}
	} else {
		log.Println("okkkkkkkkkk")
	}

	// fmt.Println(IsPhoneNum("13521146683"))
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
