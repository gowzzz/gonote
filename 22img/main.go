package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	_ "image/gif"
    _ "image/jpeg"
	_ "image/png"
	"image"
)

var filePth = "./2.jpg"
func main(){
	f, err := os.Open(filePth)
	if err != nil {
		fmt.Println("err1:", err)
		return
	}
	defer f.Close()
	// c, s1, err := image.DecodeConfig(f)
	// fmt.Printf("%+v\n",c)
	// fmt.Println("s1 = ", s1)
	// fmt.Println("width = ", c.Width)
	// fmt.Println("height = ", c.Height)

	img, s2, err := image.Decode(f)
	fmt.Printf("%+v\n",img.ColorModel())
	fmt.Printf("%+v\n",img.Bounds())
	fmt.Println("s2 = ", s2)
}
func main2() {
	f, err := os.Open(filePth)
	if err != nil {
		fmt.Println("err1:", err)
		return
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("err2:", err)
		return
	}
	CreateFile("./a.jpeg", content)
}
func CreateFile(filename string, content []byte) {
	//创建文件
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Create File1 file err:", err)
		return
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	var num int
	num, err = w.Write(content)
	fmt.Println("------------------filename:", filename)
	fmt.Println("------------------content:", len(content))
	fmt.Println("------------------num:", num)
	if err != nil {
		log.Println("Create File2 file err:", err)
		return
	}
	return
}
