package tasks

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zemzale/ubiquitest/storage"
)

type Store struct {
	db *sqlx.DB
	// TOOD Fix the name typo and name
	insertTask *storage.TaksRepository
	userRepo   *storage.UserRepository
}

func NewStore(db *sqlx.DB, insertTask *storage.TaksRepository, userRepo *storage.UserRepository) *Store {
	return &Store{db: db, insertTask: insertTask, userRepo: userRepo}
}

func (s *Store) Run(task Task) error {
	// TODO: Use storage for checking user stuff
	userExists, err := s.userRepo.Exists(task.CreatedBy)
	if err != nil {
		return err
	}

	if !userExists {
		return fmt.Errorf("user doesn't exist")
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
