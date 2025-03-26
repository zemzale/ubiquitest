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
	Cost        uint             `db:"cost"`
	TotalCost   uint             `db:"total_cost"`
}

type TaksRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaksRepository {
	return &TaksRepository{db: db}
}

func (r *TaksRepository) Create(todo Task) error {
	query := `INSERT INTO tasks 
		(id, title, created_by, completed, completed_by, parent_id, cost, total_cost)
	VALUES 
		(:id, :title, :created_by, :completed, :completed_by, :parent_id, :cost, :total_cost)`
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
			tasks.cost,
			tasks.total_cost,
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
		var cost uint
		var totalCost uint
		if err := rows.Scan(&id, &title, &createdBy, &parentID, &completed, &cost, &totalCost, &username); err != nil {
			return nil, fmt.Errorf("failed to scan tasks: %w", err)
		}

		tasks = append(tasks, &Task{
			ID:        id,
			Title:     title,
			CreatedBy: createdBy,
			Completed: completed,
			ParentID:  parentID,
			Cost:      cost,
			TotalCost: totalCost,
		})
	}

	return tasks, nil
}

func (s *TaksRepository) ListChildren(parentID string) ([]*Task, error) {
	const query = `
		SELECT 
			tasks.id, 
			tasks.title, 
			tasks.created_by, 
			tasks.parent_id, 
			tasks.completed, 
			tasks.cost,
			users.username 
		FROM tasks 
		LEFT JOIN users ON users.id = tasks.created_by
		WHERE tasks.parent_id = ?
	`
	rows, err := s.db.Queryx(query, parentID)
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
		var cost uint
		if err := rows.Scan(&id, &title, &createdBy, &parentID, &completed, &cost, &username); err != nil {
			return nil, fmt.Errorf("failed to scan tasks: %w", err)
		}

		tasks = append(tasks, &Task{
			ID:        id,
			Title:     title,
			CreatedBy: createdBy,
			Completed: completed,
			ParentID:  parentID,
			Cost:      cost,
		})
	}

	return tasks, nil
}

func (s *TaksRepository) Find(id string) (*Task, error) {
	var task Task

	err := s.db.Get(&task, "SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (s *TaksRepository) UpdateTotalCost(parentID string, cost uint) error {
	query := `UPDATE tasks SET total_cost = total_cost + ? WHERE id = ?`
	_, err := s.db.Exec(query, cost, parentID)
	return err
}
