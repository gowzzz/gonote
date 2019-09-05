package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Animal struct {
	Id   int
	Name string
	Age  int
}
type Animals []Animal

func (a Animals) Len() int { return len(a) }
func (s Animals) Less(i, j int) bool {
	a, _ := UTF82GBK(s[i].Name)
	b, _ := UTF82GBK(s[j].Name)
	bLen := len(b)
	for idx, chr := range a {
		if idx > bLen-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}
func (a Animals) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

//UTF82GBK : transform UTF8 rune into GBK byte array
func UTF82GBK(src string) ([]byte, error) {
	GB18030 := simplifiedchinese.All[0]
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), GB18030.NewEncoder()))
}

//GBK2UTF8 : transform  GBK byte array into UTF8 string
func GBK2UTF8(src []byte) (string, error) {
	GB18030 := simplifiedchinese.All[0]
	bytes, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), GB18030.NewDecoder()))
	return string(bytes), err
}

type AlarmMsg struct {
	AlarmUuid    string `json:"uuid"`
	Clusterid    int    `json:"clusterid"`
	CameraUuid   string `json:"cameraUuid"`
	CameraName   string `json:"cameraName"`
	CameraIp     string `json:"cameraIp"`
	AreaName     string `json:"areaName"`
	AreaType     string `json:"areaType"`
	Posstr       string `json:"posstr"`
	Nvruuid      string `json:"nvruuid"`
	Nvrchannel   int32  `json:"nvrchannel"`
	OriginalPath string `json:"originalPath"`
	ImgPath      string `json:"imagePath"`
	AlarmType    string `json:"category"`
	Comment      string `json:"comment"`
	AlarmTime    int    `json:"logtime"`
}

func main() {
	// aMsg := AlarmMsg{
	// 	AlarmUuid: "dsadasdsa",
	// }
	// msg, _ := json.Marshal(aMsg)
	// fmt.Printf("msg:%+v\n", string(msg))

	// return

	// fmt.Println(pinyin.LazyConvert("1奥啊2", nil))

	// animals := []string{"cat", "bird", "zebra", "fox"}
	// // Sort by strings.
	// sort.Strings(animals)
	// fmt.Println(animals) //[bird cat fox zebra]

	//sort by len
	an := Animals{
		Animal{Id: 1, Name: "请求", Age: 11},
		Animal{Id: 2, Name: "当当", Age: 22},
		Animal{Id: 3, Name: "呃呃", Age: 33},
		Animal{Id: 4, Name: "z奥啊", Age: 44},
		Animal{Id: 5, Name: "宝宝z", Age: 55},
		Animal{Id: 6, Name: "宝宝a", Age: 6},
	}
	sort.Sort(an)
	fmt.Printf("an:%+v\n", an) //[cat fox bird zebra]
	return

	data, err := ioutil.ReadFile("./1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	str := base64.StdEncoding.EncodeToString(data)
	fmt.Println("str:", str)

	return
	toBeCharge := "20190815155056"                                         //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "20060102150405"                                         //转化所需模板
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, time.Local) //使用模板在对应时区转化为time.time类型
	sr := strconv.Itoa(int(theTime.Unix()))                                //转化为时间戳 类型是int64
	fmt.Println("theTime:", theTime)                                       //打印输出theTime 2015-01-01 15:15:00 +0800 CST
	fmt.Println("sr:", sr)                                                 //打印输出时间戳 1420041600

	// //时间戳转日期
	// dataTimeStr := time.Unix(sr, 0).Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	// fmt.Println(dataTimeStr)
	return
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	line := make(chan int, 1)
	exit := make(chan bool, 1)
	go Producer(line)
	go WorkFunc(ctx, line, exit)
	<-sigs
	cancel()
	<-exit
}

// 工作协程
func WorkFunc(ctx context.Context, line chan int, exit chan bool) {
	for {
		select {
		case n := <-line:
			log.Println("work start:", n)
			time.Sleep(1 * time.Second)
			log.Println("work done:", n)
		case <-ctx.Done():
			log.Println("exit")
			goto EXIT
		}
	}
EXIT:
	exit <- true
}

// 生产协程
func Producer(line chan int) {
	for i := 0; i < 10; i++ {
		line <- i
		time.Sleep(time.Second)
	}
}
