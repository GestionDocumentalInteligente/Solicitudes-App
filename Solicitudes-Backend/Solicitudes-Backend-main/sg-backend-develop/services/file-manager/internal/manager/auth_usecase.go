package manager

import (
	"context"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/auth"
)

type authUseCase struct {
	client auth.Client
}

type AuthUseCase interface {
	Login(context.Context, auth.Credentials) (auth.TokenResponse, error)
}

func NewAuthUseCase(client auth.Client) AuthUseCase {
	return &authUseCase{
		client: client,
	}
}

func (a *authUseCase) Login(ctx context.Context, credentials auth.Credentials) (auth.TokenResponse, error) {
	resp, err := a.client.Login(ctx, credentials)
	return *resp, err
}
