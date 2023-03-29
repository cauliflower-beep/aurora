package anet

import (
	"aurora/aconf"
	"aurora/alog"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"aurora/aiface"
	"aurora/apack"
)

var zinxLogo = `                                        
   █████  ██    ██ ██████   ██████  ██████   █████ 
  ██   ██ ██    ██ ██   ██ ██    ██ ██   ██ ██   ██
  ███████ ██    ██ ██████  ██    ██ ██████  ███████
  ██   ██ ██    ██ ██   ██ ██    ██ ██   ██ ██   ██
  ██   ██  ██████  ██   ██  ██████  ██   ██ ██   ██
                                        `
var topLine = `┌──────────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└──────────────────────────────────────────────────────┘`

// Server
// @Description: 接口实现，定义一个Server服务类
type Server struct {
	Name      string //服务器的名称
	IPVersion string //tcp4 or other
	IP        string //服务绑定的IP地址
	Port      int    //服务绑定的端口

	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler aiface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr aiface.IConnManager
	//该Server的连接创建时Hook函数
	onConnStart func(conn aiface.IConnection)
	//该Server的连接断开时的Hook函数
	onConnStop func(conn aiface.IConnection)
	//数据报文封包方式
	packet aiface.IDataPack
	//异步捕获链接关闭状态
	exitChan chan struct{}
	//断粘包解码器
	decoder aiface.IDecoder
	//心跳检测器
	hc aiface.IHeartbeatChecker
}

// NewServer 创建一个服务器句柄
func NewServer(opts ...Option) aiface.IServer {
	printLogo()

	s := &Server{
		Name:       aconf.GlobalObject.Name,
		IPVersion:  "tcp",
		IP:         aconf.GlobalObject.Host,
		Port:       aconf.GlobalObject.TCPPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		exitChan:   nil,
		//默认使用zinx的TLV封包方式
		packet:  apack.Factory().NewPack(aiface.ZinxDataPack),
		decoder: apack.NewTLVDecoder(), //默认使用TLV的解码方式
	}

	for _, opt := range opts {
		opt(s)
	}

	//提示当前配置信息
	aconf.GlobalObject.Show()

	return s
}

func NewUserConfServer(config *aconf.Config, opts ...Option) aiface.IServer {
	//打印logo
	printLogo()

	s := &Server{
		Name:       config.Name,
		IPVersion:  "tcp4",
		IP:         config.Host,
		Port:       config.TCPPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		exitChan:   nil,
		packet:     apack.Factory().NewPack(aiface.ZinxDataPack),
		decoder:    apack.NewTLVDecoder(), //默认使用TLV的解码方式
	}
	//更替打包方式
	for _, opt := range opts {
		opt(s)
	}
	//刷新用户配置到全局配置变量
	aconf.UserConfToGlobal(config)

	//提示当前配置信息
	aconf.GlobalObject.Show()

	return s
}

