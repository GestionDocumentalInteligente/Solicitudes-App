package outbound

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"

	sdkpg "github.com/teamcubation/sg-backend/pkg/databases/sql/postgresql/pgxpool"
	sdkdefs "github.com/teamcubation/sg-backend/pkg/databases/sql/postgresql/pgxpool/defs"

	transport "github.com/teamcubation/sg-backend/services/requests/internal/request/adapters/outbound/transport"
	domain "github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
	ports "github.com/teamcubation/sg-backend/services/requests/internal/request/core/ports"
)

const taskType = "tasks"

type PostgreSQL struct {
	repository sdkdefs.Repository
}

func NewPostgreSQL() (ports.Repository, error) {
	r, err := sdkpg.Bootstrap()
	if err != nil {
		return nil, fmt.Errorf("bootstrap error: %w", err)
	}

	return &PostgreSQL{
		repository: r,
	}, nil
}

func (r *PostgreSQL) GetAllRequestsByCuil(ctx context.Context, cuil string) ([]domain.Request, error) {

	userID, err := r.findUserIDFromCuil(ctx, cuil)
	if err != nil {
		return nil, err
	}

	reqs, err := r.GetAllRequestsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (r *PostgreSQL) GetAllRequestsByUserID(ctx context.Context, userID int64) ([]domain.Request, error) {
	const query = `
        SELECT 
            r.id,
            r.request_type_id,
            r.property_id,
            r.status_id,
            rs.name as status_name,
            r.description,
            r.verification_complete,
            r.verification_date,
            r.created_at,
            r.updated_at,
            COALESCE(r.file_number, '') as file_number,
            COALESCE(r.abl_debt, '') as abl_debt,
            COALESCE(r.selected_activities, ARRAY[]::integer[]) as selected_activities,
            COALESCE(r.estimated_time, 0) as estimated_time,
            COALESCE(r.insurance, FALSE) as insurance
        FROM 
            public.requests r
            LEFT JOIN public.request_status rs ON r.status_id = rs.id
        WHERE 
            r.user_id = $1
        ORDER BY 
            r.created_at DESC`

	rows, err := r.repository.Pool().Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying requests: %w", err)
	}
	defer rows.Close()

	var requests []transport.RequestDataModel
	for rows.Next() {
		var req transport.RequestDataModel
		if err := rows.Scan(
			&req.ID,
			&req.RequestTypeID,
			&req.PropertyID,
			&req.StatusID,
			&req.StatusName,
			&req.Description,
			&req.VerificationComplete,
			&req.VerificationDate,
			&req.CreatedAt,
			&req.UpdatedAt,
			&req.FileNumber,
			&req.ABLDebt,
			&req.SelectedActivities,
			&req.EstimatedTime,
			&req.Insurance,
		); err != nil {
			return nil, fmt.Errorf("error scanning request: %w", err)
		}
		req.UserID = userID
		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transport.ToRequestDomainList(requests), nil
}

func (r *PostgreSQL) GetSuggestions(ctx context.Context, addrName string, addrNum, ablNum int64) ([]domain.Suggestion, error) {
	// Inicializar patrones con "%"
	addrNamePattern := "%"
	addrNumPattern := "%"
	ablNumPattern := "%"

	// Construir patrones basados en los parámetros proporcionados
	if addrName != "" {
		addrNamePattern += addrName + "%"
	}
	if addrNum > 0 {
		addrNumPattern = fmt.Sprintf("%d%%", addrNum)
	}
	if ablNum > 0 {
		ablNumPattern = fmt.Sprintf("%d%%", ablNum)
	}

	// Variables para la consulta y los argumentos
	var (
		query       string
		args        []interface{}
		orderBy     string
		whereClause string
	)

	if addrName != "" {
		whereClause = `
			regexp_replace(unaccent(lower(p.street)), '[^a-z]', '', 'g')
				ILIKE '%' || regexp_replace(unaccent(lower($1)), '[^a-z]', '', 'g') || '%' AND
			REGEXP_REPLACE(CAST(p.number AS TEXT), '[^0-9]', '', 'g') ILIKE '%' || REGEXP_REPLACE($2, '[^0-9]', '', 'g') || '%' AND
			REGEXP_REPLACE(CAST(a.abl_number AS TEXT), '[^0-9]', '', 'g') ILIKE '%' || REGEXP_REPLACE($3, '[^0-9]', '', 'g') || '%'
			`
		args = append(args, addrNamePattern, addrNumPattern, ablNumPattern)
		orderBy = `
			a.abl_number ASC, 
			p.street, 
			p.number`
	} else if ablNum > 0 {
		// Caso 2: addrName está vacío pero ablNum está presente
		whereClause = `
            REGEXP_REPLACE(CAST(a.abl_number AS TEXT), '[^0-9]', '', 'g') ILIKE '%' || REGEXP_REPLACE($1, '[^0-9]', '', 'g') || '%'`
		args = append(args, ablNumPattern)
		orderBy = `
            a.abl_number ASC, 
            p.street, 
            p.number`
	} else if addrNum > 0 {
		// Caso 3: Tanto addrName como ablNum están vacíos, pero addrNum está presente
		whereClause = `
            REGEXP_REPLACE(CAST(p.number AS TEXT), '[^0-9]', '', 'g') ILIKE '%' || REGEXP_REPLACE($1, '[^0-9]', '', 'g') || '%'`
		args = append(args, addrNumPattern)
		orderBy = `
            a.abl_number ASC, 
            p.number, 
            p.street`
	} else {
		// Caso 4: Ningún parámetro está presente
		return nil, fmt.Errorf("no se proporcionaron parámetros de búsqueda")
	}

	// Construir la consulta SQL completa
	query = fmt.Sprintf(`
        SELECT
            TRIM(p.street) AS street,
            p.number AS number, 
            a.abl_number AS abl_number, 
            p.property_id
        FROM properties p
        LEFT JOIN abl a ON p.abl_id = a.abl_id
        WHERE %s
        ORDER BY %s
        LIMIT 100;`, whereClause, orderBy)

	// Ejecutar la consulta
	var suggestions []transport.SuggestionDataModel
	if err := r.repository.SelectContext(ctx, &suggestions, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get suggestions from repository: %w", err)
	}

	return transport.ToSuggestionDomainList(suggestions), nil
}

func (r *PostgreSQL) CheckAblOwnership(ctx context.Context, cuil string, ablNumb int) (bool, error) {
	const query = `
		WITH person_info AS (
			SELECT  
				TRIM(last_name || ' ' || first_name) as full_name,
				cuil
			FROM persons 
			-- WHERE cuil = '20345678901'
			WHERE cuil = $1
		), 
		abl_info AS (
			SELECT abl_id, abl_number
			FROM public.abl 
			WHERE abl_number = $2
		)
		SELECT EXISTS (
			SELECT 1
			FROM person_info p
			CROSS JOIN abl_info ai
			JOIN public.owners o ON o.abl_id = ai.abl_id
			WHERE 
				REGEXP_REPLACE(LOWER(UNACCENT(TRIM(p.full_name))), '[^a-z]', '', 'g') = REGEXP_REPLACE(LOWER(UNACCENT(TRIM(o.primary_owner))), '[^a-z]', '', 'g')
				OR REGEXP_REPLACE(LOWER(UNACCENT(TRIM(p.full_name))), '[^a-z]', '', 'g') = REGEXP_REPLACE(LOWER(UNACCENT(TRIM(o.secondary_owner))), '[^a-z]', '', 'g')
		);`

	var exists bool
	err := r.repository.QueryRowContext(ctx, query, cuil, ablNumb).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking ABL ownership: %w", err)
	}

	return exists, nil
}

