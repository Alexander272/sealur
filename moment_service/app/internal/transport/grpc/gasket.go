package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type GasketHandlers struct {
	service service.Gasket
}

func NewGasketService(service service.Gasket) *GasketHandlers {
	return &GasketHandlers{service: service}
}

func (h *GasketHandlers) GetGasket(ctx context.Context, req *moment_proto.GetGasketRequest) (*moment_proto.GasketResponse, error) {
	gasket, err := h.service.GetGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.GasketResponse{Gasket: gasket}, nil
}

func (h *GasketHandlers) CreateGasket(ctx context.Context, gasket *moment_proto.CreateGasketRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasket(ctx context.Context, gasket *moment_proto.UpdateGasketRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteGasket(ctx context.Context, gasket *moment_proto.DeleteGasketRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

//---

func (h *GasketHandlers) GetGasketType(ctx context.Context, req *moment_proto.GetGasketTypeRequest) (*moment_proto.GasketTypeResponse, error) {
	typeGasket, err := h.service.GetTypeGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.GasketTypeResponse{GasketType: typeGasket}, nil
}

func (h *GasketHandlers) CreateGasketType(ctx context.Context, gasket *moment_proto.CreateGasketTypeRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateTypeGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasketType(ctx context.Context, gasket *moment_proto.UpdateGasketTypeRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketType(ctx context.Context, gasket *moment_proto.DeleteGasketTypeRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

//---

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

//---

func (h *GasketHandlers) CreateGasketData(ctx context.Context, data *moment_proto.CreateGasketDataRequest) (*moment_proto.Response, error) {
	if err := h.service.CreateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) UpdateGasketData(ctx context.Context, data *moment_proto.UpdateGasketDataRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketData(ctx context.Context, data *moment_proto.DeleteGasketDataRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
