package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"strings"
	"fmt"
	"go-min-chat/cache"
	"strconv"
	"go-min-chat/const"
)

func CreateRooms(conn net.Conn, rcvContent *protobuf.Content) {
	rooms := strings.Split(rcvContent.ParamString, ",")
	var success_room []string
	var fail_room []string
	redis := cache.GetReis().Conn
	//todo 事务， 多条命令一起执行
	for _, roomName := range rooms {
		iSmember, _ := redis.SIsMember(_const.ALLROOM, roomName).Result()
		if (iSmember) { // 存在房间了
			fail_room = append(fail_room, roomName)
		} else {
			incr, _ := redis.Incr(_const.ROOMIDINCR).Result()
			redis.SAdd(_const.ALLROOM, roomName)
			redis.Set("RoomName:"+roomName, incr, 0)
			redis.Set("RoomId:"+strconv.Itoa(int(incr)), roomName, 0)
			success_room = append(success_room, roomName)
		}
	}
	ret1 := fmt.Sprintf("%s %s ", strings.Join(fail_room, ","), "existing")
	ret2 := fmt.Sprintf("%s %s", strings.Join(success_room, ","), "success")

	if (len(fail_room) > 0) {
		SendFailMessage(conn, ret1)
	}
	if (len(success_room) > 0) {
		SendSuccessMessage(conn, ret2)
	}
}
