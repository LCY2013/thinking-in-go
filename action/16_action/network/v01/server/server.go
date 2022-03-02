package main

import (
	"fmt"
	"net"
)

/*
socket 监听（listen）与接收连接（accept）

socket 编程的核心在于服务端，而服务端有着自己一套相对固定的套路：Listen+Accept。在这套固定套路的基础上，我们的服务端程序通常采用一个 Goroutine 处理一个连接。

$netstat -an|grep 8888
*/

func handleConn(c net.Conn) {
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			fmt.Printf("conn close fail: [%+v]\n", err)
		}
	}(c)
	for {
		// read from the connection
		// ... ...
		// write to the connection
		//... ...
	}
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Printf("listen network fail: [%+v]\n", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("listen accept fail: [%+v]\n", err)
			break
		}

		// start a new goroutine to handle
		// the new connection.
		go handleConn(conn)
	}
}
