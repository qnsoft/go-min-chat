package msg

import (
	"go-min-chat/room"
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"strings"
)

const UNKNOW = 0
const SHOWROOMS = 1
const CREATEROOM = 2

const OK = "OK"

func DoMsg(msgContent []byte, ret *string) {
	rcvContent := &protobuf.Content{}
	proto.Unmarshal(msgContent, rcvContent)

	switch rcvContent.Id {
	case SHOWROOMS:
		doShowRooms(ret)
	case CREATEROOM:
		doCreateRooms(rcvContent, ret)
	}

}

func doShowRooms(ret *string) {
	MinChatSer := ser.GetMinChatSer()
	var innerRet string = "1)"
	rooms := MinChatSer.Rooms
	for v, r := range rooms {
		if (v == 0) {
			innerRet = fmt.Sprintf("%d)%s", v+1, r.Name)
		} else {
			innerRet = fmt.Sprintf("%s\n%d)%s", innerRet, v+1, r.Name)
		}
	}
	*ret = innerRet
}

func doCreateRooms(rcvContent *protobuf.Content, ret *string) {
	MinChatSer := ser.GetMinChatSer()
	for _, v := range MinChatSer.Rooms { // 房间名字已经存在
		if (strings.EqualFold(v.Name, rcvContent.Param)) {
			*ret = fmt.Sprintf("%s room is existing", rcvContent.Param)
			*ret = "OK"
			goto Loop
		}
	}
	MinChatSer.Rooms = append(MinChatSer.Rooms, room.BuildRoom(12, rcvContent.Param))
	*ret = "OK"
Loop:
}
