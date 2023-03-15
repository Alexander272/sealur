package grpc

import (
	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
)

type User interface {
	user_api.UserServiceServer
}

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	User
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,

		User: NewUserHandler(service.User),
	}
}

// func (h *Handler) Ping(ctx context.Context, req *user_api.PingRequest) (*user_api.PingResponse, error) {
// 	return &user_api.PingResponse{Ping: "pong"}, nil
// }
