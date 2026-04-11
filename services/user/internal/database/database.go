package database

import (
	"database/sql"
	"fmt"
)

func NewDB(connStr string) (*sql.DB, error) {
	db, err := connect(connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return db, nil
}
