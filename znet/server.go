package znet

import (
	"fmt"
	"github.com/honey-yogurt/Zinx/utils"
	"github.com/honey-yogurt/Zinx/ziface"
	"net"
)

// Server implement interface IServer
type Server struct {
	// server name
	Name string
	// tcp4 or other
	IPVersion string
	// 服务绑定的 IP 地址
	IP string
	// 服务绑定的 端口
	Port int
	// 当前 server 的消息管理模块，用来绑定 MsgID 和对应的处理业务 API 关系
	MsgHandler ziface.IMsgHandle
	// 连接管理器
	ConnMgr ziface.IConnManager
	// 创建连接后自动调用的 Hook 函数 OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// 销毁连接前自动调用的 Hook 函数 OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

// NewServer 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

var _ ziface.IServer = &Server{}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	// 开启一个 go 去做服务端 listener 业务
	go func() {
		// 0 开启消息队列及 worker 工作池
		s.MsgHandler.StartWorkerPool()
		// 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server  ", s.Name, " success, now listening...")
		var cid uint32
		cid = 0
		// 启动 server 网络连接服务
		for {
			// 阻塞,等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			//  设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			fmt.Println(">>>>>>ConnMgr Len = ", s.ConnMgr.Len(), ", MAX = ", utils.GlobalObject.MaxConn)
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("====> Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			// TODO Server.Start() 处理该新连接请求的 业务 方法，此时应该有 handler 和 conn 是绑定的

			// 将处理新连接的业务方法和conn进行绑定，得到我们定义的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将一些服务器的资源、状态、或者 已经开辟的连接信息进行停止或回收
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// 在 start 中阻塞的话，serve 就没有用了
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	// 阻塞，否则主 go 退出，listener 的 go 将会退出
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
