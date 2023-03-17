package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/google/uuid"
)

type PositionServiceNew struct {
	repo repository.Position
	snp  PositionSnp
}

func NewPositionService_New(repo repository.Position, snp PositionSnp) *PositionServiceNew {
	return &PositionServiceNew{
		repo: repo,
		snp:  snp,
	}
}

func (s *PositionServiceNew) Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error) {
	snpPosition, err := s.snp.Get(ctx, orderId)
	if err != nil {
		return nil, err
	}

	positions = append(positions, snpPosition...)

	return positions, nil
}

func (s *PositionServiceNew) GetFull(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error) {
	positions, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions. error: %w", err)
	}

	snpId := make([]string, 0, len(positions))
	snpIndex := make(map[string]int, 0)

	for i, op := range positions {
		if op.Type == position_model.PositionType_Snp {
			snpId = append(snpId, op.Id)
			snpIndex[op.Id] = i
		}
	}

	snpPositions, err := s.snp.GetFull(ctx, snpId)
	if err != nil {
		return nil, err
	}

	for _, ops := range snpPositions {
		index := snpIndex[ops.Main.PositionId]
		positions[index].SnpData = ops
	}

	return positions, nil
}

func (s *PositionServiceNew) Create(ctx context.Context, position *position_model.FullPosition) (string, error) {
	id, err := s.repo.Create(ctx, position)
	if err != nil {
		return "", fmt.Errorf("failed to create position. error: %w", err)
	}

	position.Id = id
	if position.Type == position_model.PositionType_Snp {
		if err := s.snp.Create(ctx, position); err != nil {
			return "", err
		}
	}
	return id, nil
}

func (s *PositionServiceNew) CreateSeveral(ctx context.Context, positions []*position_model.FullPosition, orderId string) error {
	var positionSnp []*position_model.FullPosition

	for _, p := range positions {
		id := uuid.New()
		p.Id = id.String()
		p.OrderId = orderId

		if p.Type == position_model.PositionType_Snp {
			positionSnp = append(positionSnp, p)
		}
	}

	if err := s.repo.CreateSeveral(ctx, positions); err != nil {
		return fmt.Errorf("failed to create several positions. error: %w", err)
	}

	if err := s.snp.CreateSeveral(ctx, positionSnp); err != nil {
		return err
	}

	return nil
}

func (s *PositionServiceNew) Update(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update position. error: %w", err)
	}

	if position.Type == position_model.PositionType_Snp {
		if err := s.snp.Update(ctx, position); err != nil {
			return err
		}
	}

	return nil
}

func (s *PositionServiceNew) Delete(ctx context.Context, positionId string) error {
	if err := s.repo.Delete(ctx, positionId); err != nil {
		return fmt.Errorf("failed to delete position. error: %w", err)
	}

	if err := s.snp.Delete(ctx, positionId); err != nil {
		return err
	}

	return nil
}
