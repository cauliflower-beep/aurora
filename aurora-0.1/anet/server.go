package anet

import (
	"aurora-v0.1/aiface"
	"fmt"
	"net"
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

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求,如果有客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2 链接建立，做一个最基本的内容回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					//若读取成功，则回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
					}
				}
			}()
		}
	}()
}

//Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Aurora server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

//Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

//NewServer 创建一个服务器句柄
func NewServer(name string) aiface.IServer {
	printLogo()

	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}

func printLogo() {
	fmt.Println(auroraLogo)
	fmt.Println(topLine)
	fmt.Println(bottomLine)
}
