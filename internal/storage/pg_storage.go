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
		UpdateData(context.Context, *models.UserModel, *payloads.CreateUserPayload) (*models.UserModel, error)
	}
	PostStore interface {
		Create(context.Context, *payloads.CreatePostPayload, *[]models.FileModel, *[]models.TagModel, *models.UserModel) (*models.PostModel, error)
		GetById(context.Context, int64) (*models.PostModel, error)
	}
	FileStore interface {
		Create(context.Context, int64, []payloads.CreateFilePayload) ([]models.FileModel, error)
	}
	TagStore interface {
		Create(context.Context, int64, []string) ([]models.TagModel, error)
	}
	BanStore interface {
		Create(context.Context, *models.UserModel, *payloads.CreateBanPayload) error
		Delete(context.Context, string) error
		GetBanByEmail(context.Context, string) (*models.BanModel, error)
	}
	CommentStore interface {
		Create(context.Context, int64, *payloads.CreateCommentPayload) error
		Delete(context.Context, int64) error
	}
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{
		UserStore:    &userStore{db: db},
		PostStore:    &postStore{db: db},
		FileStore:    &fileStore{db: db},
		TagStore:     &tagStore{db: db},
		BanStore:     &banStore{db: db},
		CommentStore: &commentStore{db: db},
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

func removeDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
