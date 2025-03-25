package sdksmtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"sync"

	defs "github.com/teamcubation/sg-mailing-api/pkg/mailing/smtp/defs"
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

	htmlBody := fmt.Sprintf(`<html>
	<body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f2f2f2;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 5px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
      <div style="background-color: #4a5e45; padding: 20px; text-align: center;">
        <h1 style="color: white; font-size: 24px; margin: 0;">SAN ISIDRO</h1>
      </div>
      <div style="padding: 20px 30px;">
        <h2 style="color: #333; font-size: 22px;">Bienvenido/a a solicitudes online San Isidro</h2>
        <p style="color: #555;">Hola, <strong>%s</strong></p>
        <p style="color: #555;">Confirmá tu dirección de correo electrónico para completar tu registro, haciendo clic en el siguiente enlace:</p>
        <div style="text-align: center; margin: 20px 0;">
          <a href="%s" style="background-color: #4a5e45; color: white; padding: 15px 25px; text-decoration: none; border-radius: 5px; font-size: 18px; display: inline-block;">Verificar mi e-mail</a>
        </div>
		<p style="color: #666; font-size: 16px;">Si el botón no funciona, copiá y pegá el siguiente enlace en tu navegador:</p>
        <p style="color: #666; font-size: 16px; word-wrap: break-word;">
          <a href="%s" style="color: #4a5e45;">%s</a>
        </p>
        <p style="color: #666; font-size: 16px;">Si no solicitaste esta verificación o crees que se trata de un error, podés ignorar este mensaje.</p>
        <p style="color: #666; font-size: 16px;">Saludos,<br>El equipo de la Municipalidad de San Isidro.</p>
      </div>
    </div>
  </body>
    </html>`, data.Name, verificationURL, verificationURL, verificationURL)

	// Prepare the email message in the correct format
	msg := []byte(fmt.Sprintf("%s%s%s", data.Subject, data.BodyTemplate, htmlBody))

	return s.sendEmail(msg, data.Email)
}

func (s *service) SendNewRequestEmail(ctx context.Context, code string, data *defs.EmailData) error {
	verificationURL := "https://sgsanisidro.gob.ar/admin/requests"

	htmlBody := fmt.Sprintf(`<html>
	<body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f2f2f2;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 5px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
      <div style="background-color: #4a5e45; padding: 20px; text-align: center;">
        <h1 style="color: white; font-size: 24px; margin: 0;">SAN ISIDRO</h1>
      </div>
      <div style="padding: 20px 30px;">
        <h2 style="color: #333; font-size: 22px;">Gracias por completar tu solicitud de Aviso de Obra</h2>
        <p style="color: #555;">Hola,</p>
        <p style="color: #555;">Ya recibimos toda la información y estamos revisándola para asegurarnos de que todo esté en orden.</p>
		<p style="color: #555;">Tu número de solicitud es %s</p>
        <div style="text-align: center; margin: 20px 0;">
          <a href="%s" style="background-color: #4a5e45; color: white; padding: 15px 25px; text-decoration: none; border-radius: 5px; font-size: 18px; display: inline-block;">Ver solicitud</a>
        </div>
		<p style="color: #666; font-size: 16px;">Si el botón no funciona, copiá y pegá el siguiente enlace en tu navegador:</p>
        <p style="color: #666; font-size: 16px; word-wrap: break-word;">
          <a href="%s" style="color: #4a5e45;">%s</a>
        </p>
        <p style="color: #666; font-size: 16px;">Si no solicitaste esta verificación o crees que se trata de un error, podés ignorar este mensaje.</p>
        <p style="color: #666; font-size: 16px;">Saludos,<br>El equipo de la Municipalidad de San Isidro.</p>
      </div>
    </div>
  </body>
    </html>`, code, verificationURL, verificationURL, verificationURL)

	// Prepare the email message in the correct format
	msg := []byte(fmt.Sprintf("%s%s%s", data.Subject, data.BodyTemplate, htmlBody))

	return s.sendEmail(msg, data.Email)
}

func formatObs(obs string) string {
	obs = strings.ReplaceAll(obs, "</br>", "<br>")
	obs = strings.ReplaceAll(obs, " <br>", "<br>")
	obs = strings.ReplaceAll(obs, "<br> ", "<br>")

	return obs
}

func (s *service) SendUpdateRequestByCodeMessage(ctx context.Context, code, obs string, data *defs.EmailData) error {
	verificationURL := "https://sgsanisidro.gob.ar/ingresar"

	htmlBody := fmt.Sprintf(`
    <html>
      <body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f2f2f2;">
        <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 5px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
          <div style="background-color: #4a5e45; padding: 20px; text-align: center;">
            <h1 style="color: white; font-size: 24px; margin: 0;">SAN ISIDRO</h1>
          </div>
          <div style="padding: 20px 30px;">
            <p style="color: #555; font-size: 18px;">
              ¡Hola %s!
            </p>
            <p style="color: #555; font-size: 16px;">
              Te informamos que tu solicitud tiene observaciones:
            </p>
            <p style="color: #555; font-size: 16px; font-weight: bold;">
              %s
            </p>
            <p style="color: #555; font-size: 16px;">
              Por favor, asegurate de revisar y cargar los documentos solicitados.
            </p>
            <div style="text-align: center; margin: 20px 0;">
              <a href="%s"
                 style="background-color: #4a5e45; color: white; padding: 15px 25px;
                        text-decoration: none; border-radius: 5px; font-size: 16px;
                        display: inline-block;">
                Modificar solicitud
              </a>
            </div>
            <p style="color: #555; font-size: 16px;">
              Recordá que tenés un plazo de 30 días, de lo contrario, la solicitud expirará.
            </p>
            <p style="color: #555; font-size: 16px;">
              Si tenés alguna consulta, no dudes en escribirnos a 
              <a href="mailto:mesadigital@sanisidro.gob.ar" style="color: #4a5e45;">
                mesadigital@sanisidro.gob.ar
              </a>.
            </p>
            <p style="color: #555; font-size: 16px;">
              ¡Muchas gracias!<br>
              Saludos,<br>
              El equipo de San Isidro
            </p>
          </div>
        </div>
      </body>
    </html>
    `, data.Name, formatObs(obs), verificationURL)

	msg := []byte(fmt.Sprintf("%s%s%s", data.Subject, data.BodyTemplate, htmlBody))

	return s.sendEmail(msg, data.Email)
}

