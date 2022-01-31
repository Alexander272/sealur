package server

import (
	"context"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(conf *config.Config, handler http.Handler) *Server {
	//TODO добавить автоматическое получение сертификата или что-то в этом роде
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + conf.Http.Port,
			Handler:        handler,
			ReadTimeout:    conf.Http.ReadTimeout,
			WriteTimeout:   conf.Http.WriteTimeout,
			MaxHeaderBytes: conf.Http.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServeTLS("localhost.cert", "localhost.key")
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
