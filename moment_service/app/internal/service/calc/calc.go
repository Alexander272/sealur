package calc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/calc/cap"
	cap_old "github.com/Alexander272/sealur/moment_service/internal/service/calc/cap_old"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/dev_cooling"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_circle"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_rectangle"
	calc_flange "github.com/Alexander272/sealur/moment_service/internal/service/calc/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/float"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/gas_cooling"
	"github.com/Alexander272/sealur/moment_service/internal/service/device"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
)

type Flange interface {
	CalculationFlange(context.Context, *calc_api.FlangeRequest) (*calc_api.FlangeResponse, error)
}

type Cap interface {
	CalculationCap(context.Context, *calc_api.CapRequest) (*calc_api.CapResponse, error)
}

type CapOld interface {
	CalculationCapOld(context.Context, *calc_api.CapRequestOld) (*calc_api.CapResponseOld, error)
}

type Float interface {
	CalculationFloat(context.Context, *calc_api.FloatRequest) (*calc_api.FloatResponse, error)
}

type DevCooling interface {
	CalculateDevCooling(context.Context, *calc_api.DevCoolingRequest) (*calc_api.DevCoolingResponse, error)
}

type GasCooling interface {
	CalculateGasCooling(context.Context, *calc_api.GasCoolingRequest) (*calc_api.GasCoolingResponse, error)
}

type ExCircle interface {
	CalculateExCircle(context.Context, *calc_api.ExpressCircleRequest) (*calc_api.ExpressCircleResponse, error)
}

type ExRect interface {
	CalculateExRect(context.Context, *calc_api.ExpressRectangleRequest) (*calc_api.ExpressRectangleResponse, error)
}

type CalcService struct {
	Flange
	Cap
	CapOld
	Float
	DevCooling
	GasCooling
	ExCircle
	ExRect
}

func NewCalcServices(flange *flange.FlangeService, gasket *gasket.GasketService, materials *materials.MaterialsService, device *device.DeviceService) *CalcService {
	graphic := graphic.NewGraphicService()

	return &CalcService{
		Flange:     calc_flange.NewFlangeService(graphic, flange, gasket, materials),
		Cap:        cap.NewCapService(graphic, flange, gasket, materials),
		CapOld:     cap_old.NewCapService(graphic, flange, gasket, materials),
		Float:      float.NewFloatService(graphic, flange, gasket, materials),
		DevCooling: dev_cooling.NewCoolingService(graphic, flange, gasket, materials),
		GasCooling: gas_cooling.NewCoolingService(graphic, flange, gasket, materials, device),
		ExCircle:   express_circle.NewExCircleService(graphic, flange, gasket, materials),
		ExRect:     express_rectangle.NewExRectServiceService(graphic, flange, gasket, materials),
	}
}
