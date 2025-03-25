package transport

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

// Custom types for JSON marshaling/unmarshaling
type DocumentTypeIDJson int
type UserTypeJSON string

// FileJson represents the structure for files in the request
type FileJson struct {
	Type    DocumentTypeIDJson `json:"type"`
	Name    string             `json:"name"`
	Content string             `json:"content"`
}

// RequestJson represents the JSON structure for requests
type RequestJson struct {
	Cuil               string
	PropertyID         string     `json:"property_id"`
	AddressStreet      string     `json:"address_street"`
	AddressNumber      string     `json:"address_number"`
	AddressABLNumber   string     `json:"address_abl_number"`
	ABLDebt            string     `json:"ablDebt"`
	CommonZone         string     `json:"commonZone"`
	UserType           string     `json:"userType"`
	SelectedActivities string     `json:"selectedActivity"`
	ProjectDesc        string     `json:"projectDescription"`
	EstimatedTime      string     `json:"estimatedTime"`
	Insurance          string     `json:"insurance"`
	Files              []FileJson `json:"files"`
}

type VerifiedRequestJson struct {
	FileNumber                string
	Cuil                      string
	VerificationType          string   `json:"verificationType"` //tasks. //property
	Observations              string   `json:"observations"`
	FinalVerificationDocument string   `json:"finalVerificationDocument"`
	Reference                 string   `json:"reference"`
	AssignedTask              []string `json:"assignedTask"`
}

type ValidateRequestJson struct {
	FileNumber  string
	Cuil        string
	IsValid     bool   `json:"is_valid"`
	FileContent string `json:"authorizationDocument"`
}

// ToRequestDomain converts a RequestJson to a domain Request
func ToRequestDomain(reqJson *RequestJson) *domain.Request {
	if reqJson == nil {
		return nil
	}

	// Convert ablNumber from string to int64
	ablNumber, _ := strconv.ParseInt(reqJson.AddressABLNumber, 10, 64)

	// Convert insurance string to bool
	insurance := strings.ToLower(reqJson.Insurance) == "true"

	// Convert commonZone string to bool
	commonZone := strings.ToLower(reqJson.CommonZone) == "true"

	// Convert estimatedTime to int
	estimatedTime, _ := strconv.ParseInt(reqJson.EstimatedTime, 10, 64)

	propertyID, _ := strconv.ParseInt(reqJson.PropertyID, 10, 64)

	var activities []int
	if err := json.Unmarshal([]byte(reqJson.SelectedActivities), &activities); err != nil {
		fmt.Println(err.Error())
	}

	return &domain.Request{
		Address: domain.Address{
			Street:    reqJson.AddressStreet,
			Number:    reqJson.AddressNumber,
			ABLNumber: ablNumber,
		},
		Cuil:               reqJson.Cuil,
		PropertyID:         propertyID,
		ABLDebt:            reqJson.ABLDebt,
		CommonZone:         commonZone,
		UserType:           domain.UserType(reqJson.UserType),
		SelectedActivities: activities,
		ProjectDesc:        reqJson.ProjectDesc,
		EstimatedTime:      estimatedTime,
		Insurance:          insurance,
		Documents:          toDocumentRequestList(reqJson.Files),
	}
}

func ToVerifiedRequestDomain(req *VerifiedRequestJson) *domain.VerifiedRequest {
	intArr := make([]int, len(req.AssignedTask))

	for i, s := range req.AssignedTask {
		act, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println(err.Error())
		}
		intArr[i] = act
	}

	return &domain.VerifiedRequest{
		FileNumber:                req.FileNumber,
		Cuil:                      req.Cuil,
		VerificationType:          req.VerificationType,
		Observations:              req.Observations,
		FinalVerificationDocument: req.FinalVerificationDocument,
		Reference:                 req.Reference,
		SelectedActivities:        intArr,
	}
}

func ToValidateRequestDomain(req *ValidateRequestJson) *domain.ValidateRequest {
	return &domain.ValidateRequest{
		FileNumber:  req.FileNumber,
		Cuil:        req.Cuil,
		IsValid:     req.IsValid,
		FileContent: req.FileContent,
	}
}

// toDocumentRequestList converts FileJson slice to domain DocumentRequests
func toDocumentRequestList(files []FileJson) []domain.DocumentRequest {
	if len(files) == 0 {
		return []domain.DocumentRequest{}
	}

	result := make([]domain.DocumentRequest, 0, len(files))
	for _, file := range files {
		result = append(result, domain.DocumentRequest{
			Name:    file.Name,
			Type:    domain.DocumentTypeID(file.Type),
			Content: file.Content,
		})
	}
	return result
}

