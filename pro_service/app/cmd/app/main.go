package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander272/sealur/pro_service/internal/config"
	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/service"
	handlers "github.com/Alexander272/sealur/pro_service/internal/transport/grpc"
	"github.com/Alexander272/sealur/pro_service/pkg/database/postgres"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/Alexander272/sealur_proto/api/user_api"
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
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&authEmail),
	}

	//* подключение к сервису
	connectEmail, err := grpc.Dial(conf.Services.EmailService.Url, optsEmail...)
	if err != nil {
		logger.Fatalf("failed connection to email service. error: %w", err)
	}
	emailClient := email_api.NewEmailServiceClient(connectEmail)

	//* данные для аутентификации
	authFile := models.Authentication{
		ServiceName: conf.Services.FileService.AuthName,
		Password:    conf.Services.FileService.AuthPassword,
	}

	//* опции grpc
	optsFile := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&authFile),
	}

	//* подключение к сервису
	connectFile, err := grpc.Dial(conf.Services.FileService.Url, optsFile...)
	if err != nil {
		logger.Fatalf("failed connection to file service. error: %w", err)
	}
	fileClient := file_api.NewFileServiceClient(connectFile)

	//* данные для аутентификации
	authUser := models.Authentication{
		ServiceName: conf.Services.UserService.AuthName,
		Password:    conf.Services.UserService.AuthPassword,
	}

	//* опции grpc
	optsUser := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&authUser),
	}

	//* подключение к сервису
	connectUser, err := grpc.Dial(conf.Services.UserService.Url, optsUser...)
	if err != nil {
		logger.Fatalf("failed connection to user service. error: %w", err)
	}
	userClient := user_api.NewUserServiceClient(connectUser)

	//* Services, Repos & API Handlers

	repos := repository.NewRepo(db)
	services := service.NewServices(repos, emailClient, fileClient, userClient)
	handlers := handlers.NewHandler(services, conf.Api)

	//TODO надо посмотреть как это по нормальному пишется

	//* GRPC Server

	// cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	// if err != nil {
	// 	logger.Fatalf("failed to load certificate. error: %w", err)
	// }

	opts := []grpc.ServerOption{
		// grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		// grpc.Creds(insecure.NewCredentials()),
		// grpc.UnaryInterceptor(handlers.UnaryInterceptor),
	}

	server := grpc.NewServer(opts...)
	pro_api.RegisterProServiceServer(server, handlers)
	flange_standard_api.RegisterFlangeStandardServiceServer(server, handlers.FlangeStandard)
	flange_type_snp_api.RegisterFlangeTypeSnpServiceServer(server, handlers.FlangeTypeSnp)
	material_api.RegisterMaterialServiceServer(server, handlers.Material)
	mounting_api.RegisterMountingServiceServer(server, handlers.Mounting)
	snp_data_api.RegisterSnpDataServiceServer(server, handlers.SnpData)
	snp_filler_api.RegisterSnpFillerServiceServer(server, handlers.SnpFiller)
	snp_material_api.RegisterSnpMaterialServiceServer(server, handlers.SnpMaterial)
	snp_standard_api.RegisterSnpStandardServiceServer(server, handlers.SnpStandard)
	snp_type_api.RegisterSnpTypeServiceServer(server, handlers.SnpType)
	standard_api.RegisterStandardServiceServer(server, handlers.Standard)
	temperature_api.RegisterTemperatureServiceServer(server, handlers.Temperature)
	snp_size_api.RegisterSnpSizeServiceServer(server, handlers.SnpSize)
	snp_api.RegisterSnpDataServiceServer(server, handlers.Snp)
	order_api.RegisterOrderServiceServer(server, handlers.Order)

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
