package trncreate

import (
	"time"

	"github.com/teamcubation/sg-users/internal/core/entities"
	person "github.com/teamcubation/sg-users/internal/person/core/entities"
)

type UserResponse struct {
	ID                   int64    `json:"id"`
	PersonID             *int64   `json:"person_id,omitempty"`
	EmailValidated       bool     `json:"email_validated"`
	AcceptsNotifications bool     `json:"accepts_notifications"`
	Roles                []string `json:"roles"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at,omitempty"`
	DeletedAt            string   `json:"deleted_at,omitempty"`
	Cuil                 string   `json:"cuil"`
	Dni                  string   `json:"dni"`
	FirstName            string   `json:"first_name"`
	LastName             string   `json:"last_name"`
	Email                string   `json:"email"`
	Phone                string   `json:"phone"`
}

// ToUserResponse converts the User entity to a response DTO
func ToUserResponse(user *entities.User) *UserResponse {
	// Convert roles to a slice of strings (role names)
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	return &UserResponse{
		ID:             user.ID,
		PersonID:       user.PersonID,
		EmailValidated: user.EmailValidated,
		Roles:          roles,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
		Cuil:           user.Person.Cuil,
		Dni:            user.Person.Dni,
		FirstName:      user.Person.FirstName,
		LastName:       user.Person.LastName,
		Email:          user.Person.Email,
		Phone:          user.Person.Phone,
		UpdatedAt: func() string {
			if user.UpdatedAt != nil {
				return user.UpdatedAt.Format(time.RFC3339)
			}
			return ""
		}(),
		DeletedAt: func() string {
			if user.DeletedAt != nil {
				return user.DeletedAt.Format(time.RFC3339)
			}
			return ""
		}(),
	}
}

type PersonResponse struct {
	ID        int64  `json:"id"`
	Cuil      string `json:"cuil"`
	Dni       string `json:"dni"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func ToPersonResponse(user *person.Person) *PersonResponse {
	return &PersonResponse{
		ID:        user.ID,
		Cuil:      user.Cuil,
		Dni:       user.Dni,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}
}
