package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	// Start 启动连接，让当前的链接准备开始工作
	Start()
	// Stop 停止连接，结束当前连接的工作
	Stop()
	// GetTCPConnection 获取当前连接绑定的 socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接的连接 ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的 TCP 状态 IP Port
	RemoteAddr() net.Addr
	// SendMsg 发送数据给远程客户端
	SendMsg(msgId uint32, data []byte) error
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}
