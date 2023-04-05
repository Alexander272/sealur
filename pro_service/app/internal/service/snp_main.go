package service

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro/models/mounting_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_model"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
)

type SnpService struct {
	filler      SnpFiller
	material    SnpMaterial
	snpType     SnpType
	mounting    Mounting
	standard    Standard
	snpStandard SnpStandard
	snpData     SnpData
	snpSize     SnpSize
}

func NewSnpService(filler SnpFiller, material SnpMaterial, snpType SnpType, mounting Mounting, standard Standard, snpData SnpData,
	snpStandard SnpStandard, snpSize SnpSize,
) *SnpService {
	return &SnpService{
		filler:      filler,
		material:    material,
		snpType:     snpType,
		mounting:    mounting,
		standard:    standard,
		snpData:     snpData,
		snpStandard: snpStandard,
		snpSize:     snpSize,
	}
}

func (s *SnpService) Get(ctx context.Context, req *snp_api.GetSnp) (snp *snp_api.Snp, err error) {
	snpData, err := s.snpData.Get(ctx, &snp_data_api.GetSnpData{TypeId: req.SnpTypeId})
	if err != nil {
		return nil, err
	}

	snpSize, err := s.snpSize.Get(ctx, &snp_size_api.GetSnpSize{TypeId: req.SnpTypeId, HasD2: req.HasD2})
	if err != nil {
		return nil, err
	}

	snp = &snp_api.Snp{
		Snp:   snpData,
		Sizes: snpSize,
	}

	return snp, nil
}

func (s *SnpService) GetData(ctx context.Context, req *snp_api.GetSnpData) (snpData *snp_model.SnpData, err error) {
	var mounting []*mounting_model.Mounting
	snpData = &snp_model.SnpData{}

	if req.StandardId == "" {
		standard, err := s.snpStandard.GetDefault(ctx)
		if err != nil {
			return nil, err
		}
		req.StandardId = standard.Standard.Id
		req.SnpStandardId = standard.Id

		mounting, err = s.mounting.GetAll(ctx, &mounting_api.GetAllMountings{})
		if err != nil {
			return nil, err
		}

	}

	snpData.Mounting = mounting

	snpTypes, err := s.snpType.GetWithFlange(ctx, req)
	if err != nil {
		return nil, err
	}

	fillers, err := s.filler.GetAll(ctx, &snp_filler_api.GetSnpFillers{StandardId: req.StandardId})
	if err != nil {
		return nil, err
	}

	materials, err := s.material.Get(ctx, &snp_material_api.GetSnpMaterial{StandardId: req.StandardId})
	if err != nil {
		return nil, err
	}

	snpData.FlangeTypes = snpTypes
	snpData.Fillers = fillers
	snpData.Materials = materials

	return snpData, nil
}
