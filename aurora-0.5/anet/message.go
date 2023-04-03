package anet

// Message 消息体 包含 data/id/dataLen 3个基本属性
type Message struct {
	Id      uint32 //消息的ID
	DataLen uint32 //消息的长度
	Data    []byte //消息的内容
}

// NewMsgPackage 创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetDataId 获取消息Id
func (msg *Message) GetDataId() uint32 {
	return msg.Id
}

// GetDataLen 获取消息数据段的长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

// GetData 获取消息数据
func (msg *Message) GetData() []byte {
	return msg.Data
}

// SetDataLen 设置消息数据长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

// SetDataId 设置消息Id
func (msg *Message) SetDataId(msgId uint32) {
	msg.Id = msgId
}

// SetData 设置消息数据段
func (msg *Message) SetData(msgData []byte) {
	msg.Data = msgData
}
