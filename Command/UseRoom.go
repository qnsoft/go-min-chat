package Command

import "go-min-chat/room"

// 进入房间 user RoomName
type UseRoom struct {
	Common
	Room room.Room
}
