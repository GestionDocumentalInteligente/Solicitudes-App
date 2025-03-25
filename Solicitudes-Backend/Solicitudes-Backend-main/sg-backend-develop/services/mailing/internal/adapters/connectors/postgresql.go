package mailconn

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	sdkpg "github.com/teamcubation/sg-mailing-api/pkg/databases/sql/postgresql/pgxpool"
	sdkpgports "github.com/teamcubation/sg-mailing-api/pkg/databases/sql/postgresql/pgxpool/defs"

	"github.com/teamcubation/sg-mailing-api/internal/core/entities"
	ports "github.com/teamcubation/sg-mailing-api/internal/core/ports"
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

func (r *PostgreSQL) GetuserByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT email_validated FROM users
		INNER JOIN persons
		ON users.person_id = persons.id WHERE persons.email = $1
	`

	row := r.repository.Pool().QueryRow(ctx, query, email)
	var user entities.User

	err := row.Scan(&user.EmailValidated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

func (r *PostgreSQL) UpdateUser(ctx context.Context, email string) error {
	query := `
		UPDATE users
		SET email_validated = true, updated_at = CURRENT_TIMESTAMP
		FROM persons
		WHERE users.person_id = persons.id AND persons.email = $1
	`

	_, err := r.repository.Pool().Exec(ctx, query, email)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
