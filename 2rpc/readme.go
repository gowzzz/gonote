package rpc

import "net/rpc"

// 规则制定者

// 服务要实现的方法
type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}

//服务的名字
const HelloServiceName = "path/to/pkg.HelloService"

// 注册该服务的函数
func RegisterHelloService(srv HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, srv)
}

// server开发者

// client开发者

// --------------------------------
type HelloServiceClient struct {
	*rpc.Client
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}
func (p *HelloServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(HelloServiceName+".Hello", request, reply)
}
