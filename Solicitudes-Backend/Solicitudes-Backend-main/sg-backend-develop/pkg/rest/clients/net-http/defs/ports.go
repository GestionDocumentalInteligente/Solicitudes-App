package defs

import (
	"context"
	"net/http"
	"net/url"
)

// TokenResponse representa la respuesta de un token de acceso
type TokenResponse interface {
	GetAccessToken() string
}

// Client define la interfaz para un cliente OAuth genérico
type Client interface {
	GetAccessToken(ctx context.Context, endpoint string, params url.Values) (TokenResponse, error)
	Do(req *http.Request) (*http.Response, error)
	AddInterceptor(interceptor Interceptor)
}

type Interceptor interface {
	Before(req *http.Request) (*http.Request, error)
	After(resp *http.Response, err error) (*http.Response, error)
}

// Config define la interfaz para la configuración del cliente
type Config interface {
	GetTokenEndpoint() string
	GetClientID() string
	GetClientSecret() string
	GetAdditionalParams() url.Values
	Validate() error
}

