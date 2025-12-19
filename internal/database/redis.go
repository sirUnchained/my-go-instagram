package database

import "github.com/redis/go-redis/v9"

func NewRedisClient(addr, password string, dbNumber int) redis.Client {
	database := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       dbNumber,
		Password: password,
	})
	return *database
}
