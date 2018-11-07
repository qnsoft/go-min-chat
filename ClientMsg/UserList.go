package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"fmt"
)

func UserList(backContent *protobuf.BackContent) {
	fmt.Println(backContent.Showroom.RoomsAndIds)
}
