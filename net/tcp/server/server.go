package main

import (
	"fmt"
	"net"
)

func main() {
	// tcpAddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	// tcpListen, _ := net.ListenTCP("tcp", tcpAddr)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	for {
		// tcpListen.AcceptTCP()
		fmt.Println("accepting")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept err:", err)
			continue
		}
		go Deal(conn)
	}
}
func Deal(conn net.Conn) {
	fmt.Println(conn.RemoteAddr())
	conn.Write([]byte(conn.RemoteAddr().String()))
	err := conn.Close()
	if err != nil {
		fmt.Println("conn.Close().Error:", err)

	}
	fmt.Println("conn writer ok!")
	return
}
