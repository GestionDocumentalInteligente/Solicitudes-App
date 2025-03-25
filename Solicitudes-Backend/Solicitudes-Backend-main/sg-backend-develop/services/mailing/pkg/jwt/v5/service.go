package sdkjwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/teamcubation/sg-mailing-api/pkg/jwt/v5/defs"
)

type service struct {
	secret     []byte
	expiration time.Duration
}

func newService(c defs.Config) (defs.Service, error) {
	return &service{
		secret:     []byte(c.GetSecretKey()),
		expiration: c.GetExpirationMinutes(),
	}, nil
}

func (s *service) GenerateToken(ctx context.Context, subject string) (string, error) {
	claims := defs.Claims{
		Subject: subject,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("error signing the token: %w", err)
	}
	return signedToken, nil
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
