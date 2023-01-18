package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	s := NewServer("[Zinx V0.1]")

	go ClientTest()

	s.Serve()
}

// 模拟客户端
func ClientTest() {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("hello Zinx"))
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}
		fmt.Printf(" server call back : %s, cnt = %d\n", buf[:cnt], cnt)

		time.Sleep(1 * time.Second)
	}

}
