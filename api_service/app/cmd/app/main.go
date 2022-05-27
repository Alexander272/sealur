package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/repository"
	"github.com/Alexander272/sealur/api_service/internal/server"
	"github.com/Alexander272/sealur/api_service/internal/service"
	transport "github.com/Alexander272/sealur/api_service/internal/transport/http"
	"github.com/Alexander272/sealur/api_service/pkg/auth"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
)

// @title Sealur
// @version 0.5
// @description API Service for Sealur

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// if err := gotenv.Load("../../.env"); err != nil {
	// 	logger.Fatalf("error loading env variables: %s", err.Error())
	// }
	conf, err := config.Init("configs")
	if err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}
	logger.Init(os.Stdout, conf.Environment)

	//* Dependencies

	// client, err := redis.NewRedisClient(redis.Config{
	// 	Host:     conf.Redis.Host,
	// 	Port:     conf.Redis.Port,
	// 	DB:       conf.Redis.DB,
	// 	Password: conf.Redis.Password,
	// })
	// if err != nil {
	// 	logger.Fatalf("failed to initialize redis %s", err.Error())
	// }

	tokenManager, err := auth.NewManager(conf.Auth.JWT.Key)
	if err != nil {
		logger.Fatalf("failed to initialize token manager: %s", err.Error())
	}

	//* Services, Repos & API Handlers

	repos := repository.NewRepo()
	// repos := repository.NewRepo(client)
	services := service.NewServices(service.Deps{
		Repos:           repos,
		TokenManager:    tokenManager,
		AccessTokenTTL:  conf.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL: conf.Auth.JWT.RefreshTokenTTL,
		Domain:          conf.Http.Domain,
	})
	handlers := transport.NewHandler(services)

	//* HTTP Server

	srv := server.NewServer(conf, handlers.Init(conf))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	logger.Infof("Application started on port: %s", conf.Http.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
