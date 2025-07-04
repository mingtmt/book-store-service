package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mingtmt/book-store/internal/users/domain"
	"github.com/mingtmt/book-store/internal/users/infrastructure/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
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

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	return jwt.GenerateToken(user.ID)
}
