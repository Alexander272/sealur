package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
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

func (s *PositionSnpService) CreateSeveral(ctx context.Context, positions []*position_model.FullPosition) error {
	if err := s.repo.CreateSeveral(ctx, positions); err != nil {
		return fmt.Errorf("failed to create several position snp. error: %w", err)
	}
	return nil
}
