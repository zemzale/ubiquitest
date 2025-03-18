package tasks

import (
	"database/sql"
	"fmt"
	"log"

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
	rows, err := l.db.Query("SELECT id, title, created_by, parent_id, completed FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	// TODO: Move this to some sort of cache
	users := make(map[uint]string)

	tasks := make([]Task, 0)
	for rows.Next() {
		var id string
		var title string
		var createdBy uint
		var parentID sql.NullString
		var completed bool
		if err := rows.Scan(&id, &title, &createdBy, &parentID, &completed); err != nil {
			return nil, fmt.Errorf("failed to scan tasks: %w", err)
		}

		username, ok := users[createdBy]
		if !ok {
			result := l.db.QueryRow("SELECT username FROM users WHERE id = ?", createdBy)
			if err := result.Scan(&username); err != nil {
				return nil, fmt.Errorf("failed to scan user: %w", err)
			}
			users[createdBy] = username
		}

		parentUUID := uuid.Nil
		if parentID.Valid {
			parentUUID, err = uuid.Parse(parentID.String)
			if err != nil {
				log.Println("failed to parse parent id:", err)
				parentUUID = uuid.Nil
			}
		}

		tasks = append(tasks, Task{
			ID:        uuid.MustParse(id),
			Title:     title,
			CreatedBy: createdBy,
			Completed: completed,
			ParentID:  parentUUID,
		})
	}

	return tasks, nil
}
