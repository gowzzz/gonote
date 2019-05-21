package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// 获取在线图片
func GetOnlineImage() {
	// var url = "https://pic2016.ytqmx.com/img/20180704/qfvzsyartbg.jpg"
	var url = "https://img3.doubanio.com/view/group_topic/large/public/p104935104.jpg"
	// var url = "https://gss3.bdstatic.com/-Po3dSag_xI4khGkpoWK1HF6hhy/baike/c0%3Dbaike150%2C5%2C5%2C150%2C50/sign=389ad835d158ccbf0fb1bd6878b1d75b/fd039245d688d43f0cd0e0757d1ed21b0ef43b3f.jpg"
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Fail exec http client Do,err:%s\n", err.Error())
	}
	fmt.Println("resp.Body:", resp.ContentLength)
	f, err := os.Create("wz.jpg")
	if err != nil {
		panic(err)
	}
	wlen, err := io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("wlen:", wlen)
	f.Close()
}

var filePth = "./2.jpg"

func main() {
	GetOnlineImage()

	return
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
	fmt.Printf("%+v\n", img.ColorModel())
	fmt.Printf("%+v\n", img.Bounds())
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
