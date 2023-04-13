package core

import (
	"fmt"
	"testing"
)

// TestNewAOIManager
// @Description: go test -run TestNewAOIManager -v
func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIMgr(100, 300, 4, 200, 450, 5)
	fmt.Println(aoiMgr)
}

// TestAOIManagerSuroundGridsByGid
// @Description: go test -run TestAOIManagerSuroundGridsByGid -v
func TestAOIManagerSuroundGridsByGid(t *testing.T) {
	aoiMgr := NewAOIMgr(0, 250, 5, 0, 250, 5)

	for k, _ := range aoiMgr.grids {
		//得到当前格子周边的九宫格
		grids := aoiMgr.GetSurroundGridsByGid(k)
		//得到九宫格所有的IDs
		fmt.Println("gid : ", k, " grids len = ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Printf("grid ID: %d, surrounding grid IDs are %v\n", k, gIDs)
	}
}
