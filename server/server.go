package main

import (
	"net"
	"fmt"
	"os"
	"flag"
	"go-min-chat/msg"
	"io"
	"go-min-chat/user"
	"go-min-chat/server/ser"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func ConnReadCheckError(err error, conn net.Conn) int {
	if err != nil {
		if (err == io.EOF) {
			fmt.Printf("client %s is close!\n", conn.RemoteAddr().String())
			return -1
		} else {
			fmt.Println("Error: %s", err.Error())
			return 0
		}
	}
	return 1
}

type Server struct {
	host string
	port int
}

var S *Server

func init() {
	var host string
	flag.StringVar(&host, "h", "127.0.0.1", "is port")
	var port int
	flag.IntVar(&port, "p", 8080, "is port")
	flag.Parse()
	S = &Server{host, port}
}

func main() {
	addr := fmt.Sprintf("%s:%d", S.host, S.port)
	listen, err := net.Listen("tcp", addr)
	checkError(err)
	defer listen.Close()
	fmt.Println("Ready to accept connections")
	MinChatSer := ser.GetMinChatSer()
	var u *user.User
	for {
		newConn, err := listen.Accept()
		// 连接上了，就把这个连接赋予一个未登录的用户
		u = user.BuildUser(0, "", 0, false)
		MinChatSer.AllUser[newConn] = u
		fmt.Println(newConn.RemoteAddr())
		checkError(err)
		ch := make(chan []byte)
		go recvConnMsg(newConn, ch)
		go sendConnMsg(newConn, ch)
	}
}

// 服务端接受消息
func recvConnMsg(conn net.Conn, ch chan []byte) {
	buf := make([]byte, 50)
Loop:
	for {
		fmt.Println("-----");
		n, err := conn.Read(buf)
		ret := ConnReadCheckError(err, conn)
		if (ret == 0) { // 读取时, 发生了错误
			os.Exit(1)
		} else if (ret == -1) { // 客户端断开了连接
			break Loop
		}
		ch <- buf[:n]
		fmt.Println(string(buf[0:n]))
	}
}

// 服务端发送消息
func sendConnMsg(conn net.Conn, ch chan []byte) {
	for {
		fmt.Println("-----")
		content, _ := <-ch
		var ret string
		msg.DoMsg(conn, content)
		fmt.Println(ret)
	}
}
