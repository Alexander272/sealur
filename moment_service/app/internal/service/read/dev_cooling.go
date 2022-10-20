package read

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
)

type DevCoolingService struct {
	materials *materials.MaterialsService
	gasket    *gasket.GasketService
}

func NewDevCoolingService(materials *materials.MaterialsService, gasket *gasket.GasketService) *DevCoolingService {
	return &DevCoolingService{
		materials: materials,
		gasket:    gasket,
	}
}

func (s *DevCoolingService) GetDevCooling(ctx context.Context, req *read_api.GetDevCoolingtRequest) (*read_api.GetDevCoolingResponse, error) {
	gasket, err := s.gasket.GetGasketWithThick(ctx, &gasket_api.GetGasketRequest{TypeGasket: []gasket_api.TypeGasket{gasket_api.TypeGasket_Soft}})
	if err != nil {
		return nil, err
	}

	env, err := s.gasket.GetEnv(ctx, &gasket_api.GetEnvRequest{})
	if err != nil {
		return nil, err
	}

	materials, err := s.materials.GetMaterials(ctx, &material_api.GetMaterialsRequest{})
	if err != nil {
		return nil, err
	}

	res := &read_api.GetDevCoolingResponse{
		Gaskets:   gasket,
		Materials: materials,
		Env:       env,
	}

	return res, nil
}
