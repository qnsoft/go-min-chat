package main

import (
	"net"
	"fmt"
	"os"
	"sync"
	"bufio"
	"flag"
	"go-min-chat/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"strings"
	"go-min-chat/msg"
	"go-min-chat/cli"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	cliSing := cli.GetCli()
	flag.StringVar(&(cliSing.Port), "p", "8080", "port")
	flag.StringVar(&(cliSing.Host), "h", "127.0.0.1", "host")
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cliSing.Host, cliSing.Port))
	checkError(err)
	defer conn.Close()
	fmt.Print(getPre())
	var wg sync.WaitGroup
	ch := make(chan []byte)
	wg.Add(3)
	go readFromStdio(ch)
	go readFromConn(conn)
	go sendMsg(conn, ch)
	//go heartBeat(conn)
	wg.Wait()
}

//func heartBeat(conn net.Conn) {
//	time1 := time2.NewTicker(5 * time2.Second)
//	var content string
//	for {
//		select {
//		case <-time1.C:
//			content = "iamok"
//			_, err := conn.Write([]byte(content))
//			checkError(err)
//			fmt.Print(pre)
//		}
//	}
//}

func readFromStdio(ch chan []byte) {
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		data = []byte(strings.Trim(string(data), " "))
		p1 := &protobuf.Content{}
		data_str := string(data)
		data_str_upper := strings.ToUpper(data_str)
		param := strings.Split(data_str, " ")
		if (strings.HasPrefix(data_str_upper, "SHOW ROOMS")) {
			p1.Id = msg.RCV_SHOWROOMS
		} else if (strings.HasPrefix(data_str_upper, "CREATE ROOM")) {
			if (len(param) < 3) {
				fmt.Println(fmt.Sprintf("(error) ERR unknown command '%s'", data_str_upper))
				fmt.Printf(getPre())
				continue
			}
			p1.Id = msg.RCV_CREATEROOM
			p1.Param = param[2]
		} else if (strings.HasPrefix(data_str_upper, "USE")) {
			p1.Id = msg.RCV_USEROOM
			p1.Param = param[1]
		} else {
			fmt.Println(fmt.Sprintf("(error) ERR unknown command '%s'", data_str_upper))
			fmt.Printf(getPre())
			continue
		}
		d, _ := proto.Marshal(p1)
		ch <- d
	}
}

func readFromConn(conn net.Conn) {
	for {
		buf := make([]byte, 500)
		n, err := conn.Read(buf)
		checkError(err)
		backContent := &protobuf.BackContent{}
		proto.Unmarshal(buf[:n], backContent)
		if (backContent.Id == msg.SEND_USEROOM) {
			param_arr := strings.Split(backContent.Param, " ")
			cli := cli.GetCli()
			cli.RoomId, err = strconv.Atoi(param_arr[0])
			cli.RoomName = param_arr[1]
		}
		fmt.Println(backContent.Param)
		fmt.Print(getPre())
	}
}

func sendMsg(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		_, err := conn.Write(content)
		checkError(err)
	}
}

func getPre() string {
	cliSing := cli.GetCli()
	pre := fmt.Sprintf("%s:%s> ", cliSing.Host, cliSing.Port)
	if (cliSing.RoomId != 0) {
		pre = fmt.Sprintf("%s[%s(%d)] ", pre, cliSing.RoomName, cliSing.RoomId)
	}
	return pre
}
