package users

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FindByID struct {
	db *sqlx.DB
}

func NewFindById(db *sqlx.DB) *FindByID {
	return &FindByID{db: db}
}

func (r *FindByID) Run(userID uint) (User, error) {
	// TODO Move this to DB layer
	var user struct {
		Id       uint   `db:"id"`
		Username string `db:"username"`
	}
	err := r.db.Get(&user, "SELECT * FROM users where id=?", userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user with id %d: %w", userID, err)
	}

	return User{ID: user.Id, Username: user.Username}, nil
}
