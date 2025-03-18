package tasks

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/zemzale/ubiquitest/storage"
)

type Task struct {
	ID        uuid.UUID
	Title     string
	CreatedBy uint
	Completed bool
	ParentID  uuid.UUID
}

func mapNewTaskToDB(task Task) storage.Task {
	parentID := uuid.Nil
	if task.ParentID != uuid.Nil {
		parentID = task.ParentID
	}

	return storage.Task{
		ID:        task.ID.String(),
		Title:     task.Title,
		CreatedBy: task.CreatedBy,
		Completed: task.Completed,
		ParentID:  sql.Null[string]{V: parentID.String(), Valid: true},
	}
}
