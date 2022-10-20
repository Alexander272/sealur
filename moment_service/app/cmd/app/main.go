package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander272/sealur/moment_service/internal/config"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur/moment_service/internal/service"
	handlers "github.com/Alexander272/sealur/moment_service/internal/transport/grpc"
	"github.com/Alexander272/sealur/moment_service/pkg/database/postgres"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
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

	//* Services, Repos & API Handlers

	repos := repository.NewRepo(db)
	services := service.NewServices(repos)
	handlers := handlers.NewHandler(services, conf.Api)

	//* GRPC Server

	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		logger.Fatalf("failed to load certificate. error: %s", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		// grpc.UnaryInterceptor(handlers.UnaryInterceptor),
	}

	server := grpc.NewServer(opts...)
	moment.RegisterPingServiceServer(server, handlers.Ping)
	material_api.RegisterMaterialsServiceServer(server, handlers.Materials)
	gasket_api.RegisterGasketServiceServer(server, handlers.Gasket)
	flange_api.RegisterFlangeServiceServer(server, handlers.Flange)
	read_api.RegisterReadServiceServer(server, handlers.Read)
	calc_api.RegisterCalcServiceServer(server, handlers.Calc)

	listener, err := net.Listen("tcp", ":"+conf.Http.Port)
	if err != nil {
		logger.Fatalf("failed to create grpc listener: %s", err.Error())
	}

	go func() {
		if err = server.Serve(listener); err != nil {
			logger.Fatalf("failed to start server: %s", err.Error())
		}
	}()
	logger.Infof("Application started on port: %s", conf.Http.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	server.GracefulStop()
}
