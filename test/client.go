package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client start..")
	time.Sleep(1 * time.Second)

	// 1 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 2 调用Write往连接写数据
		_, err := conn.Write([]byte("Hello Zinx V0.1.."))
		if err != nil {
			fmt.Println("write conn err ", err)
			return
		}
		buf := make([]byte, 512)
		// 调用Read从连接读数据
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ", err)
		}
		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		// CPU阻塞， 每隔1s进行连接
		time.Sleep(1 * time.Second)
	}
}
