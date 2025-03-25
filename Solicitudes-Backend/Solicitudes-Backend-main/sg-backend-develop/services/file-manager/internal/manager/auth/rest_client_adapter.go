package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/teamcubation/sg-file-manager-api/internal/platform/restclient"
)

type restClient struct {
	client *restclient.Client
}

func NewRestClient(client *restclient.Client) Client {
	return &restClient{
		client: client,
	}
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

func (rc *restClient) Login(ctx context.Context, credentials Credentials) (*TokenResponse, error) {
	loginReq := LoginRequest(credentials)

	params := restclient.RequestParams{
		Headers: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: loginReq,
	}

	request := rc.client.CreateRequest(params)
	request.SetContext(ctx)

	var loginResp LoginResponse
	var authError ErrorResponse

	resp, err := request.
		SetResult(&loginResp).
		SetError(&authError).
		Post("/login")

	if err != nil {
		return nil, fmt.Errorf("error making login request: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == http.StatusUnauthorized {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("login failed: %s", authError.Error)
	}

	return (*TokenResponse)(&loginResp), nil
}
