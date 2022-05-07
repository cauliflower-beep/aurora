package anet

import "aurora/aiface"

type Request struct {
	conn aiface.IConnection //已经和客户端建立好的 链接
	msg  aiface.IMessage    //客户端请求的数据
}

/*
请求层的抽象比较简单，没必要再做一个New出来，直接新建一个Request结构体调方法就行。
当然New一个也可以
*/
//获取请求连接信息
func (r *Request) GetConnection() aiface.IConnection {
	return r.conn
}

//获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
