package Msg

import (
	"net"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"fmt"
	"go-min-chat/const"
	"encoding/binary"
	"bytes"
)

func SendBackMessage(conn net.Conn, id int32, msgType int32, param string) {
	p1 := &protobuf.BackContent{Id: id,}
	data, _ := proto.Marshal(p1)
	n, err := conn.Write(data)
	fmt.Println("data send leng:", n, err)
}

func SendMessage(conn net.Conn, data []byte) {
	headSize := len(data)
	var headBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	var buffer bytes.Buffer

	buffer.Write(headBytes)
	buffer.Write(data)
	b3 := buffer.Bytes() //得到了b1+b2的结果
	n, err := conn.Write(b3)
	fmt.Println("data send leng:", n, err)
}

func SendFailMessage(conn net.Conn, msg string) {
	sendSuccessFailMessage(conn, msg, _const.RCV_FAIL)
}

func SendSuccessMessage(conn net.Conn, msg string) {
	sendSuccessFailMessage(conn, msg, _const.RCV_SUCCESS)
}

func sendSuccessFailMessage(conn net.Conn, msg string, msg_id int32) {
	p1 := &protobuf.BackContent{}
	p1.Id = msg_id
	p1.Msg = msg
	data, _ := proto.Marshal(p1)
	SendMessage(conn, data)
}
