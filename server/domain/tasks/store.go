package tasks

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/zemzale/ubiquitest/storage"
)

type Store struct {
	updateParentCost *UpdateParentCost

	taskRepo *storage.TaksRepository
	userRepo *storage.UserRepository
}

func NewStore(updateParentCost *UpdateParentCost, insertTask *storage.TaksRepository, userRepo *storage.UserRepository) *Store {
	return &Store{updateParentCost: updateParentCost, taskRepo: insertTask, userRepo: userRepo}
}

func (s *Store) Run(task Task) error {
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

	if err := s.taskRepo.Create(mapNewTaskToDB(task)); err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	if err := s.updateParentCost.Run(task.ParentID, task.Cost); err != nil {
		return fmt.Errorf("failed to update parent cost: %w", err)
	}

	return nil
}

func (s *Store) checkIfParentExists(parentID uuid.UUID) error {
	if parentID == uuid.Nil {
		return nil
	}

	err := s.taskRepo.CheckIfParentExists(parentID.String())
	if err != nil {
		return fmt.Errorf("failed to get parent id: %w", err)
	}
	return nil
}
