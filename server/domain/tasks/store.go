package tasks

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Run(task Task) error {
	err := s.checkUserExists(task.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to get user id: %w", err)
	}
	result, err := s.db.Exec(
		"INSERT INTO todos (id, title, created_by) VALUES (?, ?, ?)",
		task.ID.String(), task.Title, task.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if res == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (s *Store) checkUserExists(userID uint) error {
	var id uint
	err := s.db.Get(&id, "SELECT id FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to get user id: %w", err)
	}
	return nil
}
