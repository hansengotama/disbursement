package redishelper

import (
	"github.com/hansengotama/disbursement/internal/lib/env"
	"github.com/redis/go-redis/v9"
)

var redisConn *redis.Client

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     env.GetRedisAddress(),
		Username: env.GetRedisUser(),
		Password: env.GetRedisPassword(),
		DB:       0,
	})

	redisConn = client
}

func GetConnection() *redis.Client {
	return redisConn
}

func CloseConnection() {
	err := redisConn.Close()
	if err != nil {
		// logging
		return
	}
}
