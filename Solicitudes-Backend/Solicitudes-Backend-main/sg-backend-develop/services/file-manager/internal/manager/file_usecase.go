package manager

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
	"github.com/teamcubation/sg-file-manager-api/pkg/log"
)

type fileUseCase struct {
	fileClient     file.Client
	fileRepository file.Repository
	// queue          file.MessageQueue
}

type FileUseCase interface {
	CreateEERecord(ctx context.Context, user user.User, id int64, documents []file.Document) (string, error)
	ProcessDocuments(ctx context.Context, code string, user user.User, memo file.Memo) error
	LinkDocuments(ctx context.Context, code string, userType user.UserType, insurance bool) error
	SendDocumentToRecord(ctx context.Context, code, content, reference string, docType file.DocumentTypeID, username string) error
	UploadDocuments(ctx context.Context, code string, user user.User, documents []file.Document, memo file.Memo, updateTasks bool) error
	UpdateRequestStatus(ctx context.Context, id int64) error
}

func NewFileUseCase(fileClient file.Client, fileRepository file.Repository) FileUseCase {
	return &fileUseCase{
		fileClient:     fileClient,
		fileRepository: fileRepository,
	}
}

func (a *fileUseCase) CreateEERecord(ctx context.Context, user user.User, id int64, documents []file.Document) (string, error) {
	logger := log.FromContext(ctx)

	var wg sync.WaitGroup
	errorChan := make(chan error, len(documents))
	signedDocuments := make(chan file.SignedDocument, len(documents))

	for _, doc := range documents {
		wg.Add(1)
		go func(doc file.Document) {
			defer wg.Done()
			if doc.Content == "" {
				errorChan <- fmt.Errorf("empty doc content")
				return
			}

			fileSigned, err := a.fileClient.CreateGEDO(ctx, doc)
			if err != nil {
				errorChan <- fmt.Errorf("error sending doc to GEDO: %w", err)
				return
			}

			if fileSigned.Number == "" {
				errorChan <- fmt.Errorf("error sending document: api response error")
				return
			}

			base64Signed, err := a.fileClient.DownloadGEDO(ctx, fileSigned.Number)
			if err != nil {
				errorChan <- fmt.Errorf("error getting doc from api: %w", err)
			}
			fileSigned.Content = base64Signed

			signedDocuments <- fileSigned
		}(doc)
	}

	wg.Wait()
	close(errorChan)
	close(signedDocuments)

	var collectedErrors []string
	for err := range errorChan {
		collectedErrors = append(collectedErrors, err.Error())
	}

	if len(collectedErrors) > 0 {
		return "", fmt.Errorf("encountered errors: %s", strings.Join(collectedErrors, "; "))
	}

	code, err := a.fileClient.CreateRecord(ctx, user, id)
	if err != nil {
		var customErr *file.CustomError
		if ok := errors.As(err, &customErr); ok {
			logger.Info(fmt.Sprintf("custom error occurred: %s", customErr.Message))
			return "", customErr
		}

		logger.Info(fmt.Sprintf("unexpected error: %v", err))
		return "", file.InternalServerError("unexpected error creating record")
	}

	logger.Info(fmt.Sprintf("The request was created. Record number: %s", code))

	for signed := range signedDocuments {
		err = a.fileRepository.SaveDocument(ctx, code, signed)
		if err != nil {
			a.fileRepository.UpdateRequestStatus(ctx, id, 2)
			return code, fmt.Errorf("error saving doc in repository: %w", err)
		}
		logger.Info(fmt.Sprintf("the document %s (id: %s) was added to request : %s", signed.Filename, signed.Number, code))
	}

	return code, nil
}

