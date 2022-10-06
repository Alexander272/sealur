package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment"
)

type PingHandlers struct {
	moment.UnimplementedPingServiceServer
}

func NewPingHandlers() *PingHandlers {
	return &PingHandlers{}
}

func (h *PingHandlers) Ping(ctx context.Context, req *moment.PingRequest) (*moment.PingResponse, error) {
	return &moment.PingResponse{Ping: "pong"}, nil
}
