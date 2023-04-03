package utils

import (
	"aurora-v0.5/aiface"
	"encoding/json"
	"os"
)

/*
存储一切有关aurora框架的全局参数，供其他模块使用
一些参数可以通过 用户根据 aurora.json来配置
*/

type GlobalObj struct {
	TcpServer aiface.IServer //当前aurora的全局server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	Version   string         //当前aurora版本号

	MaxPacketSize uint32 //都需数据包的最大值
	MaxConn       int    //当前服务器主机允许的最大链接个数
}

// GlobalObject 定义一个全局指针对象，目的是让其他模块都可以访问/修改到里面的参数
var GlobalObject *GlobalObj

func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:          "AuroraServerApp",
		Version:       "v0.4",
		TcpPort:       8999,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	//从配置文件中加载配置参数
	GlobalObject.ReloadConf()
}

/*
ReloadConf
读取用户配置文件
原文 ioutil.ReadFile 函数，在Go1.16版本中被标记为过时，在Go1.17版本中已经被删除
*/
func (g *GlobalObj) ReloadConf() {
	data, err := os.ReadFile("conf/aurora.json") //os.ReadFile可以自动打开并关闭的文件，不需要经过open处理
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	//fmt.Printf("json :%s\n",data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
