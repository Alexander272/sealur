package read

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/device"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
)

type Flange interface {
	GetFlange(context.Context, *read_api.GetFlangeRequest) (*read_api.GetFlangeResponse, error)
}

type Float interface {
	GetFloat(context.Context, *read_api.GetFloatRequest) (*read_api.GetFloatResponse, error)
}

type DevCooling interface {
	GetDevCooling(context.Context, *read_api.GetDevCoolingtRequest) (*read_api.GetDevCoolingResponse, error)
}

type GasCooling interface {
	GetAVO(context.Context, *read_api.GetAVORequest) (*read_api.GetAVOResponse, error)
}

type ReadService struct {
	Flange
	Float
	DevCooling
	GasCooling
}

func NewReadService(
	flange *flange.FlangeService,
	materials *materials.MaterialsService,
	gasket *gasket.GasketService,
	device *device.DeviceService,
) *ReadService {
	return &ReadService{
		Flange:     NewFlangeService(flange, materials, gasket),
		Float:      NewFloatService(materials, gasket),
		DevCooling: NewDevCoolingService(materials, gasket),
		GasCooling: NewGasCoolingService(materials, gasket, device),
	}
}
