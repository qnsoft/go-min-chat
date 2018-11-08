package Command

import "go-min-chat/room"

type CreateRoom struct {
	Common
	Room []room.Room
}
