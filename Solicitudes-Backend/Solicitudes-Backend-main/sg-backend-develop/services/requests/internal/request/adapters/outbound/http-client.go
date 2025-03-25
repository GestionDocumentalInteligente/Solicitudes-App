package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	sdkhc "github.com/teamcubation/sg-backend/pkg/rest/clients/net-http"
	sdkdefs "github.com/teamcubation/sg-backend/pkg/rest/clients/net-http/defs"

	transport "github.com/teamcubation/sg-backend/services/requests/internal/request/adapters/outbound/transport"
	domain "github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
	ports "github.com/teamcubation/sg-backend/services/requests/internal/request/core/ports"
)

type HttpClient struct {
	client sdkdefs.Client
}

func NewHttpClient() (ports.HttpClient, error) {
	c, err := sdkhc.Bootstrap("", "", "", nil)
	if err != nil {
		return nil, fmt.Errorf("bootstrap error: %w", err)
	}

	return &HttpClient{
		client: c,
	}, nil
}

func (h *HttpClient) SendCreatedRequest(ctx context.Context, req *domain.Request) (string, error) {
	jsonBody, err := json.Marshal(transport.ToRequestPayload(req))
	if err != nil {
		return "", fmt.Errorf("error marshaling json: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/file-manager/record", os.Getenv("FILE_MANAGER_HOST")),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	code, ok := result["code"].(string)
	if !ok {
		if result["code"] == nil {
			return "", fmt.Errorf("response code is nil")
		}
		return "", fmt.Errorf("response code is not a string: %v", result["code"])
	}

	return code, nil
}

func (h *HttpClient) SendUpdateRequest(ctx context.Context, req *domain.Request) error {
	jsonBody, err := json.Marshal(transport.ToRequestPayload(req))
	if err != nil {
		return fmt.Errorf("error marshaling json: %w", err)
	}

	httpReq, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/api/v1/file-manager/record/%s/documents", os.Getenv("FILE_MANAGER_HOST"), req.FileNumber),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

type docContent struct {
	Content   string `json:"content"`
	Reference string `json:"reference"`
	Type      int    `json:"type"`
	Username  string `json:"username"`
}

func (h *HttpClient) SendVerificationDocument(ctx context.Context, content, code, reference, verificationType, name string) error {
	docType := 16
	if verificationType == "property" {
		docType = 17
	}

	jsonBody, err := json.Marshal(docContent{
		Content:   content,
		Reference: reference,
		Type:      docType,
		Username:  name,
	})
	if err != nil {
		return fmt.Errorf("error marshaling json: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/file-manager/record/%s/documents", os.Getenv("FILE_MANAGER_HOST"), code),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err)
	}

	return nil
}

func (h *HttpClient) SendValidationDocument(ctx context.Context, code, content, name string) error {
	jsonBody, err := json.Marshal(docContent{
		Content:   content,
		Reference: "Documento de Validación",
		Type:      18,
		Username:  name,
	})
	if err != nil {
		return fmt.Errorf("error marshaling json: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/file-manager/record/%s/documents", os.Getenv("FILE_MANAGER_HOST"), code),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err)
	}

	return nil
}

func (h *HttpClient) SendEmail(ctx context.Context, code, email string) error {
	req := transport.EmailRequestPayload{
		Code:  code,
		Email: email,
	}

	jsonBody, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("error marshaling json: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/mailing/new-request", os.Getenv("MAILING_HOST")),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

type EmailRequestPayload struct {
	Code         string `json:"code" binding:"required"`
	Email        string `json:"email"`
	Observations string `json:"observations"`
}

func (h *HttpClient) sendEmail(ctx context.Context, endpoint, code, email, observations string) error {
	req := EmailRequestPayload{
		Code:         code,
		Email:        email,
		Observations: observations,
	}

	jsonBody, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("error marshaling json: %w", err)
	}

	url := fmt.Sprintf("%s%s", os.Getenv("MAILING_HOST"), endpoint)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (h *HttpClient) SendEmailUpdateRequest(ctx context.Context, code, email string) error {
	return h.sendEmail(ctx, "/api/v1/mailing/update-request", code, email, "")
}

func (h *HttpClient) SendEmailUpdate(ctx context.Context, code, email, observations string) error {
	return h.sendEmail(ctx, "/api/v1/mailing/update-request-code", code, email, observations)
}

func (h *HttpClient) SendEmailValidateRequest(ctx context.Context, code, email string) error {
	return h.sendEmail(ctx, "/api/v1/mailing/validate-request", code, email, "")
}
