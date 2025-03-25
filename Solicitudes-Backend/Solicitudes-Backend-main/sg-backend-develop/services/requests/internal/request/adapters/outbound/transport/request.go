package transport

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

type DocumentDataModel struct {
	ID        int64     `db:"id"`
	RequestID int64     `db:"request_id"`
	Name      string    `db:"name"`
	Type      int       `db:"type"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type RequestDataModel struct {
	ID                   int64               `db:"id"`
	UserID               int64               `db:"user_id"`
	UserType             sql.NullString      `db:"user_type"`
	RequestTypeID        int                 `db:"request_type_id"`
	PropertyID           int64               `db:"property_id"`
	StatusID             int                 `db:"status_id"`
	StatusName           string              `db:"status_name"`
	Description          string              `db:"description"`
	VerificationComplete bool                `db:"verification_complete"`
	VerificationDate     *time.Time          `db:"verification_date"`
	CreatedAt            time.Time           `db:"created_at"`
	UpdatedAt            time.Time           `db:"updated_at"`
	FileNumber           string              `db:"file_number"`
	ABLDebt              string              `db:"abl_debt"`
	SelectedActivities   []int               `db:"selected_activity"`
	Activities           []string            `db:"activities"`
	EstimatedTime        int64               `db:"estimated_time"`
	Insurance            bool                `db:"insurance"`
	Observations         sql.NullString      `db:"observations"`
	ObservationsTasks    sql.NullString      `db:"observations_tasks"`
	Documents            []DocumentDataModel `db:"-"` // Usamos db:"-" para indicar que no es una columna en la tabla requests
}

func ToCreateRequestDataModel(req *domain.Request) *RequestDataModel {
	documents := make([]DocumentDataModel, len(req.Documents))
	for i, doc := range req.Documents {
		documents[i] = DocumentDataModel{
			Name:    doc.Name,
			Type:    int(doc.Type),
			Content: doc.Content,
		}
	}

	return &RequestDataModel{
		UserID: req.UserID,
		UserType: sql.NullString{
			String: string(req.UserType),
			Valid:  true,
		},
		PropertyID:         req.PropertyID,
		FileNumber:         req.FileNumber,
		ABLDebt:            req.ABLDebt,
		SelectedActivities: req.SelectedActivities,
		EstimatedTime:      req.EstimatedTime,
		Insurance:          req.Insurance,
		Description:        req.ProjectDesc,
		Documents:          documents,
	}
}

func ToRequestDomainList(dataModels []RequestDataModel) []domain.Request {
	domainRequests := make([]domain.Request, len(dataModels))
	for i, dm := range dataModels {
		domainRequests[i] = domain.Request{
			UserID:             dm.UserID,
			StatusName:         dm.StatusName,
			CreatedAt:          dm.CreatedAt,
			PropertyID:         dm.PropertyID,
			Address:            domain.Address{}, // Necesita venir de la tabla de direcciones
			ABLDebt:            dm.ABLDebt,
			CommonZone:         false,               // No está en RequestDataModel
			UserType:           domain.UserTypeUser, // Asignar un valor por defecto
			SelectedActivities: dm.SelectedActivities,
			ProjectDesc:        dm.Description, // Asumiendo que ProjectDesc corresponde a Description
			EstimatedTime:      dm.EstimatedTime,
			Insurance:          dm.Insurance,
			FileNumber:         dm.FileNumber,
			Documents:          ToDocumentRequestList(dm.Documents),
		}
	}

	return domainRequests
}

func ToDocumentRequestList(dataModels []DocumentDataModel) []domain.DocumentRequest {
	if len(dataModels) == 0 {
		return nil
	}

	domainDocuments := make([]domain.DocumentRequest, len(dataModels))
	for i, dm := range dataModels {
		domainDocuments[i] = domain.DocumentRequest{
			Name:    dm.Name,
			Type:    domain.DocumentTypeID(dm.Type),
			Content: dm.Content,
		}
	}

	return domainDocuments
}

// RequestDataModel represents the database model for request data
type RequestPersonDataModel struct {
	Cuil      string `db:"cuil"` // Changed to string to match DB schema
	Dni       string `db:"dni"`  // Changed to string to match DB schema
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Phone     string `db:"phone"` // Changed to string to match DB schema
}

// ToRequestDomain converts a RequestDataModel to domain.Request
func RequestPersonDataModelToRequestDomain(model *RequestPersonDataModel) *domain.Request {
	if model == nil {
		return nil
	}

	return &domain.Request{
		Cuil:      model.Cuil, // Already string in both
		Dni:       model.Dni,  // Already string in both
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Email:     model.Email,
		Phone:     model.Phone, // Already string in both
		// Set default values for required fields
		UserType:           domain.UserTypeUser,
		Address:            domain.Address{},
		ABLDebt:            "",
		CommonZone:         false,
		SelectedActivities: []int{},
		ProjectDesc:        "",
		EstimatedTime:      0,
		Insurance:          false,
		Documents:          []domain.DocumentRequest{},
	}
}

type AddressPayload struct {
	ABLNumber int64  `json:"abl_number"`
	ZipCode   int64  `json:"zc"`
	Street    string `json:"street"`
	Number    string `json:"number"`
}

type DocumentPayload struct {
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

type EmailRequestPayload struct {
	Code  string `json:"code" binding:"required"`
	Email string `json:"email"`
}

type RequestPayload struct {
	ID                 int64             `json:"id" binding:"required"`
	Cuil               int64             `json:"cuil"`
	DocumentNumber     int64             `json:"dni"`
	FirstName          string            `json:"first_name"`
	LastName           string            `json:"last_name"`
	Email              string            `json:"email"`
	Phone              int64             `json:"phone"`
	UserType           string            `json:"user_type"`
	Address            AddressPayload    `json:"address"`
	SelectedActivities []int             `json:"tasks"`
	Activities         []string          `json:"activities"`
	ProjectDesc        string            `json:"description"`
	EstimatedTime      int64             `json:"time"`
	Insurance          bool              `json:"insurance"`
	Documents          []DocumentPayload `json:"documents"`
	Observations       string            `json:"observations"`
	ObservationsTasks  string            `json:"observations_tasks"`
}

func ToRequestPayload(req *domain.Request) *RequestPayload {
	if req == nil {
		return nil
	}

	// Convert documents
	docs := make([]DocumentPayload, 0, len(req.Documents))
	for _, doc := range req.Documents {
		docs = append(docs, DocumentPayload{
			Name:    doc.Name,
			Type:    int(doc.Type),
			Content: doc.Content,
		})
	}

	// Convert string numbers to int64
	cuil, _ := strconv.ParseInt(req.Cuil, 10, 64)
	dni, _ := strconv.ParseInt(req.Dni, 10, 64)
	phone, _ := strconv.ParseInt(req.Phone, 10, 64)
	ablNumber, _ := strconv.ParseInt(strconv.FormatInt(int64(req.Address.ABLNumber), 10), 10, 64)

	return &RequestPayload{
		ID:                 req.ID,
		Cuil:               cuil,
		DocumentNumber:     dni,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Email:              req.Email,
		Phone:              phone,
		UserType:           strings.ToLower(string(req.UserType)),
		SelectedActivities: req.SelectedActivities,
		Activities:         req.Activities,
		ProjectDesc:        req.ProjectDesc,
		EstimatedTime:      req.EstimatedTime,
		Observations:       req.Observations,
		ObservationsTasks:  req.ObservationsTasks,
		Address: AddressPayload{
			ABLNumber: ablNumber,
			ZipCode:   1642,
			Street:    req.Address.Street,
			Number:    req.Address.Number,
		},
		Insurance: req.Insurance,
		Documents: docs,
	}
}
