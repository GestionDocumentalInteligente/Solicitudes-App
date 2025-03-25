package transport

import "database/sql"

type DocumentModel struct {
	ID          int64          `json:"id"`
	Code        string         `json:"code"`
	FileID      string         `json:"file_id"`
	Filename    sql.NullString `json:"filename"`
	Description sql.NullString `json:"description"`
	Content     string         `json:"content"`
	Type        int            `json:"document_type_id"`
	Activities  []string       `json:"activities"`
	Status      string         `json:"status"`
}
