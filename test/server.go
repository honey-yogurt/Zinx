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

// PreHandle Test
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle Test
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

// PostHandle Test PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("call back after ping error")
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
