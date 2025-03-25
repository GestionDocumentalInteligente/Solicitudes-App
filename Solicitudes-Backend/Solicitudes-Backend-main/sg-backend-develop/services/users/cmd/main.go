package main

import (
	"fmt"
	"log"
	"os"

	sdkcnfldr "github.com/teamcubation/sg-users/pkg/config/config-loader"

	userconn "github.com/teamcubation/sg-users/internal/adapters/connectors"
	usergtw "github.com/teamcubation/sg-users/internal/adapters/gateways"
	personconn "github.com/teamcubation/sg-users/internal/person/adapters/connectors"

	user "github.com/teamcubation/sg-users/internal/core"
	person "github.com/teamcubation/sg-users/internal/person/core"
)

func init() {
	if os.Getenv("SCOPE") == "" {
		fmt.Println("Dev mode!!")
		if err := sdkcnfldr.LoadConfig("config/.env", "config/.env.local"); err != nil {
			log.Fatalf("Viper Service error: %v", err)
		}
	} else {
		if err := sdkcnfldr.LoadConfig(".env"); err != nil {
			log.Fatalf("Viper Service error: %v", err)
		}
	}
}

// NOTE: no pude implementar wire todavia, dan errores que no entiendo, mirar mas adelante
func main() {
	personRepo, err := personconn.NewPostgreSQL()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	personUseCases := person.NewUseCases(personRepo)

	userRepo, err := userconn.NewPostgreSQL()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	usersUseCases := user.NewUseCases(userRepo, personUseCases)

	userHandler, err := usergtw.NewGinHandler(usersUseCases)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Iniciar el servidor de Gin
	err = userHandler.Start()
	if err != nil {
		log.Fatalf("Gin Server error at start: %v", err)
	}
}
