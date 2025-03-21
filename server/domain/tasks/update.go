package tasks

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Update struct {
	db *sqlx.DB
}

func NewUpdate(db *sqlx.DB) *Update {
	return &Update{db: db}
}

func (u *Update) Run(task Task, userID uint) error {
	if task.Completed {
		return u.completeTask(task, userID)
	}

	result, err := u.db.Exec(
		"UPDATE tasks SET title = ?, completed = ?, completed_by = ?, cost = ? WHERE id = ?",
		task.Title, task.Completed, nil, task.ID.String(), task.Cost,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
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

func (u *Update) completeTask(task Task, userID uint) error {
	result, err := u.db.Exec(
		"UPDATE tasks SET title = ?, completed = ?, completed_by = ? WHERE id = ?",
		task.Title, task.Completed, userID, task.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
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
