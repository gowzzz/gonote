package main

import (
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

func RunProducer(addr, topic, message string) {
	// strIP1 := "127.0.0.1:4150"
	producer1, err := initProducer(addr)
	if err != nil {
		log.Fatal("init producer1 error:", err)
	}
	defer producer1.Stop()

	//读取控制台输入
	count := 0
	for {
		err := producer1.public("test1", message)
		if err != nil {
			log.Fatal("producer1 public error:", err)
		}
		count++
	}
}

type nsqProducer struct {
	*nsq.Producer
}

//初始化生产者
func initProducer(addr string) (*nsqProducer, error) {
	fmt.Println("init producer address:", addr)
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return &nsqProducer{producer}, nil
}

//发布消息
func (np *nsqProducer) public(topic, message string) error {
	err := np.Publish(topic, []byte(message))
	if err != nil {
		log.Println("nsq public error:", err)
		return err
	}
	return nil
}
