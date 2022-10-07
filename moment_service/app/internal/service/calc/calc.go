package calc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/calc/cap"
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

type CalcService struct {
	Flange
	Cap
	Float
}

func NewCalcServices(flange *flange.FlangeService, gasket *gasket.GasketService, materials *materials.MaterialsService) *CalcService {
	graphic := graphic.NewGraphicService()

	return &CalcService{
		Flange: calc_flange.NewFlangeService(graphic, flange, gasket, materials),
		Cap:    cap.NewCapService(graphic, flange, gasket, materials),
		Float:  float.NewFloatService(graphic, flange, gasket, materials),
	}
}
