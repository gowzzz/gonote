package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	// --------nsqd-----
	// "time"
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/nsqio/go-nsq"

	// "strconv"
	"encoding/base64"
	"errors"
	"io/ioutil"
)

type Watch struct {
	puppywatch *fsnotify.Watcher
	thirdwatch *fsnotify.Watcher
}

var NSQ_ADDR1, REDIS_HOST string
var permission, rename, delete, write, create int
var puppychan, thirdchan = make(chan bool, 100), make(chan bool, 100)

const (
	PUPPY = "puppy"
	THIRD = "third"
)

func main() {
	var wg = &sync.WaitGroup{}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	var puppypath, thirdpath string
	flag.StringVar(&NSQ_ADDR1, "nsqd", "127.0.0.1:4150", "monitor dir(default:127.0.0.1:4150)")
	flag.StringVar(&REDIS_HOST, "redis", "127.0.0.1:3306", "monitor dir(default:127.0.0.1:3306)")
	flag.StringVar(&puppypath, "puppypath", "./puppypath", "monitor dir(default:./puppypath)")
	flag.StringVar(&thirdpath, "thirdpath", "./thirdpath", "monitor dir(default:./thirdpath)")
	flag.IntVar(&create, "create", 1, "monitor file create ,1=open 0=close(default:1)   ")
	flag.IntVar(&permission, "permission", 0, "monitor file permission ,1=open 0=close(default:0)")
	flag.IntVar(&rename, "rename", 0, "monitor file rename ,1=open 0=close(default:0)")
	flag.IntVar(&delete, "delete", 0, "monitor file delete ,1=open 0=close(default:0)")
	flag.IntVar(&write, "write", 0, "monitor file write ,1=open 0=close(default:0)")
	flag.Parse()
	fmt.Println("now in monitor puppypath:", puppypath)
	fmt.Println("now in monitor thirdpath:", thirdpath)
	w := Watch{}
	var err error
	w.puppywatch, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("puppywatch:", err)
	}
	w.thirdwatch, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("thirdwatch:", err)
	}
	// w := Watch{
	// 	puppywatch: puppywatch,
	// 	thirdwatch: thirdwatch,
	// }
	wg.Add(2)
	go w.watchDir(ctx1, wg, puppypath, PUPPY)
	go w.watchDir(ctx2, wg, thirdpath, THIRD)
	<-sigs
	cancel1()
	cancel2()
	wg.Wait()
}

//监控目录
func (w *Watch) watchDir(ctx context.Context, wg *sync.WaitGroup, dir string, flag string) {
	var watch *fsnotify.Watcher
	var ch chan bool
	var processfunc func(string, chan bool)
	if flag == PUPPY {
		watch = w.puppywatch
		ch = puppychan
		processfunc = processFile
	} else if flag == THIRD {
		watch = w.thirdwatch
		ch = thirdchan
		processfunc = processThirdFile
	}
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("monitor dir : ", path)
		}
		return nil
	})
	for {
		select {
		case <-ctx.Done():
			{
				fmt.Println("watchpuppyDir exit")
				wg.Done()
				return
			}
		case ev := <-watch.Events:
			{
				if 1 == create {
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("create file : ", ev.Name)
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watch.Add(ev.Name)
							fmt.Println("add monitor : ", ev.Name)
						} else {
							go processfunc(ev.Name, ch)
						}
					}
				}
				if 1 == write {
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("write file : ", ev.Name)
					}
				}
				if 1 == delete {
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("delete file : ", ev.Name)
						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watch.Remove(ev.Name)
							fmt.Println("delete dir : ", ev.Name)
						}
					}
				}
				if 1 == rename {
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("rename file: ", ev.Name)
						//如果重命名文件是目录，则移除监控
						//注意这里无法使用os.Stat来判断是否是目录了
						//因为重命名后，go已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						watch.Remove(ev.Name)
					}
				}
				if 1 == permission {
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						fmt.Println("edit permission : ", ev.Name)
					}
				}
			}
		case err := <-watch.Errors:
			{
				fmt.Println("error : ", err)
				return
			}
		}
	}
}

