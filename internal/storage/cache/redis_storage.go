package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type RedisStorage struct {
	UserCache interface {
		Set(context.Context, *models.UserModel) error
		Get(context.Context, int64) (*models.UserModel, error)
	}

	CommentCache interface {
		Set(context.Context, *models.UserModel) error
		Get(context.Context, int64) (*models.UserModel, error)
	}
}

func NewRedisStorage(client redis.Client) *RedisStorage {
	return &RedisStorage{
		UserCache: &userCache{client: client},
	}
}
