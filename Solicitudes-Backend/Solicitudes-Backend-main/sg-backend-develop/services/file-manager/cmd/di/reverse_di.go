package di

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/teamcubation/sg-file-manager-api/cmd/api/handler/ablhdl"
	"github.com/teamcubation/sg-file-manager-api/cmd/api/handler/authhdl"
	"github.com/teamcubation/sg-file-manager-api/cmd/api/handler/filehdl"
	"github.com/teamcubation/sg-file-manager-api/internal/manager"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/abl"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/auth"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/platform/db"
	"github.com/teamcubation/sg-file-manager-api/internal/platform/env"
	"github.com/teamcubation/sg-file-manager-api/internal/platform/restclient"
)

type cmds struct {
	AuthHandler *authhdl.AuthHandler
	FileHandler *filehdl.FileHandler
	AblHandler  *ablhdl.ABLHandler
}

func ConfigReverseDI() cmds {
	ctx := context.Background()
	if os.Getenv("SCOPE") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		os.Setenv("SSL_MODE", "disable")
	}

	env.LoadConfigs()

	httpClient := restclient.NewHTTPClient(restclient.Credentials{
		Email:    env.GetEmail(),
		Password: env.GetPassword(),
	},
		restclient.WithBaseURL(env.GetBaseURLGDE()),
		restclient.WithTimeout(360*time.Second),
		restclient.WithRetries(3),
		restclient.WithRetryCondition(func(response *resty.Response, err error) bool {
			return err != nil || response.StatusCode() >= 500
		}),
	)

	token, err := httpClient.Login(ctx)
	if err != nil {
		log.Fatalf("Failed to login to external API: %v", err)
	}
	httpClient.SetToken(token)

	authClient := auth.NewRestClient(httpClient)
	fileClient := file.NewRestClient(httpClient)

	ablHTTPClient := restclient.NewHTTPClient(restclient.Credentials{
		Email:    env.GetEmail(),
		Password: env.GetPassword(),
	},
		restclient.WithBaseURL("https://apis-dev.gestionmsi.gob.ar"),
		restclient.WithTimeout(360*time.Second),
		restclient.WithRetries(3),
		restclient.WithRetryCondition(func(response *resty.Response, err error) bool {
			return err != nil || response.StatusCode() >= 500
		}),
	)

	ablClient := abl.NewRestClient(ablHTTPClient)

	db, err := db.NewPGConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fileRepository := file.NewFileRepository(db)

	// Initialize use cases
	authUC := manager.NewAuthUseCase(authClient)
	fileUC := manager.NewFileUseCase(fileClient, fileRepository)
	ablUC := manager.NewAblUseCase(ablClient)

	// Initialize HTTP handlers
	authHandler := authhdl.NewAuthHandler(authUC)
	fileHandler := filehdl.NewFileHandler(fileUC)
	ablHandler := ablhdl.NewABLHandler(ablUC)

	return cmds{
		AuthHandler: authHandler,
		FileHandler: fileHandler,
		AblHandler:  ablHandler,
	}
}
