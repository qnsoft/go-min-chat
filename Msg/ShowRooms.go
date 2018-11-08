package Msg

import (
	"net"
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/protobuf/proto"
	"go-min-chat/const"
	"github.com/golang/protobuf/proto"
)

func ShowRooms(conn net.Conn) {
	MinChatSer := ser.GetMinChatSer()
	rooms := MinChatSer.AllRoomKeyRoomId
	var innerRet string
	if (len(rooms) == 0) {
		innerRet = "no room"
	} else {
		for v, r := range rooms {
			if (v == 1) {
				innerRet = fmt.Sprintf("%d)roomName:%s (roomId:%d)", v, r.Name, r.Id)
			} else {
				innerRet = fmt.Sprintf("%s\n%d)roomName:%s (roomId:%d)", innerRet, v, r.Name, r.Id)
			}
			if (r.Id == MinChatSer.AllUser[conn].RoomId) {
				innerRet += "*"
			}
		}
	}
	p1 := &protobuf.BackContent{}
	p1.Id = _const.RCV_SHOW_ROOMS
	sR := &protobuf.ShowRoom{}
	sR.Count = int32(len(rooms))
	sR.RoomsAndIds = innerRet
	p1.Showroom = sR
	data, _ := proto.Marshal(p1)
	SendMessage(conn, data)
}
