package tasks

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zemzale/ubiquitest/storage"
)

type Store struct {
	db         *sqlx.DB
	insertTask *storage.TaksRepository
}

func NewStore(db *sqlx.DB, insertTask *storage.TaksRepository) *Store {
	return &Store{db: db, insertTask: insertTask}
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

	if err := s.insertTask.Create(mapNewTaskToDB(task)); err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
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
	err := s.db.Get(&id, "SELECT id FROM tasks WHERE id = ?", parentID)
	if err != nil {
		return fmt.Errorf("failed to get parent id: %w", err)
	}
	return nil
}
