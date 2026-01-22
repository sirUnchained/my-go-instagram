package storage

import (
	"database/sql"
)

type reportStore struct {
	db *sql.DB
}
