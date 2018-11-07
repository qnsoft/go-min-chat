package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/room"
)

func CreateRooms(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	all_room := MinChatSer.AllRoomKeyRoomName
	if exist_room, ok := all_room[rcvContent.ParamString]; ok { // 存在房间了
		param := fmt.Sprintf("%s room is existing", exist_room.Name)
		SendSuccessFailMessage(conn, param)
	} else { // 房间不存在
		// 创建了当前用户的房间信息
		user := MinChatSer.AllUser[conn]
		isAuth := CheckAuth(conn);
		if (!isAuth) {
			return
		}
		rootSing := room.GetRoom()
		roomId := int(room.GetRoomNo(rootSing))
		roomName := rcvContent.ParamString
		newRome := room.BuildRoom(roomId, roomName)
		newRome.CreateUid = user.Uid
		// 把room添加到chatSer保存
		ser.AddRooms(newRome)
		SendSuccessFailMessage(conn, "OK")
	}
}
