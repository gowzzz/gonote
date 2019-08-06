package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/extrame/xls"

	"crypto/tls"
	"encoding/base64"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
)
/*
一定要先修改全部单元格的大小，在存入图片，否则会导致部分缩放的图片因为图片的拉伸而拉伸。
*/ 
func 
// 该方法只读取Sheet1
func ReadExcel(content []byte) (datas [][]string, err error) {
	// var datas [][]string

	if f, err := excelize.OpenReader(bytes.NewReader(content)); err == nil {
		sheets := f.GetSheetMap()
		if sheets == nil {
			return nil, errors.New("xlsx sheet is null")
		}
		fmt.Printf("Sheet:%+v\n", sheets)
		fmt.Println("ReadExcel:excelize")
		// 获取 Sheet1 上所有单元格
		datas = f.GetRows(sheets[1])
		// for _, row := range rows {
		// 	for _, colCell := range row {
		// 		fmt.Print(colCell, "\t")
		// 	}
		// 	fmt.Println()
		// }
		return datas, nil
	} else {
		fmt.Println("ReadExcel excelize:err:", err)
		if xlsFile, err := xls.OpenReader(bytes.NewReader(content), "utf-8"); err == nil {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("datas:%+v\n", datas)
				}
			}()
			fmt.Println("ReadExcel:xls")
			sheet1 := xlsFile.GetSheet(0)
			for i := 0; i <= int(sheet1.MaxRow); i++ {
				var rowdata []string
				row := sheet1.Row(i)
				for index := row.FirstCol(); index < row.LastCol(); index++ {
					rowdata = append(rowdata, row.Col(index))
				}
				datas = append(datas, rowdata)
			}
			return datas, nil
		} else {
			fmt.Println("ReadExcel xls:err:", err)
			fmt.Println("ReadExcel:csv")
			if datas, err := csv.NewReader(bytes.NewReader(content)).ReadAll(); err == nil {
				return datas, nil
			} else {
				fmt.Println("ReadExcel csv:err:", err)
				return nil, err
			}
		}
	}
}
func AddImgBase64ToExcel(xlsx *excelize.File, sheet, location string, width, height float64, imgbase64, imgext, realname string) error {
	// tmpfile := realname + "." + imgext
	// os.Remove(tmpfile)
	content, err := base64.StdEncoding.DecodeString(imgbase64)
	if err != nil {
		fmt.Println("err1:", err)
		return err
	}

	img, err := jpeg.Decode(bytes.NewReader(content))
	if err != nil {
		fmt.Println("jpeg.decode  err  :", err)
		return err
	}
	// pixelH := height / 2.54 * 90 //像素
	// pixelW := width / 2.54 * 90  //像素
	// pixelSize := pixelH
	// if pixelW<pixelH{
	// 	pixelSize=pixelW
	// }

	var m image.Image
	if img.Bounds().Dx() > img.Bounds().Dy() {
		m = resize.Resize(90, 0, img, resize.Lanczos3)
	} else {
		m = resize.Resize(0, 90, img, resize.Lanczos3)
	}
	buffer := bytes.NewBuffer(nil)
	jpeg.Encode(buffer, m, nil)

	format := `{"lock_aspect_ratio": true}`
	err = xlsx.AddPictureFromBytes(sheet, location, format, "xx", ".jpg", buffer.Bytes())
	if err != nil {
		fmt.Println("AddPicture err  :", err)
		return err
	}
	return nil

}
func AddPNGBase64ToExcel(xlsx *excelize.File, sheet, location string, width, height float64, imgbase64, imgext, realname string) error {
	// tmpfile := realname + "." + imgext
	// os.Remove(tmpfile)
	content, err := base64.StdEncoding.DecodeString(imgbase64)
	if err != nil {
		fmt.Println("err1:", err)
		return err
	}

	img, err := png.Decode(bytes.NewReader(content))
	if err != nil {
		fmt.Println("png.decode  err  :", err)
		return err
	}
	// pixelH := height / 2.54 * 90 //像素
	// pixelW := width / 2.54 * 90  //像素
	// pixelSize := pixelH
	// if pixelW<pixelH{
	// 	pixelSize=pixelW
	// }
	var m image.Image
	if img.Bounds().Dx() > img.Bounds().Dy() {
		m = resize.Resize(90, 0, img, resize.Lanczos3)
	} else {
		m = resize.Resize(0, 90, img, resize.Lanczos3)
	}

	buffer := bytes.NewBuffer(nil)
	png.Encode(buffer, m)

	format := `{"lock_aspect_ratio": true}`
	err = xlsx.AddPictureFromBytes(sheet, location, format, "xx", ".jpg", buffer.Bytes())
	if err != nil {
		fmt.Println("AddPicture err  :", err)
		return err
	}
	return nil

}
func AddPNGToExcel(xlsx *excelize.File, sheet, location string, width, height float64, imgpath string) error {
	var f io.Reader
	if strings.HasPrefix(imgpath, "http") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(imgpath)
		if err != nil {
			fmt.Println("client.Get err  :", err)
			return err
		}
		f = resp.Body
	} else {
		file, err := os.Open(imgpath)
		if err != nil {
			fmt.Println("os.Open err  :", err)
			return err
		}

		defer file.Close()
		f = file
	}

	img, err := png.Decode(f)
	if err != nil {
		fmt.Println("png.Decode err  :", err)
		return err
	}
	// pixelH := height / 2.54 * 90 //像素
	// pixelW := width / 2.54 * 90  //像素
	// pixelSize := pixelH
	// if pixelW<pixelH{
	// 	pixelSize=pixelW
	// }

	var m image.Image
	if img.Bounds().Dx() > img.Bounds().Dy() {
		m = resize.Resize(90, 0, img, resize.Lanczos3)
	} else {
		m = resize.Resize(0, 90, img, resize.Lanczos3)
	}

	// write new image to file
	buffer := bytes.NewBuffer(nil)
	png.Encode(buffer, m)

	format := `{"lock_aspect_ratio": true, "locked": true, "positioning": "oneCell"}` //absolute
	// err = xlsx.AddPicture(sheet, location, outname, `{"lock_aspect_ratio": true, "locked": true, "positioning": "absolute"}`)//oneCell
	err = xlsx.AddPictureFromBytes(sheet, location, format, "xx", ".jpg", buffer.Bytes())
	if err != nil {
		fmt.Println("AddPicture err  :", err)
		return err
	}
	return nil
}
func AddImgToExcel(xlsx *excelize.File, sheet, location string, width, height float64, imgpath string) error {
	var f io.Reader
	if strings.HasPrefix(imgpath, "http") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(imgpath)
		if err != nil {
			fmt.Println("client.Get err  :", err)
			return err
		}
		f = resp.Body
	} else {
		file, err := os.Open(imgpath)
		if err != nil {
			fmt.Println("os.Open err  :", err)
			return err
		}

		defer file.Close()
		f = file
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		fmt.Println("jpeg.Decode err  :", err)
		return err
	}
	// pixelH := height / 2.54 * 90 //像素
	// pixelW := width / 2.54 * 90  //像素
	// pixelSize := pixelH
	// if pixelW<pixelH{
	// 	pixelSize=pixelW
	// }

	var m image.Image
	if img.Bounds().Dx() > img.Bounds().Dy() {
		m = resize.Resize(90, 0, img, resize.Lanczos3)
	} else {
		m = resize.Resize(0, 90, img, resize.Lanczos3)
	}

	// write new image to file
	buffer := bytes.NewBuffer(nil)
	jpeg.Encode(buffer, m, nil)

	format := `{"lock_aspect_ratio": true, "locked": true, "positioning": "oneCell"}`
	// err = xlsx.AddPicture(sheet, location, outname, `{"lock_aspect_ratio": true, "locked": true, "positioning": "absolute"}`)//oneCell
	err = xlsx.AddPictureFromBytes(sheet, location, format, "xx", ".jpg", buffer.Bytes())
	if err != nil {
		fmt.Println("AddPicture err  :", err)
		return err
	}
	return nil
}