func processFile(name string, ch chan bool) {
	time.Sleep(100 * time.Millisecond)
	ch <- true
	defer func() {
		os.Remove(name)
		<-ch
		fmt.Println("-----------------------")
	}()
	fmt.Println("processFile")
	if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".jpeg") {
		fmt.Println("processFile only process jpg or jpeg")
		return
	}
	// 最后一个是文件名 10.58.122.114_20190815155056857_1.jpg
	tmp1 := strings.Split(name, "/")
	fmt.Printf("path:%+v\n", tmp1)
	path := tmp1[len(tmp1)-1]
	//  10.58.122.114_20190815155056857_1.jpg => 10.58.122.114 20190815155056857 1.jpg
	tmp2 := strings.Split(path, "_")
	fmt.Printf("filename:%+v\n", tmp2)
	if len(tmp2) < 3 {
		fmt.Printf("filename is error:%+v,split tmp2:%+v\n", path, tmp2)
		return
	}
	ip := tmp2[0]
	timestring := tmp2[1][:14]
	// 20190815155056857
	timeLayout := "20060102150405"                                         //转化所需模板
	theTime, _ := time.ParseInLocation(timeLayout, timestring, time.Local) //使用模板在对应时区转化为time.time类型
	timestampstring := strconv.Itoa(int(theTime.Unix()))                   //转化为时间戳 类型是int64

	filename := tmp2[2]
	// filename[2] 解析出1.jpg
	tmp3 := strings.Split(filename, ".")
	fmt.Printf("flagname:%+v\n", tmp3)
	if len(tmp3) < 2 {
		fmt.Printf("filename is error:%+v,split tmp3:%+v\n", path, tmp3)
		return
	}
	flagname := tmp3[0]
	if flagname == "" {
		fmt.Println("filename(should be YT or number) is error:%+v", name)
		return
	} else if flagname == "YT" {
		fmt.Println("this is yuantu:", name)
	} else {
		// 人脸图
	}
	// CheckCameraByDeviceID
	if CheckCameraByIp(ip) {
		// GetCameraInfoByDeviceID
		cameraList, err := GetCameraInfoByIp(ip)
		if err != nil {
			fmt.Println("GetCameraInfoByIp err:", err)
			return
		}
		var r = new(record)
		//文件名字中截取出来的信息
		r.RequestTime = time.Now()
		r.CameraList = *cameraList
		r.CameraUuid = cameraList.UUID
		r.DayNight = "day"
		r.FaceBody = "face"
		r.ImageName = path
		data, err := ioutil.ReadFile(name)
		if err != nil {
			fmt.Println("ImageContent ReadFile err:", err)
			return
		}
		r.ImageContent = base64.StdEncoding.EncodeToString(data)
		r.Time = timestampstring
		r.sendToMq()
		fmt.Printf("record img:%+v\n", len(r.ImageContent))
		// fmt.Printf("record:%+v\n", r)
	}
}

// 第三方相机
/*
Face_设备ID_人脸ID_时间戳毫秒_图片质量_人脸置信度.jpg
TG_设备ID_人脸ID_时间戳毫秒.jpg
if strings.HasSuffix(name,".jpg"){}
*/
func processThirdFile(name string, ch chan bool) {
	ch <- true
	time.Sleep(100 * time.Millisecond)
	fmt.Println("processThirdFile")
	defer func() {
		os.Remove(name)
		<-ch
		fmt.Println("-----------------------")
	}()
	if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".jpeg") {
		fmt.Println("processThirdFile only process jpg or jpeg")
		return
	}
	tmp1 := strings.Split(name, "/")
	fmt.Printf("path:%+v\n", tmp1)
	// Face_设备ID_人脸ID_时间戳毫秒_图片质量_人脸置信度.jpg
	path := tmp1[len(tmp1)-1]
	// Face_设备ID_人脸ID_时间戳毫秒_图片质量_人脸置信度.jpg => Face 设备ID 人脸ID 时间戳毫秒 图片质量 人脸置信度.jpg
	tmp2 := strings.Split(path, "_")
	if len(tmp2) < 4 {
		fmt.Printf("filename is error:%+v,split result:%+v\n", name, tmp2)
		return
	}
	flagname := tmp2[0]
	devicename := tmp2[1]
	facename := tmp2[2]
	timestring:=""
	if len(tmp2[3])>10{
		timestring = tmp2[3][:10] //直接是时间戳
	}
	fmt.Println("flagname:", flagname)
	fmt.Println("devicename:", devicename)
	fmt.Println("facename:", facename)
	imagename := path //Face_设备ID_人脸ID_时间戳毫秒_图片质量_人脸置信度.jpg  Face_设备ID_人脸ID_时间戳毫秒_图片质量_人脸置信度.jpg
	if flagname == "Face" {
		imagename =tmp2[1]+"_"+tmp2[2]+"_"+tmp2[3] +"_" + facename + ".jpg"
		fmt.Println("imagename:",imagename)
	} else if flagname == "TG" {
		end:=tmp2[3]
		if len(end)>4{
			end=end[:(len(end)-4)]
		}
		imagename = tmp2[1]+"_"+tmp2[2]+"_"+ end + "_YT.jpg"
		fmt.Println("imagenametg:",imagename)
	} else {
		fmt.Printf("filename is error:%+v,flagname result:%+v\n", name, flagname)
	}
	//   CheckCameraByIp
	if CheckCameraByDeviceID(devicename) {
		// GetCameraInfoByIp
		cameraList, err := GetCameraInfoByDeviceID(devicename)
		if err != nil {
			fmt.Println("GetCameraInfoByIp err:", err)
			return
		}
		var r = new(record)
		//文件名字中截取出来的信息
		r.RequestTime = time.Now()
		r.CameraList = *cameraList
		r.CameraUuid = cameraList.UUID
		r.DayNight = "day"
		r.FaceBody = "face"
		r.ImageName = imagename
		data, err := ioutil.ReadFile(name)
		if err != nil {
			fmt.Println("ImageContent ReadFile err:", err)
			return
		}
		r.ImageContent = base64.StdEncoding.EncodeToString(data)
		r.Time = timestring
		r.sendToMq()
		fmt.Printf("record img:%+v\n", len(r.ImageContent))
		// fmt.Printf("record:%+v\n", r)
	} else {
		fmt.Println("camera not in redis:", name)
	}
}

