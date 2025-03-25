package sdkjwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/teamcubation/sg-backend/pkg/jwt/v5/defs"
)

type service struct {
	config            defs.Config
	secret            []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func newService(c defs.Config) (defs.Service, error) {
	return &service{
		config:            c,
		secret:            []byte(c.GetSecretKey()),
		accessExpiration:  c.GetAccessExpiration(),
		refreshExpiration: c.GetRefreshExpiration(),
	}, nil
}

func (s *service) GenerateTokens(ctx context.Context, subject string) (string, string, error) {
	// Generación del access token con expiración corta
	accessClaims := defs.Claims{
		Subject: subject,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString(s.secret)
	if err != nil {
		return "", "", fmt.Errorf("error signing the access token: %w", err)
	}

	// Generación del refresh token con expiración más larga
	refreshClaims := defs.Claims{
		Subject: subject,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return "", "", fmt.Errorf("error signing the refresh token: %w", err)
	}

	return signedAccessToken, signedRefreshToken, nil
}

func (s *service) ValidateToken(ctx context.Context, tokenString string) (*defs.TokenClaims, error) {
	claims := &defs.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error validating the token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	tokenClaims := &defs.TokenClaims{
		Subject:   claims.Subject,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}

	return tokenClaims, nil
}

func (c *service) GetAccessExpiration() time.Duration {
	return c.config.GetAccessExpiration()
}

// GetRefreshExpiration devuelve la duración de expiración del token de refresco
func (c *service) GetRefreshExpiration() time.Duration {
	return c.config.GetRefreshExpiration()
}

func (s *service) ValidateTokenAllowExpired(ctx context.Context, tokenString string) (*defs.TokenClaims, error) {
	claims := &defs.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		// Check if the error is due to expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Token is expired but otherwise valid; proceed to extract claims
			return &defs.TokenClaims{
				Subject:   claims.Subject,
				ExpiresAt: claims.ExpiresAt.Time,
				IssuedAt:  claims.IssuedAt.Time,
			}, nil
		}
		// Other errors related to validation
		return nil, fmt.Errorf("error validating the token: %w", err)
	}

	// Check if the token is valid, even if not expired
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &defs.TokenClaims{
		Subject:   claims.Subject,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}, nil
}
