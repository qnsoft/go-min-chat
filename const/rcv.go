package _const

import (
	"fmt"
	"strings"
	"go-min-chat/Command"
)

const RCV_UNKNOW = 1
const RCV_AUTH = 2
const RCV_SHOW_ROOMS = 3
const RCV_CREATE_ROOM = 4
const RCV_USE_ROOM = 5
const RCV_GROUP_MSG = 6
const RCV_USER_LIST = 7
const RCV_WHOAMI = 8
const RCV_SUCCESS = 9
const RCV_FAIL = 10

//const RCV_SHOW_ROOMS = 3
//const RCV_CREATE_ROOM = 4
//const RCV_USE_ROOM = 5

//const RCV_USER_LIST = 7

//const RCV_WHOAMI = 9

//var All = map[int]interface{}{
//	RCV_SHOW_ROOMS:  Command.ShowRoom{},
//	RCV_CREATE_ROOM: Command.CreateRoom{},
//	RCV_USE_ROOM:    Command.UseRoom{},
//	RCV_USER_LIST:   Command.UserList{},
//	RCV_WHOAMI:      Command.WhoAmI{},
//}

func GetAllCommand(param []string) (int, interface{}) {
	param_leg := len(param)
	fmt.Println("param_leg:", param_leg, param[0])
	var msgId int
	var Cmd interface{}
	if param_leg >= 3 {
		p := strings.ToUpper(param[0]) + strings.ToUpper(param[1])
		if (p == "CREATE ROOM") {
			Cmd = Command.CreateRoom{}
		}
	} else if param_leg == 2 {
		p := strings.ToUpper(param[0] + strings.ToUpper(param[1]))
		if (p == "SHOW ROOMS") {
			Cmd = Command.ShowRoom{}
		} else if (p == "USER LIST") {
			Cmd = Command.UserList{}
		}
		if (strings.ToUpper(param[0]) == "USE") {
			Cmd = Command.UseRoom{}
		}
	} else {
		if (strings.ToUpper(param[0]) == "WHOAMI") {
			Cmd = Command.WhoAmI{}
		}
	}
	return msgId, Cmd
	//switch param_leg {
	//case 1:
	//	msgId, isContinue = OneParam[strings.ToUpper(param[0])]
	//	if (!isContinue) { // 没有找到
	//		msgId = RCV_UNKNOW
	//	}
	//case 2:
	//	msgId, isContinue = TwoParam[strings.ToUpper(param[0])]
	//	if (!isContinue) {
	//		msgId, isContinue = TwoParam[strings.ToUpper(param[0])+strings.ToUpper(param[1])]
	//		if (!isContinue) {
	//			msgId = RCV_UNKNOW
	//		} else {
	//			isContinue = false
	//			Util.EchoLine("(error) ERR " + strings.Join(param, " ")+
	//				" command missing parameter", 2)
	//			fmt.Printf(ClientUtil.GetPre())
	//		}
	//	} else {
	//		if (msgId == RCV_WHOAMI) {
	//			a = Command.WhoAmI{}
	//		} else if (msgId == RCV_USE_ROOM) {
	//			b := Command.UseRoom{}
	//			b.RoomName = param[1]
	//			a = b
	//		}
	//	}
	//default:
	//	msgId, isContinue = TwoParam[strings.ToUpper(param[0])]
	//	if (!isContinue) {
	//		msgId, isContinue = TwoParam[strings.ToUpper(param[0])+strings.ToUpper(param[1])]
	//		if (!isContinue) {
	//			msgId = RCV_UNKNOW
	//		} else {
	//			isContinue = false
	//			Util.EchoLine("(error) ERR " + strings.Join(param, " ")+
	//				" command missing parameter", 2)
	//			fmt.Printf(ClientUtil.GetPre())
	//		}
	//	} else {
	//		if (msgId == RCV_WHOAMI) {
	//			a = Command.WhoAmI{}
	//		} else if (msgId == RCV_USE_ROOM) {
	//			b := Command.UseRoom{}
	//			b.RoomName = param[1]
	//			a = b
	//		}
	//	}
	//}

}

//key := strings.Join(param, " ")
//
//v, ok := m[key];
//if ok {
//	return v, true
//}
//return v, false
