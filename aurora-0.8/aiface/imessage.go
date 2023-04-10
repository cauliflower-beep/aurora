package aiface

// IMessage 抽象层接口
type IMessage interface {
	GetDataLen() uint32 //获取消息数据段长度
	GetMsgId() uint32   //获取消息ID
	GetData() []byte    //获取消息内容

	SetDataLen(uint322 uint32) //设计消息ID
	SetMsgId(uint322 uint32)   //设计消息内容
	SetData([]byte)            //设置消息数据段长度
}
