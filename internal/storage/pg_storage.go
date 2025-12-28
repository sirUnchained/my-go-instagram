package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	FileStore interface {
		Create(context.Context, int64, []payloads.CreateFilePayload) ([]models.FileModel, error)
	}
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{
		UserStore: &userStore{db: db},
		PostStore: &postStore{db: db},
		FileStore: &fileStore{db: db},
	}
}

func executeTransaction(ctx context.Context, db *sql.DB, fnc func(ctx context.Context, tx *sql.Tx) error) error {
	if fnc == nil {
		return errors.New("transaction function cannot be nil")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fnc(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			return fmt.Errorf("original error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
