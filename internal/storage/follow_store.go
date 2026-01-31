package storage

import (
	"database/sql"

	"context"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type followStore struct {
	db *sql.DB
}

func (ls *followStore) GetFollowers(ctx context.Context, userid, limit, offset int64) ([]models.UserModel, error) {
	query := `
	SELECT u.id, u.username, p.fullname
	FROM follows 	AS f
	JOIN users 		AS u ON f.following = u.id
	JOIN profiles 	AS p ON p.id 		= u.profile
	WHERE f.following = $1
	LIMIT $2 OFFSET $3;
	`

	users := []models.UserModel{}

	rows, err := ls.db.QueryContext(ctx, query, userid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, global_varables.NOT_FOUND_ROW
		default:
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.UserModel{Profile: &models.ProfileModel{}}
		err := rows.Scan(&user.Id, &user.Username, &user.Profile.Fullname)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil

}

func (ls *followStore) GetFollowings(ctx context.Context, userid, limit, offset int64) ([]models.UserModel, error) {
	query := `
	SELECT u.id, u.username, p.fullname
	FROM follows 	AS f
	JOIN users 		AS u ON f.follower = u.id
	JOIN profiles 	AS p ON p.id = u.profile
	WHERE f.follower = $1
	LIMIT $2 OFFSET $3;
	`

	users := []models.UserModel{}

	rows, err := ls.db.QueryContext(ctx, query, userid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, global_varables.NOT_FOUND_ROW
		default:
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.UserModel{Profile: &models.ProfileModel{}}
		err := rows.Scan(&user.Id, &user.Username, &user.Profile.Fullname)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil

}

func (ls *followStore) Create(ctx context.Context, currentUser, targetUser int64) error {
	query := `INSERT INTO follows (follower, following) VALUES ($1, $2);`

	if _, err := ls.db.ExecContext(ctx, query, currentUser, targetUser); err != nil {

		return err

	}

	return nil
}

func (ls *followStore) Delete(ctx context.Context, currentUser, targetUser int64) error {
	query := `DELETE FROM follows WHERE follower = $1 AND following = $2;`

	if _, err := ls.db.ExecContext(ctx, query, currentUser, targetUser); err != nil {
		return err
	}

	return nil
}