// verificacion ciudadano: potestad inmuble propietario
func (r *PostgreSQL) RequestsVerifications(ctx context.Context, cuil string) (*domain.Verification, error) {
	const query = `
		SELECT DISTINCT 
			COALESCE(o.primary_owner, 'Sin propietario') as primary_owner,
			COALESCE(p.street, 'Sin calle') as street,
			COALESCE(p."number"::text, 'Sin número') as number,
			COALESCE(per.cuil, 'Sin CUIL') as cuil,
			COALESCE(per.dni, 'Sin DNI') as dni,
			COALESCE(rt.name, 'Sin tipo de solicitud') as request_type,
			COALESCE(rs.name, 'Sin estado') as request_status,
			COALESCE(d.id::text, 'Sin ID de documento') as document_id,
			COALESCE(d.document_type_id::text, 'Sin tipo de documento') as document_type_id
		FROM persons per
		LEFT JOIN users u ON u.person_id = per.id
		LEFT JOIN requests r ON r.user_id = u.id
		LEFT JOIN request_types rt ON rt.id = r.request_type_id
		LEFT JOIN request_status rs ON rs.id = r.status_id
		LEFT JOIN properties p ON p.property_id = r.property_id
		LEFT JOIN abl a ON a.abl_id = p.abl_id
		LEFT JOIN owners o ON o.abl_id = a.abl_id
		LEFT JOIN documents d ON d.code LIKE CONCAT(r.id, '%')
		WHERE per.cuil = $1
		LIMIT 1`

	var dataModel transport.VerificactionDataModel
	err := r.repository.QueryRowContext(ctx, query, cuil).Scan(
		&dataModel.PropertyOwner,
		&dataModel.AddrStreet,
		&dataModel.AddrNumber,
		&dataModel.CUIL,
		&dataModel.DNI,
		&dataModel.RequestType,
		&dataModel.RequestStatus,
		&dataModel.DocumentID,
		&dataModel.DocumentType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no property owner found for CUIL %s", cuil)
		}
		return nil, fmt.Errorf("error fetching property owner: %w", err)
	}

	return transport.ToVerificationDomain(&dataModel), nil
}

