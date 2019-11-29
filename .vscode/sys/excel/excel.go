package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/smtp"
	"regexp"
	"sort"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jordan-wright/email"
)

func main() {
	resp := createExcel()
	decodeStr, err := base64.StdEncoding.DecodeString(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = SendEmail(bytes.NewBuffer(decodeStr))
	if err != nil {
		fmt.Println(err)
		return
	}
}
func createExcel() string {
	sheet := "Sheet1"
	f := excelize.NewFile()
	// col3DCylinder
	// colStacked  colPercentStacked col3DStacked col3DPercentStacked
	// col3D col3DClustered

	// -第二部分
	var datas1 = map[string]int{"0-2": 262, "2-4": 760, "4-6": 999, "6-8": 511, "8-10": 122, "10-12": 522, "12-14": 266, "14-16": 277, "16-18": 777, "18-20": 511, "20-22": 300, "22-0": 200}
	var excelDatas1 ExcelDatas
	for camera, num := range datas1 {
		excelDatas1 = append(excelDatas1, ExcelData{
			Name:  camera,
			Value: num,
		})
	}
	sort.SliceStable(excelDatas1, func(i, j int) bool {
		a, _ := strconv.Atoi(regexp.MustCompile(`[0-9]+`).FindString(excelDatas1[i].Name))
		b, _ := strconv.Atoi(regexp.MustCompile(`[0-9]+`).FindString(excelDatas1[j].Name))
		// return a > b // 降序
		return a < b // 升序
	})
	f, _ = CreateChart(sheet, "相机抓拍", 1, excelDatas1, f)

	// -------------------------------

	var datas2 = map[string]int{"192.168.1.1": 100, "192.168.1.2": 90, "192.168.1.3": 80, "192.168.1.4": 20, "192.168.1.5": 120, "192.168.1.6": 90}
	var excelDatas2 ExcelDatas
	for camera, num := range datas2 {
		excelDatas2 = append(excelDatas2, ExcelData{
			Name:  camera,
			Value: num,
		})
	}
	sort.SliceStable(excelDatas2, func(i, j int) bool {
		return excelDatas2[i].Value > excelDatas2[j].Value // 降序
		// return lstPerson[i].Age < lstPerson[j].Age  // 升序
	})

	f, _ = CreateChart(sheet, "抓拍统计", 21, excelDatas2, f)
	// 保存工作簿
	buffer := bytes.NewBuffer(nil)

	num, err := f.WriteTo(buffer) //写到缓存或者任何地方
	// err := f.SaveAs("./test3.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("num:", num)
	fmt.Println("buffer:", base64.StdEncoding.EncodeToString(buffer.Bytes()))
	// file, err := os.Create("./111.xlsx")
	// if err != nil {
	// 	fmt.Println("os open:", err)
	// }
	// buffer.WriteTo(file)
	// SendEmail(buffer)
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}
func SendEmail(r io.Reader) error {
	fromUser := "golang<371600645@qq.com>"
	toUser := "810169879@qq.com"
	subject := "hello,world"
	// NewEmail返回一个email结构体的指针
	e := email.NewEmail()
	// 发件人
	e.From = fromUser
	// 收件人(可以有多个)
	e.To = []string{toUser}
	// 邮件主题
	e.Subject = subject
	e.Attach(r, "1.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// 从缓冲中将内容作为附件到邮件中
	// e.Attach(r, "1.xlsx", "application/vnd.ms-excel")

	return e.Send("smtp.qq.com:587", smtp.PlainAuth("用户名", "371600645@qq.com", "kfhjashvuvxrcabj", "smtp.qq.com"))
}

type ExcelData struct {
	Name  string
	Value int
}
type ExcelDatas []ExcelData

func CreateChart(sheet, titleName string, startline int, excelDatas ExcelDatas, f *excelize.File) (*excelize.File, error) {
	// sheet := "Sheet1"
	// // 标题放在string(65)/A，那么数据只能从string(65+1)/B开始，key在第一行，value就在第 line+1 行
	// titleName := "相机抓拍"
	// startline := 1
	titleCol, titleRow := 65, strconv.Itoa(startline+1) //A2
	// 图表位置 A10
	chartLine := string(titleCol) + strconv.Itoa(startline+2)
	// 数据  key在XX1  valueXX2 从B开始
	dataColStart, dataColEnd := titleCol+1, titleCol+len(excelDatas) //起始位置 结束位置 B~xx
	dataKeyRow := strconv.Itoa(startline)                            //B1~XX1
	dataValRow := strconv.Itoa(startline + 1)                        //B2~XX2

	keyInit := 65 + 1                                           //B
	f.SetCellValue(sheet, string(titleCol)+titleRow, titleName) //在数据前些标题
	for i := 0; i < len(excelDatas); i++ {
		col := string(keyInit + i)
		f.SetCellValue(sheet, col+dataKeyRow, excelDatas[i].Name)
		f.SetCellValue(sheet, col+dataValRow, excelDatas[i].Value)
	}
	// 图表 J10
	err := f.AddChart(sheet, chartLine, `{
		"type":"col",
		"series":[
			{
				"name":"`+sheet+`!$`+string(titleCol)+`$`+titleRow+`",
				"categories":"`+sheet+`!$`+string(dataColStart)+`$`+dataKeyRow+`:$`+string(dataColEnd)+`$`+dataKeyRow+`",
				"values":"`+sheet+`!$`+string(dataColStart)+`$`+dataValRow+`:$`+string(dataColEnd)+`$`+dataValRow+`"
			}
		],
		"format":{
			"x_scale":1.0,
			"y_scale":1.0,
			"x_offset":15,
			"y_offset":10,
			"print_obj":true,
			"lock_aspect_ratio":false,"locked":false
		},
		"legend":{
			"position":"left",
			"show_legend_key":false
		},
		"title":{
			"name":"`+titleName+`"
		},
		"plotarea":{
			"show_bubble_size":true,
			"show_cat_name":false,
			"show_leader_lines":false,
			"show_percent":true,
			"show_series_name":false,
			"show_val":true
		},
		"show_blanks_as":"zero"
	}`)
	if err != nil {
		fmt.Println(err)
	}
	return f, err
}

// func CreateCameraCaptureNum(f *excelize.File) (*excelize.File, error) {
// 	sheet := "Sheet1"
// 	// 标题放在string(65)/A，那么数据只能从string(65+1)/B开始，key在第一行，value就在第 line+1 行
// 	titleName := "相机抓拍"
// 	titleCol, titleRow := 65, "2" //A2
// 	// 图表位置 A10
// 	chartLine := "A10"
// 	// 数据  key在XX1  valueXX2 从B开始
// 	dataColStart, dataColEnd := titleCol+1, titleCol+len(cameras) //起始位置 结束位置 B~xx
// 	dataKeyRow := "1"                                             //B1~XX1
// 	dataValRow := "2"                                             //B2~XX2

// 	var excelDatas ExcelDatas
// 	for camera, num := range cameras {
// 		excelDatas = append(excelDatas, ExcelData{
// 			Name:  camera,
// 			Value: num,
// 		})
// 	}
// 	sort.SliceStable(excelDatas, func(i, j int) bool {
// 		return excelDatas[i].Value > excelDatas[j].Value // 降序
// 		// return lstPerson[i].Age < lstPerson[j].Age  // 升序
// 	})

// 	keyInit := 65 + 1                                           //B
// 	f.SetCellValue(sheet, string(titleCol)+titleRow, titleName) //在数据前些标题
// 	for i := 0; i < len(excelDatas); i++ {
// 		col := string(keyInit + i)
// 		f.SetCellValue(sheet, col+"4", excelDatas[i].Name)
// 		f.SetCellValue(sheet, col+"5", excelDatas[i].Value)
// 		// excelDatas[i].KCol = col + "4"
// 		// // B5
// 		// excelDatas[i].VCol = col + "5"
// 		// fmt.Println(col)
// 	}
// 	// 图表 J10
// 	err := f.AddChart(sheet, chartLine, `{
// 		"type":"col",
// 		"series":[
// 			{
// 				"name":"`+sheet+`!$`+string(titleCol)+`$`+titleRow+`",
// 				"categories":"`+sheet+`!$`+string(dataColStart)+`$`+dataKeyRow+`:$`+string(dataColEnd)+`$`+dataKeyRow+`",
// 				"values":"`+sheet+`!$`+string(dataColStart)+`$`+dataValRow+`:$`+string(dataColEnd)+`$`+dataValRow+`"
// 			}
// 		],
// 		"format":{
// 			"x_scale":1.0,
// 			"y_scale":1.0,
// 			"x_offset":15,
// 			"y_offset":10,
// 			"print_obj":true,
// 			"lock_aspect_ratio":false,"locked":false
// 		},
// 		"legend":{
// 			"position":"left",
// 			"show_legend_key":false
// 		},
// 		"title":{
// 			"name":"`+titleName+`"
// 		},
// 		"plotarea":{
// 			"show_bubble_size":true,
// 			"show_cat_name":false,
// 			"show_leader_lines":false,
// 			"show_percent":true,
// 			"show_series_name":false,
// 			"show_val":true
// 		},
// 		"show_blanks_as":"zero"
// 	}`)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return f, err
// }
// func CreateTimeCaptureNum(f *excelize.File) (*excelize.File, error) {

// 	var excelDatas ExcelDatas
// 	for data, num := range datas {
// 		excelDatas = append(excelDatas, ExcelData{
// 			Name:  data,
// 			Value: num,
// 		})
// 	}
// 	sort.SliceStable(excelDatas, func(i, j int) bool {
// 		a, _ := strconv.Atoi(regexp.MustCompile(`[0-9]+`).FindString(excelDatas[i].Name))
// 		b, _ := strconv.Atoi(regexp.MustCompile(`[0-9]+`).FindString(excelDatas[j].Name))
// 		// return a > b // 降序
// 		return a < b // 升序
// 	})

// 	// 起始位置 结束位置
// 	start, end := "B", string(65+len(excelDatas))
// 	sheet := "Sheet1"
// 	keyInit := 65 + 1 //B
// 	f.SetCellValue(sheet, "A2", "抓拍统计")
// 	for i := 0; i < len(excelDatas); i++ {
// 		col := string(keyInit + i)
// 		f.SetCellValue(sheet, col+"1", excelDatas[i].Name)
// 		f.SetCellValue(sheet, col+"2", excelDatas[i].Value)
// 		// // B4
// 		// excelDatas[i].KCol = col + "4"
// 		// // B5
// 		// excelDatas[i].VCol = col + "5"
// 		// fmt.Println(col)
// 	}

// 	err := f.AddChart(sheet, "A10", `{
// 		"type":"col",
// 		"series":[
// 			{
// 				"name":"`+sheet+`!$`+start+`$2",
// 				"categories":"`+sheet+`!$`+start+`$1:$`+end+`$1",
// 				"values":"`+sheet+`!$`+start+`$2:$`+end+`$2"
// 			}
// 		],
// 		"format":{
// 			"x_scale":1.0,
// 			"y_scale":1.0,
// 			"x_offset":15,
// 			"y_offset":10,
// 			"print_obj":true,
// 			"lock_aspect_ratio":false,"locked":false
// 		},
// 		"legend":{
// 			"position":"left",
// 			"show_legend_key":false
// 		},
// 		"title":{
// 			"name":"抓拍统计"
// 		},
// 		"plotarea":{
// 			"show_bubble_size":true,
// 			"show_cat_name":false,
// 			"show_leader_lines":false,
// 			"show_percent":true,
// 			"show_series_name":false,
// 			"show_val":true
// 		},
// 		"show_blanks_as":"zero"
// 	}`)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return f, err
// }
