package authconn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	sdkpg "github.com/teamcubation/sg-users/pkg/databases/sql/postgresql/pgxpool"
	sdkpgports "github.com/teamcubation/sg-users/pkg/databases/sql/postgresql/pgxpool/defs"

	entities "github.com/teamcubation/sg-users/internal/core/entities"
	ports "github.com/teamcubation/sg-users/internal/core/ports"
)

type PostgreSQL struct {
	repository sdkpgports.Repository
}

func NewPostgreSQL() (ports.Repository, error) {
	r, err := sdkpg.Bootstrap("USERS_DB")
	if err != nil {
		return nil, fmt.Errorf("bootstrap error: %w", err)
	}

	return &PostgreSQL{
		repository: r,
	}, nil
}

// Crear usuario en la base de datos con hash de contraseña
func (r *PostgreSQL) CreateUser(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (
			person_id, email_validated, accepts_notifications, created_at
		) VALUES ($1, $2, $3, CURRENT_TIMESTAMP) RETURNING id
	`

	err := r.repository.Pool().QueryRow(ctx, query,
		user.PersonID,
		user.EmailValidated,
		user.AcceptsNotifications,
	).Scan(&user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Código para unique_violation
			return errors.New("user already exists")
		}
		return err
	}

	return nil
}

func (r *PostgreSQL) FindUserByUserID(ctx context.Context, userID int64) (*entities.User, error) {
	query := `SELECT id, person_id, email_validated, accepts_notifications, created_at, updated_at, deleted_at FROM users WHERE id = $1`

	return r.findUserByQuery(ctx, query, userID)
}

func (r *PostgreSQL) UpdateUser(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET email_validated = $1, accepts_notifications= $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.repository.Pool().Exec(ctx, query,
		user.EmailValidated,
		user.AcceptsNotifications,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (r *PostgreSQL) FindUserByPersonID(ctx context.Context, personID int64) (*entities.User, error) {
	query := `SELECT id, person_id, email_validated, accepts_notifications, created_at, updated_at, deleted_at FROM users WHERE person_id = $1`

	return r.findUserByQuery(ctx, query, personID)
}

// Función auxiliar para encontrar un usuario usando una consulta SQL
func (r *PostgreSQL) findUserByQuery(ctx context.Context, query string, identifier int64) (*entities.User, error) {
	// Ejecutar la consulta
	row := r.repository.Pool().QueryRow(ctx, query, identifier)

	// Mapear el resultado a una entidad User
	var user entities.User
	var personID sql.NullInt64

	err := row.Scan(&user.ID, &personID, &user.EmailValidated, &user.AcceptsNotifications, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found for identifier: %d", identifier)
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	user.PersonID = getInt64FromNullInt64(personID)

	return &user, nil
}
