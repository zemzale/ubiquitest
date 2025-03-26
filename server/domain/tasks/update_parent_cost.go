package tasks

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/zemzale/ubiquitest/storage"
)

type UpdateParentCost struct {
	findAllParents *FindAllParents
	repo           *storage.TaksRepository
}

func NewUpdateParentCost(findAllParents *FindAllParents, repo *storage.TaksRepository) *UpdateParentCost {
	return &UpdateParentCost{findAllParents: findAllParents, repo: repo}
}

func (u *UpdateParentCost) Run(parentID uuid.UUID, cost uint) error {
	if parentID == uuid.Nil {
		return nil
	}

	if cost == 0 {
		return nil
	}

	parents, err := u.findAllParents.Run(parentID)
	if err != nil {
		return fmt.Errorf("failed to find parents: %w", err)
	}

	for _, parent := range parents {
		log.Println("updating parent cost ", parent.ID.String(), cost)
		u.repo.UpdateTotalCost(parent.ID.String(), cost)
	}

	return nil
}
