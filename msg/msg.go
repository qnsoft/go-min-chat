package msg

import (
	"go-min-chat/server/ser"
	"fmt"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"strings"
	"go-min-chat/room"
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
			innerRet = fmt.Sprintf("%d)%s(%d)", v+1, r.Name, r.Id)
		} else {
			innerRet = fmt.Sprintf("%s\n%d)%s(%d)", innerRet, v+1, r.Name, r.Id)
		}
	}
	*ret = innerRet
}

func doCreateRooms(rcvContent *protobuf.Content, ret *string) {
	MinChatSer := ser.GetMinChatSer()
	var isExist = false
	for _, v := range MinChatSer.Rooms { // 房间名字已经存在
		if (strings.EqualFold(v.Name, rcvContent.Param)) {
			*ret = fmt.Sprintf("%s room is existing", rcvContent.Param)
			*ret = "OK"
			isExist = true
			goto Loop
		}
	}
Loop:
	if (!isExist) {
		rootSing := room.GetRoom()
		MinChatSer.Rooms = append(MinChatSer.Rooms, room.BuildRoom(int(room.GetRoomNo(rootSing)), rcvContent.Param))
		*ret = "OK"
	}
}
