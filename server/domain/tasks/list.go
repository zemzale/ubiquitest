package tasks

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zemzale/ubiquitest/storage"
)

type List struct {
	db       *sqlx.DB
	taskRepo *storage.TaksRepository
}

func NewList(db *sqlx.DB, taskRepo *storage.TaksRepository) *List {
	return &List{db: db, taskRepo: taskRepo}
}

func (l *List) Run() ([]Task, error) {
	tasksRecords, err := l.taskRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	tasks := make([]Task, 0, len(tasksRecords))
	for _, taskRecord := range tasksRecords {
		parentID := uuid.Nil
		if taskRecord.ParentID.Valid {
			parentID = uuid.MustParse(taskRecord.ParentID.V)
		}

		tasks = append(tasks, Task{
			ID:        uuid.MustParse(taskRecord.ID),
			Title:     taskRecord.Title,
			CreatedBy: taskRecord.CreatedBy,
			Completed: taskRecord.Completed,
			ParentID:  parentID,
			Cost:      taskRecord.Cost,
		})
	}

	return tasks, nil
}
