package Msg

import (
	"net"
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/cache"
	"strconv"
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
	redis := cache.GetReis().Conn
	a, _ := redis.SMembers("Set:RoomName:" + u.RoomName).Result()
	fmt.Println("a:", a)
	var allUserStr string
	for _, v := range a {
		b, _ := strconv.Atoi(v)
		if (u.Uid == b) { // 如果是自己就加个*
			allUserStr += u.Nick + "(*)\n"
		} else {
			allUserStr += u.Nick + "\n"
		}
	}
	SendSuccessMessage(conn, strings.TrimSuffix(allUserStr, "\n"))
}
