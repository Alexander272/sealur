package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type PingHandlers struct {
	moment_api.UnimplementedPingServiceServer
}

func NewPingHandlers() *PingHandlers {
	return &PingHandlers{}
}

func (h *PingHandlers) Ping(ctx context.Context, req *moment_api.PingRequest) (*moment_api.PingResponse, error) {
	return &moment_api.PingResponse{Ping: "pong"}, nil
}
