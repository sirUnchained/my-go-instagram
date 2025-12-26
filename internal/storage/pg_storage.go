package storage

import (
	"context"
	"database/sql"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type PgStorage struct {
	UserStore interface {
		Create(context.Context, *payloads.CreateUserPayload) (*models.UserModel, error)
		GetById(context.Context, int64) (*models.UserModel, error)
		GetByEmail(context.Context, string) (*models.UserModel, error)
	}
	PostStore interface {
		Create(context.Context, *payloads.CreatePostPayload, *models.UserModel) (*models.PostModel, error)
	}
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{
		UserStore: &userStore{db: db},
		PostStore: &postStore{db: db},
	}
}
