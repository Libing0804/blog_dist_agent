package daors

import (
	"github.com/go-redis/redis"
)

// 定义一个全局变量
var DBrs *redis.Client

func InitRedis() (err error) {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "47.100.99.174:6379", // 指定
		Password: "",
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	_, err = redisdb.Ping().Result()
	DBrs = redisdb
	return
}
