package mysql

import (
	"net"
	"database/sql"
	"fmt"
)

type User struct {
	Uid      int
	Nick     string
	Password string
	Sex      string
	Age      int
	RoomId   int
	RoomName string
	Online   bool
	IsAuth   bool
	WaitPass bool
	Conn     net.Conn
}

func BuildUser(uid int, nick string, age int, isAuth bool) *User {
	var user = &User{}
	user.Uid = uid
	user.Nick = nick
	user.Age = age
	user.IsAuth = isAuth

	return user
}

func SelectUserById(nick string, DB *sql.DB) (User) {
	var user User
	rows, _ := DB.Query("select * from user where nick = ?", nick)
	for rows.Next() {
		rows.Scan(&user.Uid, &user.Nick, &user.Password, &user.Sex, &user.Age)
		fmt.Println("get data, id: ", user.Uid)
	}
	return user
}
