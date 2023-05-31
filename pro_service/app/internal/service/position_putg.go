package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
)

type PositionPutgService struct {
	repo repository.PositionPutg
}

func NewPositionPutgService(repo repository.PositionPutg) *PositionPutgService {
	return &PositionPutgService{
		repo: repo,
	}
}

func (s *PositionPutgService) Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error) {
	positions, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed get putg positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionPutgService) GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionPutg, error) {
	positions, err := s.repo.GetFull(ctx, positionsId)
	if err != nil {
		return nil, fmt.Errorf("failed to get full putg positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionPutgService) Create(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Create(ctx, position); err != nil {
		return fmt.Errorf("failed to create putg position. error: %w", err)
	}
	return nil
}

func (s *PositionPutgService) Update(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update putg position. error: %w", err)
	}
	return nil
}

func (s *PositionPutgService) Copy(ctx context.Context, targetId string, position *position_api.CopyPosition) (string, error) {
	drawing, err := s.repo.Copy(ctx, targetId, position)
	if err != nil {
		return "", fmt.Errorf("failed to copy putg position. error: %w", err)
	}
	return drawing, nil
}

func (s *PositionPutgService) Delete(ctx context.Context, positionId string) error {
	if err := s.repo.Delete(ctx, positionId); err != nil {
		return fmt.Errorf("failed to delete putg position. error: %w", err)
	}
	return nil
}
