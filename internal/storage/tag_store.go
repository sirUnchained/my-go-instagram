package storage

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type tagStore struct {
	db *sql.DB
}

func (fs *tagStore) Create(ctx context.Context, userid int64, tags []string) ([]models.TagModel, error) {
	if len(tags) <= 0 {
		return []models.TagModel{}, nil
	}

	query := `
	    WITH inserted_tags AS (
            INSERT INTO tags (name)
            SELECT unnest($1::text[])
            ON CONFLICT (name) DO NOTHING
            RETURNING id, name, created_at
        )
        SELECT id, name, created_at FROM inserted_tags
        UNION ALL
        SELECT id, name, created_at FROM tags
        WHERE name = ANY($1::text[])
        AND NOT EXISTS (SELECT 1 FROM inserted_tags);
	`
	rows, err := fs.db.QueryContext(ctx, query,
		pq.Array(tags),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newTags []models.TagModel
	for rows.Next() {
		var newTag models.TagModel
		err := rows.Scan(
			&newTag.Id,
			&newTag.Name,
			&newTag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		newTags = append(newTags, newTag)
	}

	return newTags, nil

}
