package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mingtmt/book-store/internal/auth/domain"
	"github.com/mingtmt/book-store/pkg/errors"
	"github.com/mingtmt/book-store/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	RegisterUser(ctx context.Context, user *domain.Auth) (*domain.Auth, error)
	FindByUsername(ctx context.Context, username string) (*domain.Auth, error)
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenStr string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenStr string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password string) (*domain.Auth, string, string, error) {
	existing, err := s.repo.FindByUsername(ctx, username)
	if err != nil && err != errors.ErrUserNotFound {
		return nil, "", "", err
	}
	if existing != nil {
		return nil, "", "", errors.ErrUserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	user := &domain.Auth{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashed),
	}

	newUser, err := s.repo.RegisterUser(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	// Generate access and refresh tokens, and store the refresh token
	accessToken, refreshToken, refreshExp, err := token.GenerateTokenPair(newUser.ID)
	if err != nil {
		return nil, "", "", err
	}

	err = s.repo.CreateRefreshToken(ctx, uuid.MustParse(newUser.ID), refreshToken, refreshExp)
	if err != nil {
		return nil, "", "", err
	}

	return newUser, accessToken, refreshToken, nil
}

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.ErrInvalidPassword
	}

	accessToken, refreshToken, expiresAt, err := token.GenerateTokenPair(user.ID)
	if err != nil {
		return "", "", err
	}

	// Store refresh token in the database and delete expired tokens
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return "", "", err
	}
	if err = s.repo.CreateRefreshToken(ctx, userUUID, refreshToken, expiresAt); err != nil {
		return "", "", err
	}
	if err := s.repo.DeleteExpiredRefreshTokens(ctx); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, tokenStr string) (string, string, error) {
	rt, err := s.repo.GetRefreshToken(ctx, tokenStr)
	if err != nil {
		return "", "", err
	}
	if rt.ExpiresAt.Before(time.Now()) || rt.Revoked {
		return "", "", err
	}

	// Revoke old token
	_ = s.repo.RevokeRefreshToken(ctx, tokenStr)

	// Generate new token pair
	access, refresh, refreshExp, err := token.GenerateTokenPair(rt.UserID)
	if err != nil {
		return "", "", err
	}

	// Store new refresh token
	userUUID, err := uuid.Parse(rt.UserID)
	if err != nil {
		return "", "", err
	}
	err = s.repo.CreateRefreshToken(ctx, userUUID, refresh, refreshExp)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
