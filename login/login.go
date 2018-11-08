package login

import (
	"go-min-chat/mysql"
	"strings"
)

func CheckLog(user1 *mysql.User, user2 mysql.UserForLog) bool {
	//db 操作
	db := mysql.InitDB()
	//fmt.Println(user2.Nick, db)
	user12 := mysql.SelectUserById(user2.Nick, db)
	if (user12.Uid > 0 && strings.EqualFold(user2.Password, user12.Password)) { //  找到该用户 && 密码正确
		user1.Uid = user12.Uid
		user1.IsAuth = true
		user1.Age = user12.Age
		user1.Nick = user12.Nick

		return true
	} else {
		return false
	}
}
