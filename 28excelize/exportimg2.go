package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"

	// "io/ioutil"

	"bytes"
	"net/http"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)
/*
一定要先修改全部单元格的大小，在存入图片，否则会导致部分缩放的图片因为图片的拉伸而拉伸。
*/ 
func openfile() (f *excelize.File) {
	f, err := excelize.OpenFile("filename")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return f
}

// AddPicture(sheet, cell, picture, format string)
//f.AddPictureFromBytes(sheet, cell, format, name, "jpeg", file)    file, _ := ioutil.ReadFile(picture)
func AddImgBase64ToExcel2(xlsx *excelize.File, sheet, cell string, width, height float64, imgbase64, imgext, realname string) error {
	// tmpfile := realname + "." + imgext
	// os.Remove(tmpfile)
	content, err := base64.StdEncoding.DecodeString(imgbase64)
	if err != nil {
		fmt.Println("err1:", err)
		return err
	}

	c, _, err := image.DecodeConfig(bytes.NewReader(content))
	// c, _, err := image.DecodeConfig(f)
	if err != nil {
		fmt.Println("err2:", err)
		return err
	}
	wmultiple, hmultiple := ImageZoom(width, height, float64(c.Width), float64(c.Height))
	multiple := wmultiple
	if hmultiple < wmultiple {
		multiple = hmultiple
	}
	multiplestr := strconv.FormatFloat(multiple, 'f', 3, 64) //保留3位小数
	format := `{"x_scale": ` + multiplestr + `, "y_scale": ` + multiplestr + `, "lock_aspect_ratio": true}`
	err = xlsx.AddPictureFromBytes(sheet, cell, format, realname, imgext, content)
	if err != nil {
		fmt.Println("AddPicture err  :", err)
		return err
	}
	fmt.Println("okkk")

	return nil

	// err = ioutil.WriteFile(tmpfile, content, os.ModePerm)
	// if err != nil {
	// 	fmt.Println("err2:", err)
	// 	return err
	// }
	// // defer os.Remove(tmpfile)
	// _, err = os.Open(tmpfile)
	// if err != nil {
	// 	fmt.Println("err3:", err)
	// 	return err
	// }
}

func AddOnlineImgToExcel2(xlsx *excelize.File, sheet, location string, width, height float64, imgpath string) error {
	// var f io.Reader
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp2, err := client.Get(imgpath)
	if err != nil {
		fmt.Println("client.Get  err  :", err)
		return err
	}

	c, _, err := image.DecodeConfig(resp2.Body)
	// c, _, err := image.DecodeConfig(f)
	if err != nil {
		fmt.Println("err2:", err)
		return err
	}
	wmultiple, hmultiple := ImageZoom(width, height, float64(c.Width), float64(c.Height))
	multiple := wmultiple
	if hmultiple < wmultiple {
		multiple = hmultiple
	}
	multiplestr := strconv.FormatFloat(multiple, 'f', 3, 64) //保留3位小数
	format := `{"x_scale": ` + multiplestr + `, "y_scale": ` + multiplestr + `, "lock_aspect_ratio": true}`

	resp, err := client.Get(imgpath)
	if err != nil {
		fmt.Println("client.Get  err  :", err)
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(" ioutil.ReadAll file err:", err)
		return err
	}
	err = xlsx.AddPictureFromBytes(sheet, location, format, imgpath, ".jpg", content)
	if err != nil {
		fmt.Println("AddPicture err  :", err)
		return err
	}
	return nil
}

func AddImgToExcel2(xlsx *excelize.File, sheet, location string, width, height float64, imgpath string) error {
	var f io.Reader

	file, err := os.Open(imgpath)
	if err != nil {
		fmt.Println("Open imgpath err  :", err)
		return err
	}
	f = file
	c, _, err := image.DecodeConfig(f)
	if err != nil {
		fmt.Println("SetCellHeight err  :", err)
		return err
	}
	wmultiple, hmultiple := ImageZoom(width, height, float64(c.Width), float64(c.Height))
	multiple := wmultiple
	if hmultiple < wmultiple {
		multiple = hmultiple
	}
	multiplestr := strconv.FormatFloat(multiple, 'f', 3, 64) //保留3位小数
	err = xlsx.AddPicture(sheet, location, imgpath, `{"x_scale": `+multiplestr+`, "y_scale": `+multiplestr+`, "lock_aspect_ratio": true}`)
	if err != nil {
		fmt.Println("AddPicture err  :", err)
		return err
	}
	return nil
}

// 像素（72磅=1英寸=2.54厘米=96像素）
/*
//像素(x) dpi(d) 英寸(y)  y=x/dpi     cm=y*2.54     2=144/72
//行高 磅=pixelH*0.6 - 0.05
//列宽 字符=(pixelW-7)/9 + 7.0/9

excel以dpi=120为准
1英寸=像素/120
1英寸=2.54cm
行高单位：磅=像素*0.6-0.05；
列宽单位：字符=(像素-7)/9；

行高有个很奇怪的舍入：0.08~0.12都记作0.1   0.13~0.17都记作0.15 ，舍入单位为0.05。所有为了全入起见，原大小+0.05
*/

// 1yc=72bang=120pixel=2.54cm
//cell cm   image pixel
// excel以120为dpi，图片以png保存在excel中，以dpi为96操作 像素需要*120/96 相当于变成了原来的1.25倍
func ImageZoom(width, height, imgWidth, ImgHeight float64) (wmultiple, hmultiple float64) {
	pixelH := height / 2.54 * 120 //像素
	pixelW := width / 2.54 * 120  //像素
	// 用像素比计算出缩放比例,宽高中选择缩放小的一方
	var pixelImgW float64 = float64(imgWidth * 120 / 96)
	var pixelImgH float64 = float64(ImgHeight * 120 / 96)
	wmultiple = pixelW / pixelImgW
	hmultiple = pixelH / pixelImgH
	return
}

// 列宽 12.56字符=120pixel
// 行号 72磅=120pixel
// 1yc=72bang=120pixel=2.54cm
func SetCellHeight(f *excelize.File, sheet string, row int, height float64) error {
	// 传入cm，写成像素
	// 设置单元格长宽
	//行高有个很奇怪的舍入：0.08~0.12都记作0.1   0.13~0.17都记作0.15 ，舍入单位为0.05。所有为了全入起见，原大小+0.05
	pixelH := height / 2.54 * 120 //像素
	ch := pixelH*0.6 - 0.05 + 0.05
	fmt.Println("ch:", ch)
	f.SetRowHeight(sheet, row, ch)
	return nil
}
func SetCellWtdth(f *excelize.File, sheet, startcol, endcol string, width float64) error {
	// 传入cm，写成像素
	// 设置单元格长宽
	pixelW := width / 2.54 * 120 //像素
	cw := (pixelW-7)/9 + 7.0/9

	fmt.Println("cw:", cw)
	f.SetColWidth(sheet, startcol, endcol, cw)

	return nil
}
