package aiface

/*
IRequest
把客户端请求的连接信息和请求的数据，放在一个叫Request的请求类中，
好处是可以从Request里得到全部客户端的请求信息，也为之后拓展框架有一定的作用。
一旦客户端有额外含义的数据，都可以放在Request里面。
*/
type IRequest interface {
	GetConnection() IConnection //获取请求连接信息
	GetData() []byte            //获取请求消息的数据
	GetMsgId() uint32           //获取请求的消息ID
	//随着框架功能的丰富，应该继续添加新的成员进来
}
