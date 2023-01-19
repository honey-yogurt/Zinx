package main

import (
	"fmt"
	"github.com/honey-yogurt/Zinx/ziface"
	"github.com/honey-yogurt/Zinx/znet"
)

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle Test
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写 ping..ping..ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1 创建一个server句柄, 使用Zinx的api
	s := znet.NewServer("[Zinx V0.4]")
	// 2 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3 启动server
	s.Serve()
}
