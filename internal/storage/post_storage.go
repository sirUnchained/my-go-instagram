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
	SELECT p.id, p.description, p.created_at, p.updated_at, u.id, u.username, u.is_private
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
		&post.Creator.IsPrivate,
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
	defer rows.Close()
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
	defer rows.Close()
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

func (ps *postStore) GetFeed(ctx context.Context, limit, offset, userid int64) ([]models.PostModel, error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// we first get posts
	postQuery := `
	SELECT p.id, p.description, p.created_at, u.id, u.username
	FROM posts AS p
	JOIN users AS u ON p.creator = u.id
	WHERE (u.is_private = FALSE) AND (u.id != $1)
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3;
	`
	posts := []models.PostModel{}
	rows, err := tx.QueryContext(ctx, postQuery, userid, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postIDs := []int64{}
	postMap := make(map[int64]*models.PostModel)

	for rows.Next() {
		post := models.PostModel{Creator: models.UserModel{}}
		err := rows.Scan(
			&post.Id,
			&post.Description,
			&post.CreatedAt,
			&post.Creator.Id,
			&post.Creator.Username,
		)
		if err != nil {
			return nil, err
		}
		postIDs = append(postIDs, post.Id)
		postMap[post.Id] = &post
		posts = append(posts, post)
	}

	// if no posts found, commit and return
	if len(posts) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return posts, nil
	}

	// now fetch files for posts using the junction table
	filesQuery := `
	SELECT pf.post, f.id, f.filename, f.filepath, f.size_bytes
	FROM  posts_files AS pf
	JOIN  files 	AS f ON pf.file = f.id
	WHERE pf.post = ANY($1)
	ORDER BY pf.post, pf.file;
	`
	fileRows, err := tx.QueryContext(ctx, filesQuery, pq.Array(postIDs))
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var postID int64
		file := models.FileModel{}

		err := fileRows.Scan(
			&postID,
			&file.Id,
			&file.Filename,
			&file.Filepath,
			&file.SizeBytes,
		)
		if err != nil {
			return nil, err
		}

		// Add file to the corresponding post
		if post, exists := postMap[postID]; exists {
			post.Files = append(post.Files, file)
		}
	}

	// finally we fetch tags
	tagsQuery := `
	SELECT pt.post, t.id, t.name
	FROM posts_tags AS pt
	JOIN tags AS t ON pt.tag = t.id
	WHERE pt.post = ANY($1)
	ORDER BY pt.post;
	`
	tagRows, err := tx.QueryContext(ctx, tagsQuery, pq.Array(postIDs))
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var postID int64
		tag := models.TagModel{}

		err := tagRows.Scan(
			&postID,
			&tag.Id,
			&tag.Name,
		)
		if err != nil {
			return nil, err
		}

		// add tag to the corresponding post
		if post, exists := postMap[postID]; exists {
			post.Tags = append(post.Tags, tag)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return posts, nil
}
