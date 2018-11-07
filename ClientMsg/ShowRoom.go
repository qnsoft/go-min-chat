package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"fmt"
)

func ShowRoom(backContent *protobuf.BackContent) {
	fmt.Println(backContent.Showroom.RoomsAndIds)
}
