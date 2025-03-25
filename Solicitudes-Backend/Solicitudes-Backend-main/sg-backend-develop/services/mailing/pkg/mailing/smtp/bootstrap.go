package sdksmtp

import (
	"fmt"
	"os"

	defs "github.com/teamcubation/sg-mailing-api/pkg/mailing/smtp/defs"
)

func Bootstrap() (defs.Service, error) {
	config := newConfig(
		os.Getenv("SMTP_SERVER"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_IDENTITY"),
		os.Getenv("VERIFICATION_URL"),
	)

	// Validar la configuración
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("SMTP config error: %w", err)
	}

	// Crear el servicio SMTP con la configuración
	return newService(config)
}
