package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/user_service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey int

const (
	clientIDKey contextKey = iota
)

func (h *Handler) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	clientName, err := h.authenticateClient(ctx)
	if err != nil {
		return nil, err
	}
	logger.Infof("query: %s, client: %s", info.FullMethod, clientName)

	ctx = context.WithValue(ctx, clientIDKey, clientName)
	return handler(ctx, req)
}

func (h *Handler) authenticateClient(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["service-name"], "")
		clientPassword := strings.Join(md["password"], "")
		if clientLogin != h.conf.Name || clientPassword != h.conf.Password {
			return "", fmt.Errorf("unknown user %s", clientLogin)
		}
		return clientLogin, nil
	}
	return "", fmt.Errorf("missing credentials")
}
