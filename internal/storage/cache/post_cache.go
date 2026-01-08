package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type postCache struct {
	client redis.Client
}

func (uc *postCache) Set(ctx context.Context, post *models.PostModel) error {
	postKey := fmt.Sprintf("post-%d", post.Id)
	postMarshal, err := json.Marshal(*post)
	if err != nil {
		return err
	}

	return uc.client.Set(ctx, postKey, postMarshal, EXP_TIME).Err()
}

func (uc *postCache) Get(ctx context.Context, postid int64) (*models.PostModel, error) {
	postKey := fmt.Sprintf("user-%d", postid)

	data, err := uc.client.Get(ctx, postKey).Result()
	if err != nil {
		return nil, err

	}

	if data == "" {
		return nil, nil
	}

	post := &models.PostModel{}
	err = json.Unmarshal([]byte(data), post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
