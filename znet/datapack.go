package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/honey-yogurt/Zinx/utils"
	"github.com/honey-yogurt/Zinx/ziface"
)

// DataPack 封包，拆包的具体模块
type DataPack struct{}

var _ ziface.IDataPack = (*DataPack)(nil)

// NewDataPack 拆包、封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32（4字节）+ ID uint32（4字节）= 8字节
	return 8
}

// Pack 封包方法
//
//	/dataLen/msgId/data/
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 将dataLen写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法，将包的Head信息读出来，再根据Head信息里的data长度，再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息，得到dataLen和msgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil
}
