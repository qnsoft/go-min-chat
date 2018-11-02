package room

import "go-min-chat/user"

type Room struct {
	Id        int
	Name      string
	CreateUid int // 创建这个房间的用户id
	AllUser   map[int]*user.User
}

func BuildRoom(id int, name string) *Room {
	var room *Room
	room = &Room{}
	room.Id = id
	room.Name = name
	room.AllUser = make(map[int]*user.User)

	return room
}
