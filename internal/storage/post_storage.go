package storage

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type postStore struct {
	db *sql.DB
}

func (ps *postStore) Create(ctx context.Context, postP *payloads.CreatePostPayload, files *[]models.FileModel, tags *[]models.TagModel, user *models.UserModel) (*models.PostModel, error) {
	// saving post
	postQuery := `INSERT INTO posts (description, creator) VALUES ($1, $2) RETURNING id, created_at;`
	post := &models.PostModel{
		Description: postP.Description,
		Creator:     *user,
	}
	err := ps.db.QueryRowContext(ctx, postQuery, postP.Description, postP.Creator).Scan(&post.Id, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	// saving post_files relation
	n := len(*files)
	filesId := make([]int64, n)
	for i, v := range *files {
		filesId[i] = v.Id
	}
	postFilesQuery := `INSERT INTO posts_files (post, file) SELECT $1, unnest($2::bigint[]);`
	if err := ps.db.QueryRowContext(ctx, postFilesQuery, post.Id, pq.Array(filesId)).Err(); err != nil {
		return nil, err
	}

	// saving post_tags relation
	if len(*tags) > 0 {
		m := len(*tags)
		tagsId := make([]int64, m)
		for i, v := range *tags {
			tagsId[i] = v.Id
		}
		tagsId = RemoveDuplicates(tagsId)
		postTagsQuery := `INSERT INTO posts_tags (post, tag) SELECT $1, unnest($2::bigint[]);`
		if err := ps.db.QueryRowContext(ctx, postTagsQuery, post.Id, pq.Array(tagsId)).Err(); err != nil {
			return nil, err
		}
	}
	return post, nil
}

func RemoveDuplicates[T comparable](slice []T) []T {
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
