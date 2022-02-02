package main

import (
	"os"

	"github.com/Alexander272/sealur/pro_service/internal/config"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur/pro_service/pkg/database/postgres"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/subosito/gotenv"
)

func main() {
	if err := gotenv.Load("../../.env"); err != nil {
		logger.Fatalf("error loading env variables: %s", err.Error())
	}
	conf, err := config.Init("configs")
	if err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}
	logger.Init(os.Stdout, conf.Environment)

	//* Dependencies

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     conf.Postgres.Host,
		Port:     conf.Postgres.Port,
		Username: conf.Postgres.Username,
		Password: conf.Postgres.Password,
		DBName:   conf.Postgres.DbName,
		SSLMode:  conf.Postgres.SSLMode,
	})
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	//* Services, Repos & API Handlers

	repos := repository.NewRepo(db)
	_ = service.NewServices(repos)
	// services := service.NewServices(repos)
	//TODO дописать
}
