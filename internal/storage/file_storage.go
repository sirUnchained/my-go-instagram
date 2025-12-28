package storage

import (
	"context"
	"database/sql"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type FileStore struct {
	db *sql.DB
}

func (fs *FileStore) Create(ctx context.Context, userid int64, files []payloads.CreateFilePayload) {
	query := `INSERT INTO files (filename, filepath, size_bytes, creator) VALUES ($1, $2, $3, $4) RETURNING id;`

	var newFiles []models.FileModel

	fs.db.QueryRowContext(ctx, query, files).Scan()

}
