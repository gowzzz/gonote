package excel

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

/*
 */
// 成交房间 客户姓名 客户手机号 首次人脸抓拍时间 报备时间 刷证时间 成交时间(2019/6/26) 经纪人姓名 经纪人号码 渠道身份 异常风险指标
type TradingRecord struct {
	Id                  int
	TradeRoomName       string    // 成交房间
	Name                string    // 客户姓名
	Phone               string    // 客户手机号
	CustomerIdNo        string    // 客户身份证号(15/18位)PID
	FirstCapturedAt     time.Time // 首次人脸抓拍时间
	RecommendedAt       time.Time // 报备时间(2017-11-07 11:11:11)
	ValidatedAt         time.Time // 刷证时间
	TradedAt            time.Time // 成交时间(2017-11-07 11:11:11)
	ReferrerName        string    // 经纪人姓名
	ReferrerMobile      string    // 经纪人手机号
	ReferrerChannelName string    // 所属渠道公司
	PropertyConsultant  string    // 置业顾问
	RiskTag             int       // 疑似风险交易 0非风险 1风险 //suspected_face
	DecideRiskTag       int       // 判定风险交易 0非风险 1风险
	Comment             string    // 判定风险交易时候的备注备注
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

/*
房号：house No
成交时间 Deal time
报备时间Filing time
经纪人 broker
客户： Customer
置业顾问(Property consultant)
渠道公司 Channel company
字段
id
*/
/*
第一行：标题
第二行：项目名称 ：项目:惠州常春藤
第三行：列名： 8列
客户姓名 客户身份证号(15/18位)PID 置业顾问 经纪人 所属渠道公司 报备时间(2017-11-07 11:11:11) 成交时间(2017-11-07 11:11:11) 房号

*/
// A 		B 					C  		 D 		E 		   F 		                    G 		                      H
// 客户姓名 客户身份证号(15/18位) 置业顾问 经纪人 所属渠道公司 报备时间(2017-11-07 11:11:11) 成交时间(2017-11-07 11:11:11) 房号
func NewExpErrPeopleFile(saveName string, title []string) (*excelize.File, int64) {
	PathExists(saveName)
	st, _ := FileExists(saveName)
	if st {
		os.Remove(saveName)
	}
	// 创建excel表格
	xlsx := excelize.NewFile()
	// 生成标题
	for i := 0; i < len(title); i++ {
		//65==A   A1~N1
		xlsx.SetCellValue("Sheet1", string(65+i)+"1", title[i])
	}
	return xlsx, 0
}

func ExportErrPeople(xlsx *excelize.File, num int, row []string, imgcontent []byte, errinfo, errpath string) {
	// 第一个num是1
	rownum := num + 1
	sheet := "Sheet1"
	// A 		B 					C  		 D 		E 		   F 		                    G 		                      H		I
	// 客户姓名 客户身份证号(15/18位) 置业顾问 经纪人 所属渠道公司 报备时间(2017-11-07 11:11:11) 成交时间(2017-11-07 11:11:11) 房号	 错误信息
	xlsx.SetCellValue(sheet, "A"+strconv.Itoa(rownum), row[0])  //从第二行开始
	xlsx.SetCellValue(sheet, "B"+strconv.Itoa(rownum), row[2])  //从第二行开始
	xlsx.SetCellValue(sheet, "C"+strconv.Itoa(rownum), row[2])  //从第二行开始
	xlsx.SetCellValue(sheet, "D"+strconv.Itoa(rownum), row[3])  //从第二行开始
	xlsx.SetCellValue(sheet, "E"+strconv.Itoa(rownum), row[4])  //从第二行开始
	xlsx.SetCellValue(sheet, "F"+strconv.Itoa(rownum), row[5])  //从第二行开始
	xlsx.SetCellValue(sheet, "G"+strconv.Itoa(rownum), row[6])  //从第二行开始
	xlsx.SetCellValue(sheet, "H"+strconv.Itoa(rownum), row[7])  //从第二行开始
	xlsx.SetCellValue(sheet, "I"+strconv.Itoa(rownum), errinfo) //从第二行开始
}
func ReadExcelFromBase64(str string) {
	/*xlsx data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,*/
	/*doc data:application/vnd.openxmlformats-officedocument.wordprocessingml.document;base64,*/
	/*pdf data:application/pdf;base64,*/
	/*xls data:application/vnd.ms-excel;base64,*/
	/*rtf data:application/msword;base64,*/
	/*data:application/vnd.ms-excel;base64,*/
	if len(str) == 0 {
		return
	}

	splitstr := strings.Split(str, ",")
	if len(splitstr) != 2 {
		panic(len(splitstr))
		return
	}
	base64Content := splitstr[1]
	content, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		panic(err)
		return
	}
	// A 		B 					C  		 D 		E 		   F 		                    G 		                      H
	// 客户姓名 客户身份证号(15/18位) 置业顾问 经纪人 所属渠道公司 报备时间(2017-11-07 11:11:11) 成交时间(2017-11-07 11:11:11) 房号
	var errtitile []string
	reader := bytes.NewReader(content)
	// if xlsx, err := excelize.OpenFile(fileName); err != nil {
	if xlsx, err := excelize.OpenReader(reader); err != nil {
		panic(err)
		return
	} else {
		// 导出错误
		errfileName := "渠道风控导入成交客户失败信息_" + time.Now().Format("20060102_150405") + ".csv"
		saveErrName := errfileName
		// ErrFilePath := errfileName
		rows, _ := xlsx.GetRows("Sheet1")
		errtitile = append(rows[0], "错误信息")
		// errxlsx, _ := NewExpErrPeopleFile(saveErrName, errtitile)
		NewExpErrPeopleFile(saveErrName, errtitile)

		// Failure := make(map[int64]string)
		// Success := make(map[int64]string)
		// errnum := 0
		for k := 1; k < len(rows); k++ {
			fmt.Printf("rows:%+v", rows[k])
			// row := rows[k]
			// if row[4] != "" {
			// 	age, err = strconv.Atoi(row[4])
			// 	if err != nil {
			// 		errinfo := "年龄必须是0-150的数字"
			// 		Failure[int64(k)] = errinfo
			// 		errnum++
			// 		ExportErrPeople(errxlsx, errnum, row, content, errinfo, "")
			// 		continue
			// 	}
			// }

			// people := &TradingRecord{
			// 	// Name:    row[2],
			// 	// Gender:  row[3],
			// 	// Age:     int64(age),
			// 	// IdNo:    row[5],
			// 	// Tel:     row[6],
			// 	// Comment: row[7],
			// }
			// imgids := []int64{img.ID}
			// var groupids []int64
			// if groupid != 0 {
			// 	groupids = []int64{groupid}
			// }
			// if _, _, _, code := people.AddPeople(imgids, groupids); code != 0 {
			// 	errinfo := "新增人员时报错:" + GetMsg(code)
			// 	Failure[int64(k)] = errinfo
			// 	errnum++
			// 	ExportErrPeople(errxlsx, errnum, row, content, errinfo, "")
			// 	continue
			// }
			// fmt.Printf("k:%+v\n", k)
			// fmt.Printf("row:%+v\n", row)
			// Success[int64(k)] = "ok"
		}
		// errxlsx.SaveAs(saveErrName)
		// time.AfterFunc(24*time.Hour, func() {
		// 	os.RemoveAll(saveErrName)
		// })
	}
}

