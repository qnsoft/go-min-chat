package ser

import (
	"sync"
	"go-min-chat/room"
)

type MinChatSer struct {
	Host  string
	Port  int
	Rooms []room.Room
}

var ins *MinChatSer
var once sync.Once

func GetMinChatSer() *MinChatSer {
	once.Do(func() {
		ins = &MinChatSer{}
	})
	return ins
}

func AddRooms(room room.Room) {
	singleMinChatSer := GetMinChatSer()
	singleMinChatSer.Rooms = append(singleMinChatSer.Rooms, room)
}
