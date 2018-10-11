package cli

import (
	"sync"
)

type Cli struct {
	RoomId   int    // 当前进入的roomId
	RoomName string // 当前进入的roomName

	Host string
	Port string
}

var ins *Cli
var once sync.Once

func GetCli() *Cli {
	once.Do(func() {
		ins = &Cli{}
	})
	return ins
}
