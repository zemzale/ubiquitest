package tasks

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type List struct {
	db *sqlx.DB
}

func NewList(db *sqlx.DB) *List {
	return &List{db: db}
}

func (l *List) Run() ([]Task, error) {
	rows, err := l.db.Query("SELECT id, title, created_by FROM todos")
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}

	// TODO: Move this to some sort of cache
	users := make(map[uint]string)

	todos := make([]Task, 0)
	for rows.Next() {
		var id string
		var title string
		var createdBy uint
		if err := rows.Scan(&id, &title, &createdBy); err != nil {
			return nil, fmt.Errorf("failed to scan todos: %w", err)
		}

		username, ok := users[createdBy]
		if !ok {
			result := l.db.QueryRow("SELECT username FROM users WHERE id = ?", createdBy)
			if err := result.Scan(&username); err != nil {
				return nil, fmt.Errorf("failed to scan user: %w", err)
			}
			users[createdBy] = username
		}

		todos = append(todos, Task{ID: uuid.MustParse(id), Title: title, CreateBy: username})
	}

	return todos, nil
}
