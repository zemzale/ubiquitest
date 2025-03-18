package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Exists(userID uint) (bool, error) {
	var id uint
	err := r.db.Get(&id, "SELECT id FROM users WHERE id = ?", userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user id: %w", err)
	}

	return true, nil
}
