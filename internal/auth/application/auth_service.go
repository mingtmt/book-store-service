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

type AuthService struct {
	authRepo    domain.AuthRepository
	rfTokenRepo domain.RefreshTokenRepository
}

func NewAuthService(authRepo domain.AuthRepository, rfTokenRepo domain.RefreshTokenRepository) *AuthService {
	return &AuthService{
		authRepo:    authRepo,
		rfTokenRepo: rfTokenRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password string) (string, string, string, error) {
	existing, err := s.authRepo.FindByUsername(ctx, username)
	if err != nil && err != errors.ErrUserNotFound {
		return "", "", "", err
	}
	if existing != nil {
		return "", "", "", errors.ErrUserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", err
	}

	user := &domain.Auth{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashed),
	}

	userID, err := s.authRepo.RegisterUser(ctx, user)
	if err != nil {
		return "", "", "", err
	}

	// Generate access and refresh tokens, and store the refresh token
	accessToken, refreshToken, refreshExp, err := token.GenerateTokenPair(userID)
	if err != nil {
		return "", "", "", err
	}

	err = s.rfTokenRepo.CreateRefreshToken(ctx, uuid.MustParse(userID), refreshToken, refreshExp)
	if err != nil {
		return "", "", "", err
	}

	return userID, accessToken, refreshToken, nil
}

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.authRepo.FindByUsername(ctx, username)
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
	if err = s.rfTokenRepo.CreateRefreshToken(ctx, userUUID, refreshToken, expiresAt); err != nil {
		return "", "", err
	}
	if err := s.rfTokenRepo.DeleteExpiredRefreshTokens(ctx); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, tokenStr string) (string, string, error) {
	rt, err := s.rfTokenRepo.GetRefreshToken(ctx, tokenStr)
	if err != nil {
		return "", "", err
	}
	if rt.ExpiresAt.Before(time.Now()) || rt.Revoked {
		return "", "", err
	}

	// Revoke old token
	_ = s.rfTokenRepo.RevokeRefreshToken(ctx, tokenStr)

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
	err = s.rfTokenRepo.CreateRefreshToken(ctx, userUUID, refresh, refreshExp)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, tokenStr string) error {
	return s.rfTokenRepo.RevokeRefreshToken(ctx, tokenStr)
}
