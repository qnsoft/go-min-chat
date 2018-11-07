package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"go-min-chat/server/ser"
	"go-min-chat/const"
	"fmt"
	"go-min-chat/mysql"
	"go-min-chat/login"
	"github.com/golang/protobuf/proto"
)

func Auth(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	u := MinChatSer.AllUser[conn]
	p1 := &protobuf.BackContent{}
	p1.Id = _const.RCV_AUTH
	fmt.Println("有人来登录", rcvContent.Nick, rcvContent.Password)
	logUser := mysql.UserForLog{Nick: rcvContent.Nick, Password: rcvContent.Password}
	isAuth := login.CheckLog(u, logUser)
	if (isAuth) {
		u.Conn = conn
		// 组装数据给客户端返回
		userinfo := &protobuf.Userinfo{}
		userinfo.Nick = u.Nick
		userinfo.Uid = int32(u.Uid)

		auth := &protobuf.Auth{}
		auth.IsOk = true
		auth.UseInfo = userinfo
		p1.Auth = auth
	} else {
		// 组装数据给客户端返回
		auth := &protobuf.Auth{}
		auth.IsOk = false
		auth.Msg = "用户名或密码错误"
		p1.Auth = auth
	}
	data, _ := proto.Marshal(p1)
	SendMessage(conn, data)
}
