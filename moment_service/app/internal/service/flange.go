package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type FlangeService struct {
	repo repository.Flange
}

func NewFlangeService(repo repository.Flange) *FlangeService {
	return &FlangeService{repo: repo}
}

func (s *FlangeService) GetFlangeSize(ctx context.Context, req *moment_api.GetFlangeSizeRequest) (models.FlangeSize, error) {
	size, err := s.repo.GetFlangeSize(ctx, req)
	if err != nil {
		return models.FlangeSize{}, fmt.Errorf("failed to get flange size. error: %w", err)
	}
	size.Area = math.Round(size.Area*1000) / 1000

	return size, nil
}

func (s *FlangeService) GetBasisFlangeSize(ctx context.Context, req *moment_api.GetBasisFlangeSizeRequest) (*moment_api.BasisFlangeSizeResponse, error) {
	res := &moment_api.BasisFlangeSizeResponse{}

	if req.IsUseRow {
		reqSize1 := models.GetBasisSize{IsUseRow: req.IsUseRow, StandId: req.StandId, Row: 0}
		sizes1, err := s.repo.GetBasisFlangeSizes(ctx, reqSize1)
		if err != nil {
			return nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		sizeRow := []*moment_api.BasisFlangeSize{}

		curD := math.Inf(-1)
		for _, fs := range sizes1 {
			fs.Pn = math.Round(fs.Pn*1000) / 1000
			if curD != fs.D {
				curD = fs.D
				sizeRow = append(sizeRow, &moment_api.BasisFlangeSize{
					Dn: fs.D,
				})
				sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, fs.Pn)
			} else {
				sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, fs.Pn)
			}
		}

		res.SizeRow1 = sizeRow

		reqSize2 := models.GetBasisSize{IsUseRow: req.IsUseRow, StandId: req.StandId, Row: 1}
		sizes2, err := s.repo.GetBasisFlangeSizes(ctx, reqSize2)
		if err != nil {
			return nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		sizeRow = []*moment_api.BasisFlangeSize{}

		curD = math.Inf(-1)
		for _, fs := range sizes2 {
			fs.Pn = math.Round(fs.Pn*1000) / 1000
			if curD != fs.D {
				curD = fs.D
				sizeRow = append(sizeRow, &moment_api.BasisFlangeSize{
					Dn: fs.D,
				})
				sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, fs.Pn)
			} else {
				sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, fs.Pn)
			}
		}

		res.SizeRow2 = sizeRow
	} else {
		reqSize := models.GetBasisSize{IsUseRow: req.IsUseRow, StandId: req.StandId, Row: 0}
		sizes, err := s.repo.GetBasisFlangeSizes(ctx, reqSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		sizeRow := []*moment_api.BasisFlangeSize{}

		curD := math.Inf(-1)
		for _, fs := range sizes {
			fs.Pn = math.Round(fs.Pn*1000) / 1000
			if curD != fs.D {
				curD = fs.D
				sizeRow = append(sizeRow, &moment_api.BasisFlangeSize{
					// Id: fs.Id,
					Dn: fs.D,
				})
				sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, fs.Pn)
			} else {
				sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, fs.Pn)
			}
		}

		res.SizeRow1 = sizeRow
	}

	return res, nil
}

func (s *FlangeService) GetFullFlangeSize(ctx context.Context, size *moment_api.GetFullFlangeSizeRequest) (*moment_api.FullFlangeSizeResponse, error) {
	sizes1, err := s.repo.GetFullFlangeSize(ctx, size, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get size for row 0. error: %w", err)
	}
	sizes2, err := s.repo.GetFullFlangeSize(ctx, size, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get size for row 1. error: %w", err)
	}

	sizeRow1, sizeRow2 := []*moment_api.FullFlangeSize{}, []*moment_api.FullFlangeSize{}
	for _, fsd := range sizes1 {
		sizeRow1 = append(sizeRow1, &moment_api.FullFlangeSize{
			Id:      fsd.Id,
			StandId: fsd.StandId,
			Pn:      fsd.Pn,
			D:       fsd.D,
			D6:      fsd.D6,
			DOut:    fsd.DOut,
			H:       fsd.H,
			S0:      fsd.S0,
			S1:      fsd.S1,
			Length:  fsd.Length,
			Count:   fsd.Count,
			BoltId:  fsd.BoltId,
		})
	}
	if len(sizes2) > 0 {
		for _, fsd := range sizes2 {
			sizeRow2 = append(sizeRow2, &moment_api.FullFlangeSize{
				Id:      fsd.Id,
				StandId: fsd.StandId,
				Pn:      fsd.Pn,
				D:       fsd.D,
				D6:      fsd.D6,
				DOut:    fsd.DOut,
				H:       fsd.H,
				S0:      fsd.S0,
				S1:      fsd.S1,
				Length:  fsd.Length,
				Count:   fsd.Count,
				BoltId:  fsd.BoltId,
			})
		}
	}

	res := &moment_api.FullFlangeSizeResponse{
		SizeRow1: sizeRow1,
		SizeRow2: sizeRow2,
	}

	return res, nil
}

func (s *FlangeService) CreateFlangeSize(ctx context.Context, size *moment_api.CreateFlangeSizeRequest) error {
	if err := s.repo.CreateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to create flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) UpdateFlangeSize(ctx context.Context, size *moment_api.UpdateFlangeSizeRequest) error {
	if err := s.repo.UpdateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to update flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteFlangeSize(ctx context.Context, size *moment_api.DeleteFlangeSizeRequest) error {
	if err := s.repo.DeleteFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to delete flange size. error: %w", err)
	}
	return nil
}
