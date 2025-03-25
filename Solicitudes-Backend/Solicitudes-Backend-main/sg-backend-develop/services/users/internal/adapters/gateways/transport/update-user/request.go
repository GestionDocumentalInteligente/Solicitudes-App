package tranport

import (
	dto "github.com/teamcubation/sg-users/internal/core/dto"
)

type User struct {
	EmailValidated bool    `json:"validated-email"`
	Person         *Person `json:"person"`
	Roles          []Role  `json:"roles,omitempty"`
}

// Person representa la estructura de datos para una persona
type Person struct {
	Cuil      string  `json:"cuil" binding:"required"`
	Dni       *string `json:"dni,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
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
		personDto = &dto.PersonDto{
			Cuil:      req.Person.Cuil,
			Dni:       req.Person.Dni,
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
		EmailValidated: req.EmailValidated,
		Person:         personDto,
		Roles:          rolesDto,
	}
}
