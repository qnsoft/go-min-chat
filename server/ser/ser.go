package ser

import (
	"sync"
	"go-min-chat/room"
	"net"
	"go-min-chat/mysql"
)

type MinChatSer struct {
	Host               string
	Port               int
	AllRoomKeyRoomId   map[int]*room.Room
	AllRoomKeyRoomName map[string]*room.Room
	AllUser            map[net.Conn]*mysql.User
}

var ins *MinChatSer
var once sync.Once

func GetMinChatSer() *MinChatSer {
	once.Do(func() {
		ins = &MinChatSer{}
		ins.AllUser = make(map[net.Conn]*mysql.User)
		ins.AllRoomKeyRoomId = make(map[int]*room.Room)
		ins.AllRoomKeyRoomName = make(map[string]*room.Room)
	})
	return ins
}

func AddRooms(room *room.Room) {
	singleMinChatSer := GetMinChatSer()
	singleMinChatSer.AllRoomKeyRoomId[room.Id] = room
	singleMinChatSer.AllRoomKeyRoomName[room.Name] = room
}
