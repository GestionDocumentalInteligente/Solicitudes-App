package sdksmtp

import (
	"fmt"
	"net/smtp"

	defs "github.com/teamcubation/sg-mailing-api/pkg/mailing/smtp/defs"
)

type config struct {
	smtpServer      string
	auth            smtp.Auth
	from            string
	port            string
	verificationURL string
}

func newConfig(smtpServer, port, from, password, identity, verificationURL string) defs.Config {
	auth := smtp.PlainAuth(identity, from, password, smtpServer)

	return &config{
		smtpServer:      smtpServer,
		auth:            auth,
		from:            from,
		port:            port,
		verificationURL: verificationURL,
	}
}

// GetSMTPServer devuelve la dirección del servidor SMTP con el puerto configurado
func (c *config) GetSMTPServer() string {
	return c.smtpServer
}

// GetAuth devuelve la autenticación SMTP configurada
func (c *config) GetAuth() smtp.Auth {
	return c.auth
}

// GetFrom devuelve la dirección de correo del remitente
func (c *config) GetFrom() string {
	return c.from
}

func (c *config) GetPort() string {
	return c.port
}

func (c *config) GetVerificationURL() string {
	return c.verificationURL
}

// Validate verifica que la configuración sea válida
func (c *config) Validate() error {
	if c.smtpServer == "" {
		return fmt.Errorf("SMTP server is not configured")
	}
	if c.auth == nil {
		return fmt.Errorf("SMTP auth is not configured")
	}
	if c.from == "" {
		return fmt.Errorf("SMTP from address is not configured")
	}
	if c.port == "" {
		return fmt.Errorf("SMTP port is not configured")
	}
	if c.verificationURL == "" {
		return fmt.Errorf("verification URL is not configured")
	}
	return nil
}
