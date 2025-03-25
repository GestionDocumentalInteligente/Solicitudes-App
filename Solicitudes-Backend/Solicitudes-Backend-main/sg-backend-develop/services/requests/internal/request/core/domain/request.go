package domain

import "time"

type DocumentTypeID int
type UserType string

// DocumentRequest represents a document in the domain
type DocumentRequest struct {
	Name    string
	Type    DocumentTypeID
	Content string
}

// Address represents an address in the domain
type Address struct {
	Street    string
	Number    string
	ABLNumber int64
}

type VerifiedRequest struct {
	FileNumber                string
	Cuil                      string
	VerificationType          string
	Observations              string
	FinalVerificationDocument string
	Reference                 string
	SelectedActivities        []int
}

type ValidateRequest struct {
	FileNumber  string
	Cuil        string
	IsValid     bool
	FileContent string
}

// Request represents the main request entity in the domain
type Request struct {
	ID                 int64
	UserID             int64
	PropertyID         int64
	Cuil               string
	Dni                string
	FirstName          string
	LastName           string
	Email              string
	Phone              string
	Address            Address
	ABLDebt            string
	CommonZone         bool
	UserType           UserType
	SelectedActivities []int
	Activities         []string
	ProjectDesc        string
	EstimatedTime      int64
	Insurance          bool
	FileNumber         string
	Documents          []DocumentRequest
	StatusName         string
	CreatedAt          time.Time
	Observations       string
	ObservationsTasks  string
	VerifyBy           string
	VerifyByTasks      string
	VerifyDate         time.Time
	VerifyDateTask     time.Time
}

// User type constants
const (
	UserTypeAdmin    UserType = "Admin"
	UserTypeUser     UserType = "User"
	UserTypeOwner    UserType = "Owner"
	UserTypeOccupant UserType = "Occupant"
)

/////

type Suggestion struct {
	ID         int64
	AddrStreet string
	AddrNum    int64
	AblNum     int64
}
