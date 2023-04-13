package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int          //格子ID
	MinX      int          //格子左边界坐标
	MaxX      int          //格子右边界坐标
	MinY      int          //格子上边界坐标
	MaxY      int          //格子下边界坐标
	playerIDs map[int]bool //当前格子内的玩家或者物体成员ID
	pIdLock   sync.RWMutex //playerIDs的map保护锁
}

// NewGrid
// @Description: 初始化一个格子
func NewGrid(gId, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gId,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// AddPlayer
// @Description: 向当前格子内添加一个玩家
func (g *Grid) AddPlayer(playerID int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	g.playerIDs[playerID] = true
}

// Remove
// @Description: 从当前格子中移除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlayerIDs
// @Description: 得到当前格子中所有的玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIdLock.RLock()
	g.pIdLock.RUnlock()

	for playerId := range g.playerIDs {
		playerIDs = append(playerIDs, playerId)
	}
	return
}

func (g *Grid) Info() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
