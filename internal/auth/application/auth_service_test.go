package application

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mingtmt/book-store/internal/auth/domain"
	"github.com/mingtmt/book-store/pkg/errors"
	"github.com/mingtmt/book-store/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) RegisterUser(ctx context.Context, user *domain.Auth) (*domain.Auth, error) {
	args := m.Called(user)
	if createdUser, ok := args.Get(0).(*domain.Auth); ok {
		return createdUser, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepo) FindByUsername(ctx context.Context, username string) (*domain.Auth, error) {
	args := m.Called(username)
	if user, ok := args.Get(0).(*domain.Auth); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepo) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockAuthRepo) GetRefreshToken(ctx context.Context, tokenStr string) (*domain.RefreshToken, error) {
	args := m.Called(tokenStr)
	if refreshToken, ok := args.Get(0).(*domain.RefreshToken); ok {
		return refreshToken, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepo) RevokeRefreshToken(ctx context.Context, tokenStr string) error {
	args := m.Called(ctx, tokenStr)
	return args.Error(0)
}

func (m *MockAuthRepo) DeleteExpiredRefreshTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestRegisterUser_Success(t *testing.T) {
	root, _ := os.Getwd()
	if filepath.Base(root) == "application" {
		root = filepath.Dir(filepath.Dir(filepath.Dir(root)))
	}
	os.Setenv("KEY_PATH", filepath.Join(root, "pkg/token/keys"))

	err := token.InitKeys()
	require.NoError(t, err)

	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	validID := uuid.New().String()

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.ErrUserNotFound)
	mockRepo.On("RegisterUser", mock.AnythingOfType("*domain.Auth")).Return(&domain.Auth{
		ID:       validID,
		Username: "testuser",
	}, nil)
	mockRepo.On("CreateRefreshToken",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
	).Return(nil)

	user, accessToken, refreshToken, err := service.RegisterUser(context.Background(), "testuser", "testpassword")

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, validID, user.ID)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	existingUser := &domain.Auth{Username: "testuser"}
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(existingUser, nil)

	user, accessToken, refreshToken, err := service.RegisterUser(
		context.Background(), "testuser", "testpassword",
	)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	mockRepo.AssertNotCalled(t, "RegisterUser", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_WithError(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.ErrInternal)

	user, accessToken, refreshToken, err := service.RegisterUser(context.Background(), "testuser", "testpassword")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	mockRepo.AssertNotCalled(t, "RegisterUser", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "nouser").Return(nil, errors.ErrUserNotFound)

	accessToken, refresToken, err := service.LoginUser(context.Background(), "nouser", "irrelevant")
	assert.Error(t, err)
	assert.Equal(t, "", accessToken)
	assert.Equal(t, "", refresToken)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := &domain.Auth{ID: "testid", Username: "testuser", Password: string(hashed)}

	mockRepo.On("FindByUsername", "testuser").Return(user, nil)

	accessToken, refreshToken, err := service.LoginUser(context.Background(), "testuser", "wrongpass")
	assert.Error(t, err)
	assert.Equal(t, "", accessToken)
	assert.Equal(t, "", refreshToken)
	assert.Equal(t, errors.ErrInvalidPassword, err)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	root, _ := os.Getwd()
	if filepath.Base(root) == "application" {
		root = filepath.Dir(filepath.Dir(filepath.Dir(root)))
	}
	os.Setenv("KEY_PATH", filepath.Join(root, "pkg/token/keys"))
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	err := token.InitKeys()
	require.NoError(t, err)

	validID := uuid.New().String()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("secretpass"), bcrypt.DefaultCost)
	user := &domain.Auth{ID: validID, Username: "testuser", Password: string(hashed)}

	mockRepo.On("FindByUsername", "testuser").Return(user, nil)
	mockRepo.On("CreateRefreshToken",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
	).Return(nil)
	mockRepo.On("DeleteExpiredRefreshTokens", mock.Anything).Return(nil)

	accessToken, refreshToken, err := service.LoginUser(context.Background(), "testuser", "secretpass")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	mockRepo.AssertExpectations(t)
}
