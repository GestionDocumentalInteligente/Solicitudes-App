package filehdl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/teamcubation/sg-file-manager-api/internal/manager"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
	"github.com/teamcubation/sg-file-manager-api/pkg/log"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	useCase manager.FileUseCase
}

func NewFileHandler(uc manager.FileUseCase) *FileHandler {
	return &FileHandler{
		useCase: uc,
	}
}

func (h *FileHandler) CreateRecord(c *gin.Context) {
	ctx := log.Context(c.Request)
	logger := log.FromContext(ctx)
	logger.Info("Entering FileHandler: CreateRecord()")

	startTime := time.Now()

	var requestData RequestDataDTO

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := requestData.UserType.IsValid(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := user.User{
		Type:           user.UserType(requestData.UserType),
		Cuil:           requestData.Cuil,
		DocumentNumber: requestData.DocumentNumber,
		FirstName:      requestData.FirstName,
		LastName:       requestData.LastName,
		Email:          requestData.Email,
		Phone:          requestData.Phone,
		Address: user.Address{
			ABLNumber: requestData.Address.ABLNumber,
			Street:    requestData.Address.Street,
			Number:    requestData.Address.Number,
			ZipCode:   requestData.Address.ZipCode,
		},
	}

	documents, err := h.validateAndConvertDocuments(requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.useCase.CreateEERecord(ctx, user, requestData.ID, documents)
	if err != nil {
		var customErr *file.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.StatusCode, gin.H{"error": customErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	goroutineCtx, cancel := context.WithTimeout(context.Background(), 4800*time.Second)

	go func(ctx context.Context, dataDTO RequestDataDTO, recordNumber string) {
		defer func() {
			if r := recover(); r != nil {
				logger.Warn(fmt.Sprintf("panic in goroutine: %v", r))
			}
			logger.Printf("document processing completed! Took %s", time.Since(startTime))
			cancel()
		}()

		memo := file.Memo{
			Insurance:   dataDTO.Insurance,
			Description: dataDTO.ProjectDesc,
			Tasks:       dataDTO.SelectedActivities,
			Time:        dataDTO.EstimatedTime,
		}

		if err := h.useCase.ProcessDocuments(ctx, recordNumber, user, memo); err != nil {
			logger.Error(fmt.Sprintf("error processing documents for code %s: %v", recordNumber, err))
			return
		}

		if err := h.useCase.LinkDocuments(ctx, recordNumber, user.Type, memo.Insurance); err != nil {
			logger.Error(fmt.Sprintf("error linking documents for code %s: %v", recordNumber, err))
			return
		}

		if err := h.useCase.UpdateRequestStatus(ctx, dataDTO.ID); err != nil {
			logger.Error(fmt.Sprintf("error updating request status %s: %v", recordNumber, err))
		}
	}(goroutineCtx, requestData, code)

	logger.Info("record created!")

	c.JSON(http.StatusCreated, gin.H{"code": code})
}

func (h *FileHandler) SendDocumentToRecord(c *gin.Context) {
	ctx := log.Context(c.Request)
	logger := log.FromContext(ctx)
	logger.Info("Entering FileHandler: SendDocumentToRecord()")

	type docContent struct {
		Content   string              `json:"content"`
		Reference string              `json:"reference"`
		Type      file.DocumentTypeID `json:"type"`
		UserName  string              `json:"username"`
	}

	var requestData docContent
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestData.Content == "" || requestData.Reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imcomplete data"})
		return
	}

	code := c.Param("id")
	err := h.useCase.SendDocumentToRecord(ctx, code, requestData.Content, requestData.Reference, requestData.Type, requestData.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "document created"})
}

func (h *FileHandler) UpdateDocumentsInRecord(c *gin.Context) {
	ctx := log.Context(c.Request)
	logger := log.FromContext(ctx)
	logger.Info("Entering FileHandler: UpdateDocumentsInRecord()")

	code := c.Param("id")
	var requestData RequestDataDTO

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("updating request %d (code: %s)", requestData.ID, code)

	var documents []file.Document
	for _, doc := range requestData.Documents {
		if _, exists := file.DocumentTypeDescriptionMap[doc.Type]; !exists {
			fmt.Printf("document Type %d is not valid. Valid types are: 9, 10, 11\n", doc.Type)
			continue
		}

		if doc.Type == file.CoOwnership ||
			doc.Type == file.PropertyTitle ||
			doc.Type == file.OwnerAuthorization ||
			doc.Type == file.AppointmentCertificate {
			if requestData.Observations == "" {
				continue
			}
		}

		if doc.Type == file.WithInsurance && requestData.ObservationsTasks == "" {
			continue
		}

		documents = append(documents, file.Document{
			TypeID:      doc.Type,
			Name:        doc.Name,
			Content:     doc.Content,
			ContentType: doc.ContentType,
			Metadata: file.DocumentMetadata{
				DocumentType: file.IfGraTypeDocument,
				Reference:    file.DocumentTypeDescriptionMap[doc.Type],
				OriginSystem: file.OriginSystem,
				FullName: user.GetFullName(user.User{
					FirstName: requestData.FirstName,
					LastName:  requestData.LastName}),
				Position:   fmt.Sprintf("%d", requestData.DocumentNumber),
				Department: "Ciudadano",
			},
		})
	}

	user := user.User{
		Type:           user.UserType(requestData.UserType),
		Cuil:           requestData.Cuil,
		DocumentNumber: requestData.DocumentNumber,
		FirstName:      requestData.FirstName,
		LastName:       requestData.LastName,
		Email:          requestData.Email,
		Phone:          requestData.Phone,
		Address: user.Address{
			ABLNumber: requestData.Address.ABLNumber,
			Street:    requestData.Address.Street,
			Number:    requestData.Address.Number,
			ZipCode:   requestData.Address.ZipCode,
		},
	}

	memo := file.Memo{
		Insurance:   requestData.Insurance,
		Description: requestData.ProjectDesc,
		Tasks:       requestData.SelectedActivities,
		Time:        requestData.EstimatedTime,
	}

	if err := h.useCase.UploadDocuments(ctx, code, user, documents, memo, requestData.ObservationsTasks != ""); err != nil {
		logger.Error(fmt.Sprintf("error processing documents for code %s: %v", code, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "updated!"})
}

func (h *FileHandler) validateAndConvertDocuments(requestData RequestDataDTO) ([]file.Document, error) {
	var documents []file.Document
	for _, doc := range requestData.Documents {
		if _, exists := file.DocumentTypeDescriptionMap[doc.Type]; !exists {
			fmt.Printf("document Type %d is not valid. Valid types are: 9, 10, 11\n", doc.Type)
			continue
		}

		if !requestData.Insurance && doc.Type == file.WithInsurance {
			continue
		}

		documents = append(documents, file.Document{
			TypeID:      doc.Type,
			Name:        doc.Name,
			Content:     doc.Content,
			ContentType: doc.ContentType,
			Metadata: file.DocumentMetadata{
				DocumentType: file.IfGraTypeDocument,
				Reference:    file.DocumentTypeDescriptionMap[doc.Type],
				OriginSystem: file.OriginSystem,
				FullName: user.GetFullName(user.User{
					FirstName: requestData.FirstName,
					LastName:  requestData.LastName}),
				Position:   fmt.Sprintf("%d", requestData.DocumentNumber),
				Department: "Ciudadano",
			},
		})
	}

	return documents, nil
}

func (h *FileHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong!"})
}
