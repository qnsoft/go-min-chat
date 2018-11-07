package mysql

import (
	"database/sql"
	"strings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//数据库配置
const (
	userName = "root"
	password = "WOAImama188"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "thomas"
)

//Db数据库连接池
var DB1 *sql.DB

//注意方法名大写，就是public
func InitDB() *sql.DB {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB1, _ := sql.Open("mysql", path)
	//设置数据库最大连接数
	DB1.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB1.SetMaxIdleConns(10)
	//验证连接
	if err := DB1.Ping(); err != nil {
		fmt.Println("opon database fail")
	}
	return DB1
}