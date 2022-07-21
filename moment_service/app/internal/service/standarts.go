package service

import (
	"context"
	"fmt"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (s *FlangeService) GetTypeFlange(ctx context.Context, req *moment_proto.GetTypeFlangeRequest) (typeFlange []*moment_proto.TypeFlange, err error) {
	data, err := s.repo.GetTypeFlange(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get types flange. error: %w", err)
	}

	for _, item := range data {
		t := moment_proto.TypeFlange(item)
		typeFlange = append(typeFlange, &t)
	}

	return typeFlange, nil
}

func (s *FlangeService) CreateTypeFlange(ctx context.Context, typeFlange *moment_proto.CreateTypeFlangeRequest) (id string, err error) {
	id, err = s.repo.CreateTypeFlange(ctx, typeFlange)
	if err != nil {
		return "", fmt.Errorf("failed to create type flange. error: %w", err)
	}

	return id, nil
}

func (s *FlangeService) UpdateTypeFlange(ctx context.Context, typeFlange *moment_proto.UpdateTypeFlangeRequest) error {
	if err := s.repo.UpdateTypeFlange(ctx, typeFlange); err != nil {
		return fmt.Errorf("failed to update type flange. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteTypeFlange(ctx context.Context, typeFlange *moment_proto.DeleteTypeFlangeRequest) error {
	if err := s.repo.DeleteTypeFlange(ctx, typeFlange); err != nil {
		return fmt.Errorf("failed to delete type flange. error: %w", err)
	}
	return nil
}

func (s *FlangeService) GetStandarts(ctx context.Context, req *moment_proto.GetStandartsRequest) (standarts []*moment_proto.Standart, err error) {
	data, err := s.repo.GetStandarts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get standarts. error: %w", err)
	}

	for _, item := range data {
		s := moment_proto.Standart(item)
		standarts = append(standarts, &s)
	}

	return standarts, nil
}

func (s *FlangeService) CreateStandart(ctx context.Context, stand *moment_proto.CreateStandartRequest) (id string, err error) {
	id, err = s.repo.CreateStandart(ctx, stand)
	if err != nil {
		return "", fmt.Errorf("failed to create standart. error: %w", err)
	}

	return id, nil
}

func (s *FlangeService) UpdateStandart(ctx context.Context, stand *moment_proto.UpdateStandartRequest) error {
	if err := s.repo.UpdateStandart(ctx, stand); err != nil {
		return fmt.Errorf("failed to update standart. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteStandart(ctx context.Context, stand *moment_proto.DeleteStandartRequest) error {
	if err := s.repo.DeleteStandart(ctx, stand); err != nil {
		return fmt.Errorf("failed to delete standart. error: %w", err)
	}
	return nil
}
