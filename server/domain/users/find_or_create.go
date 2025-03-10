package users

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type FindOrCreate struct {
	db *sqlx.DB
}

func NewFindOrCreate(db *sqlx.DB) *FindOrCreate {
	return &FindOrCreate{db: db}
}

func (u *FindOrCreate) Run(username string) (User, error) {
	result := u.db.QueryRow("SELECT id FROM users WHERE username = ?", username)
	var userID uint
	if err := result.Scan(&userID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return User{}, err
		}

		return u.insertUser(username)
	}

	return User{ID: userID, Username: username}, nil
}

func (u *FindOrCreate) insertUser(username string) (User, error) {
	result := u.db.QueryRow("INSERT INTO users (username) VALUES (?) RETURNING id", username)
	var userID uint
	if err := result.Scan(&userID); err != nil {
		return User{}, err
	}

	return User{ID: userID, Username: username}, nil
}
