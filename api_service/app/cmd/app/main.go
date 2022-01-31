package main

import (
	"os"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/subosito/gotenv"
)

// @title Sealur Prop
// @version 0.1
// @description API Service

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := gotenv.Load(); err != nil {
		logger.Fatalf("error loading env variables: %s", err.Error())
	}
	conf, err := config.Init("configs")
	if err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}
	logger.Init(os.Stdout, conf.Environment)

	//* Dependencies

	//* Services, Repos & API Handlers

	//* HTTP Server
}
