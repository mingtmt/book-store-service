package repository

import (
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

func (ur *InMemUserRepository) FindByUUID() {

}

func (ur *InMemUserRepository) Update() {

}

func (ur *InMemUserRepository) Delete() {

}

func (ur *InMemUserRepository) FindByEmail(email string) (model.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}

	return model.User{}, false
}
