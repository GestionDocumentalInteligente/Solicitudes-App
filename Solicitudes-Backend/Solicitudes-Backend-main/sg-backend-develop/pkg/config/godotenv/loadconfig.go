package sdkgodotenv

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

// LoadConfig carga múltiples archivos .env usando godotenv
func LoadConfig(filePaths ...string) error {
	if len(filePaths) == 0 {
		return errors.New("no env file paths provided")
	}

	successfullyLoaded := false // Variable para rastrear si algún archivo se cargó correctamente

	for _, filePath := range filePaths {
		if err := godotenv.Load(filePath); err != nil {
			fmt.Printf("sdkgodotenv: WARNING: Failed to load configuration file: '%s'\n", filePath)
			continue
		}
		fmt.Printf("sdkgodotenv: Configuration file successfully loaded from: %s\n", filePath)
		successfullyLoaded = true // Indica que al menos un archivo se cargó correctamente
	}

	// Si no se pudo cargar ningún archivo, devolver un error
	if !successfullyLoaded {
		return errors.New("no .env files could be loaded")
	}

	return nil
}
