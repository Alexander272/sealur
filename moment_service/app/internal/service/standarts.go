package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (s *FlangeService) GetTypeFlange(ctx context.Context, req *moment_api.GetTypeFlangeRequest) (typeFlange []*moment_api.TypeFlange, err error) {
	data, err := s.repo.GetTypeFlange(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get types flange. error: %w", err)
	}

	for _, item := range data {
		typeFlange = append(typeFlange, &moment_api.TypeFlange{
			Id:    item.Id,
			Title: item.Title,
		})
	}

	return typeFlange, nil
}

func (s *FlangeService) CreateTypeFlange(ctx context.Context, typeFlange *moment_api.CreateTypeFlangeRequest) (id string, err error) {
	id, err = s.repo.CreateTypeFlange(ctx, typeFlange)
	if err != nil {
		return "", fmt.Errorf("failed to create type flange. error: %w", err)
	}

	return id, nil
}

func (s *FlangeService) UpdateTypeFlange(ctx context.Context, typeFlange *moment_api.UpdateTypeFlangeRequest) error {
	if err := s.repo.UpdateTypeFlange(ctx, typeFlange); err != nil {
		return fmt.Errorf("failed to update type flange. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteTypeFlange(ctx context.Context, typeFlange *moment_api.DeleteTypeFlangeRequest) error {
	if err := s.repo.DeleteTypeFlange(ctx, typeFlange); err != nil {
		return fmt.Errorf("failed to delete type flange. error: %w", err)
	}
	return nil
}

func (s *FlangeService) GetStandarts(ctx context.Context, req *moment_api.GetStandartsRequest) (standarts []*moment_api.Standart, err error) {
	data, err := s.repo.GetStandarts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get standarts. error: %w", err)
	}

	for _, item := range data {
		standarts = append(standarts, &moment_api.Standart{
			Id:     item.Id,
			Title:  item.Title,
			TypeId: item.TypeId,
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
