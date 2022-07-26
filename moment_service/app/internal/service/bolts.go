package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (s *FlangeService) GetBolts(ctx context.Context, req *moment_api.GetBoltsRequest) (bolts []*moment_api.Bolt, err error) {
	data, err := s.repo.GetBolts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get bolts. error: %w", err)
	}

	for _, item := range data {
		item.Area = math.Round(item.Area*1000) / 1000

		bolts = append(bolts, &moment_api.Bolt{
			Id:       item.Id,
			Title:    item.Title,
			Diameter: item.Diameter,
			Area:     item.Area,
		})
	}

	return bolts, nil
}

func (s *FlangeService) CreateBolt(ctx context.Context, bolt *moment_api.CreateBoltRequest) error {
	if err := s.repo.CreateBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to create bolt. error: %w", err)
	}
	return nil
}

func (s *FlangeService) UpdateBolt(ctx context.Context, bolt *moment_api.UpdateBoltRequest) error {
	if err := s.repo.UpdateBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to update bolt. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteBolt(ctx context.Context, bolt *moment_api.DeleteBoltRequest) error {
	if err := s.repo.DeleteBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to delete bolt. error: %w", err)
	}
	return nil
}
