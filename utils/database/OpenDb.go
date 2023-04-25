package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./store.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
} 
