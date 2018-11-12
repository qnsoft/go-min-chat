package cache

import (
	"sync"
	"fmt"
	redis2 "github.com/garyburd/redigo/redis"
	"strconv"
)

//Db数据库连接池
var redis *Redis
var once sync.Once

type Redis struct {
	Ip   string
	Port int
	Conn redis2.Conn
}

func GetReis() *Redis {
	once.Do(func() {
		redis = &Redis{}
	})
	return redis
}

func InitCache() {
	redis := GetReis()
	c, err := redis2.Dial("tcp", redis.Ip+":"+strconv.Itoa(redis.Port))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	redis.Conn = c
}

func (redis Redis) Get(key string) interface{} {
	v, err := redis.Conn.Do("GET", key)
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get willen_key: %v \n", v)
	}
	return v
}

func (redis Redis) Set(key string, value string) {
	v, err := redis.Conn.Do("SET", key, value)
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get willen_key: %v \n", v)
	}
}

func (redis Redis) Sadd(key string, value string) {
	_, err := redis.Conn.Do("sadd", "key", value)
	if err != nil {
		fmt.Println("set add failed", err.Error())
	}
}

func (redis Redis) Sismember(key string, value string) {
	isMember, err := redis.Conn.Do("sismember", key, value)
	if err != nil {
		fmt.Println("sismember get failed", err.Error())
	} else {
		fmt.Println("foo is or not myset's member", isMember)
	}
}
