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
