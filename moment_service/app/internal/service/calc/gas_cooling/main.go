package gas_cooling

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/gas_cooling/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/gas_cooling/formulas"
	"github.com/Alexander272/sealur/moment_service/internal/service/device"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
	"github.com/goccy/go-json"
)

type CoolingService struct {
	graphic  *graphic.GraphicService
	data     *data.DataService
	formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewCoolingService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService, device *device.DeviceService) *CoolingService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	// значение зависит от поля "Условия работы"
	kp := map[bool]float64{
		true:  constants.WorkKyp,
		false: constants.TestKyp,
	}

	// значение зависит от поля "Условие затяжки"
	kz := map[string]float64{
		"uncontrollable":  constants.UncontrollableKyz,
		"controllable":    constants.ControllableKyz,
		"controllablePin": constants.ControllablePinKyz,
	}

	data := data.NewDataService(flange, materials, gasket, device, graphic)
	formulas := formulas.NewFormulasService()

	return &CoolingService{
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
		graphic:  graphic,
		data:     data,
		formulas: formulas,
	}
}

func (s *CoolingService) CalculateGasCooling(ctx context.Context, data *calc_api.GasCoolingRequest) (*calc_api.GasCoolingResponse, error) {
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := &calc_api.GasCoolingResponse{
		Data:   s.data.FormatInitData(data.Data),
		Bolts:  d.Bolt,
		Gasket: d.Gasket,
		Calc:   &gas_cooling_model.Calculated{},
	}

	result.Calc.Auxiliary = s.CalcAuxiliary(ctx, d)
	result.Calc.ForcesInBolts = s.CalcForcesInBolts(ctx, d, result.Calc.Auxiliary)
	result.Calc.Bolt = s.CalcBolts(ctx, d, result.Calc.ForcesInBolts, result.Calc.Auxiliary)
	result.Calc.Moment = s.CalcMoment(ctx, data.Data.Friction, d, result.Calc.Bolt, result.Calc.ForcesInBolts, result.Calc.Auxiliary)

	if data.IsNeedFormulas {
		result.Formulas = s.formulas.GetFormulas(data, d, result)
	}

	_, err = json.Marshal(result.Calc)
	if err != nil {
		result.Calc = &gas_cooling_model.Calculated{
			Auxiliary:     &gas_cooling_model.CalcAuxiliary{},
			ForcesInBolts: &gas_cooling_model.CalcForcesInBolts{},
			Bolt:          &gas_cooling_model.CalcBolts{},
			Moment:        &gas_cooling_model.CalcMoment{},
		}
		logger.Error("failed to marshal json. error: " + err.Error())
	}

	return result, nil
}

func (s *CoolingService) CalcAuxiliary(ctx context.Context, data models.DataGasCooling) *gas_cooling_model.CalcAuxiliary {
	aux := &gas_cooling_model.CalcAuxiliary{}

	// расчетная ширина плоской прокладки
	aux.EstimatedGasketWidth = math.Min(data.Gasket.Width, 3.87*math.Sqrt(data.Gasket.Width))
	// Bp - расчетный размер решетки в поперечном направлении
	aux.SizeTrans = data.Gasket.SizeTrans - aux.EstimatedGasketWidth
	// Lp - Расчетный размер решетки в продольном направлении
	aux.SizeLong = data.Gasket.SizeLong - aux.EstimatedGasketWidth

	return aux
}

func (s *CoolingService) CalcForcesInBolts(ctx context.Context, data models.DataGasCooling, aux *gas_cooling_model.CalcAuxiliary,
) *gas_cooling_model.CalcForcesInBolts {
	forces := &gas_cooling_model.CalcForcesInBolts{}

	// Суммарная площадь сечения болтов/шпилек
	forces.Area = float64(data.Bolt.Count) * data.Bolt.Area
	// Fв - Расчетное усилие в болтах (шпильках) в условиях эксплуатации
	forces.WorkEffort = data.Data.Pressure * (aux.SizeLong*aux.SizeTrans + 2*aux.EstimatedGasketWidth*
		data.Gasket.M*(aux.SizeLong+aux.SizeTrans))

	tmp1 := (data.Data.TestPressure / data.Data.Pressure) * forces.Area
	tmp2 := data.Data.TestPressure * (2*aux.SizeLong*aux.SizeTrans +
		2*aux.EstimatedGasketWidth*data.Gasket.M*(aux.SizeLong+aux.SizeTrans))
	// F0 - Расчетное усилие в болтах (шпильках) в условиях испытаний или монтажа
	forces.Effort = math.Max(tmp1, tmp2)

	return forces
}

