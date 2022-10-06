package flange

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/flange_model"
)

type FlangeService struct {
	repo repository.Flange
}

func NewFlangeService(repo repository.Flange) *FlangeService {
	return &FlangeService{repo: repo}
}

func (s *FlangeService) GetFlangeSize(ctx context.Context, req *flange_api.GetFlangeSizeRequest) (models.FlangeSize, error) {
	size, err := s.repo.GetFlangeSize(ctx, req)
	if err != nil {
		return models.FlangeSize{}, fmt.Errorf("failed to get flange size. error: %w", err)
	}
	size.Area = math.Round(size.Area*1000) / 1000
	size.Diameter = math.Round(size.Diameter*1000) / 1000

	return size, nil
}

func (s *FlangeService) GetBasisFlangeSize(ctx context.Context, req *flange_api.GetBasisFlangeSizeRequest) (*flange_model.BasisFlangeSizeResponse, error) {
	res := &flange_model.BasisFlangeSizeResponse{}

	reqSize1 := models.GetBasisSize{IsUseRow: req.IsUseRow, StandId: req.StandId, Row: 0, IsInch: req.IsInch}
	sizes1, err := s.repo.GetBasisFlangeSizes(ctx, reqSize1)
	if err != nil {
		return nil, fmt.Errorf("failed to get size. error: %w", err)
	}

	sizeRow := []*flange_model.BasisFlangeSize{}

	curD := math.Inf(-1)
	curDn := ""
	for _, fs := range sizes1 {
		fs.Pn = math.Round(fs.Pn*1000) / 1000
		if req.IsInch {
			if curDn != fs.Dn {
				curDn = fs.Dn
				sizeRow = append(sizeRow, &flange_model.BasisFlangeSize{
					Dn: fs.Dn,
				})
			}
		} else {
			if curD != fs.D {
				curD = fs.D
				sizeRow = append(sizeRow, &flange_model.BasisFlangeSize{
					Dn: fmt.Sprint(fs.D),
				})
			}
		}
		sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, &flange_model.BasisFlangeSize_Pn{
			Pn:       fs.Pn,
			IsEmptyD: fs.IsEmptyD,
		})
	}
	res.SizeRow1 = sizeRow

	if req.IsUseRow {
		reqSize2 := models.GetBasisSize{IsUseRow: req.IsUseRow, StandId: req.StandId, Row: 1, IsInch: req.IsInch}
		sizes2, err := s.repo.GetBasisFlangeSizes(ctx, reqSize2)
		if err != nil {
			return nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		sizeRow = []*flange_model.BasisFlangeSize{}

		curD = math.Inf(-1)
		curDn = ""
		for _, fs := range sizes2 {
			fs.Pn = math.Round(fs.Pn*1000) / 1000
			if req.IsInch {
				if curDn != fs.Dn {
					curDn = fs.Dn
					sizeRow = append(sizeRow, &flange_model.BasisFlangeSize{
						Dn: fs.Dn,
					})
				}
			} else {
				if curD != fs.D {
					curD = fs.D
					sizeRow = append(sizeRow, &flange_model.BasisFlangeSize{
						Dn: fmt.Sprint(fs.D),
					})
				}
			}
			sizeRow[len(sizeRow)-1].Pn = append(sizeRow[len(sizeRow)-1].Pn, &flange_model.BasisFlangeSize_Pn{
				Pn:       fs.Pn,
				IsEmptyD: fs.IsEmptyD,
			})
		}

		res.SizeRow2 = sizeRow
	}

	return res, nil
}

func (s *FlangeService) GetFullFlangeSize(ctx context.Context, size *flange_api.GetFullFlangeSizeRequest) (*flange_api.FullFlangeSizeResponse, error) {
	sizes1, err := s.repo.GetFullFlangeSize(ctx, size, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get size for row 0. error: %w", err)
	}
	sizes2, err := s.repo.GetFullFlangeSize(ctx, size, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get size for row 1. error: %w", err)
	}

	sizeRow1, sizeRow2 := []*flange_model.FullFlangeSize{}, []*flange_model.FullFlangeSize{}
	for _, fsd := range sizes1 {
		sizeRow1 = append(sizeRow1, &flange_model.FullFlangeSize{
			Id:      fsd.Id,
			StandId: fsd.StandId,
			Pn:      math.Round(fsd.Pn*1000) / 1000,
			Dn:      fsd.Dn,
			Dmm:     math.Round(fsd.Dmm*1000) / 1000,
			D:       math.Round(fsd.D*1000) / 1000,
			D6:      math.Round(fsd.D6*1000) / 1000,
			DOut:    math.Round(fsd.DOut*1000) / 1000,
			X:       math.Round(fsd.X*1000) / 1000,
			A:       math.Round(fsd.A*1000) / 1000,
			H:       math.Round(fsd.H*1000) / 1000,
			S0:      math.Round(fsd.S0*1000) / 1000,
			S1:      math.Round(fsd.S1*1000) / 1000,
			Length:  math.Round(fsd.Length*1000) / 1000,
			Count:   fsd.Count,
			BoltId:  fsd.BoltId,
		})
	}
	if len(sizes2) > 0 {
		for _, fsd := range sizes2 {
			sizeRow2 = append(sizeRow2, &flange_model.FullFlangeSize{
				Id:      fsd.Id,
				StandId: fsd.StandId,
				Pn:      math.Round(fsd.Pn*1000) / 1000,
				Dn:      fsd.Dn,
				Dmm:     math.Round(fsd.Dmm*1000) / 1000,
				D:       math.Round(fsd.D*1000) / 1000,
				D6:      math.Round(fsd.D6*1000) / 1000,
				DOut:    math.Round(fsd.DOut*1000) / 1000,
				X:       math.Round(fsd.X*1000) / 1000,
				A:       math.Round(fsd.A*1000) / 1000,
				H:       math.Round(fsd.H*1000) / 1000,
				S0:      math.Round(fsd.S0*1000) / 1000,
				S1:      math.Round(fsd.S1*1000) / 1000,
				Length:  math.Round(fsd.Length*1000) / 1000,
				Count:   fsd.Count,
				BoltId:  fsd.BoltId,
			})
		}
	}

	res := &flange_api.FullFlangeSizeResponse{
		SizeRow1: sizeRow1,
		SizeRow2: sizeRow2,
	}

	return res, nil
}

func (s *FlangeService) CreateFlangeSize(ctx context.Context, size *flange_api.CreateFlangeSizeRequest) error {
	if err := s.repo.CreateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to create flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) CreateFlangeSizes(ctx context.Context, size *flange_api.CreateFlangeSizesRequest) error {
	if err := s.repo.CreateFlangeSizes(ctx, size); err != nil {
		return fmt.Errorf("failed to create flange sizes. error: %w", err)
	}
	return nil
}

func (s *FlangeService) UpdateFlangeSize(ctx context.Context, size *flange_api.UpdateFlangeSizeRequest) error {
	if err := s.repo.UpdateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to update flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteFlangeSize(ctx context.Context, size *flange_api.DeleteFlangeSizeRequest) error {
	if err := s.repo.DeleteFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to delete flange size. error: %w", err)
	}
	return nil
}
