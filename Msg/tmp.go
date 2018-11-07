package Msg

import (
	"net"
	"go-min-chat/server/ser"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"go-min-chat/const"
)

func CheckAuth(conn net.Conn) bool {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	if (!user.IsAuth) { // 没有登录是不能创建房间的
		p1 := &protobuf.BackContent{}
		p1.Id = _const.RCV_AUTH
		auth := &protobuf.Auth{}
		auth.IsOk = false
		p1.Auth = auth
		data, _ := proto.Marshal(p1)
		SendMessage(conn, data)
		return false
	} else {
		return true
	}
}

func CheckRoom(conn net.Conn) bool {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	if (user.RoomId == 0) { // 没有进入房间
		SendSuccessFailMessage(conn, "请先进入房间")
		return false
	} else {
		return true
	}
}