func (a *fileUseCase) UploadDocuments(ctx context.Context, code string, userData user.User, documents []file.Document, memo file.Memo, updateTasks bool) error {
	logger := log.FromContext(ctx)

	var wg sync.WaitGroup
	errorChan := make(chan error, len(documents))
	signedDocuments := make(chan file.SignedDocument, len(documents))

	logger.Info(fmt.Sprintf("amount documents to upload: %d", len(documents)))

	for _, doc := range documents {
		wg.Add(1)
		go func(doc file.Document) {
			defer wg.Done()
			if doc.Content == "" {
				errorChan <- fmt.Errorf("empty doc content")
				return
			}

			fileSigned, err := a.fileClient.CreateGEDO(ctx, doc)
			if err != nil {
				errorChan <- fmt.Errorf("UploadDocuments: error sending doc to GEDO: %w", err)
				return
			}

			if fileSigned.Number == "" {
				errorChan <- fmt.Errorf("UploadDocuments: error sending document: api response error")
				return
			}

			base64Signed, err := a.fileClient.DownloadGEDO(ctx, fileSigned.Number)
			if err != nil {
				errorChan <- fmt.Errorf("error getting doc from api: %w", err)
			}
			fileSigned.Content = base64Signed

			signedDocuments <- fileSigned
		}(doc)
	}

	wg.Wait()
	close(errorChan)
	close(signedDocuments)

	var collectedErrors []string
	for err := range errorChan {
		collectedErrors = append(collectedErrors, err.Error())
	}

	if len(collectedErrors) > 0 {
		return fmt.Errorf("encountered errors: %s", strings.Join(collectedErrors, "; "))
	}

	filesMap := make(map[file.DocumentTypeID]string)
	var insuranceCode string
	var addressToLink []string

	for signed := range signedDocuments {
		if err := a.saveOrUpdateDocument(ctx, code, signed); err != nil {
			return err
		}

		if _, exists := file.DocumentTypeDescriptionMap[signed.TypeID]; exists && signed.TypeID != file.WithInsurance {
			filesMap[signed.TypeID] = signed.Number
			addressToLink = append(addressToLink, signed.Number)
			continue
		}

		if signed.TypeID == file.WithInsurance {
			insuranceCode = signed.Number
			continue
		}

		logger.Warn("linking? ", code, signed.TypeID)
		if err := a.linkDocumentWithRetry(code, signed.Number, 5); err != nil {
			return fmt.Errorf("error linking document: %w", err)
		}
	}

	var documentsToLink []string // for autogenerated documents
	if len(addressToLink) > 0 {
		addressDocTemplate, err := docprocessor.GetAddressDocument(userData, filesMap)
		if err != nil {
			return fmt.Errorf("error getting address document: %w", err)
		}

		addressDoc, err := docprocessor.ProcessDocument(ctx, userData, addressDocTemplate)
		if err != nil {
			return fmt.Errorf("error proccesing address document: %w", err)
		}

		fileID, err := a.updateAutogeneratedDoc(ctx, code, addressDoc)
		if err != nil {
			return fmt.Errorf("UploadDocuments: error sending new documents to GEDO: %w", err)
		}

		documentsToLink = append(documentsToLink, fileID)
		documentsToLink = append(documentsToLink, addressToLink...)
	}

	if updateTasks {
		activities, err := a.fileRepository.GetActivityNameByIDs(ctx, memo.Tasks)
		if err != nil {
			activities = "error"
		}

		insuranceDocTemplate, err := docprocessor.GetInsuranceDocument(userData, memo.Insurance, insuranceCode, activities, memo.Description, fmt.Sprintf("%d", memo.Time))
		if err != nil {
			return fmt.Errorf("error getting address document: %w", err)
		}

		insuranceDoc, err := docprocessor.ProcessDocument(ctx, userData, insuranceDocTemplate)
		if err != nil {
			return fmt.Errorf("error proccesing insurance document: %w", err)
		}

		isADifferentDocument := false
		_, _, err = a.fileRepository.GetDocumentByTypeAndCode(ctx, code, int(insuranceDocTemplate.GetTypeID()))
		if err != nil {
			logger.Info("insurance document is different")
			isADifferentDocument = true
		}

		insuranceFileID, err := a.updateAutogeneratedDoc(ctx, code, insuranceDoc)
		if err != nil {
			return fmt.Errorf("UploadDocuments: error sending new documents to GEDO: %w", err)
		}

		if isADifferentDocument {
			oldType := file.WithoutInsurance
			if insuranceDocTemplate.GetTypeID() == file.WithoutInsurance {
				oldType = file.UserInsurance
			}

			id, _, err := a.fileRepository.GetDocumentByTypeAndCode(ctx, code, int(oldType))
			if err != nil {
				return fmt.Errorf("UploadDocuments: error getting insurance doc to delete: %w", err)
			}

			err = a.fileRepository.DeleteDocumentByID(ctx, id)
			if err != nil {
				logger.Error(err)
			}

			logger.Info("the previous insurance document was deleted")
		}

		documentsToLink = append(documentsToLink, insuranceFileID)

		if insuranceCode != "" {
			documentsToLink = append(documentsToLink, insuranceCode)
		}
	}

	for _, r := range documentsToLink {
		err := a.linkDocumentWithRetry(code, r, 5)
		if err != nil {
			msg := fmt.Sprintf("UploadDocuments: error linking document: %s", r)
			logger.Error(msg)
		}
	}

	return nil
}

