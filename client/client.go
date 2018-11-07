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
	"go-min-chat/cli"
	"github.com/beego/bee/logger/colors"
	"go-min-chat/const"
	"go-min-chat/Util"
)

func main() {
	cliSing := cli.GetCli()
	flag.StringVar(&(cliSing.Nick), "u", "wang", "nick name")
	flag.StringVar(&(cliSing.Password), "p", "123456", "password")
	flag.StringVar(&(cliSing.Host), "h", "127.0.0.1", "host")
	flag.StringVar(&(cliSing.Port), "port", "8080", "port")
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cliSing.Host, cliSing.Port))

	// 啥也不干, 想登录判断
	Util.CheckError(err)
	defer conn.Close()
	fmt.Print(getPre())
	var wg sync.WaitGroup
	ch := make(chan []byte)
	wg.Add(3)
	go readFromStdio(ch)
	go readFromConn(conn)
	go sendMsg(conn, ch)
	// 发送登录的消息
	conn.Write(SendAuthMsg(cliSing.Nick, cliSing.Password))
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
		data_str := string(data)
		if (strings.EqualFold(data_str, "")) { // 直接按的回车, 不做处理
			continue
		}
		p1 := &protobuf.Content{}
		data_str_upper := strings.ToUpper(data_str)
		param := strings.Split(data_str, " ")

		if (strings.HasPrefix(data_str_upper, "SHOW ROOMS")) {
			p1.Id = _const.RCV_SHOW_ROOMS
		} else if (strings.HasPrefix(data_str_upper, "CREATE ROOM")) {
			if (len(param) < 3) {
				fmt.Println(fmt.Sprintf("(error) ERR unknown command '%s'", data_str_upper))
				fmt.Printf(getPre())
				continue
			}
			p1.Id = _const.RCV_CREATE_ROOM
			p1.ParamString = param[2]
		} else if (strings.HasPrefix(data_str_upper, "USER LIST")) {
			p1.Id = _const.RCV_USER_LIST
		} else if (strings.HasPrefix(data_str_upper, "USE")) {
			p1.Id = _const.RCV_USE_ROOM
			p1.ParamString = param[1]
		} else {
			cliSing := cli.GetCli()
			if (cliSing.RoomId != 0) { // 说明进入房间了
				p1.Id = _const.RCV_GROUP_MSG
				p1.ParamString = param[0]
			} else {
				fmt.Println(fmt.Sprintf("(error) ERR unknown command '%s'", data_str_upper))
				fmt.Printf(getPre())
				continue
			}
		}
		client := cli.GetCli()
		if (!client.IsAuth && p1.Id != _const.RCV_AUTH) {
			Util.EchoLine("请先登录", 2)
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
		Util.CheckError(err)
		backContent := &protobuf.BackContent{}
		proto.Unmarshal(buf[:n], backContent)
		switch backContent.Id {
		case _const.RCV_SUCCESS_FAIL:
			doSuccessFail(backContent)
			break
		case _const.RCV_USE_ROOM:
			useRoom(backContent)
			break
		case _const.RCV_AUTH:
			doAuth(backContent)
			break
		case _const.RCV_SHOW_ROOMS:
			doShowRoom(backContent)
			break
		case _const.RCV_USER_LIST:
			doUserList(backContent)
			break
		case _const.RCV_GROUP_MSG:
			doGroupMsg(backContent)
			break
		}
		fmt.Print(getPre())
	}
}

func doSuccessFail(backContent *protobuf.BackContent) {
	fmt.Println(colors.Green(backContent.Msg))
}

func doAuth(backContent *protobuf.BackContent) {
	cli1 := cli.GetCli()
	cli1.IsAuth = backContent.Auth.IsOk
	if (backContent.Auth.IsOk) { // 登录成功
		cli1.Nick = backContent.Auth.UseInfo.Nick
		cli1.Uid = backContent.Auth.UseInfo.Uid
		Util.EchoLine("OK", 1)
	} else { // 登录失败 or 需要去登录
		Util.EchoLine(backContent.Auth.Msg, 2)
		os.Exit(1)
	}
}

func useRoom(backContent *protobuf.BackContent) {
	cli1 := cli.GetCli()
	cli1.RoomId = int(backContent.Room.RoomId)
	cli1.RoomName = backContent.Room.RoomName
}

func doShowRoom(backContent *protobuf.BackContent) {
	fmt.Println(backContent.Showroom.RoomsAndIds)
}

func doUserList(backContent *protobuf.BackContent) {
	fmt.Println(backContent.Showroom.RoomsAndIds)
}

func doGroupMsg(backContent *protobuf.BackContent) {
	fmt.Printf("%s: %s\n", backContent.Groupmsg.Nick, backContent.Groupmsg.Content)
}

func sendMsg(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		_, err := conn.Write(content)
		Util.CheckError(err)
	}
}

func getPre() string {
	cliSing := cli.GetCli()
	var pre string
	//pre := fmt.Sprintf("%s:%s> ", cliSing.Host, cliSing.Port)
	if (cliSing.RoomId != 0) {
		room := colors.Red(fmt.Sprintf("%s", cliSing.RoomName))
		pre = fmt.Sprintf("%s%s 我: ", strings.TrimSuffix(pre, " "), room)
	}
	return pre
}

func SendAuthMsg(nick string, password string) []byte {
	p1 := &protobuf.Content{}
	p1.Id = _const.RCV_AUTH
	p1.Nick = nick
	p1.Password = password
	data, _ := proto.Marshal(p1)
	return data
}
