package ClientUtil

import (
	"strings"
	"go-min-chat/const"

	"fmt"
	"go-min-chat/protobuf/proto"
	"go-min-chat/ClientApp"
	"go-min-chat/Util"
)

func One(param []string) (int, bool) {
	if (strings.EqualFold(param[0], "whoami")) {
		return _const.RCV_WHOAMI, false
	} else {
		fmt.Println(fmt.Sprintf("(error) ERR unknown command `%s`", strings.Join(param, " ")))
		fmt.Printf(GetPre())
		return 0, true
	}
	return 0, false
}

func Two(param []string) (int, bool) {
	if (strings.EqualFold(param[0]+" "+param[1], "whoami")) {
		return _const.RCV_WHOAMI, false
	} else {
		fmt.Println(fmt.Sprintf("(error) ERR unknown command `%s`", strings.Join(param, " ")))
		fmt.Printf(GetPre())
		return 0, true
	}
	return 0, false
}

func GetMsgType(param []string) (bool, *protobuf.Content) {
	p1 := &protobuf.Content{}
	param_leg := len(param)
	var isContinue bool
	if param_leg >= 3 {
		p := Util.SliceUp(param[:2], " ")
		if (strings.EqualFold(p, "CREATE ROOM")) {
			p1.Id = _const.RCV_CREATE_ROOM
			p1.ParamString = strings.Join(param[2:], ",")
		}
	} else if param_leg == 2 {
		p := Util.SliceUp(param[:2], " ")
		if (strings.EqualFold(p, "SHOW ROOMS")) {
			p1.Id = _const.RCV_SHOW_ROOMS
		} else if (p == "USER LIST") {
			p1.Id = _const.RCV_USER_LIST
		}
		if (strings.EqualFold(param[0], "USE")) {
			p1.Id = _const.RCV_USE_ROOM
			p1.ParamString = param[1]
		}
	} else {
		if (strings.EqualFold(param[0], "WHOAMI")) {
			p1.Id = _const.RCV_WHOAMI
		}
	}

	if (p1.Id == 0) { // 任何一种消息都不是，就剩下两种可能了，1: 在房间里发送的聊天信息，2: err command
		cliSing := ClientApp.GetCli()
		if (cliSing.RoomId != 0) { // 说明进入房间了
			p1.Id = _const.RCV_GROUP_MSG
			p1.ParamString = param[0]
		} else {
			isContinue = true
			Util.EchoLine("(error) ERR unknown command `"+strings.Join(param, " ")+"`", 2)
		}
	}
	return isContinue, p1
}
