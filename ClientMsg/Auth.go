package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"go-min-chat/ClientApp"
	"go-min-chat/Util"
	"os"
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
