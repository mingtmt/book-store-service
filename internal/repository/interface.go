package repository

import "github.com/mingtmt/book-store/internal/model"

type UserRepository interface {
	FindAll()
	Create(user model.User) error
	FindByUUID()
	Update()
	Delete()
	FindByEmail(email string) (model.User, bool)
}
