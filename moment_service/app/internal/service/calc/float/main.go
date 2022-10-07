package float

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/float/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
)

type FloatService struct {
	graphic *graphic.GraphicService
	data    *data.DataService
	// formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewFloatService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService) *FloatService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	kp := map[bool]float64{
		true:  constants.WorkKyp,
		false: constants.TestKyp,
	}

	kz := map[string]float64{
		"uncontrollable":  constants.UncontrollableKyz,
		"controllable":    constants.ControllableKyz,
		"controllablePin": constants.ControllablePinKyz,
	}

	data := data.NewDataService(flange, materials, gasket, graphic)
	// formulas := formulas.NewFormulasService()

	return &FloatService{
		graphic: graphic,
		data:    data,
		// formulas: formulas,
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

func (s *FloatService) CalculationFloat(ctx context.Context, data *calc_api.FloatRequest) (*calc_api.FloatResponse, error) {
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := &calc_api.FloatResponse{
		Data:   s.data.FormatInitData(data),
		Flange: d.Flange,
		Cap:    d.Cap,
		Bolt:   d.Bolt,
		Gasket: d.Gasket,
	}

	result.Calc.B0 = d.B0
	result.Calc.Dcp = d.Dcp

	// var yp float64 = 0
	// if d.TypeGasket == "Soft" {
	// 	yp = (d.Gasket.Thickness * d.Gasket.Compression) / (d.Gasket.Epsilon * math.Pi * d.Dcp * d.Gasket.Width)
	// }

	// Lb := result.Bolt.Lenght + s.typeBolt[data.Type.String()]*d.Bolt.Diameter
	// Ab := float64(d.Bolt.Count) * d.Bolt.Area

	return result, nil
}
