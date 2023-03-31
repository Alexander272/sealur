package read

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/device"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/material_model"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
)

type GasCoolingService struct {
	materials *materials.MaterialsService
	gasket    *gasket.GasketService
	device    *device.DeviceService
}

func NewGasCoolingService(materials *materials.MaterialsService, gasket *gasket.GasketService, device *device.DeviceService) *GasCoolingService {
	return &GasCoolingService{
		materials: materials,
		gasket:    gasket,
		device:    device,
	}
}

func (s *GasCoolingService) GetAVO(ctx context.Context, req *read_api.GetAVORequest) (*read_api.GetAVOResponse, error) {
	gasket, err := s.gasket.GetGasketWithThick(ctx, &gasket_api.GetGasketRequest{TypeGasket: []gasket_api.TypeGasket{gasket_api.TypeGasket_Soft}})
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

	devices, err := s.device.GetDevices(ctx, &device_api.GetDeviceRequest{})
	if err != nil {
		return nil, err
	}

	pressure, err := s.device.GetPressure(ctx, &device_api.GetPressureRequest{})
	if err != nil {
		return nil, err
	}

	tubeCount, err := s.device.GetTubeCount(ctx, &device_api.GetTubeCountRequest{})
	if err != nil {
		return nil, err
	}

	factor, err := s.device.GetFinningFactor(ctx, &device_api.GetFinningFactorRequest{})
	if err != nil {
		return nil, err
	}

	section, err := s.device.GetSectionExecution(ctx, &device_api.GetSectionExecutionRequest{})
	if err != nil {
		return nil, err
	}

	tubeLength, err := s.device.GetTubeLength(ctx, &device_api.GetTubeLengthRequest{})
	if err != nil {
		return nil, err
	}

	numbers, err := s.device.GetNumberOfMoves(ctx, &device_api.GetNumberOfMovesRequest{})
	if err != nil {
		return nil, err
	}

	res := &read_api.GetAVOResponse{
		Gaskets:          gasket,
		BoltMaterials:    boltMaterials,
		FlangeMaterials:  flangeMaterials,
		Env:              env,
		Devices:          devices,
		Pressure:         pressure,
		TubeCount:        tubeCount,
		FinningFactor:    factor,
		SectionExecution: section,
		TubeLength:       tubeLength,
		NumberOfMoves:    numbers,
	}

	return res, nil
}
