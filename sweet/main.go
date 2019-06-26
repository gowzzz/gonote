package main

import (
	"fmt"
	"net/http"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	topic     = "mqtt"
	broker    = "tcp://118.89.137.31:1883"
	password  = ""
	user      = ""
	id        = ""
	cleansess = false
	qos       = 0
	num       = 1
	payload   = ""
	action    = "sub"
	store     = ":memory:"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {

	router := gin.Default()
	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Sweet Kite",
		})
	})

	router.GET("/websocket", func(c *gin.Context) {
		//升级get请求为webSocket协议
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()
		for {
			//读取ws中的数据
			mt, message, err := ws.ReadMessage()
			if err != nil {
				break
			}
			if len(message) != 0 {
				sendMQTT(broker, topic, string(message))
			}
			if string(message) == "ping" {
				message = []byte("pong")
			}
			//写入ws数据
			err = ws.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	})

	router.Run(":9001")
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendMQTT(addr, topic, msg string) {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(addr)
	opts.SetClientID("go-simple")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//订阅主题 aaa 并请求消息以最高qos为0交付，等待收据确认订阅
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	//在qos 1中将消息发布到/aaa，并在发送每条消息后等待来自服务器的接收
	token := c.Publish(topic, 0, false, msg)
	token.Wait()

	//unsubscribe from aaa
	if token := c.Unsubscribe("aaa"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)

	return
}
