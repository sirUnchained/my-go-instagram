package storage

import (
	"context"
	"database/sql"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type PgStorage struct {
	UserStore interface {
		Create(context.Context, *payloads.UserPayload) (*models.UserModel, error)
		GetById(context.Context, int64) (*models.UserModel, error)
		GetByEmail(context.Context, string) (*models.UserModel, error)
	}
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{
		UserStore: &userStore{db: db},
	}
}
