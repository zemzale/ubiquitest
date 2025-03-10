package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", ":memory:")
}

func CreateDB(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			created_by INTEGER NOT NULL, 
			completed BOOLEAN NOT NULL DEFAULT false,
			completed_by INTEGER NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}
	return nil
}
