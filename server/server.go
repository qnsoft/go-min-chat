package main

import (
	"net"
	"fmt"
	"os"
	"flag"
	"go-min-chat/msg"
	"io"
)

type Client struct {
	nick      string
	age       int
	online    bool
	auth      bool
	wait_pass bool
}

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

var all_user = make(map[net.Conn]*Client)

type Server struct {
	host string
	port int
}

var S *Server

func init() {
	//p1 := &protobuf.Content{
	//	Id:      1,
	//	Command: "CREATE ROOM",
	//	Param:   "abc",
	//}
	//data, _ := proto.Marshal(p1);
	//ioutil.WriteFile("./test.txt", data, os.ModePerm)
	//newTest := &protobuf.Content{}
	//err := proto.Unmarshal(data, newTest)
	//checkError(err)
	//fmt.Println(newTest.Param)
	//os.Exit(1)

	var host string
	flag.StringVar(&host, "h", "192.168.101.160", "is port")
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
	fmt.Println("connect success")
	for {
		fmt.Println("-----");
		new_conn, err := listen.Accept()
		all_user[new_conn] = &Client{}
		fmt.Println(new_conn.RemoteAddr())
		checkError(err)
		ch := make(chan []byte)
		go recvConnMsg(new_conn, ch)
		go sendConnMsg(new_conn, ch)
	}
}

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

func sendConnMsg(conn net.Conn, ch chan []byte) {
	for {
		fmt.Println("-----");
		content, _ := <-ch
		var ret string
		msg.DoMsg(content, &ret)
		_, err := conn.Write([]byte(ret))
		checkError(err)
		fmt.Println(ret)

		//if (strings.HasPrefix(content_list, "auth")) {
		//	if (all_user[conn].auth) {
		//		conn.Write([]byte("do not auth, you already login"))
		//		continue
		//	} else {
		//		name := string(content_list[5:])
		//		all_user[conn].nick = name
		//		all_user[conn].wait_pass = true
		//		all_user[conn].auth = true
		//		var s = fmt.Sprintf("login successfully!!!, welcome %s", name)
		//		_, err := conn.Write([]byte(s))
		//		checkError(err)
		//		continue
		//	}
		//} else { // 如果不是auth指令, 就判断是否auth过
		//	if (all_user[conn].auth == false) {
		//		_, err := conn.Write([]byte("please auth"))
		//		checkError(err)
		//		continue
		//	}
		//}
		//if (content_list == "list") {
		//	var all_name bytes.Buffer
		//	for k, v := range all_user {
		//		if (k != conn) {
		//			all_name.WriteString(v.nick)
		//		}
		//	}
		//	var send_str string
		//	if (strings.EqualFold(all_name.String(), "")) {
		//		send_str = fmt.Sprintf("%s(*)", all_user[conn].nick)
		//	} else {
		//		send_str = fmt.Sprintf("%s %s(*)", all_name.String(), all_user[conn].nick)
		//	}
		//	fmt.Println("content_list1:", content_list, "sendmsg:", send_str)
		//	_, err := conn.Write([]byte(send_str))
		//	checkError(err)
		//}
		////else {
		////	fmt.Println("content_list2:", content_list, "sendmsg:", string(content))
		////	_, err := conn.Write(content)
		////	checkError(err)
		////}

	}
}
