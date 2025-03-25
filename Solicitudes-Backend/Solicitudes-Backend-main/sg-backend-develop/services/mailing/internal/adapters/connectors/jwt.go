package mailconn

import (
	"context"
	"fmt"

	sdkjwt "github.com/teamcubation/sg-mailing-api/pkg/jwt/v5"
	sdkjwtdefs "github.com/teamcubation/sg-mailing-api/pkg/jwt/v5/defs"

	ports "github.com/teamcubation/sg-mailing-api/internal/core/ports"
)

type JwtService struct {
	jwtService sdkjwtdefs.Service
}

func NewJwtService() (ports.JwtService, error) {
	js, err := sdkjwt.Bootstrap("JWT_SECRET", "VERIFICATION_EXPIRATION_MINUTES")
	if err != nil {
		return nil, fmt.Errorf("jwt bootstrap error: %w", err)
	}

	return &JwtService{
		jwtService: js,
	}, nil
}

func (j *JwtService) GenerateToken(ctx context.Context, email string) (string, error) {
	token, err := j.jwtService.GenerateToken(ctx, email)
	if err != nil {
		return "", fmt.Errorf("error trying to generate token: %w", err)
	}

	return token, nil
}

func (j *JwtService) ValidateToken(ctx context.Context, token string) (*sdkjwtdefs.TokenClaims, error) {
	tokenClaims, err := j.jwtService.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("error trying to validate token: %w", err)
	}

	return tokenClaims, nil
}

func (j *JwtService) GetTokenInfo(ctx context.Context, token string) (*sdkjwtdefs.TokenClaims, error) {
	tokenClaims, err := j.jwtService.ValidateTokenAllowExpired(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("error trying to validate token: %w", err)
	}

	return tokenClaims, nil
}
