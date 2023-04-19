package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/repo"
	"github.com/Alexander272/sealur_proto/api/user/models/user_model"
)

type RegionService struct {
	repo repo.Region
}

func NewRegionService(repo repo.Region) *RegionService {
	return &RegionService{
		repo: repo,
	}
}

func (s *RegionService) GetManagerByRegion(ctx context.Context, region string) (*user_model.User, error) {
	user, err := s.repo.GetManagerByRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get manager by region. error: %w", err)
	}
	return user, nil
}
