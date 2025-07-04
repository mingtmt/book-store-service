package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mingtmt/book-store/internal/users/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password string) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashed),
	}

	return s.repo.CreateUser(ctx, user)
}
