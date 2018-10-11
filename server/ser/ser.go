package ser

import (
	"sync"
	"go-min-chat/room"
	"net"
	"go-min-chat/user"
)

type MinChatSer struct {
	Host    string
	Port    int
	Rooms   []room.Room
	AllUser map[net.Conn]user.User
}

var ins *MinChatSer
var once sync.Once

func GetMinChatSer() *MinChatSer {
	once.Do(func() {
		ins = &MinChatSer{}
		ins.AllUser = make(map[net.Conn]user.User)
	})
	return ins
}

func AddRooms(room room.Room) {
	singleMinChatSer := GetMinChatSer()
	singleMinChatSer.Rooms = append(singleMinChatSer.Rooms, room)
}