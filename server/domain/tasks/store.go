package tasks

import (
	"fmt"

	"github.com/google/uuid"
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

	err = s.checkIfParentExists(task.ParentID)
	if err != nil {
		return fmt.Errorf("failed to check if parent exists: %w", err)
	}
	result, err := s.db.Exec(
		"INSERT INTO todos (id, title, created_by, parent_id) VALUES (?, ?, ?, ?)",
		task.ID.String(), task.Title, task.CreatedBy, task.ParentID.String(),
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

func (s *Store) checkIfParentExists(parentID uuid.UUID) error {
	if parentID == uuid.Nil {
		return nil
	}

	var id uuid.UUID
	err := s.db.Get(&id, "SELECT id FROM todos WHERE id = ?", parentID)
	if err != nil {
		return fmt.Errorf("failed to get parent id: %w", err)
	}
	return nil
}
