package docprocessor

import (
	"fmt"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type TermsAndConditions struct {
	BaseDocument
}

func (d *TermsAndConditions) GetDocumentType() string {
	return file.IfTypeDocument
}

func (d *TermsAndConditions) GetTypeID() file.DocumentTypeID {
	return file.TermsAndCond
}

func (d *TermsAndConditions) GetMetadata(userData user.User) file.DocumentMetadata {
	return file.DocumentMetadata{
		DocumentType: d.GetDocumentType(),
		Reference:    "Términos y Condiciones",
		OriginSystem: file.OriginSystem,
		FullName:     user.GetFullName(userData),
		Position:     fmt.Sprintf("%d", userData.DocumentNumber),
		Department:   "Ciudadano",
	}
}
