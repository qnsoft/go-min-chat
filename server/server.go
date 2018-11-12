package main

import (
	"net"
	"fmt"
	"os"
	"go-min-chat/server/ser"
	"go-min-chat/mysql"
	"go-min-chat/Utils"
	"go-min-chat/Msg"
	"flag"
	"github.com/beego/bee/logger/colors"
	"go-min-chat/cache"
)

func init() {
	MinChatSer := ser.GetMinChatSer()
	DB := mysql.GetDB()
	Redis := cache.GetReis()
	ini_parser := Utils.IniParser{}
	conf_file_name := "conf.ini"
	if err := ini_parser.Load("../conf/conf.ini"); err != nil {
		fmt.Printf("try load config file[%s] error[%s]\n", conf_file_name, err.Error())
		return
	}

	// mysql
	DB.Ip = ini_parser.GetString("mysql", "Ip")
	DB.Username = ini_parser.GetString("mysql", "Username")
	DB.Password = ini_parser.GetString("mysql", "Password")
	DB.DbName = ini_parser.GetString("mysql", "DbName")
	DB.Charset = ini_parser.GetString("mysql", "Charset")
	DB.Port = int(ini_parser.GetInt32("mysql", "Port"))
	DB.MaxLifeTime = int(ini_parser.GetInt32("mysql", "MaxLifeTime"))
	DB.MaxIdleConns = int(ini_parser.GetInt32("mysql", "MaxIdleConns"))
	mysql.InitDB()

	// redis
	Redis.Ip = ini_parser.GetString("redis", "Ip")
	Redis.Port = int(ini_parser.GetInt32("redis", "Port"))
	cache.InitCache()
	//Redis.Set("name", "wang")
	//Redis.Sadd("room1", "wang")
	//Redis.Sismember("room1", "wang")

	// server
	MinChatSer.Host = ini_parser.GetString("test", "Ip")
	MinChatSer.Port = int(ini_parser.GetInt32("test", "Port"))
	flag.StringVar(&MinChatSer.Host, "h", MinChatSer.Host, "is port")
	flag.IntVar(&MinChatSer.Port, "p", MinChatSer.Port, "is port")
	flag.Parse()
}

const log = `                               _                    _             _
  __ _   ___        _ __ ___  (_) _ __         ___ | |__    __ _ | |_
 / _\ | / _ \  ___ | '_ \ _ \ | || '_ \  ___  / __|| '_ \  / _\ || __|
| (_| || (_) ||___|| | | | | || || | | ||___|| (__ | | | || (_| || |_
 \__, | \___/      |_| |_| |_||_||_| |_|      \___||_| |_| \__,_| \__|
 |___/
一个分布式聊天系统
`

func main() {
	MinChatSer := ser.GetMinChatSer()
	addr := fmt.Sprintf("%s:%d", MinChatSer.Host, MinChatSer.Port)
	listen, err := net.Listen("tcp", addr)
	Utils.CheckError(err)
	defer listen.Close()
	fmt.Println(colors.Red(log))
	var u *mysql.User
	for {
		newConn, err := listen.Accept()
		// 连接上了，就把这个连接赋予一个未登录的用户
		u = mysql.BuildUser(0, "", 0, false)
		MinChatSer.AllUser[newConn] = u
		fmt.Println(newConn.RemoteAddr())
		Utils.CheckError(err)
		ch := make(chan []byte)
		go recvConnMsg(newConn, ch)
		go sendConnMsg(newConn, ch)
	}
}

// 服务端接受消息
func recvConnMsg(conn net.Conn, ch chan []byte) {
	buf := make([]byte, 50)
	for {
		n, err := conn.Read(buf)
		ret := Utils.ConnReadCheckError(err, conn)
		if (ret == 0) { // 读取时, 发生了错误
			os.Exit(1)
		} else if (ret == -1) { // 客户端断开了连接
			continue
		}
		ch <- buf[:n]
	}
}

// 服务端发送消息
func sendConnMsg(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		Msg.DoAllMsg(conn, content)
	}
}