// -----------------nsqd---------------
type record struct {
	CameraUuid   string     `json:"camera_uuid" form:"camera_uuid"`
	DayNight     string     `json:"day_night" form:"day_night"`
	FaceBody     string     `json:"face_body" form:"face_body"`
	ImageName    string     `json:"image_name" form:"image_name" binding:"required"`
	ImageContent string     `json:"image" form:"image" binding:"required"`
	Time         string     `json:"time" form:"time" binding:"required"` //不用了，用RequestTime代替
	RequestTime  time.Time  `json:"request_time"`                        //服务器时间
	CameraList   CameraList //redis获取

	// Path        string    `json:"path"`  //发送来的图片地址
	// Index       string    `json:"index"`
	// FromUrl     string    `json:"from_url"`
	// ToUrl       string    `json:"to_url"`
	// AccessToken  string `json:"access_token" form:"access_token" `
}
type CameraList struct {
	UUID       string
	IP         string
	ID         int32
	Name       string
	Status     int8
	Marked     int8
	Pos        string
	AreaName   string
	AreaType   string
	NVRuuid    string
	NVRchannel int32
}

// 其中CameraUuid  DayNight FaceBody ImageName ImageContent Time RequestTime CameraList必须
// 其中 DayNight FaceBody 无法获取
func (l *record) sendToMq() error {
	b, _ := json.Marshal(l)
	nsq, err := NewNsqConn()
	if err != nil {
		return errors.New("nsq conn err:" + err.Error())
	}
	if len(b) == 0 { //不能发布空串，否则会导致error
		return errors.New("send nsq messge must be not null")
	}
	err = nsq.Publish("logImage", b) // 发布消息
	if err != nil {
		return errors.New("nsq send msg  err:" + err.Error())
	}
	return nil
}

var lock *sync.Mutex = &sync.Mutex{}
var Producer1 *nsq.Producer

func NewNsqConn() (*nsq.Producer, error) {
	var err error
	if Producer1 == nil { //加锁是为了并发，加锁前判断是为了减少操作锁的消耗
		lock.Lock()
		defer lock.Unlock()
		if Producer1 == nil {
			Producer1, err = nsq.NewProducer(NSQ_ADDR1, nsq.NewConfig())
			if err != nil {
				fmt.Println("nsq connect err:", err)
				Producer1 = nil
			}
		}
	}

	return Producer1, err
	//只有main.go中才有必要		defer xgFaceConn.Close()
}

var devicesbyip = "devicesbyip"
var devicesbysn = "devicesbysn"

func CheckCameraByIp(cameraip string) bool {
	redisc, err := redis.Dial("tcp", REDIS_HOST)
	if err != nil {
		fmt.Println("connect redis err:", err)
		return false
	}
	defer redisc.Close()
	exit, err := redis.Int64(redisc.Do("HEXISTS", devicesbyip, cameraip))
	if err != nil {
		fmt.Println("redis get failed:", err)
		return false
	}
	if exit != 1 {
		return false
	}
	return true
}