// Get All Requests

type GetAllReqPresenterDocumentTypeID int
type GetAllReqPresenterUserType string

// User type constants
const (
	GetAllReqPresenterUserTypeAdmin GetAllReqPresenterUserType = "Admin"
	GetAllReqPresenterUserTypeUser  GetAllReqPresenterUserType = "User"
)

// CustomTime para manejar el formato de fecha
type CustomTime time.Time

func (t CustomTime) MarshalJSON() ([]byte, error) {
	adjustedTime := time.Time(t).Add(-3 * time.Hour)
	formatted := adjustedTime.Format("2006-01-02 15:04")
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}

// DocumentRequest represents a document in the domain
type GetAllReqPresenterDocumentRequest struct {
	Name    string                           `json:"name"`
	Type    GetAllReqPresenterDocumentTypeID `json:"type"`
	Content string                           `json:"content"`
}

// Address represents an address in the domain
type GetAllReqPresenterAddress struct {
	Street    string `json:"street"`
	Number    string `json:"number"`
	ABLNumber int64  `json:"abl_number"`
}

// Request represents the main request entity in the domain
type GetAllReqPresenterRequest struct {
	UserID             int64                               `json:"user_id"`
	PropertyID         int64                               `json:"property_id"`
	Cuil               string                              `json:"cuil"`
	Dni                string                              `json:"dni"`
	FirstName          string                              `json:"first_name"`
	LastName           string                              `json:"last_name"`
	Email              string                              `json:"email"`
	Phone              string                              `json:"phone"`
	Address            GetAllReqPresenterAddress           `json:"address"`
	ABLDebt            string                              `json:"abl_debt"`
	CommonZone         bool                                `json:"common_zone"`
	UserType           GetAllReqPresenterUserType          `json:"user_type"`
	SelectedActivities []int                               `json:"selected_activities"`
	ProjectDesc        string                              `json:"project_desc"`
	EstimatedTime      int64                               `json:"estimated_time"`
	Insurance          bool                                `json:"insurance"`
	FileNumber         string                              `json:"file_number"`
	Documents          []GetAllReqPresenterDocumentRequest `json:"documents"`
	StatusName         string                              `json:"status"`
	CreatedAt          CustomTime                          `json:"created_at"`
	Type               string                              `json:"type"`
}

// Helper function para convertir los documentos
func toDocumentRequestListPresenter(docs []domain.DocumentRequest) []GetAllReqPresenterDocumentRequest {
	presenterDocs := make([]GetAllReqPresenterDocumentRequest, len(docs))
	for i, doc := range docs {
		presenterDocs[i] = GetAllReqPresenterDocumentRequest{
			Name:    doc.Name,
			Type:    GetAllReqPresenterDocumentTypeID(doc.Type),
			Content: doc.Content,
		}
	}
	return presenterDocs
}

// Función principal de mapeo
func ToRequestListPresenter(requests []domain.Request) []GetAllReqPresenterRequest {
	presenters := make([]GetAllReqPresenterRequest, len(requests))

	for i, req := range requests {
		fileNumber := req.FileNumber
		if fileNumber == "" {
			fileNumber = "En Proceso"
		}
		presenters[i] = GetAllReqPresenterRequest{
			UserID:     req.UserID,
			PropertyID: req.PropertyID,
			Cuil:       req.Cuil,
			Dni:        req.Dni,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Phone:      req.Phone,
			Address: GetAllReqPresenterAddress{
				Street:    req.Address.Street,
				Number:    req.Address.Number,
				ABLNumber: req.Address.ABLNumber,
			},
			ABLDebt:            req.ABLDebt,
			CommonZone:         req.CommonZone,
			UserType:           GetAllReqPresenterUserType(req.UserType),
			SelectedActivities: req.SelectedActivities,
			ProjectDesc:        req.ProjectDesc,
			EstimatedTime:      req.EstimatedTime,
			Insurance:          req.Insurance,
			FileNumber:         fileNumber,
			Documents:          toDocumentRequestListPresenter(req.Documents),
			StatusName:         req.StatusName,
			CreatedAt:          CustomTime(req.CreatedAt),
			Type:               "Aviso de Obra",
		}
	}

	return presenters
}
