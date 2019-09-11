package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("time.Duration(DEADLINE):", time.Duration(1)*time.Second)
	test()
}

type WebsocketInfo struct {
	Ws  int
	Err error
}

var Ch = make(chan WebsocketInfo)

func res() {
	res := WebsocketInfo{Ws: 1, Err: nil}
	Ch <- res
}
func test() {
	go res()
	select {
	case res := <-Ch: //如果有数据，下面打印。但是有可能ch一直没数据
		fmt.Println("res:", res)
	case <-time.After(2 * time.Second): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后3秒超时
		fmt.Println("超时")
	}
}
