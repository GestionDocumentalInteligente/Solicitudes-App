package restclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	AuthorizationHeader   = "Authorization"
	ContentTypeHeader     = "Content-Type"
	ContentTypeJSONHeader = "application/json; charset=utf-8"
	DefaultTimeout        = 360 * time.Second
	DefaultMaxRetries     = 3
	DefaultWaitTime       = 100 * time.Millisecond
	DefaultMaxWaitTime    = 2 * time.Second
)

type Client struct {
	*resty.Client
	credentials Credentials
	token       string
	mutex       sync.RWMutex
}

type Credentials struct {
	Email    string
	Password string
}

type ClientOption func(*Client)

func NewHTTPClient(credentials Credentials, options ...ClientOption) *Client {
	client := &Client{
		Client:      resty.New(),
		credentials: credentials,
	}

	client.
		SetTimeout(DefaultTimeout).
		SetRetryCount(DefaultMaxRetries).
		SetRetryWaitTime(DefaultWaitTime).
		SetRetryMaxWaitTime(DefaultMaxWaitTime).
		SetHeader(ContentTypeHeader, ContentTypeJSONHeader)

	for _, option := range options {
		option(client)
	}

	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		return err != nil || response.StatusCode() == http.StatusUnauthorized
	})

	return client
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.SetBaseURL(baseURL)
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.SetTimeout(timeout)
	}
}

func WithRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.SetRetryCount(maxRetries)
	}
}

func WithRetryCondition(condition func(*resty.Response, error) bool) ClientOption {
	return func(c *Client) {
		c.AddRetryCondition(condition)
	}
}

type RequestParams struct {
	Headers    http.Header
	Query      url.Values
	PathParams map[string]string
	Body       interface{}
}

func (c *Client) SetToken(tk string) {
	c.token = tk
}

func (c *Client) CreateRequest(params RequestParams) *resty.Request {
	request := c.R()

	c.mutex.RLock()
	if c.token != "" {
		request.SetHeader(AuthorizationHeader, fmt.Sprintf("Bearer %s", c.token))
	}
	c.mutex.RUnlock()

	if params.Body != nil {
		request.SetBody(params.Body)
	}

	for key, values := range params.Headers {
		for _, value := range values {
			request.SetHeader(key, value)
		}
	}

	if params.Query != nil {
		request.SetQueryParamsFromValues(params.Query)
	}

	if params.PathParams != nil {
		request.SetPathParams(params.PathParams)
	}

	return request
}

func (c *Client) ExecuteWithAuth(ctx context.Context, req func() (*resty.Response, error)) (*resty.Response, error) {
	resp, err := req()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		token, loginErr := c.Login(ctx)
		if loginErr != nil {
			return nil, fmt.Errorf("authentication failed: %w", loginErr)
		}

		c.mutex.Lock()
		c.token = token
		c.mutex.Unlock()

		return req()
	}

	return resp, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (c *Client) Login(ctx context.Context) (string, error) {
	loginReq := LoginRequest{
		Email:    c.credentials.Email,
		Password: c.credentials.Password,
	}

	params := RequestParams{
		Headers: http.Header{
			ContentTypeHeader: []string{ContentTypeJSONHeader},
		},
		Body: loginReq,
	}

	var loginResp LoginResponse
	var authError ErrorResponse

	resp, err := c.CreateRequest(params).
		SetContext(ctx).
		SetResult(&loginResp).
		SetError(&authError).
		Post("/login")
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("login failed: %s", authError.Error)
	}

	return loginResp.Token, nil
}