func (s *CoolingService) CalcBolts(
	ctx context.Context,
	data models.DataGasCooling,
	forces *gas_cooling_model.CalcForcesInBolts,
	aux *gas_cooling_model.CalcAuxiliary,
) *gas_cooling_model.CalcBolts {
	bolts := &gas_cooling_model.CalcBolts{}

	Kyp := s.Kyp[true]
	Kyz := s.Kyz[data.Data.Condition]
	Kyt := constants.LoadKyt

	// Расчетное напряжение в болтах/шпильках - при затяжке
	bolts.RatedStress = forces.Effort / forces.Area
	// Условия прочности болтов шпилек - при затяжке
	bolts.AllowableVoltage = 1.2 * Kyp * Kyz * Kyt * data.Bolt.SigmaAt20
	bolts.StrengthBolt = &gas_cooling_model.Condition{X: bolts.RatedStress, Y: bolts.AllowableVoltage}

	// Условие прочности прокладки
	bolts.StrengthGasket = &gas_cooling_model.Condition{
		X: math.Max(forces.WorkEffort, forces.Effort) / (2 * (aux.SizeLong + aux.SizeTrans) * data.Gasket.Width),
		Y: float64(data.Gasket.PermissiblePres),
	}

	return bolts
}

func (s *CoolingService) CalcMoment(
	ctx context.Context,
	Friction float64,
	data models.DataGasCooling,
	bolts *gas_cooling_model.CalcBolts,
	forces *gas_cooling_model.CalcForcesInBolts,
	aux *gas_cooling_model.CalcAuxiliary,
) *gas_cooling_model.CalcMoment {
	moment := &gas_cooling_model.CalcMoment{}

	ok := bolts.StrengthBolt.X <= bolts.StrengthBolt.Y && (data.TypeGasket != gas_cooling_model.GasketData_Soft ||
		(data.TypeGasket == gas_cooling_model.GasketData_Soft && bolts.StrengthGasket.X <= bolts.StrengthGasket.Y))

	if ok {
		if bolts.RatedStress > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter {
			moment.Mkp = s.graphic.CalculateMkp(data.Bolt.Diameter, bolts.RatedStress)
			moment.UseGraphic = true
		} else {
			moment.Mkp = (Friction * forces.Effort * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
			// moment.Mkp = (0.3 * forces.Effort * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		}
		moment.Friction = Friction

		if Friction == constants.DefaultFriction {
			moment.Mkp1 = 0.75 * moment.Mkp
		}

		Prek := 0.8 * forces.Area * data.Bolt.SigmaAt20
		moment.Qrek = Prek / (2 * (aux.SizeLong + aux.SizeTrans) * data.Gasket.Width)
		moment.Mrek = (Friction * Prek * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		// moment.Mrek = (0.3 * Prek * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000

		Pmax := bolts.AllowableVoltage * forces.Area
		moment.Qmax = Pmax / (2 * (aux.SizeLong + aux.SizeTrans) * data.Gasket.Width)

		if data.TypeGasket == gas_cooling_model.GasketData_Soft && moment.Qmax > data.Gasket.PermissiblePres {
			Pmax = data.Gasket.PermissiblePres * (2 * (aux.SizeLong + aux.SizeTrans) * data.Gasket.Width)
			moment.Qmax = data.Gasket.PermissiblePres
		}
		// if moment.Qmax > data.Gasket.PermissiblePres {
		// 	Pmax = data.Gasket.PermissiblePres * (2 * (aux.SizeLong + aux.SizeTrans) * data.Gasket.Width)
		// 	moment.Qmax = data.Gasket.PermissiblePres
		// }

		moment.Mmax = (Friction * Pmax * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		// moment.Mmax = (0.3 * Pmax * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000

		if moment.Mrek > moment.Mmax {
			moment.Mrek = moment.Mmax
		}
		if moment.Qrek > moment.Qmax {
			moment.Qrek = moment.Qmax
		}

	}

	return moment
}