/*
渠道成交记录_20190627101208.csv
成交房间 客户姓名 客户手机号 首次人脸抓拍时间 报备时间 刷证时间 成交时间(2019/6/26) 经纪人姓名 经纪人号码 渠道身份

疑似风险交易_20190627101212.csv  多了个 异常风险指标
成交房间 客户姓名 客户手机号 首次人脸抓拍时间 报备时间 刷证时间 成交时间(2019/6/26) 经纪人姓名 经纪人号码 渠道身份 异常风险指标

已处理交易_已判定有风险交易_20190627101215.csv
成交房间 客户姓名 客户手机号 首次人脸抓拍时间 报备时间 刷证时间 成交时间(2019/6/26) 经纪人姓名 经纪人号码 渠道身份 异常风险指标

目前不需要：刷证后未匹配到到访抓拍客户_20190627101818.csv
*/
func SaveExcelToFile(datas []TradingRecord) (string, error) {
	fileName := "人员导出信息_" + time.Now().Format("20060102_150405") + ".xlsx"
	saveName := fileName
	accessName := fileName
	// 存在的话删掉，防止已存在的内容影响导出
	PathExists(saveName)
	st, _ := FileExists(saveName)
	if st {
		os.Remove(saveName)
	}
	// 创建excel表格
	xlsx := excelize.NewFile()
	defer xlsx.SaveAs(saveName)
	// 生成标题
	title := []string{"成交房间", "客户姓名", "客户手机号", "首次人脸抓拍时间", "报备时间", "刷证时间", "成交时间", "经纪人姓名", "经纪人号码", "渠道身份", "异常风险指标"}
	for i := 0; i < len(title); i++ {
		//65==A   A1~N1
		xlsx.SetCellValue("Sheet1", string(65+i)+"1", title[i])
	}
	sheet := "Sheet1"

	for k, data := range datas {
		rownum := k + 2
		// (2019/6/26)
		// A 		B 		C  			D 			   E 		F 		G 		 H		 I		   J			K
		// 成交房间 客户姓名 客户手机号  首次人脸抓拍时间 报备时间 刷证时间 成交时间 经纪人姓名 经纪人号码 渠道身份 异常风险指标
		xlsx.SetCellValue(sheet, "A"+strconv.Itoa(rownum), data.TradeRoomName) //从第二行开始
		// xlsx.SetCellValue(sheet, "B"+strconv.Itoa(rownum), data.CustomerName)  //从第二行开始
		// xlsx.SetCellValue(sheet, "C"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "D"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "E"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "F"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "G"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "H"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "I"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "J"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
		// xlsx.SetCellValue(sheet, "K"+strconv.Itoa(rownum), data.CustomerPhone) //从第二行开始
	}
	time.AfterFunc(24*time.Hour, func() {
		os.RemoveAll(saveName)
	})
	return accessName, nil
}

//判断目录是否存在，不存在则创建，创建失败则返回错误信息
func PathExists(path string) (exist bool) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)
			if err != nil {
				log.Println("PathExists file err:", err)
				return false
			}
		}
		return false
	}
	return true
}
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, err
}
