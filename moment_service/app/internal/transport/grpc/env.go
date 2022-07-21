package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (h *GasketHandlers) GetEnv(ctx context.Context, req *moment_proto.GetEnvRequest) (*moment_proto.EnvResponse, error) {
	env, err := h.service.GetEnv(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.EnvResponse{Env: env}, nil
}

func (h *GasketHandlers) CreateEnv(ctx context.Context, env *moment_proto.CreateEnvRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateEnv(ctx, env)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateEnv(ctx context.Context, env *moment_proto.UpdateEnvRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateEnv(ctx, env); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteEnv(ctx context.Context, env *moment_proto.DeleteEnvRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteEnv(ctx, env); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

//---

func (h *GasketHandlers) CreateEnvData(ctx context.Context, data *moment_proto.CreateEnvDataRequest) (*moment_proto.Response, error) {
	if err := h.service.CreateEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) UpdateEnvData(ctx context.Context, data *moment_proto.UpdateEnvDataRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteEnvData(ctx context.Context, data *moment_proto.DeleteEnvDataRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
