package file

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type fileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) Repository {
	return &fileRepository{db: db}
}

type documentDAO struct {
	ID     int    `db:"id"`
	Code   string `db:"code"`
	FileID string `db:"file_id"`
	TypeID int    `db:"document_type_id"`
	Status string `db:"status"`
}

func (r *fileRepository) GetDocumentByTypeAndCode(ctx context.Context, code string, docType int) (int, string, error) {
	query := `
		SELECT id, code, file_id, document_type_id, status
		FROM documents
		WHERE code = $1 AND document_type_id = $2
	`

	doc := &documentDAO{}
	err := r.db.QueryRowContext(ctx, query,
		code, docType,
	).Scan(
		&doc.ID,
		&doc.Code,
		&doc.FileID,
		&doc.TypeID,
		&doc.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return doc.ID, "", NotFound("document not found")
		}
		return doc.ID, "", fmt.Errorf("error getting document: %w", err)
	}

	return doc.ID, doc.FileID, nil
}

func (r *fileRepository) SaveDocument(ctx context.Context, code string, doc SignedDocument) error {
	query := `
		INSERT INTO documents (
			code, file_id, file_url, document_type_id, content, status, original_content, filename
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		code, doc.Number, doc.URL, doc.TypeID, doc.Content, doc.Status, doc.OriginalContent, doc.Filename,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("file already exists")
		}
		return fmt.Errorf("error creating file: %w", err)
	}
	return nil
}

func (r *fileRepository) UpdateDocument(ctx context.Context, content string, id int) error {
	query := `
		UPDATE documents 
		SET content = $1 
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, content, id)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	return nil
}

func (r *fileRepository) UpdateDocumentByFileID(ctx context.Context, content, fileID string) error {
	query := `
		UPDATE documents 
		SET content = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE file_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, content, fileID)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	return nil
}

func (r *fileRepository) UpdateDocumentByTypeID(ctx context.Context, id int, doc SignedDocument) error {
	query := `
		UPDATE documents 
		SET content = $1, file_id = $2, original_content = $3, filename = $4, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, doc.Content, doc.Number, doc.OriginalContent, doc.Filename, id)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	return nil
}

func (r *fileRepository) UpdateRequestStatus(ctx context.Context, id int64, status int) error {
	query := `
		UPDATE requests 
		SET status_id = $1 
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("error updating request: %w", err)
	}

	return nil
}

func (r *fileRepository) GetActivityNameByIDs(ctx context.Context, ids []int) (string, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT name
		FROM activities
		WHERE id IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		names = append(names, name)
	}

	if len(names) == 0 {
		return "", nil
	}

	return strings.Join(names, ", "), nil
}

func (r *fileRepository) GetDocumentsByCode(ctx context.Context, code string) ([]SignedDocument, error) {
	query := `
		SELECT file_id, document_type_id
		FROM documents
		WHERE code = $1
	`

	rows, err := r.db.QueryContext(ctx, query, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFound("document not found")
		}
		return nil, fmt.Errorf("error getting document: %w", err)
	}
	defer rows.Close()

	var documents []SignedDocument
	for rows.Next() {
		var doc SignedDocument
		if err := rows.Scan(&doc.Number, &doc.TypeID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (r *fileRepository) DeleteDocumentByID(ctx context.Context, id int) error {
	query := `DELETE FROM documents WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error getting document: %w", err)
	}

	return nil
}
