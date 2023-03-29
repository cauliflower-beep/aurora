package anet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test Start")
	time.Sleep(3 * time.Second) //3秒之后发送测试请求，给服务端开启服务的机会

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("conn start err|", err, ":exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("hello, aurora v0.1! "))
		if err != nil {
			fmt.Println("data send err|", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buff err|", err)
			return
		}
		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)
		time.Sleep(time.Second)
	}
}

// TestServer
//  @Description: go test -run TestServer -v
func TestServer(t *testing.T) {
	s := NewServer("[aurora v0.1]")

	go ClientTest()

	s.Serve()
}
