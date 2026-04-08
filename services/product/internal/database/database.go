package database

import (
	"database/sql"
)

func NewDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, nil
	}

	return db, nil
}
