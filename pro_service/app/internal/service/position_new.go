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

func (s *PositionServiceNew) CreateSeveral(ctx context.Context, positions []*position_model.Position, orderId string) error {
	var positionSnp []*position_model.Position

	for _, p := range positions {
		id := uuid.New()
		p.Id = id.String()
		p.OrderId = orderId

		if p.Type == position_model.PositionType_snp {
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
