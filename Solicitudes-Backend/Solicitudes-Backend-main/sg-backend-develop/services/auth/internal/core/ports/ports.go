package ports

import (
	"context"

	"github.com/gin-gonic/gin"

	sdkjwt "github.com/teamcubation/sg-auth/pkg/jwt/v5"

	entities "github.com/teamcubation/sg-auth/internal/core/entities"
)

// UseCases define las operaciones de casos de uso para autenticación
type UseCases interface {
	Login(ctx context.Context, cuit string) (*sdkjwt.Token, error)
}

// JwtService define las operaciones del servicio JWT
type JwtService interface {
	GenerateToken(string) (*sdkjwt.Token, error)
}

type Repository interface {
	CreateEvent(context.Context) (*entities.Auth, error)
}

type HttpClient interface {
	GetAccessToken(ctx context.Context) (string, error)
	// AuthenticateUser(ctx context.Context, credentials entities.Credentials) (*entities.User, error)
	// GetUserInfo(ctx context.Context, token string) (*entities.User, error)
	// RevokeToken(ctx context.Context, token string) error
	// RefreshToken(ctx context.Context, refreshToken string) (*entities.Token, error)
}

type SessionManager interface {
	JwtToSession(*gin.Context, string, string) error
}
