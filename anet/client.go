package anet

import (
	"aurora/aconf"
	"aurora/aiface"
	"aurora/alog"
	"aurora/apack"
	"net"
	"time"
)

type Client struct {
	//目标链接服务器的IP
	Ip string
	//目标链接服务器的端口
	Port int
	//客户端链接
	conn aiface.IConnection
	//该client的连接创建时Hook函数
	onConnStart func(conn aiface.IConnection)
	//该client的连接断开时的Hook函数
	onConnStop func(conn aiface.IConnection)
	//数据报文封包方式
	packet aiface.IDataPack
	//异步捕获链接关闭状态
	exitChan chan struct{}
	//消息管理模块
	msgHandler aiface.IMsgHandle
	//断粘包解码器
	decoder aiface.IDecoder
	//心跳检测器
	hc aiface.IHeartbeatChecker
}

func NewClient(ip string, port int, opts ...ClientOption) aiface.IClient {

	c := &Client{
		Ip:         ip,
		Port:       port,
		msgHandler: NewMsgHandle(),
		packet:     apack.Factory().NewPack(aiface.ZinxDataPack), //默认使用zinx的TLV封包方式
		decoder:    apack.NewTLVDecoder(),                        //默认使用zinx的TLV解码器
	}

	//应用Option设置
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// 启动客户端，发送请求且建立链接
func (c *Client) Start() {

	c.exitChan = make(chan struct{})

	// 将解码器添加到拦截器
	if c.decoder != nil {
		c.msgHandler.AddInterceptor(c.decoder)
	}

	//客户端将协程池关闭
	aconf.GlobalObject.WorkerPoolSize = 0

	go func() {
		addr := &net.TCPAddr{
			IP:   net.ParseIP(c.Ip),
			Port: c.Port,
			Zone: "", //for ipv6, ignore
		}

		//创建原始Socket，得到net.Conn
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			//创建链接失败
			alog.Ins().ErrorF("client connect to server failed, err:%v", err)
			panic(err)
		}

		//创建Connection对象
		c.conn = newClientConn(c, conn)
		alog.Ins().InfoF("[START] Zinx Client LocalAddr: %s, RemoteAddr: %s\n", conn.LocalAddr(), conn.RemoteAddr())

		//HeartBeat心跳检测
		if c.hc != nil {
			//创建链接成功，绑定链接与心跳检测器
			c.hc.BindConn(c.conn)
		}

		//启动链接
		go c.conn.Start()

		select {
		case <-c.exitChan:
			alog.Ins().InfoF("client exit.")
		}
	}()
}

// StartHeartBeat 启动心跳检测
// interval 每次发送心跳的时间间隔
func (c *Client) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatChecker(interval)

	//添加心跳检测的路由
	c.AddRouter(checker.MsgID(), checker.Router())

	//client绑定心跳检测器
	c.hc = checker
}

// 启动心跳检测(自定义回调)
func (c *Client) StartHeartBeatWithOption(interval time.Duration, option *aiface.HeartBeatOption) {
	checker := NewHeartbeatChecker(interval)

	if option != nil {
		checker.SetHeartbeatMsgFunc(option.MakeMsg)
		checker.SetOnRemoteNotAlive(option.OnRemoteNotAlive)
		checker.BindRouter(option.HeadBeatMsgID, option.Router)
	}

	//添加心跳检测的路由
	c.AddRouter(checker.MsgID(), checker.Router())

	//client绑定心跳检测器
	c.hc = checker
}

func (c *Client) Stop() {
	alog.Ins().InfoF("[STOP] Zinx Client LocalAddr: %s, RemoteAddr: %s\n", c.conn.LocalAddr(), c.conn.RemoteAddr())
	c.conn.Stop()
	c.exitChan <- struct{}{}
	close(c.exitChan)
}

func (c *Client) AddRouter(msgID uint32, router aiface.IRouter) {
	c.msgHandler.AddRouter(msgID, router)
}

func (c *Client) Conn() aiface.IConnection {
	return c.conn
}

// 设置该Client的连接创建时Hook函数
func (c *Client) SetOnConnStart(hookFunc func(aiface.IConnection)) {
	c.onConnStart = hookFunc
}

// 设置该Client的连接断开时的Hook函数
func (c *Client) SetOnConnStop(hookFunc func(aiface.IConnection)) {
	c.onConnStop = hookFunc
}

// GetOnConnStart 得到该Server的连接创建时Hook函数
func (c *Client) GetOnConnStart() func(aiface.IConnection) {
	return c.onConnStart
}

// 得到该Server的连接断开时的Hook函数
func (c *Client) GetOnConnStop() func(aiface.IConnection) {
	return c.onConnStop
}

// 获取Client绑定的数据协议封包方式
func (c *Client) GetPacket() aiface.IDataPack {
	return c.packet
}

// 设置Client绑定的数据协议封包方式
func (c *Client) SetPacket(packet aiface.IDataPack) {
	c.packet = packet
}

func (c *Client) GetMsgHandler() aiface.IMsgHandle {
	return c.msgHandler
}

func (c *Client) AddInterceptor(interceptor aiface.Interceptor) {
	c.msgHandler.AddInterceptor(interceptor)
}

func (c *Client) SetDecoder(decoder aiface.IDecoder) {
	c.decoder = decoder
}
func (c *Client) GetLengthField() *aiface.LengthField {
	if c.decoder != nil {
		return c.decoder.GetLengthField()
	}
	return nil
}
