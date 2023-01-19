package znet

import (
	"errors"
	"fmt"
	"github.com/honey-yogurt/Zinx/ziface"
	"io"
	"net"
)

type Connection struct {
	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 连接 ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 通知当前连接停止的 channel
	ExitChan chan bool
	// 该连接处理的方法Router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

var _ ziface.IConnection = (*Connection)(nil)

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnId = ", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()

	//TODO 启动从当前连接写数据的业务
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 关闭socket连接
	c.Conn.Close()
	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 将要发送给客户端的数据，先进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	// 将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()
	// MsgDataLen/MsgId/Data
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg is = ", msgId)
		return errors.New("pack err msg")
	}
	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write mgs id ", msgId, " error: ", err)
		return errors.New("conn write error")
	}
	return nil
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, ", Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的 Msg Head 二进制流 8 个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}
		// 拆包，得到 msgID 和 msgDataLen 放到 mgs 消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}
		// 根据 dataLen 再次读取 Data，放在 msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn 的request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中，找到注册绑定的connection对应的router调用
		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}
