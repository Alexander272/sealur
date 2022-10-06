package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
)

func (h *GasketHandlers) GetEnv(ctx context.Context, req *gasket_api.GetEnvRequest) (*gasket_api.EnvResponse, error) {
	env, err := h.service.GetEnv(ctx, req)
	if err != nil {
		return nil, err
	}

	return &gasket_api.EnvResponse{Env: env}, nil
}

func (h *GasketHandlers) CreateEnv(ctx context.Context, env *gasket_api.CreateEnvRequest) (*gasket_api.IdResponse, error) {
	id, err := h.service.CreateEnv(ctx, env)
	if err != nil {
		return nil, err
	}

	return &gasket_api.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateEnv(ctx context.Context, env *gasket_api.UpdateEnvRequest) (*gasket_api.Response, error) {
	if err := h.service.UpdateEnv(ctx, env); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) DeleteEnv(ctx context.Context, env *gasket_api.DeleteEnvRequest) (*gasket_api.Response, error) {
	if err := h.service.DeleteEnv(ctx, env); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

//---
func (h *GasketHandlers) CreateManyEnvData(ctx context.Context, data *gasket_api.CreateManyEnvDataRequest) (*gasket_api.Response, error) {
	if err := h.service.CreateManyEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) CreateEnvData(ctx context.Context, data *gasket_api.CreateEnvDataRequest) (*gasket_api.Response, error) {
	if err := h.service.CreateEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) UpdateEnvData(ctx context.Context, data *gasket_api.UpdateEnvDataRequest) (*gasket_api.Response, error) {
	if err := h.service.UpdateEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) DeleteEnvData(ctx context.Context, data *gasket_api.DeleteEnvDataRequest) (*gasket_api.Response, error) {
	if err := h.service.DeleteEnvData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}
