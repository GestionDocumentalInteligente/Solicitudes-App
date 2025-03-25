package file

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
	"github.com/teamcubation/sg-file-manager-api/internal/platform/env"
	"github.com/teamcubation/sg-file-manager-api/internal/platform/restclient"
	"github.com/teamcubation/sg-file-manager-api/pkg/log"
)

type restClient struct {
	client *restclient.Client
}

type FileRequests struct {
	Type         string `json:"acronimoTipoDocumento"`
	Reference    string `json:"referencia"`
	Data         string `json:"data"`
	Source       string `json:"sistemaOrigen"`
	Username     string `json:"nombreYApellido"`
	Rol          string `json:"cargo"`
	Distribution string `json:"reparticion"`
}

type FileResponse struct {
	Number        string `json:"numero,omitempty"`
	URL           string `json:"urlArchivoGenerado,omitempty"`
	SpecialNumber string `json:"numeroEspecial,omitempty"`
	Licence       string `json:"licencia,omitempty"`
	Error         string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
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

func (rc *restClient) CreateGEDO(ctx context.Context, documentInfo Document) (SignedDocument, error) {
	logger := log.FromContext(ctx)

	fileRequest := FileRequests{
		Type:         documentInfo.Metadata.DocumentType,
		Reference:    documentInfo.Metadata.Reference,
		Data:         documentInfo.Content,
		Source:       documentInfo.Metadata.OriginSystem,
		Username:     documentInfo.Metadata.FullName,
		Rol:          documentInfo.Metadata.Position,
		Distribution: documentInfo.Metadata.Department,
	}

	url, err := getURLByDocumentType(documentInfo.Metadata.DocumentType)
	if err != nil {
		return SignedDocument{}, BadRequest(fmt.Sprintf("error in document type: %v", err))
	}

	var response FileResponse
	var errorResponse ErrorResponse

	const maxRetries = 5
	const retryDelay = 20 * time.Second
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := rc.client.ExecuteWithAuth(ctx, func() (*resty.Response, error) {
			return rc.client.CreateRequest(generateRequestParams(fileRequest)).
				SetContext(ctx).
				SetResult(&response).
				SetError(&errorResponse).
				Post(fmt.Sprintf("/api/gde/%s", url))
		})
		if err != nil {
			logger.Printf("Error making request %s on attempt %d: %v\n", url, attempt, err)
			if attempt == maxRetries {
				return SignedDocument{
					TypeID:          documentInfo.TypeID,
					Filename:        documentInfo.Name,
					Content:         documentInfo.Content,
					OriginalContent: documentInfo.Content,
					Status:          false,
				}, fmt.Errorf("error making request after %d attempts: %w", maxRetries, err)
			}

			time.Sleep(retryDelay)
			continue
		}

		if resp.IsError() {
			logger.Printf("Request failed: HTTP %d - %s on attempt %d\n", resp.StatusCode(), errorResponse.Error, attempt)

			if attempt == maxRetries {
				return SignedDocument{
					TypeID:          documentInfo.TypeID,
					Filename:        documentInfo.Name,
					Content:         documentInfo.Content,
					OriginalContent: documentInfo.Content,
					Status:          false,
				}, fmt.Errorf("request failed after %d attempts: %s", maxRetries, errorResponse.Error)
			}

			time.Sleep(retryDelay)
			continue
		}

		if response.Number == "" {
			logger.Printf("Attempt %d: lumen response: %s", attempt, resp.Body())

			if attempt == maxRetries {
				return SignedDocument{
					TypeID:          documentInfo.TypeID,
					Filename:        documentInfo.Name,
					Content:         documentInfo.Content,
					OriginalContent: documentInfo.Content,
					Status:          false,
				}, fmt.Errorf("response number is empty after %d attempts", maxRetries)
			}

			time.Sleep(retryDelay)
			continue
		}
		break
	}

	return SignedDocument{
		TypeID:          documentInfo.TypeID,
		Filename:        documentInfo.Name,
		Content:         documentInfo.Content,
		OriginalContent: documentInfo.Content,
		Number:          response.Number,
		URL:             response.URL,
		SpecialNumber:   response.SpecialNumber,
		Licence:         response.Licence,
		Status:          true,
	}, nil
}

func getURLByDocumentType(docType string) (string, error) {
	switch docType {
	case IfGraTypeDocument:
		return "generar-gedo-externo-ifgra", nil
	case IfTypeDocument:
		return "generar-gedo-externo-if?production=0", nil
	default:
		return "", errors.New("invalid type document: must be 'if' or 'ifgra'")
	}
}

