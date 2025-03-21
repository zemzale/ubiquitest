package tasks

import "github.com/google/uuid"

type CalculateCost struct{}

func NewCalculateCost() *CalculateCost { return &CalculateCost{} }

func (c *CalculateCost) Run(tasks []Task) []Task {
	tasksMap := make(map[uuid.UUID]Task)

	for _, task := range tasks {
		tasksMap[task.ID] = task
	}

	for _, task := range tasks {
		if task.ParentID == uuid.Nil {
			continue
		}

		parentTask, ok := tasksMap[task.ParentID]
		if !ok {
			continue
		}

		parentTask.Cost += task.Cost

		tasksMap[parentTask.ID] = parentTask
	}

	for i, task := range tasks {
		if task.ParentID != uuid.Nil {
			continue
		}
		tasks[i].Cost = tasksMap[task.ID].Cost
	}

	return tasks
}
