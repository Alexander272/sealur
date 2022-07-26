package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (h *GasketHandlers) GetEnv(ctx context.Context, req *moment_api.GetEnvRequest) (*moment_api.EnvResponse, error) {
	env, err := h.service.GetEnv(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.EnvResponse{Env: env}, nil
}

func (h *GasketHandlers) CreateEnv(ctx context.Context, env *moment_api.CreateEnvRequest) (*moment_api.IdResponse, error) {
	id, err := h.service.CreateEnv(ctx, env)
	if err != nil {
		return nil, err
	}

	return &moment_api.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateEnv(ctx context.Context, env *moment_api.UpdateEnvRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateEnv(ctx, env); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *GasketHandlers) DeleteEnv(ctx context.Context, env *moment_api.DeleteEnvRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteEnv(ctx, env); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

//---

func (h *GasketHandlers) CreateEnvData(ctx context.Context, data *moment_api.CreateEnvDataRequest) (*moment_api.Response, error) {
	if err := h.service.CreateEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *GasketHandlers) UpdateEnvData(ctx context.Context, data *moment_api.UpdateEnvDataRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *GasketHandlers) DeleteEnvData(ctx context.Context, data *moment_api.DeleteEnvDataRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
