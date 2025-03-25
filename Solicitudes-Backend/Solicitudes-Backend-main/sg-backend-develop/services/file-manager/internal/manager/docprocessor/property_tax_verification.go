package docprocessor

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type PropertyTaxVerification struct {
	BaseDocument
}

func (d *PropertyTaxVerification) ReplacePlaceholders(content []byte, user user.User) string {
	text := string(content)
	variables := map[string]string{
		"{ablNumber}": fmt.Sprintf("%d", user.Address.ABLNumber),
	}
	for placeholder, value := range variables {
		text = strings.ReplaceAll(text, placeholder, value)
	}
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func (d *PropertyTaxVerification) GetDocumentType() string {
	return file.IfTypeDocument
}

func (d *PropertyTaxVerification) GetTypeID() file.DocumentTypeID {
	return file.TaxVerification
}

func (d *PropertyTaxVerification) GetMetadata(user user.User) file.DocumentMetadata {
	return file.DocumentMetadata{
		DocumentType: d.GetDocumentType(),
		Reference:    "Verificación de la Situación Fiscal del Inmueble Declarado (ABL)",
		OriginSystem: file.OriginSystem,
		FullName:     "Tramites Electronicos San Isidro",
		Position:     "Administrativa",
		Department:   "Mesa Digital San Iisdro",
	}
}
