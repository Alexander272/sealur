package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander272/sealur/file_service/internal/config"
	"github.com/Alexander272/sealur/file_service/internal/repository"
	"github.com/Alexander272/sealur/file_service/internal/service"
	transport "github.com/Alexander272/sealur/file_service/internal/transport/grpc"
	proto_file "github.com/Alexander272/sealur/file_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
	"github.com/Alexander272/sealur/file_service/pkg/storage"
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

	server := grpc.NewServer(opts...)
	proto_file.RegisterFileServiceServer(server, handlers)

	listener, err := net.Listen("tcp", ":"+conf.Http.Port)
	if err != nil {
		logger.Fatalf("failed to create grpc listener:", err)
	}

	go func() {
		if err = server.Serve(listener); err != nil {
			logger.Fatalf("failed to start server:", err)
		}
	}()
	logger.Infof("Application started on port: %s", conf.Http.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	server.GracefulStop()
}
