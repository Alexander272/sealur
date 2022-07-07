package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type PingHandlers struct{}

func NewPingHandlers() *PingHandlers {
	return &PingHandlers{}
}

func (h *PingHandlers) Ping(ctx context.Context, req *moment_proto.PingRequest) (*moment_proto.PingResponse, error) {
	return &moment_proto.PingResponse{Ping: "pong"}, nil
}
