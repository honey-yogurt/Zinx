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
}

// NewServer 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
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
			// TODO Server.Start() 设置服务器最大连接控制，如果超过最大连接，那么则关闭新的连接

			// TODO Server.Start() 处理该新连接请求的 业务 方法，此时应该有 handler 和 conn 是绑定的

			// 将处理新连接的业务方法和conn进行绑定，得到我们定义的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
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
