package calc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/calc/cap"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/dev_cooling"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_circle"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_rectangle"
	calc_flange "github.com/Alexander272/sealur/moment_service/internal/service/calc/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/float"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
)

type Flange interface {
	CalculationFlange(ctx context.Context, data *calc_api.FlangeRequest) (*calc_api.FlangeResponse, error)
}

type Cap interface {
	CalculationCap(ctx context.Context, data *calc_api.CapRequest) (*calc_api.CapResponse, error)
}

type Float interface {
	CalculationFloat(ctx context.Context, data *calc_api.FloatRequest) (*calc_api.FloatResponse, error)
}

type DevCooling interface {
	CalculateDevCooling(ctx context.Context, data *calc_api.DevCoolingRequest) (*calc_api.DevCoolingResponse, error)
}

type ExCircle interface {
	CalculateExCircle(ctx context.Context, data *calc_api.ExpressCircleRequest) (*calc_api.ExpressCircleResponse, error)
}

type ExRect interface {
	CalculateExRect(ctx context.Context, data *calc_api.ExpressRectangleRequest) (*calc_api.ExpressRectangleResponse, error)
}

type CalcService struct {
	Flange
	Cap
	Float
	DevCooling
	ExCircle
	ExRect
}

func NewCalcServices(flange *flange.FlangeService, gasket *gasket.GasketService, materials *materials.MaterialsService) *CalcService {
	graphic := graphic.NewGraphicService()

	return &CalcService{
		Flange:     calc_flange.NewFlangeService(graphic, flange, gasket, materials),
		Cap:        cap.NewCapService(graphic, flange, gasket, materials),
		Float:      float.NewFloatService(graphic, flange, gasket, materials),
		DevCooling: dev_cooling.NewCoolingService(graphic, flange, gasket, materials),
		ExCircle:   express_circle.NewExCircleService(graphic, flange, gasket, materials),
		ExRect:     express_rectangle.NewExRectServiceService(graphic, flange, gasket, materials),
	}
}
