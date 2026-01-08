package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

const (
	EXP_TIME = time.Minute * 5
)

type RedisStorage struct {
	UserCache interface {
		Set(context.Context, *models.UserModel) error
		Get(context.Context, int64) (*models.UserModel, error)
	}

	CommentCache interface {
		Set(context.Context, *models.CommentModel) error
		Get(context.Context, int64) (*models.CommentModel, error)
	}
}

func NewRedisStorage(client redis.Client) *RedisStorage {
	return &RedisStorage{
		UserCache:    &userCache{client: client},
		CommentCache: &commentCache{client: client},
	}
}