func (r *PostgreSQL) GetRequestByCode(ctx context.Context, code string) (*domain.Request, error) {
	const query = `
		SELECT 
			requests.file_number,
			requests.description,
			requests.insurance,
			requests.estimated_time,
			per.first_name,
			per.last_name,
			per.cuil,	
			p.street, p.number,
			ARRAY_AGG(a.name) AS activities
		FROM requests
		JOIN
			LATERAL UNNEST(requests.selected_activities) AS activity_id ON TRUE
		JOIN
			activities a ON a.id = activity_id 
		INNER JOIN users u ON requests.user_id = u.id
		INNER JOIN persons per ON u.person_id = per.id
		INNER JOIN properties p ON p.property_id = requests.property_id 
		WHERE file_number = $1 
		GROUP BY requests.file_number, requests.description, requests.insurance, estimated_time,
			per.first_name,
			per.last_name,
			per.cuil,	
			p.street, p.number, p.locality  
		LIMIT 1`

	var dataModel transport.RequestDataModel
	var address domain.Address
	var personDataModel transport.RequestPersonDataModel
	err := r.repository.QueryRowContext(ctx, query, code).Scan(
		&dataModel.FileNumber,
		&dataModel.Description,
		&dataModel.Insurance,
		&dataModel.EstimatedTime,
		&personDataModel.FirstName,
		&personDataModel.LastName,
		&personDataModel.Cuil,
		&address.Street,
		&address.Number,
		&dataModel.Activities,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no request found for code %s", code)
		}
		return nil, fmt.Errorf("error fetching request: %w", err)
	}

	return &domain.Request{
		ProjectDesc:   dataModel.Description,
		Insurance:     dataModel.Insurance,
		EstimatedTime: dataModel.EstimatedTime,
		FirstName:     personDataModel.FirstName,
		LastName:      personDataModel.LastName,
		Cuil:          personDataModel.Cuil,
		Address: domain.Address{
			Street: address.Street,
			Number: address.Number,
		},
		Activities: dataModel.Activities,
	}, nil
}

func (r *PostgreSQL) GetAllRequestsVerifications(ctx context.Context) ([]domain.Verification, error) {
	const query = `
		SELECT  
			requests.id,
			requests.file_number,
			COALESCE(rt.description, 'Aviso de obra') as request_type,
			requests.created_at as deliveryDate,	
			COALESCE(rs.name, 'Pending') as status,		
			COALESCE(st.name, 'Pending') as status_tasks,
			COALESCE(sp.name, 'Pending') as status_property,
			per.first_name,
			per.last_name,
			per.cuil,
			p.street, p.number, p.locality
		FROM requests
		INNER JOIN users u ON requests.user_id = u.id
		INNER JOIN persons per ON u.person_id = per.id
		LEFT JOIN request_types rt ON rt.id = requests.request_type_id
		INNER JOIN request_status rs ON rs.id = requests.status_id
		LEFT JOIN request_status st ON st.id = requests.status_id_tasks
		LEFT JOIN request_status sp ON sp.id = requests.status_id_property
		LEFT JOIN properties p ON p.property_id = requests.property_id WHERE requests.status_id != 8`

	rows, err := r.repository.Pool().Query(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error querying requests: %w", err)
	}
	defer rows.Close()

	var requests []transport.VerificactionDataModel
	for rows.Next() {
		var req transport.VerificactionDataModel
		if err := rows.Scan(
			&req.ID,
			&req.FileNumber,
			&req.RequestType,
			&req.DeliveryDate,
			&req.Status,
			&req.StatusTask,
			&req.StatusProperty,
			&req.FirstName,
			&req.LastName,
			&req.CUIL,
			&req.AddrStreet,
			&req.AddrNumber,
			&req.Locality,
		); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error scanning request: %w", err)
		}
		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transport.ToVerificationListDomain(requests), nil
}

