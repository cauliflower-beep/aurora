package aiface

// IMessage
//  @Description: 提供消息的基本方法
//  @Description: 将请求的一个消息封装到message中，定义抽象层接口
type IMessage interface {
	GetDataLen() uint32 //获取消息数据段长度
	GetMsgID() uint32   //获取消息ID
	GetData() []byte    //获取消息内容
	GetRawData() []byte //获取原始数据

	SetMsgID(uint32)   //设计消息ID
	SetData([]byte)    //设计消息内容
	SetDataLen(uint32) //设置消息数据段长度
}