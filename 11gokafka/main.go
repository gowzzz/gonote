package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster" //support automatic consumer-group rebalancing and offset tracking
)

var (
	topics = "test0"
)

// consumer 消费者
func consumer() {
	groupID := "group-1"
	config := cluster.NewConfig()
	config.Group.Return.Notifications = true //如果启用，重新平衡通知将在通知通道上返回(默认禁用)。
	// config.Consumer.Offsets.CommitInterval = 1 * time.Second //提交更新偏移量的频率。默认为1。
	// config.Consumer.Offsets.Initial = sarama.OffsetNewest    //初始从最新的offset开始 应该是最新的或最迟的。默认为OffsetNewest。

	c, err := cluster.NewConsumer(strings.Split("localhost:9092", ","), groupID, strings.Split(topics, ","), config)
	if err != nil {
		glog.Errorf("Failed open consumer: %v", err)
		return
	}
	defer c.Close()
	// 这是必须的
	go func(c *cluster.Consumer) {
		errors := c.Errors()
		// 通知返回在用户重新平衡期间发生的通知通道。只有在配置的Group.Return中，通知才会通过该通道发出。通知设置为true。
		noti := c.Notifications()
		for {
			select {
			case err := <-errors:
				glog.Errorln(err)
			case <-noti:
			}
		}
	}(c)

	for msg := range c.Messages() {
		fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		// MarkOffset将提供的消息标记为已处理的消息
		c.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
	}
}

// syncProducer 同步生产者
// 并发量小时，可以用这种方式
func syncProducer() {
	config := sarama.NewConfig()
	//  config.Producer.RequiredAcks = sarama.WaitForAll
	//  config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(strings.Split("localhost:9092", ","), config)
	defer p.Close()
	if err != nil {
		glog.Errorln(err)
		return
	}

	v := "sync: " + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
	fmt.Fprintln(os.Stdout, v)
	msg := &sarama.ProducerMessage{
		Topic: topics,
		Value: sarama.ByteEncoder(v),
	}
	if _, _, err := p.SendMessage(msg); err != nil {
		glog.Errorln(err)
		return
	}
}

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式
func asyncProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewAsyncProducer(strings.Split("localhost:9092", ","), config)
	defer p.Close()
	if err != nil {
		return
	}

	//必须有这个匿名函数内容
	//循环判断哪个通道发送过来数据.
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					glog.Errorln(err)
				}
			case <-success:
			}
		}
	}(p)
	for {
		v := "async: " + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
		fmt.Fprintln(os.Stdout, v)
		msg := &sarama.ProducerMessage{
			Topic: topics,
			Value: sarama.ByteEncoder(v), //sarama.StringEncoder("test")
		}
		p.Input() <- msg
		time.Sleep(time.Second * 1)
	}
}

func main() {
	// go asyncProducer()
	go consumer()
	time.Sleep(time.Second * 10000)
}
