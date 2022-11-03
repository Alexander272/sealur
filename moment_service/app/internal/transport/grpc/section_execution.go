package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetSectionExecution(ctx context.Context, req *device_api.GetSectionExecutionRequest,
) (*device_api.SectionExecutionResponse, error) {
	section, err := h.service.GetSectionExecution(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.SectionExecutionResponse{Section: section}, nil
}

func (h *DeviceHandlers) CreateSectionExecution(ctx context.Context, section *device_api.CreateSectionExecutionRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateSectionExecution(ctx, section)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewSectionExecution(ctx context.Context, section *device_api.CreateFewSectionExecutionRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewSectionExecution(ctx, section); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateSectionExecution(ctx context.Context, section *device_api.UpdateSectionExecutionRequest) (*device_api.Response, error) {
	if err := h.service.UpdateSectionExecution(ctx, section); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteSectionExecution(ctx context.Context, section *device_api.DeleteSectionExecutionRequest) (*device_api.Response, error) {
	if err := h.service.DeleteSectionExecution(ctx, section); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
