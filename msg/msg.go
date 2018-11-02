package msg

import (
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"strings"
	"go-min-chat/room"
	"net"
	"go-min-chat/cli"
)

// server 接受信息格式
const UNKNOW = 0

const RCV_AUTH = 1
const RCV_SHOWROOMS = 2
const RCV_CREATEROOM = 3
const RCV_USEROOM = 4
const RCV_GROUP_MSG = 5

// server 发送的消息格式
const SEND_CREATEROOM = 2
const SEND_USEROOM = 3

const OK = "OK"

func DoMsg(conn net.Conn, msgContent []byte) {
	rcvContent := &protobuf.Content{}
	proto.Unmarshal(msgContent, rcvContent)
	fmt.Println("收到一个消息:", rcvContent.Id)
	switch rcvContent.Id {
	case RCV_AUTH:
		doAuth(conn, rcvContent)
		break
	case RCV_SHOWROOMS:
		doShowRooms(conn)
		break
	case RCV_CREATEROOM:
		doCreateRooms(conn, rcvContent)
		break
	case RCV_USEROOM:
		doUseRoom(conn, rcvContent)
		break
	case RCV_GROUP_MSG:
		doGroupMsg(conn, rcvContent)
		break
	}
}

func doAuth(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	u := MinChatSer.AllUser[conn]

	if (strings.EqualFold(rcvContent.ParamString, "thomas")) {
		u.Uid = 1;
		u.IsAuth = true
		u.Age = 19
		u.Nick = rcvContent.ParamString
	}
	if (strings.EqualFold(rcvContent.ParamString, "wang")) {
		u.Uid = 1;
		u.IsAuth = true
		u.Age = 19
		u.Nick = rcvContent.ParamString
	}
	SendBackMessage(conn, 1, 1, "ok")
}

func doShowRooms(conn net.Conn) {
	MinChatSer := ser.GetMinChatSer()
	var innerRet string = "1)"
	rooms := MinChatSer.AllRoom
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
	for _, v := range MinChatSer.AllRoom { // 房间名字已经存在
		if (strings.EqualFold(v.Name, rcvContent.ParamString)) {
			param := fmt.Sprintf("%s room is existing", rcvContent.ParamString)
			SendBackMessage(conn, 1, 1, param)
			isExist = true
			goto Loop
		}
	}
Loop:
	if (isExist == false) { // 不存在就创建
		// 创建了当前用户的房间信息
		user := MinChatSer.AllUser[conn]
		if (!user.IsAuth) { // 没有登录是不能创建房间的
			SendBackMessage(conn, 1, 1, "please auth first")
			return
		}
		rootSing := room.GetRoom()
		roomId := int(room.GetRoomNo(rootSing))
		roomName := rcvContent.ParamString
		newRome := room.BuildRoom(roomId, roomName)
		// 把room添加到chatSer保存
		ser.AddRooms(newRome)

		newRome.CreateUid = user.Uid
		SendBackMessage(conn, 1, 1, "OK")
		//SendBackMessage(conn, SEND_CREATEROOM, 1, fmt.Sprintf("%d %s", roomId, roomName))
	}
}

func doUseRoom(conn net.Conn, rcvContent *protobuf.Content) {
	if (cli.GetCli().RoomName == rcvContent.ParamString) { // 在当前房间
		param := fmt.Sprintf("u are already in %s room", rcvContent.ParamString)
		SendBackMessage(conn, 1, 1, param)
	} else { // 不在当前房间
		a := ser.GetMinChatSer().AllUser[conn]
		a.RoomName = rcvContent.ParamString
		//a.Uid // 用户id
		SendBackMessage(conn, 1, 1, "OK")
		SendBackMessage(conn, SEND_USEROOM, 1, fmt.Sprintf("%d %s", 1, rcvContent.ParamString))
	}
}

func doGroupMsg(conn net.Conn, rcvContent *protobuf.Content) {
	fmt.Println("group msg:", rcvContent.ParamString)
	//room := ser.GetMinChatSer().AllRoom[int(rcvContent.ParamId)]

	//if (cli.GetCli().RoomName == rcvContent.ParamString) { // 在当前房间
	//	param := fmt.Sprintf("u are already in %s room", rcvContent.ParamString)
	//	SendBackMessage(conn, 1, 1, param)
	//} else { // 不在当前房间
	//	a := ser.GetMinChatSer().AllUser[conn]
	//	a.RoomName = rcvContent.ParamString
	//	SendBackMessage(conn, 1, 1, "OK")
	//	SendBackMessage(conn, SEND_USEROOM, 1, fmt.Sprintf("%d %s", 1, rcvContent.ParamString))
	//}
}

func SendBackMessage(conn net.Conn, id int32, msgType int32, param string) {
	p1 := &protobuf.BackContent{
		Id:      id,
		MsqType: msgType,
		Param:   param,
	}
	data, _ := proto.Marshal(p1)
	n, err := conn.Write(data)
	fmt.Println("data send leng:", n, err)
}