func (r *PostgreSQL) GetAllRequestsValidations(ctx context.Context) ([]domain.Verification, error) {
	const query = `
		SELECT  
			requests.id,
			requests.file_number,
			COALESCE(rt.description, 'Aviso de obra') as request_type,
			requests.created_at as deliveryDate,	
			COALESCE(rs.name, 'Pending') as status,		
			per.first_name,
			per.last_name,
			per.cuil,
			p.street, p.number, p.locality
		FROM requests
		INNER JOIN users u ON requests.user_id = u.id
		INNER JOIN persons per ON u.person_id = per.id
		LEFT JOIN request_types rt ON rt.id = requests.request_type_id
		INNER JOIN request_status rs ON rs.id = requests.status_id
		LEFT JOIN properties p ON p.property_id = requests.property_id 
		WHERE requests.status_id_property = 3 AND requests.status_id_tasks = 3 AND requests.status_id = 1 
		AND requests.file_number IN (
			SELECT code 
			FROM documents 
			WHERE document_type_id IN (16, 17)
			GROUP BY code
			HAVING COUNT(DISTINCT document_type_id) = 2
		)`

	rows, err := r.repository.Pool().Query(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error querying requests: %w", err)
	}
	defer rows.Close()

	var requests []transport.VerificactionDataModel
	for rows.Next() {
		var req transport.VerificactionDataModel
		if err := rows.Scan(
			&req.ID,
			&req.FileNumber,
			&req.RequestType,
			&req.DeliveryDate,
			&req.Status,
			&req.FirstName,
			&req.LastName,
			&req.CUIL,
			&req.AddrStreet,
			&req.AddrNumber,
			&req.Locality,
		); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error scanning request: %w", err)
		}
		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transport.ToVerificationListDomain(requests), nil
}

func (r *PostgreSQL) GetRequestByFileNumber(ctx context.Context, id string) (*domain.Request, error) {
	const query = `
	SELECT  
		requests.file_number,
		requests.description,
		requests.insurance,
		estimated_time,
		per.first_name,
		per.last_name,
		per.cuil,
		per.email,
		perv.last_name as verify_by_name,
		pervt.last_name as verify_by_name_tasks,
		requests.verification_date_tasks,
		requests.verification_date,
		p.street, p.number,	
		ARRAY_AGG(a.name) AS activities
	FROM requests
	JOIN
		LATERAL UNNEST(requests.selected_activities) AS activity_id ON TRUE
	JOIN
		activities a ON a.id = activity_id
	INNER JOIN users u ON requests.user_id = u.id
	INNER JOIN persons per ON u.person_id = per.id
	LEFT JOIN users uv ON requests.verified_by = uv.id
	LEFT JOIN persons perv ON uv.person_id = perv.id
	LEFT JOIN users uvt ON requests.verified_by_tasks = uvt.id
	LEFT JOIN persons pervt ON uvt.person_id = pervt.id
	LEFT JOIN properties p ON p.property_id = requests.property_id 
	WHERE file_number = $1 
	GROUP BY requests.file_number, requests.description, requests.insurance, estimated_time,
		per.first_name, per.last_name, per.cuil, per.email, 
		p.street, p.number, 
		verify_by_name, verify_by_name_tasks, verification_date_tasks, verification_date
	LIMIT 1`

	var dataModel transport.RequestDataModel
	var addressStreet sql.NullString
	var addressNumber sql.NullInt64

	var verifyBy sql.NullString
	var verifyByTasks sql.NullString
	var verifyDate sql.NullTime
	var verifyDateTasks sql.NullTime

	var personDataModel transport.RequestPersonDataModel
	err := r.repository.QueryRowContext(ctx, query, id).Scan(
		&dataModel.FileNumber,
		&dataModel.Description,
		&dataModel.Insurance,
		&dataModel.EstimatedTime,
		&personDataModel.FirstName,
		&personDataModel.LastName,
		&personDataModel.Cuil,
		&personDataModel.Email,
		&verifyBy,
		&verifyByTasks,
		&verifyDate,
		&verifyDateTasks,
		&addressStreet,
		&addressNumber,
		&dataModel.Activities,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no request found for code %s", id)
		}
		return nil, fmt.Errorf("error fetching request: %w", err)
	}

	addressNumberStr := ""
	if addressNumber.Valid {
		addressNumberStr = fmt.Sprintf("%d", addressNumber.Int64)
	}

	return &domain.Request{
		FileNumber:    dataModel.FileNumber,
		ProjectDesc:   dataModel.Description,
		Insurance:     dataModel.Insurance,
		EstimatedTime: dataModel.EstimatedTime,
		FirstName:     personDataModel.FirstName,
		LastName:      personDataModel.LastName,
		Cuil:          personDataModel.Cuil,
		Email:         personDataModel.Email,
		Address: domain.Address{
			Street: addressStreet.String,
			Number: addressNumberStr,
		},
		Activities:     dataModel.Activities,
		VerifyBy:       verifyBy.String,
		VerifyByTasks:  verifyByTasks.String,
		VerifyDate:     verifyDate.Time,
		VerifyDateTask: verifyDateTasks.Time,
	}, nil
}

