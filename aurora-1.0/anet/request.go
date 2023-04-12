package anet

import "aurora-v1.0/aiface"

type Request struct {
	iConn aiface.IConnection //已经和客户端建立好的连接
	msg   aiface.IMessage    //客户端请求的数据
}

//func NewRequest() *Request {
//	return &Request{}
//}

// GetConnection 获取请求的连接信息
func (r *Request) GetConnection() aiface.IConnection {
	return r.iConn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId 获取请求的消息Id
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
