package tasks

import "github.com/google/uuid"

type Task struct {
	ID        uuid.UUID
	Title     string
	CreatedBy uint
	Completed bool
	ParentID  uuid.UUID
}
