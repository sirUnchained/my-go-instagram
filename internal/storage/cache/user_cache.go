package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

const (
	EXP_TIME = time.Minute * 5
)

type userCache struct {
	client redis.Client
}

func (uc *userCache) Set(ctx context.Context, user *models.UserModel) error {
	userKey := fmt.Sprintf("user-%d", user.Id)
	userMarshal, err := json.Marshal(*user)
	if err != nil {
		return err
	}

	return uc.client.Set(ctx, userKey, userMarshal, EXP_TIME).Err()
}

func (uc *userCache) Get(ctx context.Context, userid int64) (*models.UserModel, error) {
	userKey := fmt.Sprintf("user-%d", userid)

	data, err := uc.client.Get(ctx, userKey).Result()
	if err != nil {
		return nil, err

	}

	if data == "" {
		return nil, nil
	}

	user := &models.UserModel{Role: &models.RoleModel{}, Profile: &models.ProfileModel{}}
	err = json.Unmarshal([]byte(data), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
