package main

import "flag"
import "fmt"
import "encoding/json"

func main() {

	var a = make(map[string]int)
	var b = make(map[string]string)
	b["aaa"] = "vvv"
	a["aaa"] = 111
	res1, _ := json.Marshal(a)
	res2, _ := json.Marshal(b)
	fmt.Println("res1:", string(res1))
	fmt.Println("res2:", string(res2))
	return
	var ospath string
	flag.StringVar(&ospath, "path", "./", "被监控的文件夹，会监控当前文件夹和子文件夹")
	flag.Parse()
	fmt.Println("ospath:", ospath)

}
