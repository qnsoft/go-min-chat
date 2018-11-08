package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"go-min-chat/Util"
)

func ShowRoom(backContent *protobuf.BackContent) {
	Util.EchoLine(backContent.Showroom.RoomsAndIds, 1)
}
