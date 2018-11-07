package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"go-min-chat/ClientApp"
)

func UseRoom(backContent *protobuf.BackContent) {
	cli1 := ClientApp.GetCli()
	cli1.RoomId = int(backContent.Room.RoomId)
	cli1.RoomName = backContent.Room.RoomName
}
