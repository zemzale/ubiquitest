package storage

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Task struct {
	ID          string           `db:"id"`
	Title       string           `db:"title"`
	CreatedBy   uint             `db:"created_by"`
	Completed   bool             `db:"completed"`
	CompletedBy sql.Null[uint]   `db:"completed_by"`
	ParentID    sql.Null[string] `db:"parent_id"`
}

type TaksRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaksRepository {
	return &TaksRepository{db: db}
}

func (r *TaksRepository) Create(todo Task) error {
	query := `INSERT INTO tasks 
		(id, title, created_by, completed, completed_by, parent_id)
	VALUES 
		(:id, :title, :created_by, :completed, :completed_by, :parent_id)`
	result, err := r.db.NamedExec(query, todo)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
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

func (s *TaksRepository) CheckIfParentExists(parentID string) error {
	var id string
	err := s.db.Get(&id, "SELECT id FROM tasks WHERE id = ?", parentID)
	if err != nil {
		return fmt.Errorf("failed to get parent id: %w", err)
	}

	return nil
}

func (s *TaksRepository) List() ([]*Task, error) {
	const query = `
		SELECT 
			tasks.id, 
			tasks.title, 
			tasks.created_by, 
			tasks.parent_id, 
			tasks.completed, 
			users.username 
		FROM tasks 
		LEFT JOIN users ON users.id = tasks.created_by
	`
	rows, err := s.db.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var id string
		var title string
		var createdBy uint
		var parentID sql.Null[string]
		var completed bool
		var username string
		if err := rows.Scan(&id, &title, &createdBy, &parentID, &completed, &username); err != nil {
			return nil, fmt.Errorf("failed to scan tasks: %w", err)
		}

		tasks = append(tasks, &Task{
			ID:        id,
			Title:     title,
			CreatedBy: createdBy,
			Completed: completed,
			ParentID:  parentID,
		})
	}

	return tasks, nil
}
