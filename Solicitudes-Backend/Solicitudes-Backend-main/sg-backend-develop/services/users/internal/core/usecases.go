package user

import (
	"context"
	"fmt"
	"time"

	dto "github.com/teamcubation/sg-users/internal/core/dto"
	entities "github.com/teamcubation/sg-users/internal/core/entities"
	userports "github.com/teamcubation/sg-users/internal/core/ports"
	personports "github.com/teamcubation/sg-users/internal/person/core/ports"
)

type useCases struct {
	repository     userports.Repository
	personUseCases personports.UseCases
}

func NewUseCases(r userports.Repository, pu personports.UseCases) userports.UseCases {
	return &useCases{
		repository:     r,
		personUseCases: pu,
	}
}

func (u *useCases) CreateUser(ctx context.Context, userDto *dto.UserDto) (int64, error) {
	person := dto.ToPerson(userDto.Person)
	personFound, err := u.personUseCases.FindPersonByEmail(ctx, person.Email)
	if err != nil {
		return 0, fmt.Errorf("failed to find person by Email: %w", err)
	}

	if personFound != nil {
		return 0, entities.ErrUserAlreadyExists
	}

	personID, err := u.personUseCases.CreatePerson(ctx, person)
	if err != nil {
		return 0, fmt.Errorf("failed to create person: %w", err)
	}

	// Convertir UserDto a entidad User
	user := dto.ToUser(userDto)
	user.PersonID = &personID

	if err := u.repository.CreateUser(ctx, user); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

func (u *useCases) UpdateUserByPersonCuil(ctx context.Context, userDto *dto.UserDto) (int64, error) {
	// Actualizar la persona asociada al usuario por CUIL
	person := dto.ToPerson(userDto.Person)
	personID, err := u.personUseCases.UpdatePersonByCuil(ctx, person)
	if err != nil {
		return 0, fmt.Errorf("failed to update person by CUIL: %w", err)
	}

	existingUser, err := u.repository.FindUserByPersonID(ctx, personID)
	if err != nil {
		return 0, fmt.Errorf("failed to find user: %w", err)
	}

	// Flag para determinar si hubo cambios
	hasChanges := false

	// Actualizar campos si tienen nuevos valores
	if userDto.EmailValidated != existingUser.EmailValidated {
		existingUser.EmailValidated = userDto.EmailValidated
		hasChanges = true
	}

	// Actualizar roles si son diferentes
	if len(userDto.Roles) > 0 {
		newRoles := mapRoles(userDto.Roles)
		if !rolesAreEqual(existingUser.Roles, newRoles) {
			existingUser.Roles = newRoles
			hasChanges = true
		}
	}

	// Actualizar el campo UpdatedAt y guardar cambios si hubo modificaciones
	if hasChanges {
		now := time.Now()
		existingUser.UpdatedAt = &now

		if err := u.repository.UpdateUser(ctx, existingUser); err != nil {
			return 0, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return existingUser.ID, nil
}

func (u *useCases) FindUserByPersonCuil(ctx context.Context, cuil string) (*entities.User, error) {
	// Step 1: Find the person by CUIL using personUseCases
	person, err := u.personUseCases.FindPersonByCuil(ctx, cuil)
	if err != nil {
		return nil, fmt.Errorf("failed to find person by CUIL: %w", err)
	}
	if person == nil {
		return nil, nil
	}

	// Step 2: Find the user by the person's UUID
	user, err := u.repository.FindUserByPersonID(ctx, person.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by PersonUUID: %w", err)
	}

	user.Person = person

	return user, nil
}

func (u *useCases) FindUserByPersonID(ctx context.Context, personID int64) (*entities.User, error) {
	user, err := u.repository.FindUserByPersonID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by PersonUUID: %w", err)
	}
	return user, nil
}

// Implementación de FindUserByUserUUID
func (u *useCases) FindUserByUserID(ctx context.Context, userID int64) (*entities.User, error) {
	user, err := u.repository.FindUserByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by UserUUID: %w", err)
	}
	return user, nil
}
