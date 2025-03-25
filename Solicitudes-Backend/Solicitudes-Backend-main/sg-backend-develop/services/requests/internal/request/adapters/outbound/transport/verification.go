package transport

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

type VerificactionDataModel struct {
	ID             int64          `json:"id"`
	PropertyOwner  string         `json:"owner"`
	FileNumber     string         `json:"file_number"`
	DeliveryDate   time.Time      `json:"deliveryDate"`
	Status         string         `json:"status"`
	StatusTask     string         `json:"status_task"`
	StatusProperty string         `json:"status_property"`
	CUIL           string         `json:"cuil"`
	DNI            string         `json:"dni"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	AddrStreet     sql.NullString `json:"street"`
	AddrNumber     *int           `json:"number"`
	Locality       sql.NullString `json:"locality"`
	RequestType    string         `json:"request_type"`
	RequestStatus  string         `json:"request_status"`
	DocumentID     string         `json:"document_id"`
	DocumentType   string         `json:"document_type"`
}

func ToVerificationDomain(v *VerificactionDataModel) *domain.Verification {
	address := ""
	if v.AddrNumber != nil {
		address = fmt.Sprintf("%s %d", v.AddrStreet.String, *v.AddrNumber)
	} else {
		address = v.AddrStreet.String
	}

	return &domain.Verification{
		VerificationCase:  v.DocumentID, // Usando DocumentID como caso de verificación
		RequestType:       v.RequestType,
		DocumentType:      v.DocumentType,
		Status:            v.RequestStatus,
		RequesterFullName: v.PropertyOwner,
		RequesterCuil:     v.CUIL,
		RequesterAddress:  address,
		Documents:         []domain.Document{}, // Inicializar slice vacío para documents
	}
}

func ToVerificationListDomain(list []VerificactionDataModel) []domain.Verification {
	var verifications []domain.Verification
	for _, d := range list {
		address := ""
		if d.AddrNumber != nil {
			address = fmt.Sprintf("%s %d", d.AddrStreet.String, *d.AddrNumber)
		} else {
			address = d.AddrStreet.String
		}

		adjustedTime := time.Time(d.DeliveryDate).Add(-3 * time.Hour)

		verifications = append(verifications, domain.Verification{
			ID:                d.ID,
			RecordNumber:      d.FileNumber,
			RequestType:       d.RequestType,
			DeliveryDate:      adjustedTime.Format("2006-01-02 15:04"),
			Status:            d.Status,
			StatusTask:        d.StatusTask,
			StatusProperty:    d.StatusProperty,
			RequesterFullName: fmt.Sprintf("%s %s", d.FirstName, d.LastName),
			RequesterCuil:     d.CUIL,
			RequesterAddress:  fmt.Sprintf("%s, %s", address, d.Locality.String),
		})
	}

	return verifications
}

func ToDocumentDomainList(list []DocumentModel) []domain.Document {
	var documents []domain.Document
	for _, d := range list {
		documents = append(documents, domain.Document{
			ID:       d.ID,
			Type:     d.Type,
			Title:    d.Description.String,
			GedoCode: d.FileID,
		})
	}

	return documents
}

// FromVerificationRequest convierte de VerificationRequest al data model
// func ToVerificationDataModel(vr *domain.Verification) *VerificactionDataModel {
// 	// Separar la dirección en calle y número (si es necesario)
// 	addrParts := strings.SplitN(vr.RequesterAddress, " ", 2)
// 	street := ""
// 	number := ""
// 	if len(addrParts) > 0 {
// 		street = addrParts[0]
// 		if len(addrParts) > 1 {
// 			number = addrParts[1]
// 		}
// 	}

// 	return &VerificactionDataModel{
// 		PropertyOwner: vr.RequesterFullName,
// 		AddrStreet:    street,
// 		AddrNumber:    number,
// 		CUIL:          vr.RequesterCuil,
// 		RequestType:   vr.RequestType,
// 		DocumentType:  vr.DocumentType,
// 		RequestStatus: vr.Status,
// 		DocumentID:    vr.VerificationCase,
// 	}
// }
