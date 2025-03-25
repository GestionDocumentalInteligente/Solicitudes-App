package config

import (
	"fmt"
	"log"
	"os"
	"time"

	sdkclo "github.com/teamcubation/sg-backend/pkg/config/config-loader"
	sdkmwr "github.com/teamcubation/sg-backend/pkg/rest/middlewares/gin"
)

// Constantes del sistema
const (
	// Request Status
	ReqStatusPending  = 1
	ReqStatusAproved  = 2
	ReqStatusRejected = 3

	// Request Types
	ReqTypeConstructionNotice = 1

	// Defaults
	DefaultContextKey     = "auth"
	DefaultTokenDuration  = 24 * time.Hour
	DefaultMaxRetries     = 3
	DefaultTimeoutSeconds = 30
)

// Config estructura principal de configuración
type Config struct {
	App        AppConfig
	Auth       AuthConfig
	Middleware MiddlewareConfig
	External   ExternalServicesConfig
}

// AppConfig configuración general de la aplicación
type AppConfig struct {
	Environment string
	Debug       bool
	TimeoutSec  int
	MaxRetries  int
}

// AuthConfig configuración de autenticación
type AuthConfig struct {
	TokenDuration time.Duration
	ContextKey    string
}

// MiddlewareConfig configuración de middlewares
type MiddlewareConfig struct {
	Auth sdkmwr.Config
}

// ExternalServicesConfig configuración de servicios externos
type ExternalServicesConfig struct {
	AFIP  AFIPConfig
	MiArg MiArgConfig
}

// AFIPConfig configuración específica de AFIP
type AFIPConfig struct {
	TokenEndpoint string
	ClientSecret  string
	BaseURL       string
}

// MiArgConfig configuración específica de Mi Argentina
type MiArgConfig struct {
	ClientSecret string
	BaseURL      string
}

var (
	cfg *Config
)

// Load inicializa y carga la configuración
func Load() error {
	// Cargar archivos de configuración
	if err := sdkclo.LoadConfig("config/.env"); err != nil {
		return fmt.Errorf("error loading config files: %w", err)
	}

	// Inicializar configuración
	cfg = &Config{
		App: AppConfig{
			Environment: os.Getenv("APP_ENV"),
		},
		Auth: AuthConfig{
			ContextKey: os.Getenv("AUTH_CONTEXT_KEY"),
		},
		Middleware: MiddlewareConfig{
			Auth: sdkmwr.Config{
				SecretKey:   os.Getenv("JWT_SECRET_KEY"),
				ContextKey:  os.Getenv("AUTH_CONTEXT_KEY"),
				TokenLookup: "header:Authorization",
				TokenPrefix: "Bearer ",
			},
		},
		External: ExternalServicesConfig{
			AFIP: AFIPConfig{
				TokenEndpoint: os.Getenv("AFIP_TOKEN_ENDPOINT"),
				ClientSecret:  os.Getenv("JWT_SECRET_KEY"),
				BaseURL:       os.Getenv("AFIP_BASE_URL"),
			},
			MiArg: MiArgConfig{
				ClientSecret: os.Getenv("MIARG_CLIENT_SECRET"),
				BaseURL:      os.Getenv("MIARG_BASE_URL"),
			},
		},
	}

	// Establecer valores por defecto si no están configurados
	setDefaults()

	// Validar configuración
	if err := cfg.validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

// setDefaults establece valores por defecto para campos no configurados
func setDefaults() {
	if cfg.App.TimeoutSec == 0 {
		cfg.App.TimeoutSec = DefaultTimeoutSeconds
	}
	if cfg.App.MaxRetries == 0 {
		cfg.App.MaxRetries = DefaultMaxRetries
	}
	if cfg.Auth.TokenDuration == 0 {
		cfg.Auth.TokenDuration = DefaultTokenDuration
	}
	if cfg.Auth.ContextKey == "" {
		cfg.Auth.ContextKey = DefaultContextKey
	}
}

// validate valida la configuración
func (c *Config) validate() error {
	// Validación de AFIP
	if c.External.AFIP.TokenEndpoint == "" {
		return fmt.Errorf("AFIP_TOKEN_ENDPOINT is required")
	}
	if c.External.AFIP.ClientSecret == "" {
		return fmt.Errorf("JWT_SECRET_KEY is required")
	}

	// Validación de Mi Argentina
	if c.External.MiArg.ClientSecret == "" {
		return fmt.Errorf("MIARG_CLIENT_SECRET is required")
	}

	// Validación de configuración general
	if c.App.Environment == "" {
		return fmt.Errorf("APP_ENV is required")
	}

	return nil
}

// Getters

// GetConfig retorna la configuración completa
func GetConfig() *Config {
	return cfg
}

// GetMiddlewareConfig retorna la configuración de middleware
func GetMiddlewareConfig() MiddlewareConfig {
	return cfg.Middleware
}

// GetAFIPConfig retorna la configuración de AFIP
func GetAFIPConfig() AFIPConfig {
	return cfg.External.AFIP
}

// GetMiArgConfig retorna la configuración de Mi Argentina
func GetMiArgConfig() MiArgConfig {
	return cfg.External.MiArg
}

// GetAppConfig retorna la configuración de la aplicación
func GetAppConfig() AppConfig {
	return cfg.App
}

// GetAuthConfig retorna la configuración de autenticación
func GetAuthConfig() AuthConfig {
	return cfg.Auth
}

// IsDebug retorna si la aplicación está en modo debug
func IsDebug() bool {
	return cfg.App.Debug
}

// GetEnvironment retorna el ambiente actual
func GetEnvironment() string {
	return cfg.App.Environment
}

// MustLoad carga la configuración y termina la aplicación si hay error
func MustLoad() {
	if err := Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
}

// Helper functions

// IsDevelopment retorna si el ambiente es desarrollo
func IsDevelopment() bool {
	return cfg.App.Environment == "development"
}

// IsProduction retorna si el ambiente es producción
func IsProduction() bool {
	return cfg.App.Environment == "production"
}

// IsTest retorna si el ambiente es testing
func IsTest() bool {
	return cfg.App.Environment == "test"
}
