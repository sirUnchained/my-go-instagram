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
	profileQuery := `INSERT INTO profiles (fullname, bio, avatar) VALUES ($1, $2, $3) RETURNING id;`
	userProfile := models.ProfileModel{Fullname: userP.Fullname, Bio: userP.Bio, Avatar: userP.Avatar}
	if err := us.db.QueryRowContext(ctx, profileQuery, userProfile.Fullname, userProfile.Bio, userProfile.Avatar).Scan(&userProfile.Id); err != nil {
		return nil, err
	}

	userQuery := `INSERT INTO users (username, email, password, role, profile) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;`
	userRole := models.RoleModel{}
	user := models.UserModel{Username: userP.Username, Email: userP.Email, Password: models.Password{}, Role: userRole, Profile: userProfile}
	user.Password.Set(userP.Password)

	var userCount int
	if err := us.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	if userCount == 0 {
		user.Role.Id = 2
		user.IsVerified = true
	} else {
		user.Role.Id = 1
		user.IsVerified = false
	}

	err := us.db.QueryRowContext(ctx, userQuery,
		user.Username,
		user.Email,
		user.Password.Hash,
		user.Role.Id,
		user.Profile.Id,
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
	query := `
	SELECT u.id, u.username, u.email, u.password, u.is_verified, is_private, r.id, r.name, u.created_at, u.updated_at, p.id, p.fullname, p.bio, p.avatar
	from users AS u
	JOIN roles AS r on r.id = u.role
	JOIN profiles AS p ON p.id = u.profile
	WHERE u.id = $1
	`

	user := &models.UserModel{
		Profile: models.ProfileModel{},
		Role:    models.RoleModel{},
	}
	err := us.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.IsVerified,
		&user.IsPrivate,
		&user.Role.Id,
		&user.Role.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Profile.Id,
		&user.Profile.Fullname,
		&user.Profile.Bio,
		&user.Profile.Avatar,
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
	query := `
	SELECT u.id, u.username, u.email, u.password, u.is_verified, is_private, r.id, r.name, u.created_at, u.updated_at, p.id, p.fullname, p.bio, p.avatar
	from users AS u
	JOIN roles AS r on r.id = u.role
	JOIN profiles AS p ON p.id = u.profile
	WHERE u.email = $1
	`

	user := &models.UserModel{}
	err := us.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.IsVerified,
		&user.IsPrivate,
		&user.Role.Id,
		&user.Role.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Profile.Id,
		&user.Profile.Fullname,
		&user.Profile.Bio,
		&user.Profile.Avatar,
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
