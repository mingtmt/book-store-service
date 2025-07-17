package domain

import "context"

type Auth struct {
	ID       string
	Username string
	Password string
}

type AuthRepository interface {
	RegisterUser(ctx context.Context, user *Auth) (string, error)
	FindByUsername(ctx context.Context, username string) (*Auth, error)
}
