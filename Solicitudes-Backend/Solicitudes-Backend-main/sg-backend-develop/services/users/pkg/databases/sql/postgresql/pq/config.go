package sdkpg

import (
	"fmt"
	"os"

	"github.com/teamcubation/sg-users/pkg/databases/sql/postgresql/pq/defs"
)

type config struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
}

// newConfig crea una nueva configuración con los valores proporcionados
func newConfig(user, password, host, port, dbName string) defs.Config {
	return &config{
		Host:     host,
		User:     user,
		Password: password,
		DbName:   dbName,
		Port:     port,
	}
}

// DNS genera la cadena de conexión para PostgreSQL
func (c *config) DNS() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DbName, os.Getenv("SSL_MODE"))
}

// Validate valida que los campos necesarios estén presentes
func (c *config) Validate() error {
	if c.User == "" {
		return fmt.Errorf("POSTGRES_USERNAME environmente variable is empty")
	}
	if c.Password == "" {
		return fmt.Errorf("POSTGRES_PASSWORD environmente variable is empty")
	}
	if c.Host == "" {
		return fmt.Errorf("POSTGRES_HOST environmente variable is empty")
	}
	if c.DbName == "" {
		return fmt.Errorf("POSTGRES_DB environmente variable is empty")
	}
	return nil
}
