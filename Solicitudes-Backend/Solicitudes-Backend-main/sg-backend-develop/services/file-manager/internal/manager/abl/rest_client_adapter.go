package abl

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/teamcubation/sg-file-manager-api/internal/platform/restclient"
	"github.com/teamcubation/sg-file-manager-api/pkg/log"
)

type restClient struct {
	client *restclient.Client
}

func NewRestClient(client *restclient.Client) Client {
	return &restClient{
		client: client,
	}
}

func generateRequestParams(body interface{}) restclient.RequestParams {
	return restclient.RequestParams{
		Headers: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: body,
	}
}

type ablRequest struct {
	Number string `json:"cuenta"`
	Type   int    `json:"tipo"`
}

type response struct {
	Owner   string `json:"Titular"`
	Address string `json:"calle"`
	Number  string `json:"Altura"`
}

func (rc *restClient) ValidateABLData(ctx context.Context, ablNumber string, ablType int) (bool, error) {
	logger := log.FromContext(ctx)

	request := ablRequest{
		Number: ablNumber,
		Type:   ablType,
	}

	var response []response

	resp, err := rc.client.ExecuteWithAuth(ctx, func() (*resty.Response, error) {
		return rc.client.CreateRequest(generateRequestParams(request)).
			SetContext(ctx).
			SetResult(&response).
			Post("/consultaDeudaInicial")
	})
	if err != nil {
		logger.Error(err.Error())
		return false, fmt.Errorf("error: %w", err)
	}
	if resp.IsError() {
		logger.Error(resp.String())
		return false, fmt.Errorf("error: %s", resp.String())
	}

	return len(response) > 0, nil
}
