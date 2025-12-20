package storage

import "database/sql"

type PgStorage struct {
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{}
}
