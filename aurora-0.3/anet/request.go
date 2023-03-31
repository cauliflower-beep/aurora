package anet

import "aurora-v0.3/aiface"

type Request struct {
	iConn aiface.IConnection //已经和客户端建立好的连接
	data  []byte             //客户端请求的数据
}

//func NewRequest() *Request {
//	return &Request{}
//}

// GetConnection
//  @Description: 获取请求的连接信息
func (r *Request) GetConnection() aiface.IConnection {
	return r.iConn
}

// GetData
//  @Description: 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.data
}
