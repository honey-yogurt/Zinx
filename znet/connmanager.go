package znet

import (
	"errors"
	"fmt"
	"github.com/honey-yogurt/Zinx/ziface"
	"sync"
)

// ConnManager 连接管理模块
type ConnManager struct {
	// 管理的连接集合
	connections map[uint32]ziface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

var _ ziface.IConnManager = (*ConnManager)(nil)

// NewConnManager 创建连接
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 将 conn 加入到 ConnManager
	c.connections[conn.GetConnID()] = conn

	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除连接信息
	delete(c.connections, conn.GetConnID())

	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源 map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	// 保护共享资源 map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除 conn 并停止 conn 的工作
	for connID, conn := range c.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(c.connections, connID)
	}

	fmt.Println("Clear All connections succ! conn num = ", c.Len())
}
