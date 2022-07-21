package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type FlangeService struct {
	repo repository.Flange
}

func NewFlangeService(repo repository.Flange) *FlangeService {
	return &FlangeService{repo: repo}
}

func (s *FlangeService) GetFlangeSize(ctx context.Context, req *moment_proto.GetFlangeSizeRequest) (models.FlangeSize, error) {
	size, err := s.repo.GetFlangeSize(ctx, req)
	if err != nil {
		return models.FlangeSize{}, fmt.Errorf("failed to get flange size. error: %w", err)
	}

	return size, nil
}

func (s *FlangeService) CreateFlangeSize(ctx context.Context, size *moment_proto.CreateFlangeSizeRequest) error {
	if err := s.repo.CreateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to create flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) UpdateFlangeSize(ctx context.Context, size *moment_proto.UpdateFlangeSizeRequest) error {
	if err := s.repo.UpdateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to update flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteFlangeSize(ctx context.Context, size *moment_proto.DeleteFlangeSizeRequest) error {
	if err := s.repo.DeleteFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to delete flange size. error: %w", err)
	}
	return nil
}
