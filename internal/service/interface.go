package service

import "github.com/mingtmt/book-store/internal/model"

type UserService interface {
	GetAllUsers() ([]model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
