package sdkmwr

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Constantes del paquete
const (
	// Claves de contexto
	DefaultContextKey = "token"
	ClaimsSuffix      = "_claims"

	// Headers
	authHeaderName = "Authorization"
	bearerPrefix   = "Bearer "

	// Mensajes de error
	errMissingAuthHeader       = "authorization header required"
	errInvalidSigningMethod    = "unexpected signing method"
	errBearerPrefixRequired    = "authorization header must start with Bearer"
	errInvalidToken            = "invalid token"
	errExpiredToken            = "token has expired"
	errTokenNotFound           = "token not found in context"
	errClaimNotFound           = "claim not found in token"
	errInvalidClaimType        = "invalid claim type"
	errInvalidTokenLookup      = "invalid token lookup config"
	errUnsupportedLookupMethod = "unsupported token lookup method"
)

// Config permite configurar el comportamiento del middleware JWT
type Config struct {
	SecretKey    string // Clave secreta para HMAC
	PublicKeyPEM string // Clave pública para RSA
	TokenLookup  string // Formato: "header:Authorization" o "query:token"
	TokenPrefix  string // Ejemplo: "Bearer "
	ContextKey   string // Clave para almacenar el token en el contexto
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() Config {
	return Config{
		TokenLookup: "header:" + authHeaderName,
		TokenPrefix: bearerPrefix,
		ContextKey:  DefaultContextKey,
	}
}

// Validate middleware para validación de JWT
func Validate(config Config) gin.HandlerFunc {
	if config.TokenLookup == "" {
		config = DefaultConfig()
	}

	var rsaPublicKey *rsa.PublicKey
	if config.PublicKeyPEM != "" {
		key, err := parseRSAPublicKey(config.PublicKeyPEM)
		if err != nil {
			panic(fmt.Sprintf("failed to parse RSA public key: %v", err))
		}
		rsaPublicKey = key
	}

	return func(c *gin.Context) {
		token, err := extractToken(c, config)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		// Parsear sin validar para determinar el algoritmo
		unverifiedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, fmt.Sprintf("%s: %v", errInvalidToken, err))
			return
		}

		keyFunc := selectKeyFunc(unverifiedToken, config.SecretKey, rsaPublicKey)
		if keyFunc == nil {
			abortWithError(c, http.StatusUnauthorized, errInvalidSigningMethod)
			return
		}

		// Validar el token
		parsedToken, err := jwt.Parse(token, keyFunc)
		if err != nil || !parsedToken.Valid {
			abortWithError(c, http.StatusUnauthorized, fmt.Sprintf("%s: %v", errInvalidToken, err))
			return
		}

		// Guardar token y claims en el contexto
		contextKey := config.ContextKey
		if contextKey == "" {
			contextKey = DefaultContextKey
		}

		c.Set(contextKey, parsedToken)
		c.Set(GetClaimsKey(contextKey), parsedToken.Claims)

		c.Next()
	}
}

// ExtractClaim extrae un claim específico del token JWT
func ExtractClaim(c *gin.Context, claimKey, contextKey string) (string, error) {
	if contextKey == "" {
		contextKey = DefaultContextKey
	}

	tokenInterface, exists := c.Get(contextKey)
	if !exists {
		return "", fmt.Errorf(errTokenNotFound)
	}

	token, ok := tokenInterface.(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("invalid token type in context")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims type")
	}

	claim, exists := claims[claimKey]
	if !exists {
		return "", fmt.Errorf("%s: %s", errClaimNotFound, claimKey)
	}

	return formatClaimValue(claim)
}

// ExtractUserID helper para extraer el user_id del token
func ExtractUserID(c *gin.Context) (string, error) {
	return ExtractClaim(c, "user_id", "")
}

// GetClaimsKey retorna la clave de contexto para los claims
func GetClaimsKey(tokenKey string) string {
	if tokenKey == "" {
		tokenKey = DefaultContextKey
	}
	return tokenKey + ClaimsSuffix
}

// Funciones auxiliares privadas

func extractToken(c *gin.Context, config Config) (string, error) {
	parts := strings.Split(config.TokenLookup, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf(errInvalidTokenLookup)
	}

	switch parts[0] {
	case "header":
		return extractFromHeader(c, parts[1], config.TokenPrefix)
	case "query":
		return extractFromQuery(c, parts[1])
	default:
		return "", fmt.Errorf(errUnsupportedLookupMethod)
	}
}

func extractFromHeader(c *gin.Context, header, prefix string) (string, error) {
	auth := c.GetHeader(header)
	if auth == "" {
		return "", fmt.Errorf(errMissingAuthHeader)
	}
	if prefix != "" && !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf(errBearerPrefixRequired)
	}
	return strings.TrimPrefix(auth, prefix), nil
}

func extractFromQuery(c *gin.Context, param string) (string, error) {
	token := c.Query(param)
	if token == "" {
		return "", fmt.Errorf(errMissingAuthHeader)
	}
	return token, nil
}

func selectKeyFunc(token *jwt.Token, secretKey string, rsaKey *rsa.PublicKey) jwt.Keyfunc {
	switch token.Method.(type) {
	case *jwt.SigningMethodHMAC:
		if secretKey == "" {
			return nil
		}
		return func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		}
	case *jwt.SigningMethodRSA:
		if rsaKey == nil {
			return nil
		}
		return func(token *jwt.Token) (interface{}, error) {
			return rsaKey, nil
		}
	default:
		return nil
	}
}

func parseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaKey, nil
}

func formatClaimValue(v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case float64:
		return fmt.Sprintf("%.0f", val), nil
	case int:
		return fmt.Sprintf("%d", val), nil
	case int64:
		return fmt.Sprintf("%d", val), nil
	case bool:
		return fmt.Sprintf("%v", val), nil
	case nil:
		return "", fmt.Errorf("%s: claim is nil", errInvalidClaimType)
	default:
		return fmt.Sprintf("%v", val), nil
	}
}

func abortWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
