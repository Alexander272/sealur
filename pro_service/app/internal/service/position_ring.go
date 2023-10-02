package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
)

type PositionRingService struct {
	repo repository.PositionRing
}

func NewPositionRingService(repo repository.PositionRing) *PositionRingService {
	return &PositionRingService{
		repo: repo,
	}
}

type PositionRing interface {
	Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error)
	GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionRing, error)
	Create(context.Context, *position_model.FullPosition) error
	Update(context.Context, *position_model.FullPosition) error
	Copy(ctx context.Context, targetId string, position *position_api.CopyPosition) (string, error)
}

func (s *PositionRingService) Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error) {
	positions, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get ring positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionRingService) GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionRing, error) {
	positions, err := s.repo.GetFull(ctx, positionsId)
	if err != nil {
		return nil, fmt.Errorf("failed to get full ring positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionRingService) Create(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Create(ctx, position); err != nil {
		return fmt.Errorf("failed to create ring position. error: %w", err)
	}
	return nil
}

func (s *PositionRingService) Update(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update ring position. error: %w", err)
	}
	return nil
}

func (s *PositionRingService) Copy(ctx context.Context, targetId string, position *position_api.CopyPosition) (string, error) {
	drawing, err := s.repo.Copy(ctx, targetId, position)
	if err != nil {
		return "", fmt.Errorf("failed to copy ring position. error: %w", err)
	}
	return drawing, nil
}
