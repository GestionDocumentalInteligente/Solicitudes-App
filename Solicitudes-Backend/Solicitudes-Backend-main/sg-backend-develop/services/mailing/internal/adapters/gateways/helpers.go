package mailgtw

import (
	"errors"
	"os"
)

func getSecrets() (map[string]string, error) {
	// Crear un mapa para almacenar los secrets
	secrets := make(map[string]string)

	// Cargar los secrets cuando sea necesario
	jwtSecret := os.Getenv("JWT_SECRET")

	// Si los secrets están vacíos, retornamos error (por si acaso)
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is missing")
	}

	// Guardar los secretos en el mapa
	secrets["jwtSecret"] = jwtSecret

	return secrets, nil
}