func (r *PostgreSQL) GetDocumentsByCode(ctx context.Context, id string) ([]domain.Document, error) {
	const query = `
		SELECT 
			documents.id as id,
			document_type_id,
			document_types.description,
			file_id
		FROM documents 
		INNER JOIN document_types ON documents.document_type_id = document_types.id
		WHERE code = $1 
		AND document_types.description IS NOT NULL AND document_types.description != '' 
		AND file_id IS NOT NULL AND file_id != ''`

	rows, err := r.repository.Pool().Query(ctx, query, id)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error querying requests: %w", err)
	}
	defer rows.Close()

	var docs []transport.DocumentModel
	for rows.Next() {
		var req transport.DocumentModel
		if err := rows.Scan(
			&req.ID,
			&req.Type,
			&req.Description,
			&req.FileID,
		); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error scanning request: %w", err)
		}
		docs = append(docs, req)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transport.ToDocumentDomainList(docs), nil
}

func (r *PostgreSQL) GetValidationDocumentsByCode(ctx context.Context, id string) ([]domain.Document, error) {
	const query = `
		SELECT 
			documents.id as id,
			document_type_id,
			document_types.description,
			file_id
		FROM documents 
		INNER JOIN document_types ON documents.document_type_id = document_types.id
		WHERE code = $1 AND document_type_id IN (1, 2, 3, 17, 13, 15, 16) 
		AND file_id IS NOT NULL AND file_id != '' 
		ORDER BY 
			CASE 
				WHEN document_type_id IN (1, 2, 3) THEN 1
				WHEN document_type_id = 17 THEN 2 
				WHEN document_type_id IN (13, 15) THEN 3 
				WHEN document_type_id = 16 THEN 4 
			END`

	rows, err := r.repository.Pool().Query(ctx, query, id)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error querying requests: %w", err)
	}
	defer rows.Close()

	var docs []transport.DocumentModel
	for rows.Next() {
		var req transport.DocumentModel
		if err := rows.Scan(
			&req.ID,
			&req.Type,
			&req.Description,
			&req.FileID,
		); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error scanning request: %w", err)
		}
		docs = append(docs, req)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transport.ToDocumentDomainList(docs), nil
}

func (r *PostgreSQL) GetReplacementIFDocumentsByCode(ctx context.Context, id string, insurance bool) ([]domain.Document, error) {
	insuranceDocQuery := ""
	if !insurance {
		insuranceDocQuery = " AND document_type_id != 12"
	}
	const query = `
		SELECT 
			documents.id as id,
			document_types.name,
			document_type_id,
			file_id
		FROM documents 
		INNER JOIN document_types ON documents.document_type_id = document_types.id
		WHERE code = $1 AND file_id IS NOT NULL AND file_id != ''`

	rows, err := r.repository.Pool().Query(ctx, query+insuranceDocQuery, id)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error querying requests: %w", err)
	}
	defer rows.Close()

	var docs []transport.DocumentModel
	for rows.Next() {
		var req transport.DocumentModel
		if err := rows.Scan(
			&req.ID,
			&req.Description,
			&req.Type,
			&req.FileID,
		); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error scanning request: %w", err)
		}
		docs = append(docs, req)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transport.ToDocumentDomainList(docs), nil
}

func (r *PostgreSQL) GetInsuranceDocumentByCode(ctx context.Context, id string, docType int) (string, error) {
	const query = `
		SELECT 
			file_id
		FROM documents 
		WHERE code = $1 AND (document_type_id = 13 OR document_type_id = 15) LIMIT 1`

	var fileNumber string
	err := r.repository.QueryRowContext(ctx, query, id).Scan(&fileNumber)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("document not found for code %s", id)
		}
		return "", fmt.Errorf("error fetching request: %w", err)
	}

	return fileNumber, nil
}

func (r *PostgreSQL) GetDocumentByID(ctx context.Context, id string) (domain.Document, error) {
	const query = `
		SELECT 
			documents.id as id,
			file_id,
			COALESCE(document_types.description, 'Detalle Actividades') as description,
			content
		FROM documents 
		INNER JOIN document_types ON documents.document_type_id = document_types.id 
		WHERE file_id = $1 LIMIT 1`

	var model transport.DocumentModel
	err := r.repository.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.FileID,
		&model.Description,
		&model.Content,
	)

	if err != nil {
		return domain.Document{}, fmt.Errorf("error getting request: %w", err)
	}

	return domain.Document{
		ID:       model.ID,
		Title:    model.Description.String,
		GedoCode: model.FileID,
		Content:  model.Content,
	}, nil
}

func (r *PostgreSQL) GetRequestPersonByCuil(ctx context.Context, cuil string) (*domain.Request, error) {
	const query = `
		SELECT cuil, dni, first_name, last_name, email, phone
		FROM persons
		WHERE cuil = $1`

	var model transport.RequestPersonDataModel
	err := r.repository.Pool().QueryRow(ctx, query, cuil).Scan(
		&model.Cuil,
		&model.Dni,
		&model.FirstName,
		&model.LastName,
		&model.Email,
		&model.Phone,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting request: %w", err)
	}

	return transport.RequestPersonDataModelToRequestDomain(&model), nil
}

