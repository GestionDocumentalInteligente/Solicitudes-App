package domain

type Verification struct {
	ID                int64
	VerificationCase  string
	RecordNumber      string
	RequestType       string
	DocumentType      string
	DeliveryDate      string
	Status            string
	StatusTask        string
	StatusProperty    string
	RequesterFullName string
	RequesterCuil     string
	RequesterAddress  string
	Documents         []Document
}

type Document struct {
	ID           int64
	Type         int
	Title        string
	GedoCode     string
	Content      string
	VerifiedBy   string
	VerifiedDate string
}
