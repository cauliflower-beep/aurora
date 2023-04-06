package anet

import (
	"aurora-v0.6/aiface"
	"aurora-v0.6/utils"
	"fmt"
	"net"
	"time"
)

var auroraLogo = `                                        
 _____  _   _   ____  ___    ____  _____ 
(____ || | | | / ___)/ _ \  / ___)(____ |
/ ___ || |_| || |   | |_| || |    / ___ |
\_____||____/ |_|    \___/ |_|    \_____|
                                        `
var topLine = `┌───────────────────────────────────────────────────┐`
var bottomLine = `└───────────────────────────────────────────────────┘`

// Server 接口实现，定义一个Server服务类
type Server struct {
	Name      string //服务器名称
	IPVersion string //tcp4 or other
	IP        string //服务绑定的IP地址
	Port      int    //服务绑定的端口

	MsgMgr aiface.IMsgMgr //当前Server由服务器开发绑定的回调router，也就是Server注册的链接对应的处理业务
}

// NewServer 创建一个服务器句柄
func NewServer() aiface.IServer {
	printLogo()

	return &Server{
		Name:      utils.GlobalObject.Name, //从全局参数获取
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,    //从全局参数获取
		Port:      utils.GlobalObject.TcpPort, //从全局参数获取
		MsgMgr:    NewMsgMgr(),                //MsgMgr 初始化
	}
}

//============== 实现 aiface.IServer 里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	//追加日志查看配置是否生效
	fmt.Printf("[Aurora] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	//开启一个goroutine去做服务端Listener业务
	go func() {
		//1 获取一个TCP的Addr解析对象
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		fmt.Println("start Aurora server  ", s.Name, " succ, now listenning...")

		// todo server.go 应该有一个自动生成ID的方法
		var cid uint32 = 0

		//3 启动server网络连接业务
		for {
			var conn *net.TCPConn
			conn, err = listener.AcceptTCP() //阻塞等待客户端建立连接请求
			if err != nil {
				fmt.Println("accept err|", err)
				continue
			}

			// todo server.start()设置服务器最大连接控制，如果超过最大连接，则关闭此新的连接

			// 处理该新连接请求的业务方法，此时应该把 Router 和 conn 绑定
			dealConn := NewConnection(conn, cid, s.MsgMgr)
			cid++
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Aurora server , name ", s.Name)

	// todo 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// AddRouter 路由功能
func (s *Server) AddRouter(msgId uint32, router aiface.IRouter) {
	s.MsgMgr.AddRouter(msgId, router)
	fmt.Println("Add Server Router succ!")
}

func printLogo() {
	fmt.Println(auroraLogo)
	fmt.Println(topLine)
	fmt.Println(bottomLine)
}
