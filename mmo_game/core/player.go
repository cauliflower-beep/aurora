package core

import (
	"aurora-v1.0/aiface"
	"auroraTags/mmo_game/pb/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

type Player struct {
	Pid  int32              //玩家Id
	Conn aiface.IConnection //当前玩家的链接
	X    float32            //平面x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标 (注意不是Y)
	V    float32            //旋转0-360度
}

/*
	PlayerId生成器
*/

var PidGen int32 = 1  //用来生成玩家的计数器
var IdLock sync.Mutex //保护PidGen的互斥机制

// NewPlayer
// @Description: 创建一个玩家
func NewPlayer(conn aiface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	defer IdLock.Unlock()
	player := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴偏移若干坐标
		Y:    0,                            //高度为0
		Z:    float32(134 + rand.Intn(17)), //随机在134坐标点 基于Y轴偏移若干坐标
		V:    0,                            //角度为0，尚未实现
	}
	return player
}

// SendMsg
// @Description: 将pb的protobuf数据序列化之后发送给客户端
// @Description: 因为player经常与客户端通信，所以提供了这个方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	fmt.Printf("before Marshal data = %+v\n", data)

	//将proto message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err|", err)
		return
	}
	fmt.Printf("after marshal data = %+v\n", msg)

	if p.Conn == nil {
		fmt.Println("conn in player is nil")
		return
	}

	//调用aurora框架的SendMsg发包
	if err = p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("player sendmsg error!")
		return
	}
	return
}

// SyncPid
// @Description: 同步当前的playerID给客户端，走MsgID-1 消息
func (p *Player) SyncPid() {
	//组建msgid-1 proto数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, data)
}

// BroadCastStartPosition
// @Description: 同步当前玩家的初始化坐标信息给客户端，走MsgID-200 消息
func (p *Player) BroadCastStartPosition() {
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, msg)
}

// Talk
// @Description: 玩家发起世界聊天
func (p *Player) Talk(content string) {
	//1.组建msgId200 proto数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //TP-1代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//2.得到当前世界所有在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//3.向所有玩家发送MsgId:200消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}
