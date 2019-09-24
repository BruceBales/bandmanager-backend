package dao

import (
	"fmt"

	"github.com/brucebales/bandmanager-backend/src/internal/config"

	"github.com/go-redis/redis"
)

func NewRedis() *redis.Client {
	conf := config.GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		Password: conf.RedisPass,
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Redis error: ", err)
	}

	return client

}
