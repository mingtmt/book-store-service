package service

import (
	"strings"

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

func (us *userService) GetAllUsers(search string, page, limit int) ([]model.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(utils.ErrCodeInternal, "Failed to fetch users", err)
	}

	var filterUsers []model.User
	if search != "" {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)

			if strings.Contains(name, search) || strings.Contains(email, search) {
				filterUsers = append(filterUsers, user)
			}
		}
	} else {
		filterUsers = users
	}

	start := (page - 1) * limit
	if start >= len(filterUsers) {
		return []model.User{}, nil
	}

	end := start + limit
	if end > len(filterUsers) {
		end = len(filterUsers)
	}

	return filterUsers[start:end], nil
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

func (us *userService) GetUserByUUID(uuid string) (model.User, error) {
	user, found := us.repo.FindByUUID(uuid)
	if !found {
		return model.User{}, utils.NewError(utils.ErrCodeNotFound, "User not found")
	}

	return user, nil
}

func (us *userService) UpdateUser(uuid string, updatedUser model.User) (model.User, error) {
	updatedUser.Email = utils.NormalizeString(updatedUser.Email)
	if u, exist := us.repo.FindByEmail(updatedUser.Email); exist && u.UUID != uuid {
		return model.User{}, utils.NewError(utils.ErrCodeConflict, "Email already exist")
	}

	currentUser, found := us.repo.FindByUUID(uuid)
	if !found {
		return model.User{}, utils.NewError(utils.ErrCodeNotFound, "User not found")
	}

	currentUser.Name = updatedUser.Name
	currentUser.Email = updatedUser.Email
	currentUser.Age = updatedUser.Age
	currentUser.Status = updatedUser.Status
	currentUser.Level = updatedUser.Level

	if updatedUser.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, utils.WrapError(utils.ErrCodeInternal, "Failed to hash password", err)
		}

		currentUser.Password = string(hashPassword)
	}

	if err := us.repo.Update(uuid, currentUser); err != nil {
		return model.User{}, utils.WrapError(utils.ErrCodeInternal, "Failed to update user", err)
	}

	return currentUser, nil
}

func (us *userService) DeleteUser() {

}
