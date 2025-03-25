package docprocessor

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type RequestStart struct {
	BaseDocument
}

func (d *RequestStart) GetTypeID() file.DocumentTypeID {
	return file.Request
}

func (d *RequestStart) GetDocumentType() string {
	return file.IfTypeDocument
}

func (d *RequestStart) ReplacePlaceholders(content []byte, userData user.User) string {
	text := string(content)
	variables := map[string]string{
		"{name}":  user.GetFullName(userData),
		"{cuil}":  fmt.Sprintf("%d", userData.Cuil),
		"{email}": userData.Email,
		"{phone}": fmt.Sprintf("%d", userData.Phone),
	}
	for placeholder, value := range variables {
		text = strings.ReplaceAll(text, placeholder, value)
	}
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func (d *RequestStart) GetMetadata(user user.User) file.DocumentMetadata {
	reference := "Inicio de solicitud Aviso de Obra"

	return file.DocumentMetadata{
		DocumentType: d.GetDocumentType(),
		Reference:    reference,
		OriginSystem: file.OriginSystem,
		FullName:     "Tramites Electronicos San Isidro",
		Position:     "Administrativa",
		Department:   "Mesa Digital San Iisdro",
	}
}
