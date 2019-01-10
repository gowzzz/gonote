package main

import (
	"net"
	"net/rpc"
)

type HelloService struct{}

// Go语言的RPC规则：方法只能有两个可序列化的参数，其中第二个参数是指针类型，并且返回一个error类型，同时必须是公开的方法
func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

// 为HelloService注册一个RPC服务
func RegisterRPC() {
	rpc.RegisterName("HelloService", new(HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic("Listen err:" + err.Error())
	}
	// 每接受一次就执行一次
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Accept err:" + err.Error())
		}
		rpc.ServeConn(conn)
	}

}

func main() {
	RegisterRPC()
}
