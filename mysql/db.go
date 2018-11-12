package mysql

import (
	"database/sql"
	"strings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"strconv"
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
var DB1 *DB
var once sync.Once

var path = strings.Join([]string{userName, ":",
	password, "@tcp(", ip, ":", port, ")/",
	dbName, "?charset=utf8"}, "")

type DB struct {
	Ip           string
	Port         int
	Username     string
	Password     string
	DbName       string
	Charset      string
	MaxLifeTime  int
	MaxIdleConns int
}

func GetDB() *DB {
	once.Do(func() {
		DB1 = &DB{}
	})
	return DB1
}

//注意方法名大写，就是public
func InitDB() *sql.DB {
	DB := GetDB()
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{DB.Username, ":", DB.Password, "@tcp(", DB.Ip, ":", strconv.Itoa(DB.Port), ")/", DB.DbName, "?charset=", DB.Charset}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB1, _ := sql.Open("mysql", path)
	//fmt.Println(path)
	//os.Exit(1)
	//设置数据库最大连接数

	DB1.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB1.SetMaxIdleConns(DB.MaxIdleConns)
	//验证连接
	if err := DB1.Ping(); err != nil {
		fmt.Println("opon database fail")
	}
	return DB1
}
