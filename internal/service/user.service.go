package service

import (
	"github.com/google/uuid"
	"github.com/mingtmt/book-store/internal/model"
	"github.com/mingtmt/book-store/internal/repository"
	"github.com/mingtmt/book-store/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers() {

}

func (us *userService) CreateUser(user model.User) (model.User, error) {
	user.Email = utils.NormalizeString(user.Email)
	if _, exist := us.repo.FindByEmail(user.Email); exist {
		return model.User{}, utils.NewError(utils.ErrCodeConflict, "Email already exist")
	}

	user.UUID = uuid.New().String()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, utils.WrapError(utils.ErrCodeInternal, "Failed to hash password", err)
	}

	user.Password = string(hashPassword)

	if err := us.repo.Create(user); err != nil {
		return model.User{}, utils.WrapError(utils.ErrCodeInternal, "Failed to create new user", err)
	}

	return user, nil
}

func (us *userService) GetUserByUUID() {

}

func (us *userService) UpdateUser() {

}

func (us *userService) DeleteUser() {

}
