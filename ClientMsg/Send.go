package ClientMsg

import (
	"net"
	"go-min-chat/Utils"
)

func Send(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		_, err := conn.Write(content)
		Utils.CheckError(err)
	}
}
