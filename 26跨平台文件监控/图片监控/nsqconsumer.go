package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

/*
err := initConsumer("test1", "test-channel1", "127.0.0.1:4161")
if err != nil {
	log.Fatal("init Consumer error")
}
*/
func RunConsumer(addr, topic, channel string) {
	err := initConsumer(addr, topic, channel)
	if err != nil {
		log.Fatal("init Consumer error")
	}
}

type nsqHandler struct {
	nsqConsumer      *nsq.Consumer
	messagesReceived int
}

//处理消息
func (nh *nsqHandler) HandleMessage(msg *nsq.Message) error {
	nh.messagesReceived++
	fmt.Printf("receive ID:%s,addr:%s,message:%s", msg.ID, msg.NSQDAddress, string(msg.Body))
	fmt.Println()
	return nil
}

func initConsumer(addr, topic, channel string) error {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 3 * time.Second
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Println("init Consumer NewConsumer error:", err)
		return err
	}

	handler := &nsqHandler{nsqConsumer: c}
	c.AddHandler(handler)

	err = c.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Println("init Consumer ConnectToNSQLookupd error:", err)
		return err
	}
	return nil
}
