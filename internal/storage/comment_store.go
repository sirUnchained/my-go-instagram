package storage

import (
	"context"
	"database/sql"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type commentStore struct {
	db *sql.DB
}

func (cs *commentStore) Create(ctx context.Context, userid int64, commentP *payloads.CreateCommentPayload) error {
	return nil
}

func (cs *commentStore) GetPostComments(ctx context.Context, postid int64) ([]models.CommentModel, error) {
	return nil, nil
}

func (cs *commentStore) GetRepliedComments(ctx context.Context, parrentid int64) ([]models.CommentModel, error) {
	return nil, nil
}

func (cs *commentStore) Delete(ctx context.Context, commentid int64) error {
	return nil
}
