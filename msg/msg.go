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

	isAuth := CheckAuth(conn)
	if (!isAuth) { // 这个自己没有登录就不能使用这个命令
		return
	}

	isRoom := CheckRoom(conn)
	if (!isRoom) { // 这个自己没有登录就不能使用这个命令
		return
	}
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
	SendSuccessFailMessage(conn, strings.TrimSuffix(allUserStr, "\n"))
}

func doAuth(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	u := MinChatSer.AllUser[conn]
	p1 := &protobuf.BackContent{}
	p1.Id = RCV_AUTH
	if (strings.EqualFold(rcvContent.Password, "123456")) { // 相等说明登录正确
		if (strings.EqualFold(rcvContent.Nick, "wang")) {
			u.Uid = 1;
			u.IsAuth = true
			u.Age = 19
			u.Nick = rcvContent.Nick
			u.Conn = conn

			// 组装数据给客户端返回
			userinfo := &protobuf.Userinfo{}
			userinfo.Nick = rcvContent.ParamString
			userinfo.Uid = int32(u.Uid)

			auth := &protobuf.Auth{}
			auth.IsOk = true
			auth.UseInfo = userinfo
			p1.Auth = auth
		}
	}
	data, _ := proto.Marshal(p1)
	SendMessage(conn, data);
}

func doShowRooms(conn net.Conn) {
	MinChatSer := ser.GetMinChatSer()
	rooms := MinChatSer.AllRoomKeyRoomId
	var innerRet string
	if (len(rooms) == 0) {
		innerRet = "no room"
	} else {
		for v, r := range rooms {
			if (v == 1) {
				innerRet = fmt.Sprintf("%d)%s(%d)", v, r.Name, r.Id)
			} else {
				innerRet = fmt.Sprintf("%s\n%d)%s(%d)", innerRet, v, r.Name, r.Id)
			}
			if (r.Id == MinChatSer.AllUser[conn].RoomId) {
				innerRet += "*"
			}
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

func doUseRoom(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	if r, ok := MinChatSer.AllRoomKeyRoomName[rcvContent.ParamString]; ok {
		if (r.Id == user.RoomId) { // 在当前房间
			SendSuccessFailMessage(conn, fmt.Sprintf("you are already in %s room", rcvContent.ParamString))
		} else { // 不在当前房间
			user.RoomName = r.Name
			user.RoomId = r.Id
			r.AllUser[user.Uid] = user

			p1 := &protobuf.BackContent{}
			room1 := &protobuf.Room{}
			room1.RoomId = int32(r.Id)
			room1.RoomName = r.Name
			p1.Id = RCV_USE_ROOM
			p1.Room = room1
			data, _ := proto.Marshal(p1)
			SendSuccessFailMessage(conn, "OK")
			SendMessage(conn, data)
			//SendSuccessFailMessage(conn, fmt.Sprintf("%d %s", 1, rcvContent.ParamString))

		}
	} else { // 不存在
		SendSuccessFailMessage(conn, "room "+rcvContent.ParamString+" is not found")
		return
	}
}

func doGroupMsg(conn net.Conn, rcvContent *protobuf.Content) {
	fmt.Println("group msg:", rcvContent.ParamString)
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn] // 当前这个用户
	Room := MinChatSer.AllRoomKeyRoomId[user.RoomId]
	p1 := &protobuf.BackContent{}
	p1.Id = RCV_GROUP_MSG
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

func SendSuccessFailMessage(conn net.Conn, msg string) {
	p1 := &protobuf.BackContent{}
	p1.Id = RCV_SUCCESS_FAIL
	p1.Msg = msg
	data, _ := proto.Marshal(p1)
	conn.Write(data)
}

func CheckAuth(conn net.Conn) bool {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	if (!user.IsAuth) { // 没有登录是不能创建房间的
		p1 := &protobuf.BackContent{}
		p1.Id = RCV_AUTH
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
