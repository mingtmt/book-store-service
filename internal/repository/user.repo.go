package repository

import (
	"fmt"
	"slices"

	"github.com/mingtmt/book-store/internal/model"
)

type InMemUserRepository struct {
	users []model.User
}

func NewInMemUserRepository() UserRepository {
	return &InMemUserRepository{
		users: make([]model.User, 0),
	}
}

func (ur *InMemUserRepository) FindAll() ([]model.User, error) {
	return ur.users, nil
}

func (ur *InMemUserRepository) Create(user model.User) error {
	ur.users = append(ur.users, user)
	return nil
}

func (ur *InMemUserRepository) FindByUUID(uuid string) (model.User, bool) {
	for _, user := range ur.users {
		if user.UUID == uuid {
			return user, true
		}
	}

	return model.User{}, false
}

func (ur *InMemUserRepository) Update(uuid string, user model.User) error {
	for i, u := range ur.users {
		if u.UUID == uuid {
			ur.users[i] = user
			return nil
		}
	}

	return fmt.Errorf("User not found")
}

func (ur *InMemUserRepository) Delete(uuid string) error {
	for i, u := range ur.users {
		if u.UUID == uuid {
			ur.users = slices.Delete(ur.users, i, i+1)
			return nil
		}
	}

	return fmt.Errorf("User not found")
}

func (ur *InMemUserRepository) FindByEmail(email string) (model.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}

	return model.User{}, false
}
