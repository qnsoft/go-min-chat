package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"go-min-chat/server/ser"
	"go-min-chat/room"
	"strings"
	"fmt"
)

func CreateRooms(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	all_room := MinChatSer.AllRoomKeyRoomName
	rooms := strings.Split(rcvContent.ParamString, ",")
	var success_room []string
	var fail_room []string
	for _, roomName := range rooms {
		if _, ok := all_room[roomName]; ok { // 存在房间了
			fail_room = append(fail_room, roomName)
			//param := fmt.Sprintf("%s room is existing", exist_room.Name)
			//SendSuccessFailMessage(conn, param)
		} else {
			rootSing := room.GetRoom()
			roomId := int(room.GetRoomNo(rootSing))
			newRome := room.BuildRoom(roomId, roomName)
			newRome.CreateUid = user.Uid
			// 把room添加到chatSer保存
			ser.AddRooms(newRome)
			success_room = append(success_room, roomName)
		}
	}
	ret1 := fmt.Sprintf("%s %s ", strings.Join(fail_room, ","), "existing")
	ret2 := fmt.Sprintf("%s %s", strings.Join(success_room, ","), "success")

	if (len(fail_room) > 0) {
		SendFailMessage(conn, ret1)
	}
	if (len(success_room) > 0) {
		SendSuccessMessage(conn, ret2)
	}
}
