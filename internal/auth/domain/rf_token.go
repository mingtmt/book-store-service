package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenStr string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenStr string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
}
