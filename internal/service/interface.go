package service

import "github.com/mingtmt/book-store/internal/model"

type UserService interface {
	GetAllUsers()
	CreateUser(user model.User) (model.User, error)
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
