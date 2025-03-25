package sdkmwr

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authHeaderName             = "Authorization"
	bearerPrefix               = "Bearer "
	errMissingAuthHeader       = "authorization header required"
	errInvalidSigningMethod    = "unexpected signing method"
	errBearerPrefixRequired    = "authorization header must start with Bearer"
	errInvalidToken            = "invalid token"
	errExpiredToken            = "token has expired"
	errInsufficientPermissions = "insufficient permissions"
)

func formatPublicKey(key string) string {
	// Añadir encabezado y pie de clave pública PEM
	var formattedKey strings.Builder
	formattedKey.WriteString("-----BEGIN PUBLIC KEY-----\n")

	// Insertar saltos de línea cada 64 caracteres
	for i := 0; i < len(key); i += 64 {
		end := i + 64
		if end > len(key) {
			end = len(key)
		}
		formattedKey.WriteString(key[i:end] + "\n")
	}

	formattedKey.WriteString("-----END PUBLIC KEY-----")
	return formattedKey.String()
}

func ValidateJwt(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Obtener el token desde el encabezado Authorization
		authHeader := c.GetHeader(authHeaderName)

		if authHeader != "" {
			// Verificar que el encabezado empiece con "Bearer"
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				abortWithError(c, http.StatusUnauthorized, errBearerPrefixRequired)
				return
			}
			// Extraer el token quitando el prefijo "Bearer "
			tokenString = strings.TrimPrefix(authHeader, bearerPrefix)
		} else {
			// 2. Si no se encuentra en el encabezado Authorization, buscar en la query string
			tokenString = c.Query("token")
			if tokenString == "" {
				abortWithError(c, http.StatusUnauthorized, errMissingAuthHeader)
				return
			}
		}

		// 3. Parsear el token sin validar para obtener el algoritmo
		unverifiedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, &jwt.RegisteredClaims{})
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, errInvalidToken+": "+err.Error())
			return
		}

		// 4. Determinar el método de firma
		var keyFunc jwt.Keyfunc

		switch unverifiedToken.Method.(type) {
		case *jwt.SigningMethodHMAC:
			// Método de firma HMAC
			if secretKey == "" {
				abortWithError(c, http.StatusUnauthorized, "secret key is not provided")
				return
			}
			keyFunc = func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			}
		case *jwt.SigningMethodRSA:
			// Método de firma RSA
			provider := c.Query("provider")
			rsaPublicKey, err := getRSAPublicKey(provider)
			if err != nil || rsaPublicKey == nil {
				abortWithError(c, http.StatusUnauthorized, "RSA public key is not provided")
				return
			}
			keyFunc = func(token *jwt.Token) (interface{}, error) {
				return rsaPublicKey, nil
			}
		default:
			// Método de firma no soportado
			abortWithError(c, http.StatusUnauthorized, errInvalidSigningMethod)
			return
		}

		// 5. Validar el token con el keyFunc adecuado
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, keyFunc)
		if err != nil || !token.Valid {
			abortWithError(c, http.StatusUnauthorized, errInvalidToken+": "+err.Error())
			return
		}

		// 6. Verificar si el token ha expirado
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			abortWithError(c, http.StatusUnauthorized, errInvalidToken)
			return
		}

		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			abortWithError(c, http.StatusUnauthorized, errExpiredToken)
			return
		}

		// 7. Guardar el token validado en el contexto para que los handlers lo utilicen
		c.Set("token", token)

		// Continuar con la siguiente función en la cadena de middlewares
		c.Next()
	}
}

func getRSAPublicKey(provider string) (*rsa.PublicKey, error) {
	var rsaPublicKeyPEM string

	switch provider {
	case "AFIP":
		rsaPublicKeyPEM = os.Getenv("RSA_PUBLIC_KEY")
	case "MI_ARGENTINA":
		rsaPublicKeyPEM = os.Getenv("RSA_PUBLIC_KEY_MIARG")
	case "ANSES":
		rsaPublicKeyPEM = os.Getenv("RSA_PUBLIC_KEY_ANSES")
	default:
		return nil, fmt.Errorf("provider not found: %s", provider)
	}

	if rsaPublicKeyPEM == "" {
		return nil, fmt.Errorf("RSA public key not found in environment")
	}

	rsaPublicKeyPEM = formatPublicKey(rsaPublicKeyPEM) // Formatear la clave pública
	var rsaPublicKey *rsa.PublicKey

	block, _ := pem.Decode([]byte(rsaPublicKeyPEM))
	if block == nil || (block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY") {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: " + err.Error())
	}
	var ok bool
	rsaPublicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPublicKey, nil
}

// abortWithError centraliza la lógica para abortar con un error
func abortWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
