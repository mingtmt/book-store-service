package domain

import "time"

type Auth struct {
	ID       string
	Username string
	Password string
}

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}
