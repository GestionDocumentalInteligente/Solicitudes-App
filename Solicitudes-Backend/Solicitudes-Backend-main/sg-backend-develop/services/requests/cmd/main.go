package main

import (
	"context"
	"log"

	config "github.com/teamcubation/sg-backend/services/requests/internal/config"

	reqinb "github.com/teamcubation/sg-backend/services/requests/internal/request/adapters/inbound"
	reqout "github.com/teamcubation/sg-backend/services/requests/internal/request/adapters/outbound"
	req "github.com/teamcubation/sg-backend/services/requests/internal/request/core"
)

func init() {
	config.Load()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repository, err := reqout.NewPostgreSQL()
	if err != nil {
		log.Fatalf("PostgreSQL error: %v", err)
	}

	httpClient, err := reqout.NewHttpClient()
	if err != nil {
		log.Fatalf("Http Client error: %v", err)
	}

	reqUsecases := req.NewUseCases(repository, httpClient)

	reqHandler, err := reqinb.NewGinHandler(reqUsecases)
	if err != nil {
		log.Fatalf("req Handler error: %v", err)
	}

	err = reqHandler.Start(ctx)
	if err != nil {
		log.Fatalf("Gin Server error at start: %v", err)
	}
}
