package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
)

type PositionSnpService struct {
	repo repository.PositionSnp
}

func NewPositionSnpService(repo repository.PositionSnp) *PositionSnpService {
	return &PositionSnpService{repo: repo}
}

func (s *PositionSnpService) Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error) {
	positions, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionSnpService) GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionSnp, error) {
	positions, err := s.repo.GetFull(ctx, positionsId)
	if err != nil {
		return nil, fmt.Errorf("failed to get full snp positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionSnpService) Create(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Create(ctx, position); err != nil {
		return fmt.Errorf("failed to create snp position. error: %w", err)
	}
	return nil
}

func (s *PositionSnpService) CreateSeveral(ctx context.Context, positions []*position_model.FullPosition) error {
	if err := s.repo.CreateSeveral(ctx, positions); err != nil {
		return fmt.Errorf("failed to create several snp position. error: %w", err)
	}
	return nil
}

func (s *PositionSnpService) Update(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update snp position. error: %w", err)
	}
	return nil
}

func (s *PositionSnpService) Copy(ctx context.Context, targetId string, position *position_api.CopyPosition) (string, error) {
	drawing, err := s.repo.Copy(ctx, targetId, position)
	if err != nil {
		return "", fmt.Errorf("failed to copy snp position. error: %w", err)
	}
	return drawing, nil
}

func (s *PositionSnpService) Delete(ctx context.Context, positionId string) error {
	if err := s.repo.Delete(ctx, positionId); err != nil {
		return fmt.Errorf("failed to delete snp position. error: %w", err)
	}
	return nil
}
