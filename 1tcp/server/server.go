package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	var buf = "wzzzz\n"
	n, err := conn.Write([]byte(buf))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
	fmt.Println(buf)
	conn.Close()
}
