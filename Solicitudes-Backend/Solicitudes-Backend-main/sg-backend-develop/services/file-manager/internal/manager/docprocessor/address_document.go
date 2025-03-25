package docprocessor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type AddressDocument struct {
	BaseDocument
	UserType  user.UserType
	Variables map[string]string
}

func (d *AddressDocument) GetDocumentType() string {
	return file.IfTypeDocument
}

func (d *AddressDocument) GetTypeID() file.DocumentTypeID {
	switch d.UserType {
	case user.Admin:
		return file.AddressAdmin
	case user.Owner:
		return file.AddressOwner
	case user.Occupant:
		return file.AddressOccupant
	default:
		return file.AddressOccupant
	}
}

func (d *AddressDocument) ReplacePlaceholders(content []byte, user user.User) string {
	text := string(content)
	for placeholder, value := range d.Variables {
		text = strings.ReplaceAll(text, placeholder, value)
	}
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func (d *AddressDocument) GetMetadata(userData user.User) file.DocumentMetadata {
	return file.DocumentMetadata{
		DocumentType: d.GetDocumentType(),
		Reference:    "Domicilio declarado y acreditación de titularidad o legitimación",
		OriginSystem: file.OriginSystem,
		FullName:     user.GetFullName(userData),
		Position:     fmt.Sprintf("%d", userData.DocumentNumber),
		Department:   "Ciudadano",
	}
}

func GetAddressDocument(userData user.User, docs map[file.DocumentTypeID]string) (DocumentTemplate, error) {
	variables := map[string]string{
		"{street}":    userData.Address.Street,
		"{number}":    userData.Address.Number,
		"{ablNumber}": fmt.Sprintf("%d", userData.Address.ABLNumber),
	}

	// Seleccionar el archivo correcto según el tipo de usuario
	var filePath string
	switch userData.Type {
	case user.Admin:
		filePath = "./assets/address/AdminZonaComun.txt"
		variables["{ifNumberRule}"] = docs[file.CoOwnership]
		variables["{ifNumberAdmin}"] = docs[file.AppointmentCertificate]
	case user.Owner:
		filePath = "./assets/address/Titular.txt"
		variables["{ifNumber}"] = docs[file.PropertyTitle]
	case user.Occupant:
		filePath = "./assets/address/InquilinoOcupanteLegitimo.txt"
		variables["{ifNumberTitle}"] = docs[file.PropertyTitle]
		variables["{ifNumberAuth}"] = docs[file.OwnerAuthorization]
	default:
		return nil, errors.New("invalid user_type: must be 'admin', 'owner', or 'occupant'")
	}

	return &AddressDocument{
		BaseDocument: BaseDocument{FilePath: filePath},
		UserType:     userData.Type,
		Variables:    variables,
	}, nil
}