func (r *PostgreSQL) CreateRequestByCuil(ctx context.Context, req *domain.Request) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}

	userID, err := r.findUserIDFromCuil(ctx, req.Cuil)
	if err != nil {
		return err
	}

	req.UserID = userID

	err = r.CreateRequestByUserID(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgreSQL) CreateRequestByUserID(ctx context.Context, req *domain.Request) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}

	tx, err := r.repository.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	reqDataModel := transport.ToCreateRequestDataModel(req)

	reqDataModel.RequestTypeID = 1
	reqDataModel.StatusID = 8 // Processing

	const query = `
        INSERT INTO requests (
            user_id,
			user_type,
            request_type_id,
            property_id,
            status_id,
            description,
            verification_complete,
            verification_date,
            created_at,
            updated_at,
            file_number,
            abl_debt,
            estimated_time,
            insurance,
            selected_activities
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            $9, $10, $11, $12, $13
        ) RETURNING id`

	var requestID int64
	err = tx.QueryRow(ctx, query,
		reqDataModel.UserID,
		reqDataModel.UserType,
		reqDataModel.RequestTypeID,
		reqDataModel.PropertyID,
		reqDataModel.StatusID,
		reqDataModel.Description,
		reqDataModel.VerificationComplete,
		reqDataModel.VerificationDate,
		reqDataModel.FileNumber,
		reqDataModel.ABLDebt,
		reqDataModel.EstimatedTime,
		reqDataModel.Insurance,
		pq.Array(reqDataModel.SelectedActivities),
	).Scan(&requestID)

	if err != nil {
		return fmt.Errorf("%w: %v", transport.ErrCreateRequest, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	req.ID = requestID

	return nil
}

func (r *PostgreSQL) UpdateRequest(ctx context.Context, id int64, code string) error {
	query := `
		UPDATE requests 
		SET file_number = $1 
		WHERE id = $2
	`

	_, err := r.repository.Pool().Exec(ctx, query, code, id)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	return nil
}

func (r *PostgreSQL) UpdateRequestWithObservations(ctx context.Context, req *domain.VerifiedRequest) (string, string, string, error) {
	userID, err := r.findUserIDFromCuil(ctx, req.Cuil)
	if err != nil {
		return "", "", "", err
	}

	// var currentGlobalStatus int
	// err = r.repository.Pool().QueryRow(ctx, `
	// 	SELECT status_id
	// 	FROM requests
	// 	WHERE file_number = $1`, req.FileNumber).Scan(&currentGlobalStatus)
	// if err != nil {
	// 	return "", "", "", fmt.Errorf("error fetching current global status: %w", err)
	// }

	newStatusID := 9
	statusID := 3
	if req.Observations != "" {
		newStatusID = 7
		statusID = 5
	}

	var set string
	if req.VerificationType == taskType {
		set = fmt.Sprintf("SET observations_tasks = $1, verified_by_tasks = $2, verification_date_tasks = $3, status_id_tasks = %d, status_id = %d", statusID, newStatusID)
	} else {
		set = fmt.Sprintf("SET observations = $1, verified_by = $2, verification_date = $3, status_id_property = %d, status_id = %d", statusID, newStatusID)
	}

	args := []interface{}{req.Observations, userID, time.Now(), req.FileNumber}
	where := "WHERE file_number = $4"
	if len(req.SelectedActivities) > 0 {
		set += ", selected_activities = $4"
		args = []interface{}{req.Observations, userID, time.Now(), pq.Array(req.SelectedActivities), req.FileNumber}
		where = "WHERE file_number = $5"
	}

	query := `
		UPDATE requests
		` + set + `
		` + where + ` RETURNING observations, observations_tasks, user_id`

	var observationsProperty sql.NullString
	var observationsTasks sql.NullString
	var userIDToMail int64
	err = r.repository.Pool().QueryRow(ctx, query, args...).Scan(&observationsProperty, &observationsTasks, &userIDToMail)
	if err != nil {
		return "", "", "", fmt.Errorf("error updating requests: %w", err)
	}

	email, err := r.GetPersonByUserID(ctx, userIDToMail)
	if err != nil {
		return "", "", "", fmt.Errorf("error updating requests: %w", err)
	}

	return observationsProperty.String, observationsTasks.String, email, nil
}

func (r *PostgreSQL) ValidateRequest(ctx context.Context, req *domain.ValidateRequest) (int64, error) {
	userID, err := r.findUserIDFromCuil(ctx, req.Cuil)
	if err != nil {
		return userID, err
	}

	set := "SET validated_by = $1, verification_complete = true, validation_date = $2, status_id = 4"
	if !req.IsValid {
		set = "SET validated_by = $1, validation_date = $2, status_id = 1, status_id_tasks = 1, status_id_property = 1"
	}

	query := `
		UPDATE requests 
		` + set + `
		WHERE file_number = $3 RETURNING user_id`

	// _, err = r.repository.Pool().Exec(ctx, query, userID, time.Now(), req.FileNumber)
	// if err != nil {
	// 	return fmt.Errorf("error validating requests: %w", err)
	// }

	var userIDToMail int64
	err = r.repository.Pool().QueryRow(ctx, query, userID, time.Now(), req.FileNumber).Scan(&userIDToMail)
	if err != nil {
		return userIDToMail, fmt.Errorf("error updating requests: %w", err)
	}

	return userIDToMail, nil
}

func (r *PostgreSQL) UpdateRequestStatus(ctx context.Context, id int64, status int) error {
	query := `
		UPDATE requests 
		SET status_id = $1 
		WHERE id = $2
	`

	_, err := r.repository.Pool().Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	return nil
}

func (r *PostgreSQL) UpdateRequestStatusByFileNumber(ctx context.Context, fileNumber string, status int) error {
	query := `
		UPDATE requests 
		SET status_id = $1 
		WHERE file_number = $2
	`

	_, err := r.repository.Pool().Exec(ctx, query, status, fileNumber)
	if err != nil {
		return fmt.Errorf("error updating request status: %w", err)
	}

	return nil
}

// helpers
func (r *PostgreSQL) findUserIDFromCuil(ctx context.Context, cuil string) (int64, error) {
	const query = `
        SELECT u.id as user_id
        FROM users u
        JOIN persons p ON u.person_id = p.id
        WHERE p.cuil = $1
        AND u.deleted_at IS NULL`

	var userID int64
	err := r.repository.Pool().QueryRow(ctx, query, cuil).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%w: cuil %s", transport.ErrUserNotFound, cuil)
		}
		return 0, fmt.Errorf("error finding user by cuil: %v", err)
	}

	return userID, nil
}

