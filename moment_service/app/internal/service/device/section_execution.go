package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetSectionExecution(ctx context.Context, req *device_api.GetSectionExecutionRequest,
) (section []*device_model.SectionExecution, err error) {
	data, err := s.repo.GetSectionExecution(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get section execution. error: %w", err)
	}

	for _, item := range data {
		section = append(section, &device_model.SectionExecution{
			Id:    item.Id,
			DevId: item.DevId,
			Value: item.Value,
		})
	}
	return section, nil
}

func (s *DeviceService) CreateSectionExecution(ctx context.Context, section *device_api.CreateSectionExecutionRequest) (id string, err error) {
	id, err = s.repo.CreateSectionExecution(ctx, section)
	if err != nil {
		return "", fmt.Errorf("failed to create section execution. error: %w", err)
	}
	return id, err
}

func (s *DeviceService) CreateFewSectionExecution(ctx context.Context, section *device_api.CreateFewSectionExecutionRequest) error {
	if err := s.repo.CreateFewSectionExecution(ctx, section); err != nil {
		return fmt.Errorf("failed to create few section exection. error: %w", err)
	}
	return nil
}

func (s *DeviceService) UpdateSectionExecution(ctx context.Context, section *device_api.UpdateSectionExecutionRequest) error {
	if err := s.repo.UpdateSectionExecution(ctx, section); err != nil {
		return fmt.Errorf("failed to update section exection. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteSectionExecution(ctx context.Context, section *device_api.DeleteSectionExecutionRequest) error {
	if err := s.repo.DeleteSectionExecution(ctx, section); err != nil {
		return fmt.Errorf("failed to delete section exection. error: %w", err)
	}
	return nil
}
