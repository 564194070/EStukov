package redis

import (
	"github.com/go-redis/redis"
	"time"
)

var (
	pool *redis.Client
	redisHost = "119.3.232.83:6379:6379"
	redisPass = "mohican123."
)

func InitRedisClient() error {
	redisDB := redis.NewClient(&redis.Options{
		Addr: redisHost,
		Password: redisPass,
		DB: 0,
		PoolSize: 50,
		MaxConnAge: 300*time.Second,
	})

	_, err := redisDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func RedisPool() *redis.Client {
	return pool
}