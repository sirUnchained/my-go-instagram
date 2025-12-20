package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgreSQL(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	database, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(maxOpenConns)
	database.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	database.SetConnMaxIdleTime(duration)

	// check connection health
	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
}
