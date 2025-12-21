package storage

import (
	"context"
	"database/sql"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
)

type PgStorage struct {
	UserStore interface {
		Create(context.Context, *payloads.UserPayload) error
	}
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{
		UserStore: &userStore{db: db},
	}
}
