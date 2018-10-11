package user

import "net"

type User struct {
	Nick     string
	Age      int
	RoomId   int
	RoomName string
	Online   bool
	IsAuth   bool
	WaitPass bool
	Conn     *net.Conn
}
