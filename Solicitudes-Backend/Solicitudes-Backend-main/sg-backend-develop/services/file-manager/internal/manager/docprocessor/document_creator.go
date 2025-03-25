package docprocessor

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type DocumentTemplate interface {
	ReadFile() ([]byte, error)
	ReplacePlaceholders(content []byte, user user.User) string
	GetDocumentType() string
	GetTypeID() file.DocumentTypeID
	GetMetadata(user user.User) file.DocumentMetadata
}

type BaseDocument struct {
	FilePath string
}

func (d *BaseDocument) ReadFile() ([]byte, error) {
	return os.ReadFile(d.FilePath)
}

func (d *BaseDocument) ReplacePlaceholders(content []byte, user user.User) string {
	return base64.StdEncoding.EncodeToString(content)
}

func ProcessDocument(ctx context.Context, user user.User, docTemplate DocumentTemplate) (file.Document, error) {
	content, err := docTemplate.ReadFile()
	if err != nil {
		return file.Document{}, fmt.Errorf("error reading document: %w", err)
	}

	encodedContent := docTemplate.ReplacePlaceholders(content, user)
	metadata := docTemplate.GetMetadata(user)

	document := file.Document{
		Name:     metadata.Reference,
		TypeID:   docTemplate.GetTypeID(),
		Metadata: metadata,
		Content:  encodedContent,
	}

	return document, nil
}
