package auth

import (
	"context"
	"fmt"

	sdkjwt "github.com/teamcubation/sg-auth/pkg/jwt/v5"

	ports "github.com/teamcubation/sg-auth/internal/core/ports"
)

type UseCases struct {
	jwtService ports.JwtService
}

func NewUseCases(js ports.JwtService) ports.UseCases {
	return &UseCases{
		jwtService: js,
	}
}

func (u *UseCases) Login(ctx context.Context, cuit string) (*sdkjwt.Token, error) {
	token, err := u.jwtService.GenerateToken(cuit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return the generated token
	return token, nil
}
