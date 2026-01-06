package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type userCache struct {
	client redis.Client
}

func (uc *userCache) Set(ctx context.Context, user *models.UserModel) error {
	userKey := fmt.Sprintf("user-%d", user.Id)
	userMarshal, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return uc.client.Set(ctx, userKey, userMarshal, time.Minute*5).Err()
}

func (uc *userCache) Get(ctx context.Context, userid int64) (*models.UserModel, error) {
	userKey := fmt.Sprintf("user-%d", userid)

	data, err := uc.client.Get(ctx, userKey).Result()
	if err != nil {
		switch {
		case err == redis.Nil:
			return nil, nil
		default:
			return nil, err
		}
	}

	user := &models.UserModel{}
	if data != "" {
		err := json.Unmarshal([]byte(data), user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
