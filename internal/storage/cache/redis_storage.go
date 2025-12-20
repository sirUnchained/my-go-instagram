package cache

import "github.com/redis/go-redis/v9"

type RedisStorage struct {
}

func NewRedisStorage(client redis.Client) *RedisStorage {
	return &RedisStorage{}
}
