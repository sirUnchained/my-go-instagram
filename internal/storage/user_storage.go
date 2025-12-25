package storage

import (
	"context"
	"database/sql"
	"fmt"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type userStore struct {
	db *sql.DB
}

func (us *userStore) Create(ctx context.Context, userP *payloads.CreateUserPayload) (*models.UserModel, error) {
	query := `INSERT INTO users (username, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;`

	user := models.UserModel{Username: userP.Username, Fullname: userP.Fullname, Email: userP.Email, Password: models.Password{}, Role: models.RoleModel{}}
	user.Password.Set(userP.Password)

	var userCount int
	err := us.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	if userCount == 0 {
		user.Role.Id = 2
		user.IsVerified = true
	} else {
		user.Role.Id = 1
		user.IsVerified = false
	}

	err = us.db.QueryRowContext(ctx, query,
		user.Username,
		user.Fullname,
		user.Email,
		user.Password.Hash,
		user.Role.Id,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"":
			return nil, global_varables.USERNAME_DUP

		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return nil, global_varables.EMAIL_DUP

		default:
			return nil, err
		}
	}

	return &user, nil
}

func (us *userStore) GetById(ctx context.Context, userId int64) (*models.UserModel, error) {
	query := `SELECT id, username, fullname, email, password, is_verified, role, created_at, updated_at from users WHERE id = $1`

	user := &models.UserModel{}
	err := us.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.Password.Hash,
		&user.IsVerified,
		&user.Role.Id,
		&user.CreatedAt,
		&user.UreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, global_varables.NOT_FOUND_ROW
		default:
			return nil, err
		}
	}

	return user, nil

}

func (us *userStore) GetByEmail(ctx context.Context, email string) (*models.UserModel, error) {
	query := `SELECT id, username, fullname, email, password, is_verified, role, created_at, updated_at from users WHERE email = $1`

	user := &models.UserModel{}
	err := us.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.Password.Hash,
		&user.IsVerified,
		&user.Role.Id,
		&user.CreatedAt,
		&user.UreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == "sql: no rows in result set":
			return nil, global_varables.NOT_FOUND_ROW
		default:
			return nil, err
		}
	}

	return user, nil
}
