package core

import "fmt"

// AOIMgr
// @Description: AOI管理模块
type AOIMgr struct {
	MinX  int           //区域左边界坐标
	MaxX  int           //区域右边界坐标
	CntX  int           //x方向格子的数量
	MinY  int           //区域上边界坐标
	MaxY  int           //区域下边界坐标
	CntY  int           //y方向格子的数量
	grids map[int]*Grid //当前区域中都有哪些格子 key->格子ID val->格子对象
}

// NewAOIMgr
// @Description: 初始化一个aoi区域
func NewAOIMgr(minX, maxX, cntX, minY, maxY, cntY int) *AOIMgr {
	aoiMgr := &AOIMgr{
		MinX:  minX,
		MaxX:  maxX,
		CntX:  cntX,
		MinY:  minY,
		MaxY:  maxY,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}

	//给AOI初始化区域中所有的格子
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			//计算格子ID
			//格子编号：id = idy * nx +idx 利用格子坐标计算格子编号
			gid := y*cntY + x

			//初始化一个格子放在AOI中的map里，key是当前格子的ID
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength(),
			)
		}
	}

	return aoiMgr
}

// gridWidth
// @Description:获取格子宽度
func (am *AOIMgr) gridWidth() int {
	return (am.MaxX - am.MinX) / am.CntX
}

// gridLength
// @Description: 获取格子长度
func (am *AOIMgr) gridLength() int {
	return (am.MaxY - am.MinY) / am.CntY
}

// Info
// @Description: aoi管理模块信息
func (am *AOIMgr) Info() string {
	s := fmt.Sprintf("AOIMgr:\nminX:%d, maxX:%d, cntX:%d, minY:%d, maxY:%d, cntY:%d\n Grids in AOIMgr:\n",
		am.MinX, am.MaxX, am.CntX, am.MinY, am.MaxY, am.CntY)
	for _, grid := range am.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// GetSurroundGridsByGid
// @Description: 求出某个格子周围都有哪些九宫格
// @Description: 基本思路是先求出这个格子左右是否有格子，再求出这一行上下是否有格子
func (am *AOIMgr) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	//判断gid是否存在
	if _, ok := am.grids[gid]; !ok {
		return
	}

	//将当前gid添加到九宫格中
	grids = append(grids, am.grids[gid])

	//根据gid得到当前格子所在的x轴编号
	idx := gid % am.CntX

	//判断当前idx左边是否还有格子
	if idx > 0 {
		grids = append(grids, am.grids[gid-1])
	}

	//判断当前idx右边是否还有格子
	if idx < am.CntX-1 {
		grids = append(grids, am.grids[gid+1])
	}

	//将x轴当前的格子都取出，进行遍历，分别判断每个格子的上下是否有格子

	//得到当前x轴的格子id集合
	gridsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridsX = append(gridsX, v.GID)
	}

	//遍历x轴格子
	for _, v := range gridsX {
		//计算该格子处于第几列
		idy := v / am.CntX

		//判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, am.grids[v-am.CntX])
		}
		//判断当前的idy下边是否还有格子
		if idy < am.CntY-1 {
			grids = append(grids, am.grids[v+am.CntX])
		}
	}
	return
}

// GetGidByPos
// @Description: 通过横纵坐标获取对应的格子Id
func (am *AOIMgr) GetGidByPos(x, y float32) int {
	gx := (int(x) - am.MinX) / am.gridWidth()
	gy := (int(y) - am.MinY) / am.gridLength()
	return gy*am.CntX + gx
}

// GetPidByPos
// @Description: 通过横纵坐标得到周边九宫格内的全部playerIds
func (am *AOIMgr) GetPidByPos(x, y float32) (playerIds []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gid := am.GetGidByPos(x, y)

	//根据格子ID得到周边九宫格的信息
	grids := am.GetSurroundGridsByGid(gid)
	for _, v := range grids {
		playerIds = append(playerIds, v.GetPlayerIDs()...)
		fmt.Printf("===> grid ID : %d, pids : %v  ====", v.GID, v.GetPlayerIDs())
	}

	return
}

// GetPidsByGid
// @Description: 通过gid获取当前格子的全部playerId
func (am *AOIMgr) GetPidsByGid(gid int) (playerIDs []int) {
	playerIDs = am.grids[gid].GetPlayerIDs()
	return
}

// RemovePidFromGrid
// @Description: 移除一个格子中的playerId
func (am *AOIMgr) RemovePidFromGrid(pid, gid int) {
	am.grids[gid].Remove(pid)
}

// AddPidToGrid
// @Description: 添加一个playerId到一个格子中
func (am *AOIMgr) AddPidToGrid(pid, gid int) {
	am.grids[gid].AddPlayer(pid)
}

// AddToGridByPos
// @Description: 通过横纵坐标添加一个player到一个格子中
func (am *AOIMgr) AddToGridByPos(pid int, x, y float32) {
	gid := am.GetGidByPos(x, y)
	am.AddPidToGrid(pid, gid)
}

// RemoveFromGridByPos
// @Description: 通过横纵坐标把一个player从对应格子中删除
func (am *AOIMgr) RemoveFromGridByPos(pid int, x, y float32) {
	gid := am.GetGidByPos(x, y)
	am.RemovePidFromGrid(pid, gid)
}