func (a *fileUseCase) ProcessDocuments(ctx context.Context, code string, user user.User, memo file.Memo) error {
	logger := log.FromContext(ctx)
	logger.Info("Documents generation started")

	err := a.generateNewDocuments(ctx, code, user, memo)
	if err != nil {
		logger.Error(fmt.Sprintf("error generating documents for ID %s: %v", code, err))
		return err
	}

	return nil
}

func (a *fileUseCase) SendDocumentToRecord(ctx context.Context, code, content, reference string, docType file.DocumentTypeID, username string) error {
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("Adding new document on record: %s", code))

	newDocument := file.Document{
		Name:    reference,
		Content: content,
		TypeID:  docType,
		Metadata: file.DocumentMetadata{
			DocumentType: "IF",
			Reference:    reference,
			OriginSystem: file.OriginSystem,
			FullName:     username,
			Position:     "Municipalidad",
			Department:   "De San Iisdro",
		},
	}

	fileSigned, err := a.fileClient.CreateGEDO(ctx, newDocument)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("error sending document: %w", err)
	}

	if fileSigned.Number == "" {
		return fmt.Errorf("error sending document: api response error")
	}

	base64Signed, err := a.fileClient.DownloadGEDO(ctx, fileSigned.Number)
	if err != nil {
		return fmt.Errorf("error getting doc from api: %w", err)
	}
	fileSigned.Content = base64Signed

	if err := a.saveOrUpdateDocument(ctx, code, fileSigned); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("linking documents on record: %s", code))
	if err := a.linkDocumentWithRetry(code, fileSigned.Number, 5); err != nil {
		return fmt.Errorf("error linking document: %w", err)
	}

	logger.Info(fmt.Sprintf("updating done: %s!", code))

	return nil
}

func (a *fileUseCase) saveOrUpdateDocument(ctx context.Context, code string, fileSigned file.SignedDocument) error {
	idDoc, _, err := a.fileRepository.GetDocumentByTypeAndCode(ctx, code, int(fileSigned.TypeID))
	if err != nil {
		err = a.fileRepository.SaveDocument(ctx, code, fileSigned)
		if err != nil {
			return fmt.Errorf("error saving doc: %w", err)
		}
	} else {
		err := a.fileRepository.UpdateDocumentByTypeID(ctx, idDoc, fileSigned)
		if err != nil {
			return fmt.Errorf("error updating doc in repository: %w", err)
		}
	}

	return nil
}

