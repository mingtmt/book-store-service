package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mingtmt/book-store/internal/auth/domain"
	"github.com/mingtmt/book-store/internal/auth/infrastructure/token"
	"github.com/mingtmt/book-store/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	RegisterUser(ctx context.Context, user *domain.Auth) (*domain.Auth, error)
	FindByUsername(ctx context.Context, username string) (*domain.Auth, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password string) (*domain.Auth, error) {
	existing, err := s.repo.FindByUsername(ctx, username)
	if err != nil && err != errors.ErrUserNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.Auth{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashed),
	}

	return s.repo.RegisterUser(ctx, user)
}

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.ErrInvalidPassword
	}

	return token.GenerateToken(user.ID)
}
