package main

import (
	"net"
	"fmt"
	"os"
	"sync"
	"bufio"
	"flag"
	time2 "time"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"strings"
	"go-min-chat/msg"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

var pre string

func main() {
	var port string
	var host string
	flag.StringVar(&port, "p", "8080", "port")
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	checkError(err)
	defer conn.Close()
	pre = fmt.Sprintf("%s:%s> ", host, port)
	fmt.Print(pre)
	var wg sync.WaitGroup
	ch := make(chan []byte)
	wg.Add(3)
	go readFromStdio(ch)
	go readFromConn(conn)
	go sendMsg(conn, ch)
	//go heartBeat(conn)
	wg.Wait()
}

func heartBeat(conn net.Conn) {
	time1 := time2.NewTicker(5 * time2.Second)
	var content string
	for {
		select {
		case <-time1.C:
			content = "iamok"
			_, err := conn.Write([]byte(content))
			checkError(err)
			fmt.Print(pre)
		}
	}
}

func readFromStdio(ch chan []byte) {
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		data = []byte(strings.Trim(string(data), " "))
		p1 := &protobuf.Content{}
		data_str := string(data)
		data_str_upper := strings.ToUpper(data_str)
		if (strings.HasPrefix(data_str_upper, "SHOW ROOMS")) {
			p1.Id = msg.SHOWROOMS
		} else if (strings.HasPrefix(data_str_upper, "CREATE ROOM")) {
			param := strings.Split(data_str, " ")
			p1.Id = msg.CREATEROOM
			p1.Param = param[2]
		}
		d, _ := proto.Marshal(p1)
		ch <- d
	}
}

func readFromConn(conn net.Conn) {
	for {
		buf := make([]byte, 50)
		n, err := conn.Read(buf)
		checkError(err)
		fmt.Print(string(buf[:n]))
		fmt.Println()
		fmt.Print(pre)
	}
}

func sendMsg(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		_, err := conn.Write(content)
		checkError(err)
	}
}