func GetCameraInfoByIp(cameraip string) (*CameraList, error) {
	redisc, err := redis.Dial("tcp", REDIS_HOST)
	if err != nil {
		fmt.Println("connect redis err:", err)
		return nil, err
	}
	defer redisc.Close()
	cameralist, err := redis.Bytes(redisc.Do("hget", devicesbyip, cameraip))
	if err != nil {
		fmt.Println("redis get camerainfo failed:", err)
		return nil, err
	}
	var cl CameraList
	if err := json.Unmarshal(cameralist, &cl); err != nil {
		if err != nil {
			fmt.Println("cameralist json Unmarshal err:", err)
			return nil, err
		}
	}
	return &cl, nil
}

func CheckCameraByDeviceID(deviceid string) bool {
	redisc, err := redis.Dial("tcp", REDIS_HOST)
	if err != nil {
		fmt.Println("connect redis err:", err)
		return false
	}
	defer redisc.Close()
	exit, err := redis.Int64(redisc.Do("HEXISTS", devicesbysn, deviceid))
	if err != nil {
		fmt.Println("redis get failed:", err)
		return false
	}
	if exit != 1 {
		return false
	}
	return true
}

func GetCameraInfoByDeviceID(deviceid string) (*CameraList, error) {
	redisc, err := redis.Dial("tcp", REDIS_HOST)
	if err != nil {
		fmt.Println("connect redis err:", err)
		return nil, err
	}
	defer redisc.Close()
	cameralist, err := redis.Bytes(redisc.Do("hget", devicesbysn, deviceid))
	if err != nil {
		fmt.Println("redis get camerainfo failed:", err)
		return nil, err
	}
	var cl CameraList
	if err := json.Unmarshal(cameralist, &cl); err != nil {
		if err != nil {
			fmt.Println("cameralist json Unmarshal err:", err)
			return nil, err
		}
	}
	return &cl, nil
}

// //监控目录
// func (w *Watch) watchthirdDir(ctx context.Context, wg *sync.WaitGroup, dir string) {
// 	watch := w.thirdwatch
// 	//通过Walk来遍历目录下的所有子目录
// 	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
// 		//这里判断是否为目录，只需监控目录即可
// 		//目录下的文件也在监控范围内，不需要我们一个一个加
// 		if info.IsDir() {
// 			path, err := filepath.Abs(path)
// 			if err != nil {
// 				return err
// 			}
// 			err = watch.Add(path)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println("monitor dir : ", path)
// 		}
// 		return nil
// 	})
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			{
// 				fmt.Println("watchthirdDir exit")
// 				wg.Done()
// 				return
// 			}
// 		case ev := <-watch.Events:
// 			{
// 				if 1 == create {
// 					if ev.Op&fsnotify.Create == fsnotify.Create {
// 						fmt.Println("create file : ", ev.Name)
// 						//这里获取新创建文件的信息，如果是目录，则加入监控中
// 						fi, err := os.Stat(ev.Name)
// 						if err == nil && fi.IsDir() {
// 							watch.Add(ev.Name)
// 							fmt.Println("add monitor : ", ev.Name)
// 						} else {
// 							go processThirdFile(ev.Name)
// 						}
// 					}
// 				}
// 				if 1 == write {
// 					if ev.Op&fsnotify.Write == fsnotify.Write {
// 						fmt.Println("write file : ", ev.Name)
// 					}
// 				}
// 				if 1 == delete {
// 					if ev.Op&fsnotify.Remove == fsnotify.Remove {
// 						fmt.Println("delete file : ", ev.Name)
// 						//如果删除文件是目录，则移除监控
// 						fi, err := os.Stat(ev.Name)
// 						if err == nil && fi.IsDir() {
// 							watch.Remove(ev.Name)
// 							fmt.Println("delete dir : ", ev.Name)
// 						}
// 					}
// 				}
// 				if 1 == rename {
// 					if ev.Op&fsnotify.Rename == fsnotify.Rename {
// 						fmt.Println("rename file: ", ev.Name)
// 						//如果重命名文件是目录，则移除监控
// 						//注意这里无法使用os.Stat来判断是否是目录了
// 						//因为重命名后，go已经无法找到原文件来获取信息了
// 						//所以这里就简单粗爆的直接remove好了
// 						watch.Remove(ev.Name)
// 					}
// 				}
// 				if 1 == permission {
// 					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
// 						fmt.Println("edit permission : ", ev.Name)
// 					}
// 				}
// 			}
// 		case err := <-watch.Errors:
// 			{
// 				fmt.Println("error : ", err)
// 				return
// 			}
// 		}
// 	}

// }
