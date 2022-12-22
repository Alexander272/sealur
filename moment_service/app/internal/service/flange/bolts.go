package flange

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/flange_model"
)

// Получение размеров болтов по id
func (s *FlangeService) GetBolt(ctx context.Context, boltId string) (bolt models.BoltSize, err error) {
	bolt, err = s.repo.GetBolt(ctx, boltId)
	if err != nil {
		return bolt, fmt.Errorf("failed to get bolt by id. error: %w", err)
	}
	bolt.Diameter = math.Round(bolt.Diameter*1000) / 1000
	bolt.Area = math.Round(bolt.Area*1000) / 1000

	return bolt, err
}

func (s *FlangeService) GetBolts(ctx context.Context, req *flange_api.GetBoltsRequest) (bolts []*flange_model.Bolt, err error) {
	data, err := s.repo.GetBolts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get bolts. error: %w", err)
	}

	for _, item := range data {
		bolts = append(bolts, &flange_model.Bolt{
			Id:       item.Id,
			Title:    item.Title,
			Diameter: math.Round(item.Diameter*1000) / 1000,
			Area:     math.Round(item.Area*1000) / 1000,
			IsInch:   item.IsInch,
		})
	}

	return bolts, nil
}

func (s *FlangeService) GetAllBolts(ctx context.Context, req *flange_api.GetBoltsRequest) (bolts []*flange_model.Bolt, err error) {
	data, err := s.repo.GetAllBolts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get bolts. error: %w", err)
	}

	for _, item := range data {
		bolts = append(bolts, &flange_model.Bolt{
			Id:       item.Id,
			Title:    item.Title,
			Diameter: math.Round(item.Diameter*1000) / 1000,
			Area:     math.Round(item.Area*1000) / 1000,
			IsInch:   item.IsInch,
		})
	}

	return bolts, nil
}

func (s *FlangeService) CreateBolt(ctx context.Context, bolt *flange_api.CreateBoltRequest) error {
	if err := s.repo.CreateBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to create bolt. error: %w", err)
	}
	return nil
}

func (s *FlangeService) CreateBolts(ctx context.Context, bolt *flange_api.CreateBoltsRequest) error {
	if err := s.repo.CreateBolts(ctx, bolt); err != nil {
		return fmt.Errorf("failed to create bolt. error: %w", err)
	}
	return nil
}

func (s *FlangeService) UpdateBolt(ctx context.Context, bolt *flange_api.UpdateBoltRequest) error {
	if err := s.repo.UpdateBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to update bolt. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteBolt(ctx context.Context, bolt *flange_api.DeleteBoltRequest) error {
	if err := s.repo.DeleteBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to delete bolt. error: %w", err)
	}
	return nil
}
