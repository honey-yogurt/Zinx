package znet

import (
	"fmt"
	"github.com/honey-yogurt/Zinx/ziface"
	"strconv"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
}

var _ ziface.IMsgHandle = &MsgHandle{}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
	}
	// 2 根据MsgID调度router对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		// id 已经注册
		panic("repeat api, msg ID = " + strconv.Itoa(int(msgID)))
	}
	// 2 添加msg和API的绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}
