package person

import (
	"context"
	"fmt"

	"github.com/teamcubation/sg-users/internal/person/core/entities"
	"github.com/teamcubation/sg-users/internal/person/core/ports"
)

type useCases struct {
	repository ports.Repository
}

func NewUseCases(r ports.Repository) ports.UseCases {
	return &useCases{
		repository: r,
	}
}

func (u *useCases) CreatePerson(ctx context.Context, person *entities.Person) (int64, error) {
	person.Cuil = keepOnlyNumbers(person.Cuil)
	person.Dni = keepOnlyNumbers(person.Dni)
	person.Phone = keepOnlyNumbers(person.Phone)

	// Intentar crear la persona en el repositorio
	if err := u.repository.CreatePerson(ctx, person); err != nil {
		return 0, fmt.Errorf("failed to create person: %w", err)
	}

	return person.ID, nil
}

func (u *useCases) FindPersonByCuil(ctx context.Context, cuil string) (*entities.Person, error) {
	cuil = keepOnlyNumbers(cuil)

	person, err := u.repository.FindPersonByCuil(ctx, cuil)
	if err != nil {
		return nil, fmt.Errorf("failed to find person by CUIL: %w", err)
	}

	return person, nil
}

func (u *useCases) FindPersonByEmail(ctx context.Context, email string) (*entities.Person, error) {
	person, err := u.repository.FindPersonByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find person by CUIL: %w", err)
	}

	return person, nil
}

func (u *useCases) UpdatePersonByCuil(ctx context.Context, person *entities.Person) (int64, error) {
	person.Cuil = keepOnlyNumbers(person.Cuil)
	person.Dni = keepOnlyNumbers(person.Dni)
	person.Phone = keepOnlyNumbers(person.Phone)

	existingPerson, err := u.FindPersonByCuil(ctx, person.Cuil)
	if err != nil {
		return 0, fmt.Errorf("failed to find person by CUIL: %w", err)
	}

	hasChanges := false

	if person.Dni != "" && person.Dni != existingPerson.Dni {
		existingPerson.Dni = person.Dni
		hasChanges = true
	}
	if person.FirstName != "" && person.FirstName != existingPerson.FirstName {
		existingPerson.FirstName = person.FirstName
		hasChanges = true
	}
	if person.LastName != "" && person.LastName != existingPerson.LastName {
		existingPerson.LastName = person.LastName
		hasChanges = true
	}
	if person.Email != "" && person.Email != existingPerson.Email {
		existingPerson.Email = person.Email
		hasChanges = true
	}
	if person.Phone != "" && person.Phone != existingPerson.Phone {
		existingPerson.Phone = person.Phone
		hasChanges = true
	}

	if hasChanges {
		if err := u.repository.UpdatePerson(ctx, existingPerson); err != nil {
			return 0, fmt.Errorf("failed to update person: %w", err)
		}
	}

	return existingPerson.ID, nil
}

func (u *useCases) FindPersonByID(ctx context.Context, id int64) (*entities.Person, error) {
	person, err := u.repository.FindPersonByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find person by UUID: %w", err)
	}
	return person, nil
}