type Record struct {
	Department     string `json:"reparticion"`
	Sector         string `json:"sector"`
	System         string `json:"sistema"`
	Reason         string `json:"motivo"`
	Description    string `json:"descripcion"`
	TreatmentCode  string `json:"selectTrataCod"`
	External       bool   `json:"externo"`
	Internal       string `json:"interno"`
	Person         string `json:"persona"`
	DocumentType   string `json:"tipoDoc"`
	DocumentNumber int64  `json:"nroDoc"`
	LastName       string `json:"apellido"`
	FirstName      string `json:"nombre"`
	Address        string `json:"domicilio"`
	Email          string `json:"email"`
	ZipCode        int64  `json:"codigoPostal"`
	Phone          int64  `json:"telefono"`
}

type RecordResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (rc *restClient) CreateRecord(ctx context.Context, user user.User, id int64) (string, error) {
	logger := log.FromContext(ctx)

	recordRequest := Record{
		Department:     "MDSI",
		Sector:         "MESADIGITAL",
		System:         OriginSystem,
		Reason:         fmt.Sprintf("%s %d", env.GetReason(), id),
		Description:    fmt.Sprintf("%s %d", env.GetReason(), id),
		TreatmentCode:  env.GetTreatmentCode(),
		Internal:       "false",
		External:       true,
		Person:         "true",
		DocumentType:   "DU",
		DocumentNumber: user.DocumentNumber,
		LastName:       user.LastName,
		FirstName:      user.FirstName,
		Email:          user.Email,
		Phone:          user.Phone,
		Address:        fmt.Sprintf("%s %s", user.Address.Street, user.Address.Number),
		ZipCode:        user.Address.ZipCode,
	}

	ctx = context.Background()
	var response RecordResponse
	startTime := time.Now()
	resp, err := rc.client.ExecuteWithAuth(ctx, func() (*resty.Response, error) {
		return rc.client.CreateRequest(generateRequestParams(recordRequest)).
			SetContext(ctx).
			SetResult(&response).
			SetError(&response).
			Post("/api/gde/generar-expediente")
	})
	logger.Printf("Request to /api/gde/generar-expediente took %s", time.Since(startTime))
	if err != nil {
		return "", InternalServerError(fmt.Sprintf("error making request: %v", err))
	}
	if resp.IsError() {
		return "", InternalServerError(fmt.Sprintf("request failed: %s", response.Message))
	}

	if response.Error {
		return "", BadRequest(fmt.Sprintf("request failed: %s", response.Message))
	}

	return response.Message, nil
}

func (rc *restClient) DownloadGEDO(ctx context.Context, fileID string) (string, error) {
	const (
		maxRetries = 3
		retryDelay = 2 * time.Second
		endpoint   = "/api/gde/descargar-gedo-pdf"
	)

	var (
		base64Response string
		errorResponse  RecordResponse
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := rc.client.ExecuteWithAuth(ctx, func() (*resty.Response, error) {
			return rc.client.CreateRequest(restclient.RequestParams{
				Query: url.Values{
					"numero":  {fileID},
					"sistema": {OriginSystem},
				},
			}).
				SetContext(ctx).
				SetResult(&base64Response).
				SetError(&errorResponse).
				Get(endpoint)
		})
		if err != nil {
			if attempt == maxRetries {
				return base64Response, InternalServerError(fmt.Sprintf("error making request: %v", err))
			}
			time.Sleep(retryDelay)
			continue
		}
		if resp.IsError() {
			if attempt == maxRetries {
				return base64Response, InternalServerError(fmt.Sprintf("request failed: %s", errorResponse.Message))
			}
			time.Sleep(retryDelay)
			continue
		}
		break
	}

	return base64Response, nil
}

type Link struct {
	Special bool   `json:"especial"`
	System  string `json:"sistema"`
}

func (rc *restClient) LinkDocument(code, fileID string) error {
	body := Link{
		Special: false,
		System:  OriginSystem,
	}

	url := fmt.Sprintf("/api/gde/vincular-gedo-a-expediente/%s/%s", code, fileID)
	ctx := context.Background()

	var response interface{}
	resp, err := rc.client.ExecuteWithAuth(ctx, func() (*resty.Response, error) {
		return rc.client.CreateRequest(generateRequestParams(body)).
			SetContext(ctx).
			SetResult(&response).
			SetError(&ErrorResponse{}).
			Post(url)
	})
	if err != nil {
		return InternalServerError(fmt.Sprintf("error making request: %v", err))
	}
	if resp.IsError() {
		fmt.Printf("Server responded with an error. Status Code: %d\n", resp.StatusCode())
		return InternalServerError(fmt.Sprintf("request error: %s", resp.String()))
	}

	switch v := response.(type) {
	case bool:
		return nil
	case map[string]interface{}:
		if errMsg, ok := v["error"].(string); ok {
			return fmt.Errorf("API error: %s", errMsg)
		}
		return fmt.Errorf("unexpected error response format: %v", v)
	default:
		return fmt.Errorf("unexpected response type: %T", v)
	}
}
