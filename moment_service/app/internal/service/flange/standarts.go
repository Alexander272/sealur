package flange

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (s *FlangeService) GetStandarts(ctx context.Context, req *moment_api.GetStandartsRequest) (standarts []*moment_api.Standart, err error) {
	data, err := s.repo.GetStandarts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get standarts. error: %w", err)
	}

	for _, item := range data {
		var rows []string
		if item.Rows != "" {
			rows = strings.Split(item.Rows, "; ")
		}

		standarts = append(standarts, &moment_api.Standart{
			Id:             item.Id,
			Title:          item.Title,
			TypeId:         item.TypeId,
			TitleDn:        item.TitleDn,
			TitlePn:        item.TitlePn,
			IsNeedRow:      item.IsNeedRow,
			IsInch:         item.IsInch,
			HasDesignation: item.HasDesignation,
			Rows:           rows,
		})
	}

	return standarts, nil
}

func (s *FlangeService) GetStandartsWithSize(ctx context.Context, req *moment_api.GetStandartsRequest) (standarts []*moment_api.StandartWithSize, err error) {
	data, err := s.repo.GetStandarts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get standarts. error: %w", err)
	}

	for _, item := range data {
		sizes, err := s.GetBasisFlangeSize(ctx, &moment_api.GetBasisFlangeSizeRequest{
			IsUseRow: item.IsNeedRow,
			StandId:  item.Id,
			IsInch:   item.IsInch,
		})
		if err != nil {
			return nil, err
		}

		var rows []string
		if item.Rows != "" {
			rows = strings.Split(item.Rows, "; ")
		}

		standarts = append(standarts, &moment_api.StandartWithSize{
			Id:             item.Id,
			Title:          item.Title,
			TypeId:         item.TypeId,
			TitleDn:        item.TitleDn,
			TitlePn:        item.TitlePn,
			IsNeedRow:      item.IsNeedRow,
			IsInch:         item.IsInch,
			HasDesignation: item.HasDesignation,
			Rows:           rows,
			Sizes:          sizes,
		})
	}

	return standarts, nil
}

func (s *FlangeService) CreateStandart(ctx context.Context, stand *moment_api.CreateStandartRequest) (id string, err error) {
	id, err = s.repo.CreateStandart(ctx, stand)
	if err != nil {
		return "", fmt.Errorf("failed to create standart. error: %w", err)
	}

	return id, nil
}

func (s *FlangeService) UpdateStandart(ctx context.Context, stand *moment_api.UpdateStandartRequest) error {
	if err := s.repo.UpdateStandart(ctx, stand); err != nil {
		return fmt.Errorf("failed to update standart. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteStandart(ctx context.Context, stand *moment_api.DeleteStandartRequest) error {
	if err := s.repo.DeleteStandart(ctx, stand); err != nil {
		return fmt.Errorf("failed to delete standart. error: %w", err)
	}
	return nil
}
