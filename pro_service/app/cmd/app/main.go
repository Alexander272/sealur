package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander272/sealur/pro_service/internal/config"
	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/service"
	handlers "github.com/Alexander272/sealur/pro_service/internal/transport/grpc"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	proto_email "github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto/email"
	proto_file "github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto/file"
	"github.com/Alexander272/sealur/pro_service/pkg/database/postgres"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	_ "github.com/lib/pq"
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

	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	if err != nil {
		logger.Fatalf("failed to load certificate. error: %w", err)
	}

	//* данные для аутентификации
	authEmail := models.Authentication{
		ServiceName: conf.Services.EmailService.AuthName,
		Password:    conf.Services.EmailService.AuthPassword,
	}

	//* опции grpc
	optsEmail := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&authEmail),
	}

	//* подключение к сервису
	connectEmail, err := grpc.Dial(conf.Services.EmailService.Url, optsEmail...)
	if err != nil {
		logger.Fatalf("failed connection to email service. error: %w", err)
	}
	emailClient := proto_email.NewEmailServiceClient(connectEmail)

	//* данные для аутентификации
	authFile := models.Authentication{
		ServiceName: conf.Services.FileService.AuthName,
		Password:    conf.Services.FileService.AuthPassword,
	}

	//* опции grpc
	optsFile := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&authFile),
	}

	//* подключение к сервису
	connectFile, err := grpc.Dial(conf.Services.FileService.Url, optsFile...)
	if err != nil {
		logger.Fatalf("failed connection to email service. error: %w", err)
	}
	fileClient := proto_file.NewFileServiceClient(connectFile)

	//* Services, Repos & API Handlers

	repos := repository.NewRepo(db)
	services := service.NewServices(repos, emailClient, fileClient)
	handlers := handlers.NewHandler(services, conf.Api)

	//TODO надо посмотреть как это по нормальному пишется

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
	proto.RegisterProServiceServer(server, handlers)

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
