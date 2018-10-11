package msg

import (
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"strings"
	"go-min-chat/room"
	"net"
)

// server 接受信息格式
const UNKNOW = 0
const RCV_SHOWROOMS = 1
const RCV_CREATEROOM = 2

// server 发送的消息格式
const SEND_CREATEROOM = 2

const OK = "OK"

func DoMsg(conn net.Conn, msgContent []byte) {
	rcvContent := &protobuf.Content{}
	proto.Unmarshal(msgContent, rcvContent)

	switch rcvContent.Id {
	case RCV_SHOWROOMS:
		doShowRooms(conn)
	case RCV_CREATEROOM:
		doCreateRooms(conn, rcvContent)
	}
}

func doShowRooms(conn net.Conn) {
	MinChatSer := ser.GetMinChatSer()
	var innerRet string = "1)"
	rooms := MinChatSer.Rooms
	for v, r := range rooms {
		if (v == 0) {
			innerRet = fmt.Sprintf("%d)%s(%d)", v+1, r.Name, r.Id)
		} else {
			innerRet = fmt.Sprintf("%s\n%d)%s(%d)", innerRet, v+1, r.Name, r.Id)
		}
	}
	SendBackMessage(conn, 1, 1, innerRet)
}

func doCreateRooms(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	var isExist = false
	for _, v := range MinChatSer.Rooms { // 房间名字已经存在
		if (strings.EqualFold(v.Name, rcvContent.Param)) {
			param := fmt.Sprintf("%s room is existing", rcvContent.Param)
			SendBackMessage(conn, 1, 1, param)
			isExist = true
			goto Loop
		}
	}
Loop:
	if (isExist == false) { // 不存在就创建
		rootSing := room.GetRoom()
		roomId := int(room.GetRoomNo(rootSing))
		roomName := rcvContent.Param
		// 把room添加到chatSer保存
		MinChatSer.Rooms = append(MinChatSer.Rooms, room.BuildRoom(roomId, roomName))

		// 创建了当前用户的房间信息
		user := ser.GetMinChatSer().AllUser[conn]
		user.RoomId = roomId
		user.RoomName = roomName

		SendBackMessage(conn, 1, 1, "OK")
		SendBackMessage(conn, SEND_CREATEROOM, 1, fmt.Sprintf("%d %s", roomId, roomName))
	}
}

func SendBackMessage(conn net.Conn, id int32, msgType int32, param string) {
	p1 := &protobuf.BackContent{
		Id:      id,
		MsqType: msgType,
		Param:   param,
	}
	data, _ := proto.Marshal(p1)
	conn.Write(data)
}