func (r *PostgreSQL) GetRequestByID(ctx context.Context, reqID int64) (*domain.Request, error) {
	const query = `
        SELECT 
            r.id,
            r.user_id,
            r.request_type_id,
            r.property_id,
            r.status_id,
            rs.name as status_name,
            r.description,
            r.verification_complete,
            r.verification_date,
            r.created_at,
            r.updated_at,
            COALESCE(r.file_number, '') as file_number,
            COALESCE(r.abl_debt, '') as abl_debt,
            COALESCE(r.selected_activities, ARRAY[]::integer[]) as selected_activities,
            COALESCE(r.estimated_time, 0) as estimated_time,
            COALESCE(r.insurance, FALSE) as insurance
        FROM 
            public.requests r
            LEFT JOIN public.request_status rs ON r.status_id = rs.id
        WHERE 
            r.id = $1`

	var req transport.RequestDataModel
	err := r.repository.Pool().QueryRow(ctx, query, reqID).Scan(
		&req.ID,
		&req.UserID,
		&req.RequestTypeID,
		&req.PropertyID,
		&req.StatusID,
		&req.StatusName,
		&req.Description,
		&req.VerificationComplete,
		&req.VerificationDate,
		&req.CreatedAt,
		&req.UpdatedAt,
		&req.FileNumber,
		&req.ABLDebt,
		&req.SelectedActivities,
		&req.EstimatedTime,
		&req.Insurance,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("request not found with ID %d", reqID)
		}
		return nil, fmt.Errorf("error getting request: %w", err)
	}

	return transport.ToRequestDomain(&req), nil
}

