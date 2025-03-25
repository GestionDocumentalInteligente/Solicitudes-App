package docprocessor

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type InsuranceDocument struct {
	BaseDocument
	Insurance bool
	Variables map[string]string
}

func (d *InsuranceDocument) GetDocumentType() string {
	return file.IfTypeDocument
}

func (d *InsuranceDocument) GetTypeID() file.DocumentTypeID {
	if d.Insurance {
		return file.UserInsurance
	}

	return file.WithoutInsurance
}

func (d *InsuranceDocument) ReplacePlaceholders(content []byte, user user.User) string {
	text := string(content)
	for placeholder, value := range d.Variables {
		text = strings.ReplaceAll(text, placeholder, value)
	}
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func GetInsuranceDocument(user user.User, insurance bool, ifNumber, tasks, memory, time string) (DocumentTemplate, error) {
	variables := map[string]string{
		"{tasks}":              tasks,
		"{projectDescription}": memory,
		"{plazo}":              time,
	}

	filePath := "./assets/insurance/SinPoliza.txt"
	if insurance {
		variables["{ifNumber}"] = ifNumber
		filePath = "./assets/insurance/Poliza.txt"
	}

	return &InsuranceDocument{
		BaseDocument: BaseDocument{FilePath: filePath},
		Insurance:    insurance,
		Variables:    variables,
	}, nil
}

func (d *InsuranceDocument) GetMetadata(userData user.User) file.DocumentMetadata {
	return file.DocumentMetadata{
		DocumentType: d.GetDocumentType(),
		Reference:    "Actividades declaradas para el otorgamiento del Aviso de Obra",
		OriginSystem: file.OriginSystem,
		FullName:     user.GetFullName(userData),
		Position:     fmt.Sprintf("%d", userData.DocumentNumber),
		Department:   "Ciudadano",
	}
}
