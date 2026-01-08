package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type commentCache struct {
	client redis.Client
}

func (uc *commentCache) Set(ctx context.Context, comment *models.CommentModel) error {
	commentKey := fmt.Sprintf("comment-%d", comment.ID)
	commentMarshal, err := json.Marshal(*comment)
	if err != nil {
		return err
	}

	return uc.client.Set(ctx, commentKey, commentMarshal, EXP_TIME).Err()
}

func (uc *commentCache) Get(ctx context.Context, commentid int64) (*models.CommentModel, error) {
	commentKey := fmt.Sprintf("user-%d", commentid)

	data, err := uc.client.Get(ctx, commentKey).Result()
	if err != nil {
		return nil, err

	}

	if data == "" {
		return nil, nil
	}

	comment := &models.CommentModel{}
	err = json.Unmarshal([]byte(data), comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