func (a *fileUseCase) sendDocsToAPIAsync(ctx context.Context, code string, documents []file.Document) error {
	logger := log.FromContext(ctx)

	var wg sync.WaitGroup
	errorChan := make(chan error, len(documents))

	for _, doc := range documents {
		wg.Add(1)
		go func(doc file.Document) {
			defer wg.Done()
			if doc.Content == "" {
				logger.Error("empty doc content...")
				return
			}

			fileSigned, err := a.fileClient.CreateGEDO(ctx, doc)
			if err != nil {
				logger.Error(err.Error())
				errorChan <- fmt.Errorf("error sending document: %w", err)
				// TODO: SQS
				// a.queue.SendMessage(ctx, code, fileSigned.TypeID)
			}

			if fileSigned.Number == "" {
				errorChan <- fmt.Errorf("GDE number empty: api response error")
			}

			if isAddressDocument(doc) || isInsuranceDocument(doc) {
				base64Signed, err := a.fileClient.DownloadGEDO(ctx, fileSigned.Number)
				if err != nil {
					errorChan <- fmt.Errorf("error getting doc from api: %w", err)
				}
				fileSigned.Content = base64Signed
			}

			err = a.fileRepository.SaveDocument(ctx, code, fileSigned)
			if err != nil {
				errorChan <- fmt.Errorf("error saving doc: %w", err)
			}
		}(doc)
	}

	wg.Wait()
	close(errorChan)

	var collectedErrors []string
	for err := range errorChan {
		collectedErrors = append(collectedErrors, err.Error())
	}

	if len(collectedErrors) > 0 {
		return fmt.Errorf("encountered errors: %s", strings.Join(collectedErrors, "; "))
	}

	return nil
}

func (a *fileUseCase) updateAutogeneratedDoc(ctx context.Context, code string, doc file.Document) (string, error) {
	logger := log.FromContext(ctx)

	if doc.Content == "" {
		logger.Error("empty doc content...")
		return "", nil
	}

	fileSigned, err := a.fileClient.CreateGEDO(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return "", fmt.Errorf("error calling GDE api: %w", err)
		// TODO: SQS
		// a.queue.SendMessage(ctx, code, fileSigned.TypeID)
	}

	if fileSigned.Number == "" {
		return "", fmt.Errorf("error sending document: api response error")
	}

	base64Signed, err := a.fileClient.DownloadGEDO(ctx, fileSigned.Number)
	if err != nil {
		return fileSigned.Number, fmt.Errorf("error getting doc from api: %w", err)
	}
	fileSigned.Content = base64Signed

	if err := a.saveOrUpdateDocument(ctx, code, fileSigned); err != nil {
		return fileSigned.Number, fmt.Errorf("error saving doc: %w", err)
	}

	return fileSigned.Number, nil
}

func (a *fileUseCase) generateNewDocuments(ctx context.Context, code string, userData user.User, memo file.Memo) error {
	logger := log.FromContext(ctx)
	logger.Info("Creating new documents on record")

	var documents []file.Document

	templates := []docprocessor.DocumentTemplate{
		&docprocessor.RequestStart{BaseDocument: docprocessor.BaseDocument{FilePath: "./assets/InicioSolicitud.txt"}},
		&docprocessor.SwornStatement{BaseDocument: docprocessor.BaseDocument{FilePath: "./assets/DecJurada.txt"}},
		&docprocessor.TermsAndConditions{BaseDocument: docprocessor.BaseDocument{FilePath: "./assets/TerminosYCondiciones.txt"}},
		&docprocessor.PropertyTaxVerification{BaseDocument: docprocessor.BaseDocument{FilePath: "./assets/VeriFiscalInmuebleDeclaradoABL.txt"}},
	}

	for _, template := range templates {
		doc, err := docprocessor.ProcessDocument(ctx, userData, template)
		if err != nil {
			msg := fmt.Sprintf("error processing document type %d", template.GetTypeID())
			logger.Error(msg)
			return fmt.Errorf("%s: %w", msg, err)
		}
		documents = append(documents, doc)
	}

	var docTypeID []file.DocumentTypeID
	switch userData.Type {
	case user.Admin:
		docTypeID = []file.DocumentTypeID{file.CoOwnership, file.AppointmentCertificate}
	case user.Owner:
		docTypeID = []file.DocumentTypeID{file.PropertyTitle}
	case user.Occupant:
		docTypeID = []file.DocumentTypeID{file.PropertyTitle, file.OwnerAuthorization}
	default:
		docTypeID = nil
	}

	filesMap := make(map[file.DocumentTypeID]string)
	for _, dType := range docTypeID {
		_, fileID, err := a.fileRepository.GetDocumentByTypeAndCode(ctx, code, int(dType))
		if err != nil {
			filesMap[dType] = "error"
		}
		filesMap[dType] = fileID
	}

	addressDocTemplate, err := docprocessor.GetAddressDocument(userData, filesMap)
	if err != nil {
		return fmt.Errorf("error getting address document: %w", err)
	}

	addressDoc, err := docprocessor.ProcessDocument(ctx, userData, addressDocTemplate)
	if err != nil {
		return fmt.Errorf("error proccesing address document: %w", err)
	}
	documents = append(documents, addressDoc)

	ifNumber := ""
	if memo.Insurance {
		_, fileID, err := a.fileRepository.GetDocumentByTypeAndCode(ctx, code, int(file.WithInsurance))
		if err != nil {
			ifNumber = "error"
		}
		ifNumber = fileID
	}

	activities, err := a.fileRepository.GetActivityNameByIDs(ctx, memo.Tasks)
	if err != nil {
		activities = "error"
	}

	insuranceDocTemplate, err := docprocessor.GetInsuranceDocument(userData, memo.Insurance, ifNumber, activities, memo.Description, fmt.Sprintf("%d", memo.Time))
	if err != nil {
		return fmt.Errorf("error getting address document: %w", err)
	}

	insuranceDoc, err := docprocessor.ProcessDocument(ctx, userData, insuranceDocTemplate)
	if err != nil {
		return fmt.Errorf("error proccesing insurance document: %w", err)
	}
	documents = append(documents, insuranceDoc)

	return a.sendDocsToAPIAsync(ctx, code, documents)
}

