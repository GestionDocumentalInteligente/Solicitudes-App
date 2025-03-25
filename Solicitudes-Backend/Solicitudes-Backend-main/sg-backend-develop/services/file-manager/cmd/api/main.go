package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teamcubation/sg-file-manager-api/cmd/api/handler/ablhdl"
	"github.com/teamcubation/sg-file-manager-api/cmd/api/handler/authhdl"
	"github.com/teamcubation/sg-file-manager-api/cmd/api/handler/filehdl"
	"github.com/teamcubation/sg-file-manager-api/cmd/di"

	"github.com/gin-gonic/gin"
)

func main() {
	cmds := di.ConfigReverseDI()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	v1 := router.Group("/api/v1/file-manager")

	authhdl.NewRouter(cmds.AuthHandler).AddRoutesV1(v1)
	filehdl.NewRouter(cmds.FileHandler).AddRoutesV1(v1)
	ablhdl.NewRouter(cmds.AblHandler).AddRoutesV1(v1)

	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = fmt.Sprintf(":%s", envPort)
	}

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received. Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully shut down.")
}
