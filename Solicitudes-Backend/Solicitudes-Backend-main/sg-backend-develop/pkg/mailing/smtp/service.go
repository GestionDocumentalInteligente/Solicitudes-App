package sdksmtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"sync"

	defs "github.com/teamcubation/sg-backend/pkg/mailing/smtp/defs"
)

var (
	instance defs.Service
	once     sync.Once
	initErr  error
)

// service representa el servicio SMTP que envía correos
type service struct {
	config defs.Config
}

// newService crea una nueva instancia del servicio SMTP usando la configuración proporcionada
func newService(config defs.Config) (defs.Service, error) {
	once.Do(func() {
		instance = &service{
			config: config,
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

func (s *service) SendVerificationEmail(ctx context.Context, data *defs.EmailData) error {
	verificationURL := fmt.Sprintf("%s?token=%s", s.config.GetVerificationURL(), data.Token)

	// Prepare the email subject and body
	body := fmt.Sprintf("%s:\n\n%s", data.BodyTemplate, verificationURL)

	// Prepare the email message in the correct format
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", data.Email, data.Subject, body))

	// Get and parse server address
	host := s.config.GetSMTPServer()
	port := s.config.GetPort()
	auth := s.config.GetAuth()
	from := s.config.GetFrom()

	// Check if we're in development mode by reading the STAGE environment variable
	stage := os.Getenv("STAGE")
	if stage == "DEV" {
		// Development mode: Do not use TLS
		client, err := smtp.Dial(fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer client.Quit()

		// Only authenticate if not MailHog
		if host != "mailhog" && host != "localhost" {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
			}
		}

		// Set the sender and recipient
		if err := client.Mail(from); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}
		if err := client.Rcpt(data.Email); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		// Write the email message
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get SMTP data writer: %w", err)
		}
		if _, err := w.Write(msg); err != nil {
			return fmt.Errorf("failed to write email message: %w", err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close email message writer: %w", err)
		}

		fmt.Printf("Verification email sent to %s (Development Mode)\n", data.Email)
		return nil
	}

	// Production mode: Use TLS
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{
		InsecureSkipVerify: true, // Only for development, ensure proper certificates for production
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.GetSMTPServer())
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// SMTP authentication
	if err := client.Auth(s.config.GetAuth()); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
	}

	// Set the sender and recipient
	if err := client.Mail(s.config.GetFrom()); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(data.Email); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Write the email message
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get SMTP data writer: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("failed to write email message: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close email message writer: %w", err)
	}

	fmt.Printf("Verification email sent to %s\n", data.Email)
	return nil
}
