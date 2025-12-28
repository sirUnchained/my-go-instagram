package storage

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type fileStore struct {
	db *sql.DB
}

func (fs *fileStore) Create(ctx context.Context, userid int64, files []payloads.CreateFilePayload) ([]models.FileModel, error) {
	n := len(files)
	filenames := make([]string, n)
	filepaths := make([]string, n)
	sizes := make([]int, n)
	creators := make([]int64, n)

	for i, file := range files {
		filenames[i] = file.Filename
		filepaths[i] = file.Filepath
		sizes[i] = file.SizeBytes
		creators[i] = userid
	}

	query := `
				INSERT INTO files (filename, filepath, size_bytes, creator) 
					SELECT unnest($1::text[]), unnest($2::text[]), unnest($3::bigint[]), unnest($4::bigint[]) 
					RETURNING id, filename, filepath, size_bytes, creator, created_at;
			`

	rows, err := fs.db.QueryContext(ctx, query,
		pq.Array(filenames),
		pq.Array(filepaths),
		pq.Array(sizes),
		pq.Array(creators))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newFiles []models.FileModel
	for rows.Next() {
		var newFile models.FileModel
		err := rows.Scan(
			&newFile.Id,
			&newFile.Filename,
			&newFile.Filepath,
			&newFile.SizeBytes,
			&newFile.Creator,
			&newFile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		newFiles = append(newFiles, newFile)
	}

	return newFiles, nil

}
