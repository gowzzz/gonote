package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("Dial err:" + err.Error())
	}
	var reply string
	// 用client.Call时，第一个参数是用点号链接的RPC服务名字和方法名字，第二和第三个参数分别我们定义RPC方法的两个参数一个请求一个响应。
	err = client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		panic("Call err:" + err.Error())
	}
	fmt.Println("reply:", reply)
}
