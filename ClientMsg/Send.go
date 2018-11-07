package ClientMsg

import (
	"net"
	"go-min-chat/Util"
)

func Send(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		_, err := conn.Write(content)
		Util.CheckError(err)
	}
}