//============== 实现 aiface.IServer 里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start() {
	alog.Ins().InfoF("[START] Server name: %s,listener at IP: %s, Port %d is starting", s.Name, s.IP, s.Port)
	s.exitChan = make(chan struct{})

	// 将解码器添加到拦截器
	if s.decoder != nil {
		s.msgHandler.AddInterceptor(s.decoder)
	}

	//开启一个go去做服务端Listener业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			alog.Ins().ErrorF("[START] resolve tcp addr err: %v\n", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		alog.Ins().InfoF("[START] start Zinx server  %s succ, now listening...", s.Name)

		var cID uint64

		go func() {
			//3 启动server网络连接业务
			for {
				//3.1 设置服务器最大连接控制,如果超过最大连接，则等待
				if s.ConnMgr.Len() >= aconf.GlobalObject.MaxConn {
					alog.Ins().InfoF("Exceeded the maxConnNum:%d, Wait:%d", aconf.GlobalObject.MaxConn, AcceptDelay.duration)
					AcceptDelay.Delay()
					continue
				}

				//3.2 阻塞等待客户端建立连接请求
				conn, err := listener.Accept()
				if err != nil {
					//Go 1.16+
					if errors.Is(err, net.ErrClosed) {
						alog.Ins().ErrorF("Listener closed")
						return
					}
					alog.Ins().ErrorF("Accept err: %v", err)
					AcceptDelay.Delay()
					continue
				}

				AcceptDelay.Reset()

				//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
				dealConn := newServerConn(s, conn, cID)
				cID++

				//HeartBeat 心跳检测
				if s.hc != nil {
					//从Server端克隆一个心跳检测器
					heartBeatChecker := s.hc.Clone()
					//绑定当前链接
					heartBeatChecker.BindConn(dealConn)
				}

				//3.4 启动当前链接的处理业务
				go dealConn.Start()
			}
		}()

		select {
		case <-s.exitChan:
			err := listener.Close()
			if err != nil {
				alog.Ins().ErrorF("listener close err: %v", err)
			}
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	alog.Ins().InfoF("[STOP] Zinx server , name %s", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
	s.exitChan <- struct{}{}
	close(s.exitChan)
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()
	//阻塞,否则主Go退出， listenner的go将会退出
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill信号
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	alog.Ins().InfoF("[SERVE] Zinx server , name %s, Serve Interrupt, signal = %v", s.Name, sig)
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router aiface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() aiface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(aiface.IConnection)) {
	s.onConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(aiface.IConnection)) {
	s.onConnStop = hookFunc
}

// GetOnConnStart 得到该Server的连接创建时Hook函数
func (s *Server) GetOnConnStart() func(aiface.IConnection) {
	return s.onConnStart
}

// 得到该Server的连接断开时的Hook函数
func (s *Server) GetOnConnStop() func(aiface.IConnection) {
	return s.onConnStop
}

func (s *Server) GetPacket() aiface.IDataPack {
	return s.packet
}

func (s *Server) SetPacket(packet aiface.IDataPack) {
	s.packet = packet
}

func (s *Server) GetMsgHandler() aiface.IMsgHandle {
	return s.msgHandler
}

// StartHeartBeat 启动心跳检测
// interval 每次发送心跳的时间间隔
func (s *Server) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatChecker(interval)

	//添加心跳检测的路由
	s.AddRouter(checker.MsgID(), checker.Router())

	//server绑定心跳检测器
	s.hc = checker
}

// StartHeartBeatWithFunc 启动心跳检测
// option 心跳检测的配置
func (s *Server) StartHeartBeatWithOption(interval time.Duration, option *aiface.HeartBeatOption) {
	checker := NewHeartbeatChecker(interval)

	if option != nil {
		checker.SetHeartbeatMsgFunc(option.MakeMsg)
		checker.SetOnRemoteNotAlive(option.OnRemoteNotAlive)
		checker.BindRouter(option.HeadBeatMsgID, option.Router)
	}

	//添加心跳检测的路由
	s.AddRouter(checker.MsgID(), checker.Router())

	//server绑定心跳检测器
	s.hc = checker
}

func (s *Server) GetHeartBeat() aiface.IHeartbeatChecker {
	return s.hc
}

func (s *Server) SetDecoder(decoder aiface.IDecoder) {
	s.decoder = decoder
}

func (s *Server) GetLengthField() *aiface.LengthField {
	if s.decoder != nil {
		return s.decoder.GetLengthField()
	}
	return nil
}

func (s *Server) AddInterceptor(interceptor aiface.Interceptor) {
	s.msgHandler.AddInterceptor(interceptor)
}

func printLogo() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/aceld                    %s", borderLine, borderLine))
	fmt.Println(fmt.Sprintf("%s [tutorial] https://www.yuque.com/aceld/npyr8s/bgftov %s", borderLine, borderLine))
	fmt.Println(fmt.Sprintf("%s [document] https://www.yuque.com/aceld/tsgooa        %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		aconf.GlobalObject.Version,
		aconf.GlobalObject.MaxConn,
		aconf.GlobalObject.MaxPacketSize)
}

func init() {}
