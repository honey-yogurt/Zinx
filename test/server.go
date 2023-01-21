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
	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的数据，再回写 ping..ping..ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

// HelloZinxRouter hello Zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Handle Test
func (h *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	// 先读取客户端的数据，再回写 ping..ping..ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(1, []byte("Hello! Welcome to Zinx!")); err != nil {
		fmt.Println(err)
	}
}

// DoConnectionBegin 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN!")); err != nil {
		fmt.Println(err)
	}
	// 给当前的一些连接设置属性
	conn.SetProperty("name", "honey-yogurt")
}

// DoConnectionLost 连接断开前执行的钩子函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionLost is Called...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost...")

	// 获取连接属性
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("name", name)
	}

}

func main() {
	// 1 创建一个server句柄, 使用Zinx的api
	s := znet.NewServer("[Zinx V1.0]")
	// 2 注册连接 Hook 钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	// 3 给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 4 启动server
	s.Serve()
}
