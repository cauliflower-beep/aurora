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

// SyncSurrounding
// @Description: 同步周围玩家并显示
func (p *Player) SyncSurrounding() {
	//1 根据自己的位置，获取周围九宫格内的玩家pid
	pids := WorldMgrObj.AoiMgr.GetPidByPos(p.X, p.Z)
	//2 根据pid得到所有玩家对象
	players := make([]*Player, 0, len(pids))
	//3 给这些玩家发送msgId:200 的消息，暴露自己的位置
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	//3.1 组建msgId200 proto数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //Tp-2代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//3.2 把自己的位置信息，发送给九宫格内的其他玩家
	for _, player := range players {
		player.SendMsg(200, msg)
	}

	//4 让周围九宫格内的玩家出现在自己的视野中
	//4.1构建 message SyncPlayers数据
	playersData := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersData = append(playersData, p)
	}

	//4.2 封装SyncPlayer protobuf数据
	SyncPlayersMsg := &pb.SyncPlayers{
		Ps: playersData[:], // 将playersData数据复制一份给Ps
	}

	//4.3 给当前玩家发送周围九宫格内的全部其他玩家数据
	p.SendMsg(202, SyncPlayersMsg)
}

// UpdatePos
// @Description: 广播玩家位置移动
func (p *Player) UpdatePos(x, y, z, v float32) {
	//更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//构建protobuf协议，发送位置给周围玩家
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取当前玩家周边全部玩家
	players := p.GetSurroundingPlayers()
	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

// GetSurroundingPlayers
// @Description: 获取当前玩家的AOI周边玩家信息
func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI区域的所有pid
	pids := WorldMgrObj.AoiMgr.GetPidByPos(p.X, p.Z)

	//将所有pid对应的player返回
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}
