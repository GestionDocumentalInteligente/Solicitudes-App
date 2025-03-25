package file

import (
	"context"
	"time"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

const IfGraTypeDocument string = "IFGRA"
const IfTypeDocument string = "IF"
const OriginSystem string = "SG-MSI"

type DocumentTypeID int

const (
	AddressAdmin           DocumentTypeID = 1  // "admin"
	AddressOwner           DocumentTypeID = 2  // "owner"
	AddressOccupant        DocumentTypeID = 3  // "occupant"
	Request                DocumentTypeID = 4  // "request_start"
	Cover                  DocumentTypeID = 5  // "cover"
	Statement              DocumentTypeID = 6  // "statement"
	TermsAndCond           DocumentTypeID = 7  // "terms_and_conditions"
	TaxVerification        DocumentTypeID = 8  // "tax_verification"
	CoOwnership            DocumentTypeID = 9  // "co_ownership_regulation_appointment certificate"
	PropertyTitle          DocumentTypeID = 10 // "property_title_or_ownership_report"
	OwnerAuthorization     DocumentTypeID = 11 // "owner_authorization"
	WithInsurance          DocumentTypeID = 12 // "insurance"
	WithoutInsurance       DocumentTypeID = 13 // "not"
	AppointmentCertificate DocumentTypeID = 14
	UserInsurance          DocumentTypeID = 15 // "insurance_signed"
	VerificationTasks      DocumentTypeID = 16
	VerificationProperty   DocumentTypeID = 17 // "insurance_signed"
	Administrative         DocumentTypeID = 18
)

var DocumentTypeDescriptionMap = map[DocumentTypeID]string{
	CoOwnership:            "Reglamento de Co-propiedad o Acta de Designación",
	PropertyTitle:          "Título de Propiedad o Informe de Dominio",
	OwnerAuthorization:     "Autorización del Propietario",
	WithInsurance:          "Poliza de Seguro",
	AppointmentCertificate: "Acta de Designación",
}

type Credentials struct {
	Email    string
	Password string
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Document struct {
	Metadata    DocumentMetadata
	TypeID      DocumentTypeID
	Name        string
	Content     string
	ContentType string
}

type DocumentMetadata struct {
	ID             string
	DocumentType   string
	Reference      string
	OriginSystem   string
	FullName       string
	Position       string
	Department     string
	CreatedAt      time.Time
	DocumentStatus string
	LastModified   time.Time
}

type SignedDocument struct {
	TypeID          DocumentTypeID
	Filename        string
	OriginalContent string
	Content         string
	Number          string
	URL             string
	SpecialNumber   string
	Licence         string
	Status          bool
}

type Memo struct {
	Insurance   bool
	Description string
	Tasks       []int
	Time        int64
}

type Repository interface {
	SaveDocument(ctx context.Context, code string, doc SignedDocument) error
	GetDocumentByTypeAndCode(ctx context.Context, code string, docType int) (int, string, error)
	GetDocumentsByCode(ctx context.Context, code string) ([]SignedDocument, error)
	GetActivityNameByIDs(ctx context.Context, id []int) (string, error)
	UpdateDocument(ctx context.Context, content string, id int) error
	UpdateDocumentByFileID(ctx context.Context, content string, id string) error
	UpdateDocumentByTypeID(ctx context.Context, idDoc int, doc SignedDocument) error
	UpdateRequestStatus(ctx context.Context, id int64, status int) error
	DeleteDocumentByID(ctx context.Context, id int) error
}

type Client interface {
	CreateGEDO(ctx context.Context, documentInfo Document) (SignedDocument, error)
	DownloadGEDO(ctx context.Context, fileID string) (string, error)
	CreateRecord(ctx context.Context, user user.User, id int64) (string, error)
	LinkDocument(code, fileID string) error
}

type MessageQueue interface {
	SendMessage(ctx context.Context, code string, typeID DocumentTypeID) error
}
