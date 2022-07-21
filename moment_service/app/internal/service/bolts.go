package service

import (
	"context"
	"fmt"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (s *FlangeService) GetBolts(ctx context.Context, req *moment_proto.GetBoltsRequest) (bolts []*moment_proto.Bolt, err error) {
	data, err := s.repo.GetBolts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get bolts. error: %w", err)
	}

	for _, item := range data {
		b := moment_proto.Bolt(item)
		bolts = append(bolts, &b)
	}

	return bolts, nil
}

func (s *FlangeService) CreateBolt(ctx context.Context, bolt *moment_proto.CreateBoltRequest) error {
	if err := s.repo.CreateBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to create bolt. error: %w", err)
	}
	return nil
}

func (s *FlangeService) UpdateBolt(ctx context.Context, bolt *moment_proto.UpdateBoltRequest) error {
	if err := s.repo.UpdateBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to update bolt. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteBolt(ctx context.Context, bolt *moment_proto.DeleteBoltRequest) error {
	if err := s.repo.DeleteBolt(ctx, bolt); err != nil {
		return fmt.Errorf("failed to delete bolt. error: %w", err)
	}
	return nil
}
