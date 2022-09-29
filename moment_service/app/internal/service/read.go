package service

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type ReadService struct {
	flange    *FlangeService
	materials *MaterialsService
	gasket    *GasketService
}

func NewReadService(flange *FlangeService, materials *MaterialsService, gasket *GasketService) *ReadService {
	return &ReadService{
		flange:    flange,
		materials: materials,
		gasket:    gasket,
	}
}

func (s *ReadService) GetFlange(ctx context.Context, req *moment_api.GetFlangeRequest) (*moment_api.GetFlangeResponse, error) {
	typeFlange, err := s.flange.GetTypeFlange(ctx, &moment_api.GetTypeFlangeRequest{})
	if err != nil {
		return nil, err
	}

	st, err := s.flange.GetStandarts(ctx, &moment_api.GetStandartsRequest{TypeId: typeFlange[0].Id})
	if err != nil {
		return nil, err
	}

	standarts := []*moment_api.StandartWithSize{}

	for _, s2 := range st {
		sizes, err := s.flange.GetBasisFlangeSize(ctx, &moment_api.GetBasisFlangeSizeRequest{
			IsUseRow: s2.IsNeedRow,
			StandId:  s2.Id,
			IsInch:   s2.IsInch,
		})
		if err != nil {
			return nil, err
		}

		standarts = append(standarts, &moment_api.StandartWithSize{
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

	gasket, err := s.gasket.GetGasketWithThick(ctx, &moment_api.GetGasketRequest{})
	if err != nil {
		return nil, err
	}

	env, err := s.gasket.GetEnv(ctx, &moment_api.GetEnvRequest{})
	if err != nil {
		return nil, err
	}

	materials, err := s.materials.GetMaterials(ctx, &moment_api.GetMaterialsRequest{})
	if err != nil {
		return nil, err
	}

	res := &moment_api.GetFlangeResponse{
		TypeFlange: typeFlange,
		Standarts:  standarts,
		Gaskets:    gasket,
		Env:        env,
		Materials:  materials,
	}

	return res, nil
}
