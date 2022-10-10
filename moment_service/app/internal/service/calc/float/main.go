package float

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/float/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/float/formulas"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

type FloatService struct {
	graphic  *graphic.GraphicService
	data     *data.DataService
	formulas *formulas.FormulasService
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
	formulas := formulas.NewFormulasService()

	return &FloatService{
		graphic:  graphic,
		data:     data,
		formulas: formulas,
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

// расчет плавающей головки теплообменного аппарата
func (s *FloatService) CalculationFloat(ctx context.Context, data *calc_api.FloatRequest) (*calc_api.FloatResponse, error) {
	// получение данных (либо из бд, либо либо их передают с клиента) для расчетов (+ там пару формул записано)
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := calc_api.FloatResponse{
		Data:     s.data.FormatInitData(data),
		Flange:   d.Flange,
		Cap:      d.Cap,
		Bolt:     d.Bolt,
		Gasket:   d.Gasket,
		Calc:     &float_model.Calculated{},
		Formulas: &float_model.Formulas{},
	}

	result.Calc.B0 = d.B0
	result.Calc.Dcp = d.Dcp

	var yp, alpha float64 = 0, 1
	if d.TypeGasket == "Soft" {
		yp = (d.Gasket.Thickness * d.Gasket.Compression) / (d.Gasket.Epsilon * math.Pi * d.Dcp * d.Gasket.Width)
	}

	result.Calc.Lb = result.Bolt.Lenght + s.typeBolt[data.Type.String()]*d.Bolt.Diameter
	result.Calc.Yb = result.Calc.Lb / (d.Bolt.EpsilonAt20 * d.Bolt.Area * float64(d.Bolt.Count))
	result.Calc.Yp = yp
	result.Calc.A = float64(d.Bolt.Count) * d.Bolt.Area

	if d.TypeGasket == "Soft" {
		alpha = 1 - (yp-(d.Cap.Y-d.Flange.B)*d.Flange.B)/(yp+result.Calc.Yb)
	}
	result.Calc.Alpha = alpha

	result.Calc.Po = 0.5 * math.Pi * d.Dcp * d.B0 * d.Gasket.Pres

	if data.Pressure >= 0 {
		// формула 7
		result.Calc.Rp = math.Pi * d.Dcp * d.B0 * d.Gasket.M * math.Abs(data.Pressure)
	}

	result.Calc.Qd = 0.785 * math.Pow(d.Dcp, 2) * float64(data.Pressure)

	result.Calc.MinB = 0.4 * result.Calc.A * d.Bolt.SigmaAt20
	result.Calc.Pb2 = math.Max(result.Calc.Po, result.Calc.MinB)
	result.Calc.Pb1 = alpha*result.Calc.Qd + result.Calc.Rp

	result.Calc.Pb = math.Max(result.Calc.Pb1, result.Calc.Pb2)
	result.Calc.Pbr = result.Calc.Pb + (1-alpha)*result.Calc.Qd

	result.Calc.SigmaB1 = result.Calc.Pb / result.Calc.A
	result.Calc.SigmaB2 = result.Calc.Pbr / result.Calc.A

	Kyp := s.Kyp[data.IsWork]
	Kyz := s.Kyz[data.Condition.String()]
	Kyt := constants.NoLoadKyt
	// формула Г.3
	result.Calc.DSigmaM = 1.2 * Kyp * Kyz * Kyt * d.Bolt.SigmaAt20
	// формула Г.4
	result.Calc.DSigmaR = Kyp * Kyz * Kyt * d.Bolt.Sigma

	if d.TypeGasket == "Soft" {
		result.Calc.Qmax = math.Max(result.Calc.Pb, result.Calc.Pbr) / (math.Pi * d.Dcp * d.Gasket.Width)
	}

	if result.Calc.SigmaB1 <= result.Calc.DSigmaM {
		result.Calc.VSigmaB1 = true
	}
	if result.Calc.SigmaB2 <= result.Calc.DSigmaR {
		result.Calc.VSigmaB2 = true
	}

	if (result.Calc.VSigmaB1 && result.Calc.VSigmaB2 && d.TypeGasket != "Soft") ||
		(result.Calc.VSigmaB1 && result.Calc.VSigmaB2 && result.Calc.Qmax <= float64(d.Gasket.PermissiblePres) && d.TypeGasket == "Soft") {
		if result.Calc.SigmaB1 > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter {
			result.Calc.Mkp = s.graphic.CalculateMkp(d.Bolt.Diameter, result.Calc.SigmaB1)
		} else {
			result.Calc.Mkp = (0.3 * result.Calc.Pb * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		}

		result.Calc.Mkp1 = 0.75 * result.Calc.Mkp

		Prek := 0.8 * result.Calc.A * d.Bolt.SigmaAt20
		result.Calc.Qrek = Prek / (math.Pi * d.Dcp * d.Gasket.Width)
		result.Calc.Mrek = (0.3 * Prek * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000

		Pmax := result.Calc.DSigmaM * result.Calc.A
		result.Calc.Qmax = Pmax / (math.Pi * d.Dcp * d.Gasket.Width)

		if d.TypeGasket == "Soft" && result.Calc.Qmax > d.Gasket.PermissiblePres {
			Pmax = float64(d.Gasket.PermissiblePres) * (math.Pi * d.Dcp * d.Gasket.Width)
			result.Calc.Qmax = d.Gasket.PermissiblePres
		}

		result.Calc.Mmax = (0.3 * Pmax * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
	}

	if data.IsNeedFormulas {
		result.Formulas = s.formulas.GetFormulas(
			data.Condition.String(),
			data.IsWork,
			d,
			result,
		)
	}

	return &result, nil
}
