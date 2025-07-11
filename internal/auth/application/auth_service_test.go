package application

import (
	"context"
	"testing"

	"github.com/mingtmt/book-store/internal/auth/domain"
	"github.com/mingtmt/book-store/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.ErrUserNotFound)
	mockRepo.On("RegisterUser", mock.AnythingOfType("*domain.Auth")).Return(&domain.Auth{
		ID:       "test-id",
		Username: "testuser",
	}, nil)

	user, err := service.RegisterUser(context.Background(), "testuser", "testpassword")

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test-id", user.ID)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	existingUser := &domain.Auth{Username: "testuser"}
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(existingUser, nil)

	user, err := service.RegisterUser(
		context.Background(), "testuser", "testpassword",
	)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertNotCalled(t, "RegisterUser", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_WithError(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.ErrInternal)

	user, err := service.RegisterUser(context.Background(), "testuser", "testpassword")

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertNotCalled(t, "RegisterUser", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_PasswordNotEncrypt(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.ErrUserNotFound)
	mockRepo.On("RegisterUser", mock.AnythingOfType("*domain.Auth")).Return(nil, errors.ErrInternal)

	user, err := service.RegisterUser(context.Background(), "testuser", "testpassword")

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
