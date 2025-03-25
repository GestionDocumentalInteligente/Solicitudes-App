package transport

import (
	"fmt"
	"strings"

	domain "github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

// type Document struct {
// 	ID       int64  `json:"id"`
// 	Title    string `json:"title"`
// 	GedoCode string `json:"gedoCode"`
// }

// type VerificationRequestPresenter struct {
// 	VerificationCase  string     `json:"verificationCase"`
// 	RecordNumber      string     `json:"recordNumber"`
// 	RequestType       string     `json:"requestType"`
// 	DocumentType      string     `json:"documentType"`
// 	DeliveryDate      string     `json:"deliveryDate"`
// 	Status            string     `json:"status"`
// 	RequesterFullName string     `json:"requesterFullName"`
// 	RequesterCuil     string     `json:"requesterCuil"`
// 	RequesterAddress  string     `json:"requesterAddress"`
// 	Documents         []Document `json:"documents"`
// }

// func ToVerificationRequestPresenter(req *domain.Verification) *VerificationRequestPresenter {
// 	if req == nil {
// 		return nil
// 	}

// 	docs := make([]Document, 0, len(req.Documents))
// 	for _, doc := range req.Documents {
// 		docs = append(docs, Document{
// 			ID:       doc.ID,
// 			Title:    doc.Title,
// 			GedoCode: doc.GedoCode,
// 		})
// 	}

// 	return &VerificationRequestPresenter{
// 		VerificationCase:  req.VerificationCase,
// 		RecordNumber:      req.RecordNumber,
// 		RequestType:       req.RequestType,
// 		DocumentType:      req.DocumentType,
// 		DeliveryDate:      req.DeliveryDate,
// 		Status:            req.Status,
// 		RequesterFullName: req.RequesterFullName,
// 		RequesterCuil:     req.RequesterCuil,
// 		RequesterAddress:  req.RequesterAddress,
// 		Documents:         docs,
// 	}
// }

type Document struct {
	ID           int64          `json:"id,string"`
	Title        string         `json:"title"`
	GedoCode     string         `json:"gedoCode,omitempty"`
	Content      string         `json:"content,omitempty"`
	Request      *RequestResume `json:"request,omitempty"`
	VerifiedBy   string         `json:"verifiedBy,omitempty"`
	VerifiedDate string         `json:"verifiedDate,omitempty"`
}

type RequestResume struct {
	DescriptiveMemory string   `json:"descriptiveMemory,omitempty"`
	AssignedTask      []string `json:"assignedTask,omitempty"`
	EstimatedTime     string   `json:"estimatedTime,omitempty"`
	RequiresInsurance string   `json:"requiresInsurance,omitempty"`
}

type DocContent struct {
	DescriptiveMemory string `json:"descriptiveMemory,omitempty"`
	AssignedTask      string `json:"assignedTask,omitempty"`
	EstimatedTime     string `json:"estimatedTime,omitempty"`
	RequiresInsurance string `json:"requiresInsurance,omitempty"`
}

type VerificationRequestPresenter struct {
	ID                int64      `json:"id"`
	VerificationCase  string     `json:"verificationCase,omitempty"`
	RecordNumber      string     `json:"recordNumber"`
	RequestType       string     `json:"requestType"`
	DocumentType      string     `json:"documentType"`
	DeliveryDate      string     `json:"deliveryDate"`
	Status            string     `json:"status"`
	StatusTask        string     `json:"status_task"`
	StatusProperty    string     `json:"status_property"`
	RequesterFullName string     `json:"requesterFullName"`
	RequesterCuil     string     `json:"requesterCuil"`
	RequesterAddress  string     `json:"requesterAddress"`
	Documents         []Document `json:"documents,omitempty"`
}

func ToVerificationRequestPresenter(req *domain.Verification) *VerificationRequestPresenter {
	if req == nil {
		return nil
	}

	docs := make([]Document, 0, len(req.Documents))
	for _, doc := range req.Documents {
		d := Document{
			Title: doc.Title,
		}

		docs = append(docs, d)
	}

	return &VerificationRequestPresenter{
		//RecordNumber:      req.FileNumber,
		RequestType:  "Aviso de obra",
		DocumentType: "Potestad sobre el inmueble",
		//DeliveryDate:      req.CreatedAt.Format("02/01/2006"),
		//Status:            req.StatusName,
		//RequesterFullName: fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		//RequesterCuil:     req.Cuil,
		//RequesterAddress:  fmt.Sprintf("%s %s", req.Address.Street, req.Address.Number),
		Documents: docs,
	}
}

func ToVerificationListPresenter(list []domain.Verification) []VerificationRequestPresenter {
	verifications := make([]VerificationRequestPresenter, len(list))
	for i, model := range list {
		verifications[i] = VerificationRequestPresenter{
			ID:                model.ID,
			RecordNumber:      model.RecordNumber,
			RequestType:       model.RequestType,
			DocumentType:      "Potestad sobre el inmueble",
			DeliveryDate:      model.DeliveryDate,
			Status:            strings.ToLower(model.Status),
			StatusTask:        strings.ToLower(model.StatusTask),
			StatusProperty:    strings.ToLower(model.StatusProperty),
			RequesterFullName: model.RequesterFullName,
			RequesterCuil:     model.RequesterCuil,
			RequesterAddress:  model.RequesterAddress,
		}
	}
	return verifications
}

func ToDocumentPresenter(doc domain.Document) Document {
	return Document{
		ID:       doc.ID,
		Title:    doc.Title,
		GedoCode: doc.GedoCode,
		Content:  doc.Content,
	}
}

func ToDocumentListPresenter(request *domain.Request, list []domain.Document, code string) []Document {
	docs := make([]Document, len(list)+1)

	var insurance string
	if request.Insurance {
		insurance = "Si"
	} else {
		insurance = "No"
	}

	docs[0] = Document{
		ID:       1,
		Title:    "Datos suministrados",
		GedoCode: code,
		Request: &RequestResume{
			DescriptiveMemory: request.ProjectDesc,
			AssignedTask:      request.Activities,
			EstimatedTime:     fmt.Sprintf("%d dias", request.EstimatedTime),
			RequiresInsurance: insurance,
		},
		VerifiedBy:   request.VerifyBy,
		VerifiedDate: request.VerifyDate.Format("2006-01-02 15:04:05"),
	}

	for i, model := range list {
		id := i + 2
		docs[i+1] = Document{
			ID:           int64(id),
			Title:        model.Title,
			GedoCode:     model.GedoCode,
			VerifiedBy:   request.VerifyBy,
			VerifiedDate: request.VerifyDate.Format("2006-01-02 15:04:05"),
		}
	}
	return docs
}

func ToDocumentList(list []domain.Document) []Document {
	docs := make([]Document, len(list))

	for i, model := range list {
		id := i + 1
		docs[i] = Document{
			ID:           int64(id),
			Title:        model.Title,
			GedoCode:     model.GedoCode,
			VerifiedBy:   model.VerifiedBy,
			VerifiedDate: model.VerifiedDate,
		}
	}
	return docs
}

func ToReplacementDocumentList(list []domain.Document) []Document {
	docs := make([]Document, 0)

	groupedCodes := []string{}
	for _, model := range list {
		if model.Type == 9 || model.Type == 10 || model.Type == 11 || model.Type == 14 {
			groupedCodes = append(groupedCodes, model.GedoCode)
		}
	}

	groupedGedoCode := ""
	if len(groupedCodes) > 0 {
		groupedGedoCode = strings.Join(groupedCodes, " y GEDO N° ")
	}

	userDocsReady := false
	id := 1
	for _, model := range list {
		title := model.Title
		gedoCode := model.GedoCode
		if model.Type == 9 || model.Type == 10 || model.Type == 11 || model.Type == 14 {
			if userDocsReady {
				continue
			}
			title = "uploaded_user_documents"
			gedoCode = groupedGedoCode
			userDocsReady = true
		}

		docs = append(docs, Document{
			ID:       int64(id),
			Title:    title,
			GedoCode: gedoCode,
		})
		id++
	}
	return docs
}
