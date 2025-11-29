package db

import (
	"github.com/SorenHQ/joern-port/env"
	"github.com/go-redis/redis/v8"
)
var redisCli *redis.Client
func InitRedis() {
			redisCli = redis.NewClient(&redis.Options{
			Addr:     env.GetRedisURI(),
			Password: "",
			PoolSize: 10,
		})

	}


func GetRedisClient()*redis.Client{
	if redisCli==nil{
		InitRedis()
	}
	return redisCli
}