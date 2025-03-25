package sdkpg

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	defs "github.com/teamcubation/sg-users/pkg/databases/sql/postgresql/pq/defs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq" // Importación de driver
)

var (
	instance  defs.Repository
	once      sync.Once
	initError error
)

type repository struct {
	db *sql.DB
}

func newRepository(c defs.Config) (defs.Repository, error) {
	once.Do(func() {
		instance = &repository{}
		initError = instance.Connect(c)
		if initError != nil {
			instance = nil
		}
	})
	return instance, initError
}

func (r *repository) Connect(c defs.Config) error {
	// Construir la cadena de conexión
	connString := c.DNS()

	// Conectar con la base de datos PostgreSQL
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Verificar la conexión
	if err = db.Ping(); err != nil {
		return fmt.Errorf("unable to ping the database: %w", err)
	}

	r.db = db
	return nil
}

func (r *repository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *repository) DB() *sql.DB {
	return r.db
}

// Aplicar migraciones a la base de datos PostgreSQL
func ApplyMigrations(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// Crear una instancia de migración
	m, err := migrate.NewWithDatabaseInstance("file:///app/migrations", dbName, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Aplicar las migraciones
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
