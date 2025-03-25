package tranport

import (
	"strconv"

	dto "github.com/teamcubation/sg-users/internal/core/dto"
)

type User struct {
	*Person
	AcceptsNotifications bool   `json:"accepts_notifications"`
	EmailValidated       bool   `json:"email_validated"`
	Roles                []Role `json:"roles,omitempty"`
}

// Person representa la estructura de datos para una persona
type Person struct {
	Cuil      string  `json:"cuit" binding:"required"`
	Dni       *string `json:"dni,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     string  `json:"email" binding:"required,email"`
	Phone     string  `json:"phone" binding:"required"`
}

// Role representa la estructura de datos del rol para la API
type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

// Perm representa la estructura de datos del permiso para la API
type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ToUserDto(req *User) *dto.UserDto {
	var personDto *dto.PersonDto
	if req.Person != nil {
		var dni string
		if _, err := strconv.Atoi(req.Person.Cuil); err == nil {
			dni = req.Person.Cuil[2:10]
		}

		personDto = &dto.PersonDto{
			Cuil:      req.Person.Cuil,
			Dni:       &dni,
			FirstName: req.Person.FirstName,
			LastName:  req.Person.LastName,
			Email:     req.Person.Email,
			Phone:     req.Person.Phone,
		}
	}

	// Mapeamos los roles y permisos
	rolesDto := make([]dto.RoleDto, len(req.Roles))
	for i, role := range req.Roles {
		permissionsDto := make([]dto.PermissionDto, len(role.Permissions))
		for j, perm := range role.Permissions {
			permissionsDto[j] = dto.PermissionDto{
				Name:        perm.Name,
				Description: perm.Description,
			}
		}

		rolesDto[i] = dto.RoleDto{
			Name:        role.Name,
			Permissions: permissionsDto,
		}
	}

	return &dto.UserDto{
		Person:         personDto,
		EmailValidated: req.EmailValidated,
		Roles:          rolesDto,
	}
}
