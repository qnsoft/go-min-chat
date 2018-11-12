package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"go-min-chat/Utils"
)

func ShowRoom(backContent *protobuf.BackContent) {
	Utils.EchoLine(backContent.Showroom.RoomsAndIds, 1)
}
