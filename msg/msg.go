package msg

import (
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"strings"
	"net"
	"go-min-chat/room"
)

// server 接受信息格式
const UNKNOW = 0

const RCV_AUTH = 1
const RCV_SHOW_ROOMS = 2
const RCV_CREATE_ROOM = 3
const RCV_USE_ROOM = 4
const RCV_GROUP_MSG = 5
const RCV_USER_LIST = 6
const RCV_SUCCESS_FAIL = 7

// server 发送的消息格式
//const SEND_CREATEROOM = 2
//const SEND_USEROOM = 3

const OK = "OK"

func DoMsg(conn net.Conn, msgContent []byte) {
	rcvContent := &protobuf.Content{}
	proto.Unmarshal(msgContent, rcvContent)
	fmt.Println("收到一个消息:", rcvContent.Id)
	switch rcvContent.Id {
	case RCV_AUTH:
		doAuth(conn, rcvContent)
		break
	case RCV_SHOW_ROOMS:
		doShowRooms(conn)
		break
	case RCV_CREATE_ROOM:
		doCreateRooms(conn, rcvContent)
		break
	case RCV_USE_ROOM:
		doUseRoom(conn, rcvContent)
		break
	case RCV_GROUP_MSG:
		doGroupMsg(conn, rcvContent)
		break
	case RCV_USER_LIST:
		doUserList(conn)
		break
	}
}

func doUserList(conn net.Conn) {
	MinChatSer := ser.GetMinChatSer()
	u := MinChatSer.AllUser[conn]
	fmt.Println("roomName:", u.RoomName, "roomId:", u.RoomId)
	allUser := MinChatSer.AllRoomKeyRoomId[u.RoomId].AllUser
	fmt.Println(allUser)
	var allUserStr string
	for _, v := range allUser {
		if (u == v) { // 如果是自己就加个*
			allUserStr += v.Nick + "(*)\n"
		} else {
			allUserStr += v.Nick + "\n"
		}
	}
	SendBackMessage(conn, 1, 1, strings.TrimSuffix(allUserStr, "\n"))
}

func doAuth(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	u := MinChatSer.AllUser[conn]
	p1 := &protobuf.BackContent{}
	p1.Id = RCV_AUTH
	if (strings.EqualFold(rcvContent.ParamString, "thomas")) {
		// 设置本地的
		u.Uid = 1;
		u.IsAuth = true
		u.Age = 19
		u.Nick = rcvContent.ParamString
		// 组装数据给客户端返回
		userinfo := &protobuf.Userinfo{}
		userinfo.Nick = rcvContent.ParamString
		userinfo.Uid = int32(u.Uid)

		auth := &protobuf.Auth{}
		auth.IsOk = true
		auth.UseInfo = userinfo

		p1.Auth = auth

	}
	if (strings.EqualFold(rcvContent.ParamString, "wang")) {
		u.Uid = 1;
		u.IsAuth = true
		u.Age = 19
		u.Nick = rcvContent.ParamString

		userinfo := &protobuf.Userinfo{}
		userinfo.Nick = rcvContent.ParamString
		userinfo.Uid = int32(u.Uid)

		auth := &protobuf.Auth{}
		auth.IsOk = true
		auth.UseInfo = userinfo

		p1.Auth = auth
	}
	data, _ := proto.Marshal(p1)
	SendMessage(conn, data);
}

func doShowRooms(conn net.Conn) {
	MinChatSer := ser.GetMinChatSer()
	var innerRet string = "1)"
	rooms := MinChatSer.AllRoomKeyRoomId

	for v, r := range rooms {
		if (v == 0) {
			innerRet = fmt.Sprintf("%d)%s(%d)", v+1, r.Name, r.Id)
		} else {
			innerRet = fmt.Sprintf("%s\n%d)%s(%d)", innerRet, v+1, r.Name, r.Id)
		}
	}
	p1 := &protobuf.BackContent{}
	p1.Id = RCV_SHOW_ROOMS
	sR := &protobuf.ShowRoom{}
	sR.Count = int32(len(rooms))
	sR.RoomsAndIds = innerRet
	p1.Showroom = sR
	data, _ := proto.Marshal(p1)
	SendMessage(conn, data)
}

func doCreateRooms(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	all_room := MinChatSer.AllRoomKeyRoomName
	if exist_room, ok := all_room[rcvContent.ParamString]; ok { // 存在房间了
		param := fmt.Sprintf("%s room is existing", exist_room.Name)
		SendBackMessage(conn, 1, 1, param)
	} else { // 房间不存在
		// 创建了当前用户的房间信息
		user := MinChatSer.AllUser[conn]
		if (!user.IsAuth) { // 没有登录是不能创建房间的
			p1 := &protobuf.BackContent{}
			p1.Id = RCV_AUTH
			auth := &protobuf.Auth{}
			auth.IsOk = false
			p1.Auth = auth
			data, _ := proto.Marshal(p1)
			SendMessage(conn, data)
			return
		}
		rootSing := room.GetRoom()
		roomId := int(room.GetRoomNo(rootSing))
		roomName := rcvContent.ParamString
		newRome := room.BuildRoom(roomId, roomName)
		newRome.CreateUid = user.Uid
		// 把room添加到chatSer保存
		ser.AddRooms(newRome)
		//SendBackMessage(conn, 1, 1, "OK")
		p2 := &protobuf.BackContent{}
		p2.Id = RCV_SUCCESS_FAIL
		p2.Msg = "OK"
		data2, _ := proto.Marshal(p2)
		SendMessage(conn, data2)
	}
}

func doUseRoom(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	if r, ok := MinChatSer.AllRoomKeyRoomName[rcvContent.ParamString]; ok {
		if (r.Id == user.RoomId) { // 在当前房间
			param := fmt.Sprintf("you are already in %s room", rcvContent.ParamString)
			SendBackMessage(conn, 1, 1, param)
		} else { // 不在当前房间
			user.RoomName = r.Name
			user.RoomId = r.Id
			r.AllUser[user.Uid] = user
			//a.Uid // 用户id
			SendBackMessage(conn, 1, 1, "OK")
			SendBackMessage(conn, RCV_USE_ROOM, 1, fmt.Sprintf("%d %s", 1, rcvContent.ParamString))
		}
	} else { // 不存在
		SendBackMessage(conn, 1, 1, "room "+rcvContent.ParamString+" is not found")
		return
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
	p1 := &protobuf.BackContent{Id: id,}
	data, _ := proto.Marshal(p1)
	n, err := conn.Write(data)
	fmt.Println("data send leng:", n, err)
}

func SendMessage(conn net.Conn, data []byte) {
	n, err := conn.Write(data)
	fmt.Println("data send leng:", n, err)
}
