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
	"go-min-chat/const"
	"go-min-chat/Utils"
	"go-min-chat/ClientMsg"
	"go-min-chat/ClientUtil"
	"go-min-chat/ClientApp"
	"go-min-chat/Network"
)

func main() {
	cliSing := ClientApp.GetCli()
	flag.StringVar(&(cliSing.Nick), "u", "wang", "nick name")
	flag.StringVar(&(cliSing.Password), "p", "123456", "password")
	flag.StringVar(&(cliSing.Host), "h", "127.0.0.1", "host")
	flag.StringVar(&(cliSing.Port), "port", "8080", "port")
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cliSing.Host, cliSing.Port))

	// 啥也不干, 想登录判断
	Utils.CheckError(err)
	defer conn.Close()
	fmt.Print(ClientUtil.GetPre())
	var wg sync.WaitGroup
	ch := make(chan []byte)
	wg.Add(3)
	go readFromStdio(ch)
	go readFromConn(conn)
	go ClientMsg.Send(conn, ch)
	// 发送登录的消息
	conn.Write(ClientMsg.SendAuth(cliSing.Nick, cliSing.Password))
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
		if (!ClientApp.GetCli().IsAuth) { // 没有登录
			Utils.EchoLine("请先登录", 2)
			continue
		}
		data = []byte(strings.Trim(string(data), " "))
		data_str := string(data)
		if (strings.EqualFold(data_str, "")) { // 直接按的回车, 不做处理
			continue
		}
		param := strings.Split(data_str, " ")
		isContinue, p1 := ClientUtil.GetMsgType(param)
		if (isContinue) {
			continue
		}
		d, _ := proto.Marshal(p1)
		ch <- d
	}
}

func readFromConn(conn net.Conn) {
	msg := make(chan string, 10)
	for {
		select {
		case msg1 := <-msg: // 从msg chan里面取数据
			doMsg(msg1)
		default: // 读到数据放到msg chan里面
			buffer := Network.NewBuffer(conn, 1024)
			buffer.ThomasRead(msg)
		}
	}
}

func doMsg(buf string) {
	backContent := &protobuf.BackContent{}
	proto.Unmarshal([]byte(buf), backContent)
	switch backContent.Id {
	case _const.RCV_FAIL:
		Utils.EchoLine(backContent.Msg, 2)
		fmt.Print(ClientUtil.GetPre())
		break
	case _const.RCV_SUCCESS:
		Utils.EchoLine(backContent.Msg, 1)
		fmt.Print(ClientUtil.GetPre())
		break
	case _const.RCV_USE_ROOM:
		ClientMsg.UseRoom(backContent)
		fmt.Print(ClientUtil.GetPre())
		break
	case _const.RCV_AUTH:
		ClientMsg.Auth(backContent)
		fmt.Print(ClientUtil.GetPre())
		break
	case _const.RCV_SHOW_ROOMS:
		ClientMsg.ShowRoom(backContent)
		fmt.Print(ClientUtil.GetPre())
		break
	case _const.RCV_USER_LIST:
		ClientMsg.UserList(backContent)
		fmt.Print(ClientUtil.GetPre())
		break
	case _const.RCV_GROUP_MSG:
		ClientMsg.GroupMsg(backContent)
		fmt.Print(ClientUtil.GetPre())
		break
	}
}
