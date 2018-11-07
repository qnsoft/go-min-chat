package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"fmt"
	"go-min-chat/server/ser"
	"go-min-chat/const"
	"github.com/golang/protobuf/proto"
)

func GroupMsg(conn net.Conn, rcvContent *protobuf.Content) {
	fmt.Println("group msg:", rcvContent.ParamString)
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn] // 当前这个用户
	Room := MinChatSer.AllRoomKeyRoomId[user.RoomId]
	p1 := &protobuf.BackContent{}
	p1.Id = _const.RCV_GROUP_MSG
	groupMsg := &protobuf.GroupMsg{}
	groupMsg.Content = rcvContent.ParamString
	groupMsg.Uid = int32(user.Uid)
	groupMsg.Nick = user.Nick
	p1.Groupmsg = groupMsg

	data, _ := proto.Marshal(p1)
	for _, v := range Room.AllUser {
		if (conn != v.Conn) {
			SendMessage(v.Conn, data)
		}
	}
}
