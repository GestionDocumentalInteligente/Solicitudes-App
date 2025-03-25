package sdkjwt

import (
	"fmt"
	"sync"

	"github.com/golang-jwt/jwt/v5"

	"github.com/teamcubation/sg-auth/pkg/jwt/v5/ports"
)

var (
	instance  ports.Service
	once      sync.Once
	initError error
)

type service struct {
	secretKey string
}

// newService inicializa el servicio JWT con la configuración proporcionada
func newService(cfg ports.Config) (ports.Service, error) {
	once.Do(func() {
		if err := cfg.Validate(); err != nil {
			initError = err
			return
		}
		instance = &service{
			secretKey: cfg.GetSecretKey(),
		}
	})
	if initError != nil {
		return nil, initError
	}

	return instance, nil
}

// GenerateToken genera un token JWT con las reclamaciones proporcionadas
func (j *service) GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateToken valida un token JWT proporcionado
func (j *service) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
