package auth

import (
	"context"
)

type Credentials struct {
	Email    string
	Password string
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Client interface {
	Login(ctx context.Context, credentials Credentials) (*TokenResponse, error)
}
