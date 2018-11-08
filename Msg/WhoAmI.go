package Msg

import (
	"net"
	"go-min-chat/server/ser"
)

func WhoAmI(conn net.Conn) {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	SendSuccessMessage(conn, "Nick:"+user.Nick)
}
