package user

import (
	"github.com/teamcubation/sg-users/internal/core/dto"
	"github.com/teamcubation/sg-users/internal/core/entities"
)

// Función auxiliar para mapear RoleDto a entities.Role
func mapRoles(roleDtos []dto.RoleDto) []entities.Role {
	roles := make([]entities.Role, len(roleDtos))
	for i, roleDto := range roleDtos {
		permissions := make([]entities.Permission, len(roleDto.Permissions))
		for j, permDto := range roleDto.Permissions {
			permissions[j] = entities.Permission{
				Name:        permDto.Name,
				Description: permDto.Description,
			}
		}
		roles[i] = entities.Role{
			Name:        roleDto.Name,
			Permissions: permissions,
		}
	}
	return roles
}

// Función auxiliar para comparar si dos slices de roles son iguales
func rolesAreEqual(a, b []entities.Role) bool {
	if len(a) != len(b) {
		return false
	}
	roleMap := make(map[string]entities.Role)
	for _, role := range a {
		roleMap[role.Name] = role
	}
	for _, role := range b {
		if existingRole, ok := roleMap[role.Name]; !ok || !permissionsAreEqual(existingRole.Permissions, role.Permissions) {
			return false
		}
	}
	return true
}

// Función auxiliar para comparar si dos slices de permisos son iguales
func permissionsAreEqual(a, b []entities.Permission) bool {
	if len(a) != len(b) {
		return false
	}
	permMap := make(map[string]entities.Permission)
	for _, perm := range a {
		permMap[perm.Name] = perm
	}
	for _, perm := range b {
		if existingPerm, ok := permMap[perm.Name]; !ok || existingPerm.Description != perm.Description {
			return false
		}
	}
	return true
}
