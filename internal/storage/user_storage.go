package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type userStore struct {
	db *sql.DB
}

func (us *userStore) Create(ctx context.Context, userP *payloads.CreateUserPayload) (*models.UserModel, error) {
	// # using transactions
	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// # first create a profile
	profileQuery := `INSERT INTO profiles (fullname, bio, avatar) VALUES ($1, $2, $3) RETURNING id;`
	userProfile := models.ProfileModel{Fullname: userP.Fullname, Bio: userP.Bio}
	if err := tx.QueryRowContext(ctx, profileQuery, userProfile.Fullname, userProfile.Bio, nil).Scan(&userProfile.Id); err != nil {
		return nil, err
	}

	// # now create user
	userQuery := `INSERT INTO users (username, email, password, role, profile) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;`
	userRole := models.RoleModel{}
	user := models.UserModel{Username: userP.Username, Email: userP.Email, Password: models.Password{}, Role: userRole, Profile: userProfile}
	user.Password.Set(userP.Password)

	// ## set user role by users count logic
	var userCount int
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}
	if userCount == 0 {
		user.Role.Id = 2
		user.IsVerified = true
	} else {
		user.Role.Id = 1
		user.IsVerified = false
	}

	err = tx.QueryRowContext(ctx, userQuery,
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

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *userStore) GetById(ctx context.Context, userId int64) (*models.UserModel, error) {
	// todo : select avatar too
	query := `
	SELECT u.id, u.username, u.email, u.password, u.is_verified, is_private, r.id, r.name, u.created_at, u.updated_at, p.id, p.fullname, p.bio
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
	// todo : select avatar too
	query := `
	SELECT u.id, u.username, u.email, u.password, u.is_verified, is_private, r.id, r.name, u.created_at, u.updated_at, p.id, p.fullname, p.bio
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

func (us *userStore) UpdateData(ctx context.Context, user *models.UserModel, userP *payloads.CreateUserPayload) (*models.UserModel, error) {
	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateUser := &models.UserModel{
		Id:         user.Id,
		Username:   userP.Username,
		Email:      userP.Email,
		IsPrivate:  user.IsPrivate,
		IsVerified: user.IsVerified,
		Password:   models.Password{},
		Profile: models.ProfileModel{
			Fullname: userP.Fullname,
			Bio:      userP.Bio,
		},
	}
	if err := updateUser.Password.Set(userP.Password); err != nil {
		return nil, err
	}

	// update profile first
	queryProfile := `UPDATE profiles SET fullname = $1, bio = $2, updated_at = NOW() WHERE id = $3`
	if _, err := tx.ExecContext(ctx, queryProfile, updateUser.Profile.Fullname, updateUser.Profile.Bio, updateUser.Id); err != nil {
		return nil, err
	}

	// if user updated email so he/she is no more verified
	emailChanged := strings.Compare(updateUser.Email, user.Email) == 0
	if emailChanged {
		queryUser := `UPDATE users SET password = $1, username = $2, updated_at = NOW() WHERE id = $3`
		_, err = tx.ExecContext(ctx, queryUser, updateUser.Password.Hash, updateUser.Username, updateUser.Id)
	} else {
		queryUser := `UPDATE users SET password = $1, email = $2, username = $3, is_verified = FALSE, updated_at = NOW() WHERE id = $4`
		updateUser.IsVerified = false
		_, err = tx.ExecContext(ctx, queryUser, updateUser.Password.Hash, updateUser.Email, updateUser.Username, updateUser.Id)
	}
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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return updateUser, nil
}
