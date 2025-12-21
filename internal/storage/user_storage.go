package storage

import (
	"context"
	"database/sql"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
)

type userStore struct {
	db *sql.DB
}

func (us *userStore) Create(ctx context.Context, userP *payloads.UserPayload) error {
	return nil
}