func (r *PostgreSQL) GetRequestByExpCode(ctx context.Context, expCode string) (*domain.Request, error) {
	const requestQuery = `
		SELECT 
			r.id, 
			r.user_id, 
			r.user_type,
			r.request_type_id, 
			r.property_id, 
			r.status_id, 
			rs.name as status_name, 
			r.description, 
			r.verification_complete, 
			r.verification_date,
			r.observations,
			r.observations_tasks,
			r.created_at, 
			r.updated_at,
			p.street, p.number, a.abl_number,
			COALESCE(r.file_number, '') as file_number,
			COALESCE(r.abl_debt, '') as abl_debt,
			COALESCE(r.selected_activities, ARRAY[]::integer[]) as selected_activities,
			COALESCE(r.estimated_time, 0) as estimated_time,
			COALESCE(r.insurance, FALSE) as insurance
		FROM public.requests r 
		LEFT JOIN properties p ON p.property_id = r.property_id 
		LEFT JOIN abl a ON a.abl_id = p.abl_id 
		LEFT JOIN public.request_status rs ON r.status_id = rs.id
		WHERE r.file_number = $1`

	var req transport.RequestDataModel
	var addressStreet sql.NullString
	var addressNumber sql.NullInt64
	var addressABL sql.NullInt64

	err := r.repository.Pool().QueryRow(ctx, requestQuery, expCode).Scan(
		&req.ID,
		&req.UserID,
		&req.UserType,
		&req.RequestTypeID,
		&req.PropertyID,
		&req.StatusID,
		&req.StatusName,
		&req.Description,
		&req.VerificationComplete,
		&req.VerificationDate,
		&req.Observations,
		&req.ObservationsTasks,
		&req.CreatedAt,
		&req.UpdatedAt,
		&addressStreet,
		&addressNumber,
		&addressABL,
		&req.FileNumber,
		&req.ABLDebt,
		&req.SelectedActivities,
		&req.EstimatedTime,
		&req.Insurance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("request not found with expedition code (file_number): %s", expCode)
		}
		return nil, fmt.Errorf("error getting request: %w", err)
	}

	const documentQuery = `
		SELECT 
			documents.id,
			document_type_id,
			file_id,
			document_types.description,
			original_content as content,
			filename
		FROM documents
		INNER JOIN document_types ON documents.document_type_id = document_types.id
		WHERE code = $1 
		AND document_types.description IS NOT NULL 
		AND document_types.description != ''`

	rows, err := r.repository.Pool().Query(ctx, documentQuery, expCode)
	if err != nil {
		return nil, fmt.Errorf("error querying documents: %w", err)
	}
	defer rows.Close()

	var documents []transport.DocumentModel
	for rows.Next() {
		var doc transport.DocumentModel
		if err := rows.Scan(&doc.ID, &doc.Type, &doc.FileID, &doc.Description, &doc.Content, &doc.Filename); err != nil {
			return nil, fmt.Errorf("error scanning document: %w", err)
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}

	addressNumberStr := ""
	if addressNumber.Valid {
		addressNumberStr = fmt.Sprintf("%d", addressNumber.Int64)
	}

	request := transport.ToRequestDomain(&req)
	request.Observations = req.Observations.String
	request.ObservationsTasks = req.ObservationsTasks.String
	request.Address.Street = addressStreet.String
	request.Address.Number = addressNumberStr
	request.Address.ABLNumber = addressABL.Int64

	request.Documents = transport.DocumentModelToDocumentRequestList(documents)

	return request, nil
}

func (r *PostgreSQL) UpdateUserRequest(ctx context.Context, req *domain.Request) (string, string, error) {
	if req == nil {
		return "", "", fmt.Errorf("nil request")
	}

	tx, err := r.repository.Pool().Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	reqDataModel := transport.ToCreateRequestDataModel(req)
	reqDataModel.StatusID = 1

	const query = `
        UPDATE requests 
        SET 
            property_id = $1,
            status_id = $2,
            description = $3,
            updated_at = CURRENT_TIMESTAMP,
            abl_debt = $4,
            estimated_time = $5,
            insurance = $6,
            selected_activities = $7
        WHERE file_number = $8 
        RETURNING observations, observations_tasks, id`

	var observations sql.NullString
	var observationsTasks sql.NullString
	var reqID int64
	err = tx.QueryRow(ctx, query,
		reqDataModel.PropertyID,
		reqDataModel.StatusID,
		reqDataModel.Description,
		reqDataModel.ABLDebt,
		reqDataModel.EstimatedTime,
		reqDataModel.Insurance,
		pq.Array(reqDataModel.SelectedActivities),
		req.FileNumber,
	).Scan(&observations, &observationsTasks, &reqID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", fmt.Errorf("%w: request not found", transport.ErrUpdateRequest)
		}
		return "", "", fmt.Errorf("%w: %v", transport.ErrUpdateRequest, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", "", fmt.Errorf("error committing transaction: %v", err)
	}

	req.ID = reqID

	return observations.String, observationsTasks.String, nil
}

func (r *PostgreSQL) UpdateVerificationStatus(ctx context.Context, fileNumber string) error {
	tx, err := r.repository.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	statusID := 1

	const query = `
        UPDATE requests 
        SET 
			status_id_tasks = CASE WHEN status_id_tasks = 5 THEN $1 ELSE status_id_tasks END,
        	status_id_property = CASE WHEN status_id_property = 5 THEN $2 ELSE status_id_property END
        WHERE file_number = $3 
        RETURNING id, user_id`

	var requestID int64
	var userID int64
	err = tx.QueryRow(ctx, query,
		statusID,
		statusID,
		fileNumber,
	).Scan(&requestID, &userID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("%w: request not found", transport.ErrUpdateRequest)
		}
		return fmt.Errorf("%w: %v", transport.ErrUpdateRequest, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (r *PostgreSQL) GetPersonByUserID(ctx context.Context, id int64) (string, error) {
	const query = `
		SELECT 
			email
		FROM persons
		INNER JOIN users ON users.person_id = persons.id
		WHERE users.id = $1 
		LIMIT 1`

	var email sql.NullString
	err := r.repository.QueryRowContext(ctx, query, id).Scan(
		&email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("no user found for id %d", id)
		}
		return "", fmt.Errorf("error fetching user: %w", err)
	}

	return email.String, nil
}
