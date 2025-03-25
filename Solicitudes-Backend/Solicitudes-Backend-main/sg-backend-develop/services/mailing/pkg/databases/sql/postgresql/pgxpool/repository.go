package sdkpostgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"

	defs "github.com/teamcubation/sg-mailing-api/pkg/databases/sql/postgresql/pgxpool/defs"
)

var (
	instance  defs.Repository
	once      sync.Once
	initError error
)

type repository struct {
	pool *pgxpool.Pool
}

func newRepository(config defs.Config) (defs.Repository, error) {
	once.Do(func() {
		instance = &repository{}
		initError = instance.Connect(config)
		if initError != nil {
			instance = nil
		}
	})
	return instance, initError
}

func (r *repository) Connect(c defs.Config) error {
	// Construcción de la cadena de conexión
	connString := c.DNS()

	// Conexión al pool de PostgreSQL
	pool, err := ConnectPool(connString)
	if err != nil {
		return err
	}
	r.pool = pool
	return nil
}

func (r *repository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *repository) Pool() *pgxpool.Pool {
	return r.pool
}

// Función que conecta al pool de PostgreSQL
func ConnectPool(connString string) (*pgxpool.Pool, error) {
	// Parsear la configuración de la cadena de conexión
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database connection string: %w", err)
	}

	// Establecer la conexión
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Realizar un ping a la base de datos para verificar la conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %w", err)
	}

	return pool, nil
}

// Función para aplicar migraciones
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
