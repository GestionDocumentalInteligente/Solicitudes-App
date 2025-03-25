package ports

import (
	"context"

	dto "github.com/teamcubation/sg-users/internal/core/dto"
	entities "github.com/teamcubation/sg-users/internal/core/entities"
)

type UseCases interface {
	CreateUser(context.Context, *dto.UserDto) (int64, error)
	UpdateUserByPersonCuil(context.Context, *dto.UserDto) (int64, error)
	FindUserByPersonCuil(context.Context, string) (*entities.User, error)
	FindUserByPersonID(context.Context, int64) (*entities.User, error)
	FindUserByUserID(context.Context, int64) (*entities.User, error)
}

type Repository interface {
	CreateUser(context.Context, *entities.User) error
	FindUserByPersonID(ctx context.Context, personID int64) (*entities.User, error)
	FindUserByUserID(ctx context.Context, userID int64) (*entities.User, error)
	UpdateUser(context.Context, *entities.User) error
}
