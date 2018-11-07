package Msg

import (
	"fmt"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"net"
	"go-min-chat/const"
)

// server 接受信息格式
func DoAllMsg(conn net.Conn, msgContent []byte) {
	rcvContent := &protobuf.Content{}
	proto.Unmarshal(msgContent, rcvContent)
	fmt.Println("收到一个消息:", rcvContent.Id)
	switch rcvContent.Id {
	case _const.RCV_AUTH:
		Auth(conn, rcvContent)
		break
	case _const.RCV_SHOW_ROOMS:
		ShowRooms(conn)
		break
	case _const.RCV_CREATE_ROOM:
		CreateRooms(conn, rcvContent)
		break
	case _const.RCV_USE_ROOM:
		UseRoom(conn, rcvContent)
		break
	case _const.RCV_GROUP_MSG:
		GroupMsg(conn, rcvContent)
		break
	case _const.RCV_USER_LIST:
		UserList(conn)
		break
	}
}

