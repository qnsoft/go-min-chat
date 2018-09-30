package main

import (
	"net"
	"fmt"
	"os"
	"flag"
	"bytes"
	"strings"
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

var all_user = make(map[net.Conn]*Client)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "is port")
	flag.Parse()

	addr := fmt.Sprintf("%s%s", "192.168.101.200:", port)
	fmt.Println(addr)
	listen, err := net.Listen("tcp", addr)
	checkError(err)
	defer listen.Close()

	fmt.Println(all_user)
	for {
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
	for {
		n, err := conn.Read(buf)
		checkError(err)
		ch <- buf[:n]
		fmt.Println(string(buf[0:n]))
	}
}

func sendConnMsg(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		content_list := string(content)
		if (strings.HasPrefix(content_list, "auth")) {
			if (all_user[conn].auth) {
				conn.Write([]byte("do not auth, you already login"))
				continue
			} else {
				name := string(content_list[5:])
				all_user[conn].nick = name
				all_user[conn].wait_pass = true
				all_user[conn].auth = true
			}
		} else { // 如果不是auth指令, 就判断是否auth过
			if (all_user[conn].auth == false) {
				_, err := conn.Write([]byte("please auth"))
				checkError(err)
				continue
			}
		}
		if (content_list == "list") {
			var all_name bytes.Buffer
			for k, v := range all_user {
				if (k != conn) {
					all_name.WriteString(v.nick)
				}
			}
			fmt.Println("content_list1:", content_list, "sendmsg:", all_name.String())
			_, err := conn.Write([]byte(all_name.String()))
			checkError(err)
		} else {
			fmt.Println("content_list2:", content_list, "sendmsg:", string(content))
			_, err := conn.Write(content)
			checkError(err)
		}
		fmt.Println("-------")
	}
}
