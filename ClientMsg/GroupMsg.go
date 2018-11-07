package ClientMsg

import (
	"go-min-chat/protobuf/proto"
	"fmt"
)

func GroupMsg(backContent *protobuf.BackContent) {
	fmt.Printf("%s: %s\n", backContent.Groupmsg.Nick, backContent.Groupmsg.Content)
}
