package storage

import (
	"context"
	"database/sql"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type postStore struct {
	db *sql.DB
}

func (ps *postStore) Create(ctx context.Context, postP *payloads.CreatePostPayload, user *models.UserModel) (*models.PostModel, error) {
	query := `INSERT INTO posts (description, creator) VALUES ($1, $2) RETURNING id, created_at;`

	post := &models.PostModel{
		Description: postP.Description,
		Creator:     *user,
	}
	err := ps.db.QueryRowContext(ctx, query, postP.Description, postP.Creator).Scan(post.Id, post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return post, nil
}
