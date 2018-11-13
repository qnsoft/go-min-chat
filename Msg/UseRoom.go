package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"go-min-chat/cache"
	"fmt"
	"go-min-chat/const"
	"github.com/golang/protobuf/proto"
	"go-min-chat/server/ser"
	"reflect"
	"strconv"
)

func UseRoom(conn net.Conn, rcvContent *protobuf.Content) {
	MinChatSer := ser.GetMinChatSer()
	user := MinChatSer.AllUser[conn]
	key := "RoomName:" + rcvContent.ParamString
	redis := cache.GetReis().Conn
	a, _ := redis.Get(key).Result()
	value, _ := strconv.Atoi(a)
	fmt.Println("reflect.TypeOf(value)", reflect.TypeOf(value))
	fmt.Println("value", value)
	if value != 0 {
		if (value == user.RoomId) { // 在当前房间
			SendFailMessage(conn, fmt.Sprintf("you are already in %s room", rcvContent.ParamString))
		} else { // 不在当前房间
			user.RoomName = rcvContent.ParamString
			user.RoomId = value
			redis.SAdd("Set:RoomName:"+rcvContent.ParamString, user.Uid)
			p1 := &protobuf.BackContent{}
			room1 := &protobuf.Room{}
			room1.RoomId = int32(value)
			room1.RoomName = rcvContent.ParamString
			p1.Id = _const.RCV_USE_ROOM
			p1.Room = room1
			data, _ := proto.Marshal(p1)
			SendSuccessMessage(conn, "OK")
			SendMessage(conn, data)
		}
	} else { // 不存在
		SendFailMessage(conn, "room "+rcvContent.ParamString+" is not found")
		return
	}
}
