package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", "./db.sqlite")
}

func CreateDB(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			created_by INTEGER NOT NULL, 
			completed BOOLEAN NOT NULL DEFAULT false,
			completed_by INTEGER NULL,
			parent_id TEXT NULL,
			cost INTEGER NOT NULL DEFAULT 0
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS task_tree (
			id TEXT PRIMARY KEY,
			parent_id TEXT NOT NULL,
			root_id TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}
