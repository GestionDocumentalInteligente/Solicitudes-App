package perconn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"

	sdkpg "github.com/teamcubation/sg-users/pkg/databases/sql/postgresql/pgxpool"
	sdkpgports "github.com/teamcubation/sg-users/pkg/databases/sql/postgresql/pgxpool/defs"

	"github.com/teamcubation/sg-users/internal/person/core/entities"
	"github.com/teamcubation/sg-users/internal/person/core/ports"
)

type PostgreSQL struct {
	repository sdkpgports.Repository
}

func NewPostgreSQL() (ports.Repository, error) {
	r, err := sdkpg.Bootstrap("USERS_DB")
	if err != nil {
		return nil, fmt.Errorf("bootstrap error: %w", err)
	}
	return &PostgreSQL{repository: r}, nil
}

func (r *PostgreSQL) CreatePerson(ctx context.Context, person *entities.Person) error {
	query := `
		INSERT INTO persons (
			cuil, dni, first_name, last_name, email, phone
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`
	err := r.repository.Pool().QueryRow(ctx, query,
		person.Cuil, person.Dni, person.FirstName, person.LastName,
		person.Email, person.Phone,
	).Scan(&person.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("person already exists")
		}
		return fmt.Errorf("error creating person: %w", err)
	}
	return nil
}

func (r *PostgreSQL) UpdatePerson(ctx context.Context, person *entities.Person) error {
	query := `
		UPDATE persons
		SET dni = $1, first_name = $2, last_name = $3, email = $4, phone = $5
		WHERE cuil = $7
	`
	result, err := r.repository.Pool().Exec(ctx, query,
		person.Dni, person.FirstName, person.LastName,
		person.Email, person.Phone, person.Cuil,
	)
	if err != nil {
		return fmt.Errorf("error updating person: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("person not found")
	}

	return nil
}

func (r *PostgreSQL) FindPersonByCuil(ctx context.Context, cuil string) (*entities.Person, error) {
	query := `
		SELECT id, cuil, dni, first_name, last_name, email, phone
		FROM persons
		WHERE cuil = $1
	`
	person := &entities.Person{}
	err := r.repository.Pool().QueryRow(ctx, query, cuil).Scan(
		&person.ID,
		&person.Cuil,
		&person.Dni,
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding person: %w", err)
	}

	return person, nil
}

func (r *PostgreSQL) FindPersonByEmail(ctx context.Context, email string) (*entities.Person, error) {
	query := `
		SELECT id, cuil, dni, first_name, last_name, email, phone
		FROM persons
		WHERE email = $1
	`
	person := &entities.Person{}
	err := r.repository.Pool().QueryRow(ctx, query, email).Scan(
		&person.ID,
		&person.Cuil,
		&person.Dni,
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.Phone,
	)
	if err != nil {
		log.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding person: %w", err)
	}

	return person, nil
}

func (r *PostgreSQL) FindPersonByID(ctx context.Context, id int64) (*entities.Person, error) {
	query := `
		SELECT id, cuil, dni, first_name, last_name, email, phone
		FROM persons
		WHERE uuid = $1
	`
	person := &entities.Person{}
	err := r.repository.Pool().QueryRow(ctx, query, id).Scan(
		&person.ID,
		&person.Cuil,
		&person.Dni,
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.Phone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("person not found")
		}
		return nil, fmt.Errorf("error finding person by UUID: %w", err)
	}

	return person, nil
}
