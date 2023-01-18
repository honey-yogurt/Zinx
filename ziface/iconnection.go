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
	// Send 发送数据给远程客户端
	Send(data []byte) error
}

// HandleFunc 定义了处理连接业务的方法
// 原生 socket 连接
// 客户端请求的数据
// 客户端请求数据的长度
type HandleFunc func(*net.TCPConn, []byte, int) error
