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

type MockRefreshTokenRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) RegisterUser(ctx context.Context, user *domain.Auth) (string, error) {
	args := m.Called(user)
	if userID, ok := args.Get(0).(string); ok {
		return userID, args.Error(1)
	}
	return "", args.Error(1)
}

func (m *MockAuthRepo) FindByUsername(ctx context.Context, username string) (*domain.Auth, error) {
	args := m.Called(username)
	if user, ok := args.Get(0).(*domain.Auth); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRefreshTokenRepo) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockRefreshTokenRepo) GetRefreshToken(ctx context.Context, tokenStr string) (*domain.RefreshToken, error) {
	args := m.Called(tokenStr)
	if refreshToken, ok := args.Get(0).(*domain.RefreshToken); ok {
		return refreshToken, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRefreshTokenRepo) RevokeRefreshToken(ctx context.Context, tokenStr string) error {
	args := m.Called(ctx, tokenStr)
	return args.Error(0)
}

func (m *MockRefreshTokenRepo) DeleteExpiredRefreshTokens(ctx context.Context) error {
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

	mockAuthRepo := new(MockAuthRepo)
	mockRfTokenRepo := new(MockRefreshTokenRepo)
	service := NewAuthService(mockAuthRepo, mockRfTokenRepo)

	validUUID := uuid.New()
	validID := validUUID.String()

	mockAuthRepo.On("FindByUsername", "testuser").Return(nil, errors.ErrUserNotFound)

	var registeredUser *domain.Auth
	mockAuthRepo.On("RegisterUser", mock.MatchedBy(func(u *domain.Auth) bool {
		registeredUser = u
		return u.Username == "testuser" && u.ID != "" && u.Password != ""
	})).Return(validID, nil)

	mockRfTokenRepo.On("CreateRefreshToken",
		mock.Anything,
		validUUID,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
	).Return(nil)

	userID, accessToken, refreshToken, err := service.RegisterUser(context.Background(), "testuser", "testpassword")

	assert.NoError(t, err)
	assert.Equal(t, validID, userID)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	require.NotNil(t, registeredUser)
	assert.Equal(t, "testuser", registeredUser.Username)
	assert.NotEqual(t, "testpassword", registeredUser.Password)

	// Password should match bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(registeredUser.Password), []byte("testpassword"))
	assert.NoError(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	mockRfTokenRepo := new(MockRefreshTokenRepo)
	existingUser := &domain.Auth{Username: "testuser"}
	service := NewAuthService(mockAuthRepo, mockRfTokenRepo)

	mockAuthRepo.On("FindByUsername", "testuser").Return(existingUser, nil)

	userID, accessToken, refreshToken, err := service.RegisterUser(
		context.Background(), "testuser", "testpassword",
	)

	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	mockAuthRepo.AssertNotCalled(t, "RegisterUser", mock.Anything)
	mockAuthRepo.AssertExpectations(t)
}

func TestRegisterUser_WithError(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	mockRfTokenRepo := new(MockRefreshTokenRepo)
	service := NewAuthService(mockAuthRepo, mockRfTokenRepo)

	mockAuthRepo.On("FindByUsername", "testuser").Return(nil, errors.ErrInternal)

	userID, accessToken, refreshToken, err := service.RegisterUser(context.Background(), "testuser", "testpassword")

	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	mockAuthRepo.AssertNotCalled(t, "RegisterUser", mock.Anything)
	mockAuthRepo.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	mockRfTokenRepo := new(MockRefreshTokenRepo)
	service := NewAuthService(mockAuthRepo, mockRfTokenRepo)

	mockAuthRepo.On("FindByUsername", "nouser").Return(nil, errors.ErrUserNotFound)

	accessToken, refresToken, err := service.LoginUser(context.Background(), "nouser", "irrelevant")
	assert.Error(t, err)
	assert.Equal(t, "", accessToken)
	assert.Equal(t, "", refresToken)
	mockAuthRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	mockRfTokenRepo := new(MockRefreshTokenRepo)
	service := NewAuthService(mockAuthRepo, mockRfTokenRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := &domain.Auth{ID: "testid", Username: "testuser", Password: string(hashed)}

	mockAuthRepo.On("FindByUsername", "testuser").Return(user, nil)

	accessToken, refreshToken, err := service.LoginUser(context.Background(), "testuser", "wrongpass")
	assert.Error(t, err)
	assert.Equal(t, "", accessToken)
	assert.Equal(t, "", refreshToken)
	assert.Equal(t, errors.ErrInvalidPassword, err)
	mockAuthRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	root, _ := os.Getwd()
	if filepath.Base(root) == "application" {
		root = filepath.Dir(filepath.Dir(filepath.Dir(root)))
	}
	os.Setenv("KEY_PATH", filepath.Join(root, "pkg/token/keys"))
	mockAuthRepo := new(MockAuthRepo)
	mockRfTokenRepo := new(MockRefreshTokenRepo)
	service := NewAuthService(mockAuthRepo, mockRfTokenRepo)

	err := token.InitKeys()
	require.NoError(t, err)

	validID := uuid.New().String()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("secretpass"), bcrypt.DefaultCost)
	user := &domain.Auth{ID: validID, Username: "testuser", Password: string(hashed)}

	mockAuthRepo.On("FindByUsername", "testuser").Return(user, nil)
	mockRfTokenRepo.On("CreateRefreshToken",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
	).Return(nil)
	mockRfTokenRepo.On("DeleteExpiredRefreshTokens", mock.Anything).Return(nil)

	accessToken, refreshToken, err := service.LoginUser(context.Background(), "testuser", "secretpass")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	mockAuthRepo.AssertExpectations(t)
}
