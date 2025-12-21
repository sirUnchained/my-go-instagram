package storage

import (
	"context"
	"database/sql"
	"fmt"

	error_messages "github.com/sirUnchained/my-go-instagram/internal/errors"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type userStore struct {
	db *sql.DB
}

func (us *userStore) Create(ctx context.Context, userP *payloads.UserPayload) error {
	query := `INSERT INTO users (username, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5);`

	user := models.UserModel{Username: userP.Username, Fullname: userP.Fullname, Email: userP.Email, Password: models.Password{}, Role: models.RoleModel{}}
	user.Password.Set(userP.Password)

	var userCount int
	err := us.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
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
	).Err()

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"":
			return error_messages.USERNAME_DUP

		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return error_messages.EMAIL_DUP

		default:
			return err
		}
	}

	return nil
}
