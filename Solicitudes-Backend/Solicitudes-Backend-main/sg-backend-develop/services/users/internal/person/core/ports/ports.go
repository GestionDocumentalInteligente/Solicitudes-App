package ports

import (
	"context"

	entities "github.com/teamcubation/sg-users/internal/person/core/entities"
)

type UseCases interface {
	CreatePerson(ctx context.Context, person *entities.Person) (int64, error)
	FindPersonByCuil(context.Context, string) (*entities.Person, error)
	FindPersonByEmail(context.Context, string) (*entities.Person, error)
	FindPersonByID(context.Context, int64) (*entities.Person, error)
	UpdatePersonByCuil(context.Context, *entities.Person) (int64, error)
}

type Repository interface {
	CreatePerson(context.Context, *entities.Person) error
	UpdatePerson(context.Context, *entities.Person) error
	FindPersonByCuil(context.Context, string) (*entities.Person, error)
	FindPersonByEmail(context.Context, string) (*entities.Person, error)
	FindPersonByID(context.Context, int64) (*entities.Person, error)
}
