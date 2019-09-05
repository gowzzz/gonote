package main
// import(
// 	"time"
// 	"github.com/nsqio/go-nsq"
// 	"encoding/json"
// 	"github.com/garyburd/redigo/redis"

// )

// type record struct {
// 	CameraUuid   string `json:"camera_uuid" form:"camera_uuid"`
// 	DayNight     string `json:"day_night" form:"day_night"`
// 	FaceBody     string `json:"face_body" form:"face_body"`
// 	ImageName    string `json:"image_name" form:"image_name" binding:"required"`
// 	ImageContent string `json:"image" form:"image" binding:"required"`
// 	Time         string `json:"time" form:"time" binding:"required"` //不用了，用RequestTime代替
// 	RequestTime time.Time `json:"request_time"` //服务器时间
// 	CameraList  CameraList  //redis获取

// 	// Path        string    `json:"path"`  //发送来的图片地址
// 	// Index       string    `json:"index"` 
// 	// FromUrl     string    `json:"from_url"`
// 	// ToUrl       string    `json:"to_url"`
// 	// AccessToken  string `json:"access_token" form:"access_token" `
// }
// type CameraList struct {
// 	UUID     string
// 	IP       string
// 	ID       int32
// 	Name     string
// 	Status   int8
// 	Marked   int8
// 	Pos      string
// 	AreaName string
// 	AreaType string
// 	NVRuuid  string
// 	NVRchannel int32
// }
// // 其中CameraUuid  DayNight FaceBody ImageName ImageContent Time RequestTime CameraList必须
// // 其中 DayNight FaceBody 无法获取
// func (l *record) sendToMq() (error) {
// 	b, _ := json.Marshal(l)
// 	nsq, err := NewNsqConn()
// 	if err != nil {
// 		return errors.New("nsq conn err:"+err.Error())
// 	}
// 	if len(b) == 0 { //不能发布空串，否则会导致error
// 		return errors.New("send nsq messge must be not null")
// 	}
// 	err = nsq.Publish("logImage", b) // 发布消息
// 	if err != nil {
// 		return  errors.New("nsq send msg  err:"+err.Error())
// 	}
// 	return nil
// }
// var Producer1 *nsq.Producer
// func NewNsqConn() (*nsq.Producer, error) {
// 	if Producer1 == nil { //加锁是为了并发，加锁前判断是为了减少操作锁的消耗
// 		lock.Lock()
// 		defer lock.Unlock()
// 		if Producer1 == nil {
// 			Producer1, err = nsq.NewProducer(NSQ_ADDR1, nsq.NewConfig())
// 			if err != nil {
// 				log.Println("nsq connect err:", err)
// 				Producer1 = nil
// 			}
// 		}
// 	}

// 	return Producer1, err
// 	//只有main.go中才有必要		defer xgFaceConn.Close()
// }

// var devicesbyip="devicesbyip"
// var devicesbysn="devicesbysn"
// func CheckCameraByIp(cameraip string)bool{
// 	redisc, err := redis.Dial("tcp", REDIS_HOST)
// 	if err != nil {
// 		fmt.Println("connect redis err:", err)
// 		return false
// 	}
// 	defer redisc.Close()
// 	exit, err := redis.Int64(redisc.Do("HEXISTS", devicesbyip, cameraip))
// 	if err != nil {
// 		fmt.Println("redis get failed:", err)
// 		return false
// 	}
// 	if exit!=1{
// 		return false
// 	}
// 	return true
// }


// func GetCameraInfoByIp(cameraip string)(*CameraList,error){
// 	redisc, err := redis.Dial("tcp", REDIS_HOST)
// 	if err != nil {
// 		fmt.Println("connect redis err:", err)
// 		return nil,err
// 	}
// 	defer redisc.Close()
// 	cameralist, err := redis.Bytes(redisc.Do("hget", devicesbyip, cameraip))
// 	if err != nil {
// 		log.Println("redis get camerainfo failed:", err)
// 		return nil,err
// 	}
// 	var cl CameraList
// 	if err := json.Unmarshal(cameralist, &cl); err != nil {
// 		if err != nil {
// 			fmt.Println("cameralist json Unmarshal err:", err)
// 			return nil,err
// 		}
// 	} 
// 	return cl,nil
// }

// func CheckCameraByDeviceID(deviceid string)bool{
// 	redisc, err := redis.Dial("tcp", REDIS_HOST)
// 	if err != nil {
// 		fmt.Println("connect redis err:", err)
// 		return false
// 	}
// 	defer redisc.Close()
// 	exit, err := redis.Int64(redisc.Do("HEXISTS", devicesbysn, deviceid))
// 	if err != nil {
// 		fmt.Println("redis get failed:", err)
// 		return false
// 	}
// 	if exit!=1{
// 		return false
// 	}
// 	return true
// }


// func GetCameraInfoByDeviceID(deviceid string)(*CameraList){
// 	redisc, err := redis.Dial("tcp", REDIS_HOST)
// 	if err != nil {
// 		fmt.Println("connect redis err:", err)
// 		return nil
// 	}
// 	defer redisc.Close()
// 	cameralist, err := redis.Bytes(redisc.Do("hget", devicesbysn, deviceid))
// 	if err != nil {
// 		log.Println("redis get camerainfo failed:", err)
// 		return nil
// 	}
// 	var cl CameraList
// 	if err := json.Unmarshal(cameralist, &cl); err != nil {
// 		if err != nil {
// 			fmt.Println("cameralist json Unmarshal err:", err)
// 			return nil
// 		}
// 	} 
// 	return cl
// }