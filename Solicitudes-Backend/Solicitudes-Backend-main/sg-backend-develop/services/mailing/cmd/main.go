package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	mailconn "github.com/teamcubation/sg-mailing-api/internal/adapters/connectors"
	mailgtw "github.com/teamcubation/sg-mailing-api/internal/adapters/gateways"
	mailing "github.com/teamcubation/sg-mailing-api/internal/core"
)

func init() {
	if os.Getenv("SCOPE") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		os.Setenv("SSL_MODE", "disable")
	}
}

func main() {
	userRepo, err := mailconn.NewPostgreSQL()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	jwtService, err := mailconn.NewJwtService()
	if err != nil {
		log.Fatalf("JWT Service error: %v", err)
	}

	smtpService, err := mailconn.NewSmtpService()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	mailingUseCases := mailing.NewUseCases(jwtService, smtpService, userRepo)

	userHandler, err := mailgtw.NewGinHandler(mailingUseCases)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	err = userHandler.Start()
	if err != nil {
		log.Fatalf("Gin Server error at start: %v", err)
	}
}
