package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/service"
	handlers "github.com/Alexander272/sealur/email_service/internal/transport/grpc"
	"github.com/Alexander272/sealur/email_service/pkg/email/smtp"
	"github.com/Alexander272/sealur/email_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"google.golang.org/grpc"
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

	logger.Debug(conf.Email)
	emailSender, err := smtp.NewSMTPSender(conf.Email.Sender, conf.Email.Password, conf.Email.Host, conf.Email.Port)
	if err != nil {
		logger.Fatalf("error initializing email sender", err)
	}

	services := service.NewServices(emailSender, conf.Recipients)
	handlers := handlers.NewHandler(services, conf.Api)

	//* GRPC Server

	// cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	// if err != nil {
	// 	logger.Fatalf("failed to load certificate. error: %w", err)
	// }

	// opts := []grpc.ServerOption{
	// 	// grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	// 	// grpc.Creds(insecure.NewCredentials()),
	// 	grpc.UnaryInterceptor(handlers.UnaryInterceptor),
	// }

	// server := grpc.NewServer(opts...)
	server := grpc.NewServer()
	email_api.RegisterEmailServiceServer(server, handlers)

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
