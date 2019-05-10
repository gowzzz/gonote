package main

import (
    "fmt"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"

    "github.com/360EntSecGroup-Skylar/excelize"
)
var filename="./Book1.xlsx"
func openfile()(f *excelize.File){
	f,err:=excelize.OpenFile(filename)
	if err!=nil{
		fmt.Println(err)
		return nil	
	}
	return  f
}

func newfile()(f *excelize.File){
	return excelize.NewFile()
}
func main() {
	var err error
    f := newfile()
	cellsize(f)
	// AddImg(f)
    err = f.SaveAs(filename)
    if err != nil {
        fmt.Println(err)
    }
}

// 像素（72磅=1英寸=2.54厘米=96像素）
/*


//像素(x) dpi(d) 英寸(y)  y=x/dpi     cm=y*2.54     2=144/72
//行高 磅=(x-7)/9
//列宽 字符=x*0.6  5/3
var imgdpi float64 = 72.00
var exceldpi float64 = 120.00
var cellpixel float64 = 2.00                   //1代表2.54厘米
var pixel float64 = exceldpi * cellpixel       //2英寸  5.08 (excel以120为dpi，dpi为72的像素存入execl 像素需要*120/72 相当于变成了原来的1.25倍)
var excelw float64 = (pixel-7)/9 + 0.78 + 0.01 //单元格稍大一点
var excelh float64 = pixel*0.6 + 0.01          //单元格稍大一点

单元格大小设置成2*2cm
单元格设置
2cm=120*2*2.54 =609.6像素
ch=609.6*0.6-0.05
cw=(609.6-7)/9+7.0/9
照片缩放
h2=h*120/72=989*1.25=1236.25
w2=w*120/72=700*1.25=875
hmultiple=609.6/1236.25=0.494
wmultiple=609.6/875=0.6967


excel以120为dpi，dpi为72的像素存入execl 像素需要*120/72 相当于变成了原来的1.25倍


ch=h2*0.6-0.05=741.7
cw=(w2-7)/9=96.44

excel以dpi=120为准
1英寸=像素/120
1英寸=2.54cm
行高单位：磅=像素*0.6-0.05；
列宽单位：字符=(像素-7)/9；

行高有个很奇怪的舍入：0.08~0.12都记作0.1   0.13~0.17都记作0.15 ，舍入单位为0.05。所有为了全入起见，原大小+0.05
*/
func cellsize(f *excelize.File){
	fmt.Println(7.0/9)
	var err error
	err = f.SetColWidth("Sheet1", "A", "H", (10+7.0/9))//19.22字符
	// err = f.SetColWidth("Sheet1", "A", "H",((609.6-7)/9+7.0/9))//19.22字符
	if err != nil {
		fmt.Println(err)
		return 
	}
	err = f.SetRowHeight("Sheet1", 1, 50.12+0.05)//50磅  12=>10 13/17=>15  18=>20
	// err = f.SetRowHeight("Sheet1", 1, (609.6*0.6) )//50磅
	if err != nil {
		fmt.Println(err)
		return 
	}

	width, err := f.GetColWidth("Sheet1", "A")
	if err != nil {
		fmt.Println(err)
		return 
	}
	height, err := f.GetRowHeight("Sheet1", 1)
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Println("width:",width)
	fmt.Println("height:",height)
}
func AddImg(f *excelize.File){
	    // 插入图片
    err := f.AddPicture("Sheet1", "A1", "./2.jpg", `{"x_scale": 0.494, "y_scale": 0.494}`)
    if err != nil {
        fmt.Println(err)
    }
    // 插入带有缩放比例和超链接的图片
    // err = f.AddPicture("Sheet1", "D2", "./2.jpg", `{"x_scale": 0.5, "y_scale": 0.5}`)
    // if err != nil {
    //     fmt.Println(err)
    // }
    // // 插入图片，并设置图片的外部超链接、打印和位置属性
    // err = f.AddPicture("Sheet1", "H2", "./3.gif", `{"x_offset": 15, "y_offset": 10, "hyperlink": "https://github.com/360EntSecGroup-Skylar/excelize", "hyperlink_type": "External", "print_obj": true, "lock_aspect_ratio": false, "locked": false, "positioning": "oneCell"}`)
    // if err != nil {
    //     fmt.Println(err)
    // }
}