package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexander272/sealur/file_service/internal/config"
	"github.com/Alexander272/sealur/file_service/internal/repository"
	"github.com/Alexander272/sealur/file_service/internal/server"
	"github.com/Alexander272/sealur/file_service/internal/service"
	transport "github.com/Alexander272/sealur/file_service/internal/transport/grpc"
	transport_http "github.com/Alexander272/sealur/file_service/internal/transport/http"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
	"github.com/Alexander272/sealur/file_service/pkg/storage"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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

	store, err := storage.NewClient(conf.MinIO)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	//* Services, Repos & API Handlers

	repo := repository.NewRepo(store)
	services := service.NewService(repo)
	handlers := transport.NewHandler(services, conf.Api)

	//* GRPC Server

	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		logger.Fatalf("failed to load certificate. error: %w", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(handlers.UnaryInterceptor),
	}

	grpcServer := grpc.NewServer(opts...)
	file_api.RegisterFileServiceServer(grpcServer, handlers)

	listener, err := net.Listen("tcp", ":"+conf.Tcp.Port)
	if err != nil {
		logger.Fatalf("failed to create grpc listener:", err)
	}

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			logger.Fatalf("failed to start server:", err)
		}
	}()
	logger.Infof("Application tcp started on port: %s", conf.Tcp.Port)

	//* API Handlers
	handler := transport_http.NewHandler(services)

	//* HTTP Server
	srv := server.NewServer(conf, handler.Init(conf))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	logger.Infof("Application started on port: %s", conf.Http.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	grpcServer.GracefulStop()

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
