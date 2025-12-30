package storage

import (
	"database/sql"
	"strings"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
	"golang.org/x/net/context"
)

type banStore struct {
	db *sql.DB
}

func (bs *banStore) Create(ctx context.Context, user *models.UserModel, whyBanned *payloads.CreateBanPayload) error {
	return executeTransaction(ctx, bs.db, func(ctx context.Context, tx *sql.Tx) error {
		deleteUserQuery := `DELETE FROM users WHERE email = $1;`
		if _, err := tx.ExecContext(ctx, deleteUserQuery, user.Email); err != nil {
			switch {
			case err.Error() == "sql: no rows in result set":
				return global_varables.NOT_FOUND_ROW
			default:
				return err
			}
		}

		banQuery := `INSERT INTO bans (email, why_banned) VALUES ($1, $2);`
		if _, err := tx.ExecContext(ctx, banQuery, whyBanned.Email, whyBanned.WhyBanned); err != nil {
			return err
		}

		return nil
	})
}

func (bs *banStore) Delete(ctx context.Context, email string) error {
	query := `DELETE FROM bans where email = $1;`
	if _, err := bs.db.ExecContext(ctx, query, email); err != nil {
		return err
	}

	return nil
}

func (bs *banStore) GetBanByEmail(ctx context.Context, email string) (*models.BanModel, error) {
	ban := &models.BanModel{}
	query := `SELECT id, email, why_bannedd, created_at FROM bans WHERE email = $1;`
	if err := bs.db.QueryRowContext(ctx, query, email).Scan(&ban.Id, &ban.Email, &ban.WhyBanned, &ban.CreatedAt); err != nil {
		switch {
		case strings.Compare(err.Error(), "sql: no rows in result set") == 0:
			return nil, global_varables.NOT_FOUND_ROW

		default:
			return nil, err
		}
	}

	return ban, nil
}
