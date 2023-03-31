package read

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/flange_model"
	"github.com/Alexander272/sealur_proto/api/moment/models/material_model"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
)

type FlangeService struct {
	flange    *flange.FlangeService
	materials *materials.MaterialsService
	gasket    *gasket.GasketService
}

func NewFlangeService(flange *flange.FlangeService, materials *materials.MaterialsService, gasket *gasket.GasketService) *FlangeService {
	return &FlangeService{
		flange:    flange,
		materials: materials,
		gasket:    gasket,
	}
}

func (s *FlangeService) GetFlange(ctx context.Context, req *read_api.GetFlangeRequest) (*read_api.GetFlangeResponse, error) {
	typeFlange, err := s.flange.GetTypeFlange(ctx, &flange_api.GetTypeFlangeRequest{})
	if err != nil {
		return nil, err
	}

	st, err := s.flange.GetStandarts(ctx, &flange_api.GetStandartsRequest{TypeId: typeFlange[0].Id})
	if err != nil {
		return nil, err
	}

	standarts := []*flange_model.StandartWithSize{}

	for _, s2 := range st {
		sizes, err := s.flange.GetBasisFlangeSize(ctx, &flange_api.GetBasisFlangeSizeRequest{
			IsUseRow: s2.IsNeedRow,
			StandId:  s2.Id,
			IsInch:   s2.IsInch,
		})
		if err != nil {
			return nil, err
		}

		standarts = append(standarts, &flange_model.StandartWithSize{
			Id:             s2.Id,
			Title:          s2.Title,
			TypeId:         s2.TypeId,
			Sizes:          sizes,
			TitleDn:        s2.TitleDn,
			TitlePn:        s2.TitlePn,
			IsNeedRow:      s2.IsNeedRow,
			Rows:           s2.Rows,
			IsInch:         s2.IsInch,
			HasDesignation: s2.HasDesignation,
		})
	}

	gasket, err := s.gasket.GetGasketWithThick(ctx, &gasket_api.GetGasketRequest{})
	if err != nil {
		return nil, err
	}

	env, err := s.gasket.GetEnv(ctx, &gasket_api.GetEnvRequest{})
	if err != nil {
		return nil, err
	}

	boltMaterials, err := s.materials.GetMaterials(ctx, &material_api.GetMaterialsRequest{Type: material_model.MaterialType_bolt})
	if err != nil {
		return nil, err
	}

	flangeMaterials, err := s.materials.GetMaterials(ctx, &material_api.GetMaterialsRequest{Type: material_model.MaterialType_flange})
	if err != nil {
		return nil, err
	}

	res := &read_api.GetFlangeResponse{
		TypeFlange:      typeFlange,
		Standarts:       standarts,
		Gaskets:         gasket,
		Env:             env,
		BoltMaterials:   boltMaterials,
		FlangeMaterials: flangeMaterials,
	}

	return res, nil
}
