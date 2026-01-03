package storage

import (
	"context"
	"database/sql"
	"strings"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type commentStore struct {
	db *sql.DB
}

func (cs *commentStore) Create(ctx context.Context, userid int64, commentP *payloads.CreateCommentPayload) error {
	quety := `INSERT INTO comments (content, creator, post, parent) VALUES ($1, $2, $3, $4);`

	_, err := cs.db.ExecContext(ctx, quety, commentP.Content, commentP.CreatorID, commentP.PostID, commentP.ParentID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "comments_post_fkey"):
			return global_varables.NOT_FOUND_ROW
		default:
			return err
		}
	}

	return nil
}

func (cs *commentStore) GetPostComments(ctx context.Context, postid, limit, offset int64) ([]models.CommentModel, error) {
	postComments := []models.CommentModel{}
	query := `
	SELECT c.content, c.parent, c.created_at, u.id, u.username
	FROM comments AS c 
	JOIN users AS u ON c.creator = u.id
	WHERE c.post = $1
	LIMIT $2 OFFSET $3;
	`

	rows, err := cs.db.QueryContext(ctx, query, postid, limit, offset)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, global_varables.NOT_FOUND_ROW
		default:
			return nil, err
		}
	}

	for rows.Next() {
		postComment := &models.CommentModel{Creator: &models.UserModel{}}
		err := rows.Scan(&postComment.Content, &postComment.ParentID, &postComment.CreatedAt, &postComment.Creator.Id, &postComment.Creator.Username)
		if err != nil {
			return nil, err
		}

		postComment.PostID = postid
		postComments = append(postComments, *postComment)
	}

	return postComments, nil
}

func (cs *commentStore) GetRepliedComments(ctx context.Context, parrentid int64) ([]models.CommentModel, error) {
	return nil, nil
}

func (cs *commentStore) Delete(ctx context.Context, commentid int64) error {
	return nil
}
