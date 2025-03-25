package authconn

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	sdkjwt "github.com/teamcubation/sg-auth/pkg/jwt/v5"
	sdkports "github.com/teamcubation/sg-auth/pkg/jwt/v5/ports"

	ports "github.com/teamcubation/sg-auth/internal/core/ports"
)

type JwtService struct {
	JwtService sdkports.Service
}

func NewJwtService() (ports.JwtService, error) {
	js, err := sdkjwt.Bootstrap()
	if err != nil {
		return nil, fmt.Errorf("bootstrap error: %w", err)
	}

	return &JwtService{
		JwtService: js,
	}, nil
}

func (j *JwtService) GenerateToken(cuit string) (*sdkjwt.Token, error) {
	claims := map[string]interface{}{
		"sub": cuit,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	// Llama al SDK para generar el token
	tokenString, err := j.JwtService.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// Crea la entidad Token con el token generado
	token := &sdkjwt.Token{
		AccessToken: tokenString,
		ExpiresAt:   time.Now().Add(time.Hour * 24),
	}

	return token, nil
}

func (j *JwtService) ValidateToken(tokenString string) (*sdkjwt.TokenClaims, error) {
	// Llama al SDK para validar el token
	token, err := j.JwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("error al validar el token: %w", err)
	}

	// Verifica que el token sea válido y extrae las claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	// Extrae las claims necesarias
	subject, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("el token no contiene 'sub'")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("el token no contiene 'exp'")
	}
	expiresAt := time.Unix(int64(expFloat), 0)

	iatFloat, ok := claims["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("el token no contiene 'iat'")
	}
	issuedAt := time.Unix(int64(iatFloat), 0)

	// Crea una estructura TokenClaims con la información extraída
	tokenClaims := &sdkjwt.TokenClaims{
		Subject:   subject,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}

	return tokenClaims, nil
}
