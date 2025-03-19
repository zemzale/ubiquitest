package users

import (
	"github.com/zemzale/ubiquitest/storage"
)

type FindByID struct {
	userRepo *storage.UserRepository
}

func NewFindById(userRepo *storage.UserRepository) *FindByID {
	return &FindByID{userRepo: userRepo}
}

func (r *FindByID) Run(userID uint) (User, error) {
	userRecord, err := r.userRepo.FindByID(userID)
	if err != nil {
		return User{}, err
	}

	return User{ID: userRecord.ID, Username: userRecord.Username}, nil
}
