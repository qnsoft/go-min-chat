package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/const"
	"github.com/golang/protobuf/proto"
)

func UseRoom(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	if r, ok := MinChatSer.AllRoomKeyRoomName[rcvContent.ParamString]; ok {
		if (r.Id == user.RoomId) { // 在当前房间
			SendFailMessage(conn, fmt.Sprintf("you are already in %s room", rcvContent.ParamString))
		} else { // 不在当前房间
			user.RoomName = r.Name
			user.RoomId = r.Id
			r.AllUser[user.Uid] = user

			p1 := &protobuf.BackContent{}
			room1 := &protobuf.Room{}
			room1.RoomId = int32(r.Id)
			room1.RoomName = r.Name
			p1.Id = _const.RCV_USE_ROOM
			p1.Room = room1
			data, _ := proto.Marshal(p1)
			SendSuccessMessage(conn, "OK")
			SendMessage(conn, data)
			//SendSuccessFailMessage(conn, fmt.Sprintf("%d %s", 1, rcvContent.ParamString))

		}
	} else { // 不存在
		SendFailMessage(conn, "room "+rcvContent.ParamString+" is not found")
		return
	}
}
