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
	if n > 0 {
		filesId := make([]int64, n)
		for i, v := range *files {
			filesId[i] = v.Id
		}
		postFilesQuery := `INSERT INTO posts_files (post, file) SELECT $1, unnest($2::bigint[]);`
		if err := ps.db.QueryRowContext(ctx, postFilesQuery, post.Id, pq.Array(filesId)).Err(); err != nil {
			return nil, err
		}
	}

	// saving post_tags relation
	m := len(*tags)
	if m > 0 {
		tagsId := make([]int64, m)
		for i, v := range *tags {
			tagsId[i] = v.Id
		}
		tagsId = removeDuplicates(tagsId)
		postTagsQuery := `INSERT INTO posts_tags (post, tag) SELECT $1, unnest($2::bigint[]);`
		if err := ps.db.QueryRowContext(ctx, postTagsQuery, post.Id, pq.Array(tagsId)).Err(); err != nil {
			return nil, err
		}
	}

	return post, nil
}

func (ps *postStore) GetById(ctx context.Context, postid int64) (*models.PostModel, error) {
	// get post
	queryPost := `
	SELECT p.id, p.description, p.created_at, p.updated_at, u.id, u.username
		FROM posts 		AS p 
		JOIN users 		AS u 	ON u.id 	= p.creator
		JOIN profiles   AS up 	ON up.id 	= u.profile
	WHERE p.id = $1;
	`
	post := &models.PostModel{Creator: models.UserModel{}}
	err := ps.db.QueryRowContext(ctx, queryPost, postid).Scan(
		&post.Id,
		&post.Description,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Creator.Id,
		&post.Creator.Username,
	)
	if err != nil {
		return nil, err
	}

	// get post files
	queryFiles := `
		SELECT f.id, f.filename, f.filepath
			FROM posts_files as pf 
			JOIN files as f ON f.id = pf.file
		WHERE pf.post = $1
	`
	rows, err := ps.db.QueryContext(ctx, queryFiles, postid)
	if err != nil {
		return nil, err
	}
	files := []models.FileModel{}
	for rows.Next() {
		var newFile models.FileModel
		err := rows.Scan(&newFile.Id, &newFile.Filename, &newFile.Filepath)
		if err != nil {
			return nil, err
		}
		files = append(files, newFile)
	}
	post.Files = files

	// get post tags
	queryTags := `
		SELECT t.id, t.name
			FROM posts_tags as pt 
			JOIN tags as t ON t.id = pt.tag
		WHERE pt.post = $1
	`
	rows, err = ps.db.QueryContext(ctx, queryTags, postid)
	if err != nil {
		return nil, err
	}
	tags := []models.TagModel{}
	for rows.Next() {
		var newTag models.TagModel
		err := rows.Scan(&newTag.Id, &newTag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, newTag)
	}
	post.Tags = tags

	return post, nil
}
