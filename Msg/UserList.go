package Msg

import (
	"net"
	"go-min-chat/server/ser"
	"fmt"
	"strings"
)

func UserList(conn net.Conn) {
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
	SendSuccessMessage(conn, strings.TrimSuffix(allUserStr, "\n"))
}
