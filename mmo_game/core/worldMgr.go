package core

import "sync"

/*
本模块管理当前世界所有玩家，
将AOI和玩家做了一层统一管理，协调其他模块。
拥有全部的当前在线玩家信息和当前世界的AOI划分规则，
方便玩家之间聊天、同步位置等功能
*/

type WorldMgr struct {
	AoiMgr  *AOIMgr           //当前世界地图的AOI规划管理器
	Players map[int32]*Player //当前在线的玩家集合
	pLock   sync.RWMutex      //保护Players的锁
}

// WorldMgrObj 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldMgr

// init 初始化WorldMgrObj
func init() {
	WorldMgrObj = &WorldMgr{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIMgr(AoiMinX, AoiMaxX, AoiCntX, AoiMinY, AoiMaxY, AoiCntY),
	}
}

// AddPlayer
// @Description: 提供一个添加玩家的功能，将玩家添加进玩家信息表Players
// @Description: 玩家上线的时候应该添加进来
func (wm *WorldMgr) AddPlayer(player *Player) {
	//将player添加到世界管理器中
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	//将player 添加到AOI网络规划中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

// RemovePlayer
// @Description: 提供一个删除玩家的功能，将玩家从玩家信息表Players中删除
func (wm *WorldMgr) RemovePlayer(player *Player) {
	wm.pLock.Lock()
	delete(wm.Players, player.Pid)
	wm.pLock.Unlock()
}

// GetPlayerByPid
// @Description: 根据玩家id获取对应玩家对象
func (wm *WorldMgr) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]
}

// GetAllPlayers
// @Description: 获取当前世界全部在线玩家对象列表
func (wm *WorldMgr) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0, len(wm.Players))
	for _, player := range wm.Players {
		players = append(players, player)
	}
	return players
}
