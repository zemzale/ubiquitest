package tasks

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/storage"
)

type FindAllParents struct {
	repo *storage.TaksRepository
}

func NewFindAllParents(repo *storage.TaksRepository) *FindAllParents {
	return &FindAllParents{repo: repo}
}

func (f *FindAllParents) Run(parentID uuid.UUID) ([]Task, error) {
	parent, err := f.repo.Find(parentID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find parent: %w", err)
	}

	parentRecords := []*storage.Task{parent}

	for {
		children, err := f.repo.ListChildren(parent.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list children: %w", err)
		}

		parentRecords = append(parentRecords, children...)

		if !parent.ParentID.Valid || parent.ParentID.V == uuid.Nil.String() {
			break
		}

		parent, err = f.repo.Find(parent.ParentID.V)
		if err != nil {
			return nil, fmt.Errorf("failed to find parent: %w", err)
		}
		parentRecords = append(parentRecords, parent)
	}

	return lo.Map(parentRecords, func(t *storage.Task, _ int) Task {
		return mapNewTaskFromDB(*t)
	}), nil
}