func (s *service) SendUpdateRequestMessage(ctx context.Context, code string, data *defs.EmailData) error {
	verificationURL := "https://sgsanisidro.gob.ar/ingresar"

	htmlBody := fmt.Sprintf(`
    <html>
      <body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f2f2f2;">
        <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 5px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
          <div style="background-color: #4a5e45; padding: 20px; text-align: center;">
            <h1 style="color: white; font-size: 24px; margin: 0;">SAN ISIDRO</h1>
          </div>
          <div style="padding: 20px 30px;">
            <p style="color: #555; font-size: 18px;">
              ¡Hola %s!
            </p>
            <p style="color: #555; font-size: 16px;">
              Gracias por actualizar la información para el aviso de obra.
              Estamos verificando los datos que enviaste, y este proceso puede demorar
              hasta 48 horas hábiles.
            </p>
            <p style="color: #555; font-size: 16px;">
              Mientras tanto, podés revisar el estado de tu solicitud accediendo a tu cuenta:
            </p>
            <div style="text-align: center; margin: 20px 0;">
              <a href="%s"
                 style="background-color: #4a5e45; color: white; padding: 15px 25px;
                        text-decoration: none; border-radius: 5px; font-size: 16px;
                        display: inline-block;">
                Ingresar a mi cuenta
              </a>
            </div>
            <p style="color: #555; font-size: 16px;">
              Si tenés alguna consulta, no dudes en escribirnos a 
              <a href="mailto:mesadigital@sanisidro.gob.ar" style="color: #4a5e45;">
                mesadigital@sanisidro.gob.ar
              </a>.
            </p>
            <p style="color: #555; font-size: 16px;">
              ¡Muchas gracias!<br>
              Saludos,<br>
              El equipo de San Isidro
            </p>
          </div>
        </div>
      </body>
    </html>
    `, data.Name, verificationURL)

	// Preparar el mensaje final combinando Subject, BodyTemplate (si la usas) y el htmlBody
	msg := []byte(fmt.Sprintf("%s%s%s", data.Subject, data.BodyTemplate, htmlBody))

	// Finalmente, se envía el correo
	return s.sendEmail(msg, data.Email)
}

func (s *service) SendValidateRequestMessage(ctx context.Context, code string, data *defs.EmailData) error {
	verificationURL := "https://sgsanisidro.gob.ar/ingresar"

	htmlBody := fmt.Sprintf(`
    <html>
      <body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f2f2f2;">
        <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 5px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
          <div style="background-color: #4a5e45; padding: 20px; text-align: center;">
            <h1 style="color: white; font-size: 24px; margin: 0;">SAN ISIDRO</h1>
          </div>
          <div style="padding: 20px 30px;">
            <p style="color: #555; font-size: 18px;">
              ¡Hola %s!
            </p>
            <p style="color: #555; font-size: 16px;">
              Buenas noticias. Tu solicitud de aviso de obra N° %s fue aprobada.
            </p>
            <p style="color: #555; font-size: 16px;">
              Ya podés descargar el certificado o revisar los detalles entrando a tu cuenta.
            </p>
            <div style="text-align: center; margin: 20px 0;">
              <a href="%s"
                 style="background-color: #4a5e45; color: white; padding: 15px 25px;
                        text-decoration: none; border-radius: 5px; font-size: 16px;
                        display: inline-block;">
                Ingresar a mi cuenta
              </a>
            </div>
            <p style="color: #555; font-size: 16px;">
              Si tenés alguna consulta, no dudes en escribirnos a 
              <a href="mailto:mesadigital@sanisidro.gob.ar" style="color: #4a5e45;">
                mesadigital@sanisidro.gob.ar
              </a>.
            </p>
            <p style="color: #555; font-size: 16px;">
              ¡Muchas gracias!<br>
              Saludos,<br>
              El equipo de San Isidro
            </p>
          </div>
        </div>
      </body>
    </html>
    `, data.Name, code, verificationURL)

	// Preparar el mensaje final combinando Subject, BodyTemplate (si la usas) y el htmlBody
	msg := []byte(fmt.Sprintf("%s%s%s", data.Subject, data.BodyTemplate, htmlBody))

	// Finalmente, se envía el correo
	return s.sendEmail(msg, data.Email)
}

func (s *service) sendEmail(msg []byte, email string) error {
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
		if err := client.Rcpt(email); err != nil {
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

		fmt.Printf("Verification email sent to %s (Development Mode)\n", email)
		return nil
	}

	// Production mode: Use TLS
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	})
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, host)
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
	if err := client.Rcpt(email); err != nil {
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

	return nil
}
