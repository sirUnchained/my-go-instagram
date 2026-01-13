package storage

import (
	"context"
	"database/sql"
	"strings"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
)

type likeStore struct {
	db *sql.DB
}

func (ls *likeStore) Create(ctx context.Context, userid, postid int64) error {
	query := `INSERT INTO likes (post, creator) VALUES ($1, $2);`

	if _, err := ls.db.ExecContext(ctx, query, postid, userid); err != nil {
		switch {
		case strings.Contains(err.Error(), "unique_post_like"):
			return global_varables.DUP_ITEM
		default:
			return err
		}
	}

	return nil
}

func (ls *likeStore) Delete(ctx context.Context, userid, postid int64) error {
	query := `DELETE FROM likes WHERE post = $1 AND creator = $2;`

	if _, err := ls.db.ExecContext(ctx, query, postid, userid); err != nil {
		return err
	}

	return nil
}
