package anet

import (
	"aurora-v0.2/aiface"
	"errors"
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

//Server 接口实现，定义一个Server服务类
type Server struct {
	Name      string //服务器名称
	IPVersion string //tcp4 or other
	IP        string //服务绑定的IP地址
	Port      int    //服务绑定的端口
}

//============== 定义当前客户端链接的handle api ===========

// CallBackToClient
//  @Description: 客户端消息回显业务
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err|", err)
		return errors.New("callBackToClient err")
	}
	return nil
}

//============== 实现 aiface.IServer 里的全部接口方法 ========

//Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
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

			// 处理该新连接请求的业务方法，此时应该把 handler 和 conn 绑定
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++
			go dealConn.Start()
		}
	}()
}

//Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Aurora server , name ", s.Name)

	// todo 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

//Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

//NewServer 创建一个服务器句柄
func NewServer(name string) aiface.IServer {
	printLogo()

	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}

func printLogo() {
	fmt.Println(auroraLogo)
	fmt.Println(topLine)
	fmt.Println(bottomLine)
}