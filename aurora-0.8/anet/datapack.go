package anet

import (
	"aurora-v0.8/aiface"
	"aurora-v0.8/utils"
	"bytes"
	"encoding/binary"
	"errors"
)

// DataPack 封包拆包类实例，暂时不需要成员
type DataPack struct{}

// NewDataPack 实例化一个封包拆包类
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// 此处暂时固定为8 Id uint32(4字节) + DadaLen uint32(4字节)
	return 8
}

// Pack 封包方法(压缩数据 将消息打包成二进制数据以便于传输)
func (dp *DataPack) Pack(msg aiface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲区
	/*
		缓冲区是一个用于存储数据的临时存储区域。
		这里因为消息的数据长度、ID和数据长度是不确定的，所以需要使用缓冲区，动态的存储这些数据
	*/
	dataBuff := bytes.NewBuffer([]byte{})

	/*
		LittleEndian是一种字节序，它指定了在多字节数据类型(如int、float等)的存储中，最低有效字节(即最右边的字节)先存储在内存中。
		这与BigEndian相反，后者将最高有效字节(即最左边的字节)存储在内存中。
	*/
	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// Unpack 拆包(解压数据)
func (dp *DataPack) Unpack(binaryData []byte) (aiface.IMessage, error) {
	//创建一个从输入获取二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	/*
		拆包的时候是分两次过程的，第二次依赖第一次的dataLen结果。
		所以Unpack只能解压出包头head的内容，得到msgId和dataLen，
		之后调用者再根据dataLen继续从io流中读取body中的数据
	*/

	//只解压head的信息，得到dataLen/msgId
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data recieved")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
