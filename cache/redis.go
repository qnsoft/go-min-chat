package cache

import (
	"sync"
	"fmt"
	redis2 "gopkg.in/redis.v4"
	"strconv"
)

//Db数据库连接池
var redis *Redis
var once sync.Once

type Redis struct {
	Ip   string
	Port int
	Db   int
	Conn *redis2.Client
}

func GetReis() *Redis {
	once.Do(func() {
		redis = &Redis{}
	})
	return redis
}

func InitCache() {
	redis := GetReis()
	client := redis2.NewClient(&redis2.Options{
		Addr:     redis.Ip + ":" + strconv.Itoa(redis.Port),
		Password: "",
		DB:       0,
	})
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	redis.Conn = client
}
