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

type FindByUsername struct {
	userRepo *storage.UserRepository
}

func NewFindByUsername(userRepo *storage.UserRepository) *FindByUsername {
	return &FindByUsername{userRepo: userRepo}
}

func (r *FindByUsername) Run(username string) (User, error) {
	userRecord, err := r.userRepo.FindByUsername(username)
	if err != nil {
		return User{}, err
	}

	return User{ID: userRecord.ID, Username: userRecord.Username}, nil
}
