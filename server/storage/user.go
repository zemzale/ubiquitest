package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       uint   `db:"id"`
	Username string `db:"username"`
}

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

func (r *UserRepository) FindByID(userID uint) (User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users where id=?", userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user with id %d: %w", userID, err)
	}

	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users where username=?", username)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user with username %s: %w", username, err)
	}

	return user, nil
}
