package room

import (
	"sync"
	"sync/atomic"
)

var ins *roomSta
var once sync.Once

type roomSta struct {
	sum int32 // 当前最大的房间号
}

func GetRoom() *roomSta {
	once.Do(func() {
		ins = &roomSta{}
	})
	return ins
}

// 获取一个最新的房间号，规则: 当前房间号 + 1
func GetRoomNo(roomSta *roomSta) int32 {
	atomic.AddInt32(&roomSta.sum, 1)
	return roomSta.sum
}

// 获取当前房间数
func GetRoomSum(roomSta *roomSta) int32 {
	return roomSta.sum
}
