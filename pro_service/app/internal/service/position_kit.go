package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
)

type PositionKitService struct {
	repo repository.PositionKit
}

func NewPositionKitService(repo repository.PositionKit) *PositionKitService {
	return &PositionKitService{
		repo: repo,
	}
}

type PositionKit interface {
	Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error)
	GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionRingsKit, error)
	Create(ctx context.Context, position *position_model.FullPosition) error
	Update(ctx context.Context, position *position_model.FullPosition) error
	Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error)
}

func (s *PositionKitService) Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error) {
	positions, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get rings kit positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionKitService) GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionRingsKit, error) {
	positions, err := s.repo.GetFull(ctx, positionsId)
	if err != nil {
		return nil, fmt.Errorf("failed to get full rings kit position. error: %w", err)
	}
	return positions, nil
}

func (s *PositionKitService) Create(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Create(ctx, position); err != nil {
		return fmt.Errorf("failed to create rings kit position. error: %w", err)
	}
	return nil
}

func (s *PositionKitService) Update(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update rings kit position. error: %w", err)
	}
	return nil
}

func (s *PositionKitService) Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error) {
	drawing, err := s.repo.Copy(ctx, targetPositionId, position)
	if err != nil {
		return "", fmt.Errorf("failed to copy rings kit position. error: %w", err)
	}
	return drawing, nil
}
