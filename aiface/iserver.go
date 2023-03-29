package aiface

import (
	"time"
)

// IServer
//  @Description: 提供Server抽象层全部接口声明，是整个aurora框架最基本的接口
type IServer interface {
	Start() //启动服务器方法
	Stop()  //停止服务器方法
	Serve() //开启业务服务方法

	AddRouter(msgID uint32, router IRouter)                   //路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	GetConnMgr() IConnManager                                 //得到链接管理
	SetOnConnStart(func(IConnection))                         //设置该Server的连接创建时Hook函数
	SetOnConnStop(func(IConnection))                          //设置该Server的连接断开时的Hook函数
	GetOnConnStart() func(IConnection)                        //得到该Server的连接创建时Hook函数
	GetOnConnStop() func(IConnection)                         //得到该Server的连接断开时的Hook函数
	GetPacket() IDataPack                                     //获取Server绑定的数据协议封包方式
	GetMsgHandler() IMsgHandle                                //获取Server绑定的消息处理模块
	SetPacket(IDataPack)                                      //设置Server绑定的数据协议封包方式
	StartHeartBeat(time.Duration)                             //启动心跳检测
	StartHeartBeatWithOption(time.Duration, *HeartBeatOption) //启动心跳检测(自定义回调)
	GetHeartBeat() IHeartbeatChecker                          //获取心跳检测器
	GetLengthField() *LengthField
	SetDecoder(IDecoder)
	AddInterceptor(Interceptor)
}
