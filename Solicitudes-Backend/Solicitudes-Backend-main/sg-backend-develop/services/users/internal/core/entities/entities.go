package entities

import (
	"errors"
	"time"

	datmod "github.com/teamcubation/sg-users/internal/adapters/connectors/data-model"
	"github.com/teamcubation/sg-users/internal/person/core/entities"
)

var ErrUserAlreadyExists = errors.New("user with this email already exists")

type UserType string

const (
	PersonType UserType = "person" // The user is a person
)

type User struct {
	ID                   int64  // Unique identifier for the user
	PersonID             *int64 // Reference to the associated person (nullable)
	EmailValidated       bool
	AcceptsNotifications bool
	Roles                []Role     // List of roles associated with the user
	CreatedAt            time.Time  // User creation date
	LoggedAt             *time.Time // Last time the user logged in (optional)
	UpdatedAt            *time.Time // Date when the user was last updated
	DeletedAt            *time.Time // Date when the user was deleted (optional)
	Person               *entities.Person
}

type Role struct {
	Name        string       // Role name (e.g., Admin, User)
	Permissions []Permission // List of permissions associated with the role
}

type Permission struct {
	Name        string // Permission name (e.g., Create, Edit)
	Description string // Permission description
}

func ToDataModel(user *User) (*datmod.User, error) {
	var dUser datmod.User

	dUser.ID = user.ID
	dUser.PersonID = user.PersonID
	dUser.EmailValidated = user.EmailValidated
	dUser.CreatedAt = user.CreatedAt
	dUser.LoggedAt = user.LoggedAt
	dUser.UpdatedAt = user.UpdatedAt
	dUser.DeletedAt = user.DeletedAt

	return &dUser, nil
}
