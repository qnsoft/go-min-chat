package Msg

import (
"net"
"go-min-chat/protobuf/proto"
"github.com/golang/protobuf/proto"
"fmt"
"go-min-chat/const"
)

func SendBackMessage(conn net.Conn, id int32, msgType int32, param string) {
	p1 := &protobuf.BackContent{Id: id,}
	data, _ := proto.Marshal(p1)
	n, err := conn.Write(data)
	fmt.Println("data send leng:", n, err)
}

func SendMessage(conn net.Conn, data []byte) {
	n, err := conn.Write(data)
	fmt.Println("data send leng:", n, err)
}

func SendSuccessFailMessage(conn net.Conn, msg string) {
	p1 := &protobuf.BackContent{}
	p1.Id = _const.RCV_SUCCESS_FAIL
	p1.Msg = msg
	data, _ := proto.Marshal(p1)
	conn.Write(data)
}