func (a *fileUseCase) LinkDocuments(ctx context.Context, code string, userType user.UserType, insurance bool) error {
	logger := log.FromContext(ctx)
	logger.Info("LinkDocuments documents started")

	order := []int{int(file.Request)}

	switch userType {
	case user.Admin:
		order = append(order, int(file.AddressAdmin), int(file.CoOwnership), int(file.AppointmentCertificate))
	case user.Owner:
		order = append(order, int(file.AddressOwner), int(file.PropertyTitle))
	case user.Occupant:
		order = append(order, int(file.AddressOccupant), int(file.PropertyTitle), int(file.OwnerAuthorization))
	}

	order = append(order, int(file.TaxVerification))

	if insurance {
		order = append(order, int(file.UserInsurance), int(file.WithInsurance))
	} else {
		order = append(order, int(file.WithoutInsurance))
	}

	order = append(order, int(file.TermsAndCond), int(file.Statement))

	ordenMap := make(map[int]int)
	for i, tipo := range order {
		ordenMap[tipo] = i
	}

	documents, err := a.fileRepository.GetDocumentsByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("error finding documents by code in repository: %w", err)
	}

	sort.SliceStable(documents, func(i, j int) bool {
		return ordenMap[int(documents[i].TypeID)] < ordenMap[int(documents[j].TypeID)]
	})

	for _, r := range documents {
		if r.Number == "" {
			logger.Info(fmt.Sprintf("error linking document. Code: %s. Type: %d", code, r.TypeID))
			continue
		}
		err := a.linkDocumentWithRetry(code, r.Number, 5)
		if err != nil {
			msg := fmt.Sprintf("error linking document: %s", r.Number)
			logger.Error(msg)
		}
	}

	return nil
}

func (a *fileUseCase) linkDocumentWithRetry(code string, number string, maxRetries int) error {
	var err error
	for attempts := 0; attempts < maxRetries; attempts++ {
		err := a.fileClient.LinkDocument(code, number)
		if err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to link document after %d retries: %w", maxRetries, err)
}

func (a *fileUseCase) UpdateRequestStatus(ctx context.Context, id int64) error {
	return a.fileRepository.UpdateRequestStatus(ctx, id, 1)
}

func isAddressDocument(doc file.Document) bool {
	if doc.TypeID == file.AddressAdmin || doc.TypeID == file.AddressOccupant || doc.TypeID == file.AddressOwner {
		return true
	}

	return false
}

func isInsuranceDocument(doc file.Document) bool {
	if doc.TypeID == file.WithoutInsurance || doc.TypeID == file.UserInsurance {
		return true
	}

	return false
}
