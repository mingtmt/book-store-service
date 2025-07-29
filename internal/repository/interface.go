package repository

import "github.com/mingtmt/book-store/internal/model"

type UserRepository interface {
	FindAll() ([]model.User, error)
	Create(user model.User) error
	FindByUUID(uuid string) (model.User, bool)
	Update(uuid string, user model.User) error
	Delete()
	FindByEmail(email string) (model.User, bool)
}
