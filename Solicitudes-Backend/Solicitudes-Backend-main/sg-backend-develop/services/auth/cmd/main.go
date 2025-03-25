package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	authconn "github.com/teamcubation/sg-auth/internal/adapters/connectors"
	authgtw "github.com/teamcubation/sg-auth/internal/adapters/gateways"
	auth "github.com/teamcubation/sg-auth/internal/core"
)

func init() {
	if os.Getenv("SCOPE") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func main() {
	jwtService, err := authconn.NewJwtService()
	if err != nil {
		log.Fatalf("JWT Service error: %v", err)
	}

	authUsecases := auth.NewUseCases(jwtService)

	authHandler, err := authgtw.NewGinHandler(authUsecases)
	if err != nil {
		log.Fatalf("Auth Handler error: %v", err)
	}

	err = authHandler.Start()
	if err != nil {
		log.Fatalf("Gin Server error at start: %v", err)
	}

}
