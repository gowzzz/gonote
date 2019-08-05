package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
	"time"

	"github.com/extrame/xls"

	"github.com/gin-gonic/gin"
)

type Import struct {
	Content string `form:"Content" json:"Content" xml:"Content"  binding:"required"`
}

func main() {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/import", func(c *gin.Context) {
			var imp Import
			if err := c.ShouldBind(&imp); err != nil {
				c.SecureJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			} else {
				ReadExcelFromBase64(imp.Content)
			}
		})
		v1.POST("/export", func(c *gin.Context) {
			c.SecureJSON(http.StatusOK, "export")
		})
	}

	// http.ListenAndServe(":8888", r)
	r.Run(":8888") // listen and serve on 0.0.0.0:8080
}

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
	var tRecords []TradingRecord
	// A 		B 					C  		 D 		E 		   F 		                    G 		                      H
	// 客户姓名 客户身份证号(15/18位) 置业顾问 经纪人 所属渠道公司 报备时间(2017-11-07 11:11:11) 成交时间(2017-11-07 11:11:11) 房号
	// var errtitile []string
	reader := bytes.NewReader(content)
	// if xlsx, err := excelize.OpenFile(fileName); err != nil {
	if xlFile, err := xls.OpenReader(reader, "utf-8"); err == nil {
		// xlFile.NumSheets() 表格个数
		sheet1 := xlFile.GetSheet(0)
		for i := 3; i <= int(sheet1.MaxRow); i++ {
			row := sheet1.Row(i)
			//时间转换的模板，golang里面只能是 "2006-01-02 15:04:05" （go的诞生时间）
			timeTemplate1 := "2006-1-2 15:4:5" //常规类型
			// for j := row.FirstCol(); j < row.LastCol(); j++ {
			// stamp, _ := time.ParseInLocation(timeTemplate1, t1, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
			var tRecord TradingRecord
			tRecord.Name = row.Col(0)                                                              // 客户姓名
			tRecord.CustomerIdNo = row.Col(1)                                                      // 客户身份证号(15/18位)PID
			tRecord.PropertyConsultant = row.Col(2)                                                // 置业顾问
			tRecord.ReferrerName = row.Col(3)                                                      // 经纪人姓名
			tRecord.ReferrerChannelName = row.Col(4)                                               // 所属渠道公司
			tRecord.RecommendedAt, _ = time.ParseInLocation(timeTemplate1, row.Col(5), time.Local) // 报备时间(2017-11-07 11:11:11)
			tRecord.TradedAt, _ = time.ParseInLocation(timeTemplate1, row.Col(6), time.Local)      // 成交时间(2017-11-07 11:11:11)
			tRecord.TradeRoomName = row.Col(7)
			tRecords = append(tRecords, tRecord) // 成交房间
			// 未导入的新
			// Phone               string    // 客户手机号
			// FirstCapturedAt     time.Time // 首次人脸抓拍时间
			// ValidatedAt         time.Time // 刷证时间
			// ReferrerMobile      string    // 经纪人手机号
			// RiskTag             int       // 疑似风险交易 0非风险 1风险 //suspected_face
			// DecideRiskTag       int       // 判定风险交易 0非风险 1风险
			// Comment             string    // 判定风险交易时候的备注备注
		}
		fmt.Printf("tRecords:%+v\n", tRecords[0])
	} else {
		panic(err)
	}
}
