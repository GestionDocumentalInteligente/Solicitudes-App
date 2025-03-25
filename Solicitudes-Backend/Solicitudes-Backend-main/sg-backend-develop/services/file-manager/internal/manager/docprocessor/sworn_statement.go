package docprocessor

import (
	"fmt"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type SwornStatement struct {
	BaseDocument
}

func (d *SwornStatement) GetDocumentType() string {
	return file.IfTypeDocument
}

func (d *SwornStatement) GetTypeID() file.DocumentTypeID {
	return file.Statement
}

func (d *SwornStatement) GetMetadata(userData user.User) file.DocumentMetadata {
	return file.DocumentMetadata{
		DocumentType: d.GetDocumentType(),
		Reference:    "Declaración Jurada sobre la solicitud de Aviso de Obra",
		OriginSystem: file.OriginSystem,
		FullName:     user.GetFullName(userData),
		Position:     fmt.Sprintf("%d", userData.DocumentNumber),
		Department:   "Ciudadano",
	}
}
