package repository

import (
	"log"

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

func (ur *InMemUserRepository) FindAll() {
	log.Println("Finding all users in memory repository")
}

func (ur *InMemUserRepository) Create() {

}

func (ur *InMemUserRepository) FindByUUID() {

}

func (ur *InMemUserRepository) Update() {

}

func (ur *InMemUserRepository) Delete() {

}
