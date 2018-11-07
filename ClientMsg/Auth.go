package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"go-min-chat/ClientApp"
	"go-min-chat/Util"
	"os"
	"go-min-chat/const"
	"github.com/golang/protobuf/proto"
)

func Auth(backContent *protobuf.BackContent) {
	cli1 := ClientApp.GetCli()
	cli1.IsAuth = backContent.Auth.IsOk
	if (backContent.Auth.IsOk) { // 登录成功
		cli1.Nick = backContent.Auth.UseInfo.Nick
		cli1.Uid = backContent.Auth.UseInfo.Uid
		Util.EchoLine("OK", 1)
	} else { // 登录失败 or 需要去登录
		Util.EchoLine(backContent.Auth.Msg, 2)
		os.Exit(1)
	}
}

func SendAuth(nick string, password string) []byte {
	p1 := &protobuf.Content{}
	p1.Id = _const.RCV_AUTH
	p1.Nick = nick
	p1.Password = password
	data, _ := proto.Marshal(p1)
	return data
}
