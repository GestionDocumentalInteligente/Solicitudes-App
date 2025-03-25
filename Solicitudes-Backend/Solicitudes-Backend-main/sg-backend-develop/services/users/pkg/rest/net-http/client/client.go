package sdkhclnt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/teamcubation/sg-users/pkg/rest/net-http/client/defs"
)

type client struct {
	config       defs.Config
	httpClient   *http.Client
	interceptors []defs.Interceptor
}

func newClient(config defs.Config) (defs.Client, error) {
	c := &client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Timeout predeterminado
		},
	}

	return c, nil
}

func (c *client) GetAccessToken(ctx context.Context, endpoint string, params url.Values) (defs.TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resp *http.Response
	err = c.retryWithBackoff(func() error {
		var err error
		resp, err = c.do(req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error obteniendo el token de acceso: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error obteniendo el token de acceso: código de estado %d", resp.StatusCode)
	}

	var tokenRes map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, fmt.Errorf("error decodificando la respuesta: %w", err)
	}

	return &defs.GenericTokenResponse{tokenRes}, nil
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	for _, interceptor := range c.interceptors {
		var err error
		req, err = interceptor.Before(req)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.Do(req)

	if err == nil {
		for i := len(c.interceptors) - 1; i >= 0; i-- {
			resp, err = c.interceptors[i].After(resp, err)
			if err != nil {
				return resp, err
			}
		}
	}

	return resp, err
}

func (c *client) retryWithBackoff(operation func() error) error {
	backoff := 100 * time.Millisecond
	for i := 0; i < 3; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return fmt.Errorf("operación fallida después de 3 intentos")
}

func (c *client) AddInterceptor(interceptor defs.Interceptor) {
	c.interceptors = append(c.interceptors, interceptor)
}
