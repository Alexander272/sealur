package express_circle

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_circle/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_circle/formulas"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
)

type ExCircleService struct {
	graphic  *graphic.GraphicService
	data     *data.DataService
	formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewExCircleService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService) *ExCircleService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	// занчение зависит от поля "Условия работы"
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

	data := data.NewDataService(flange, materials, gasket, graphic)
	formulas := formulas.NewFormulasService()

	return &ExCircleService{
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
		graphic:  graphic,
		data:     data,
		formulas: formulas,
	}
}

// экспресс оценка момента затяжки
func (s *ExCircleService) CalculateExCircle(ctx context.Context, data *calc_api.ExpressCircleRequest) (*calc_api.ExpressCircleResponse, error) {
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := calc_api.ExpressCircleResponse{
		Data:   s.data.FormatInitData(data),
		Bolts:  d.Bolt,
		Gasket: d.Gasket,
		Calc:   &express_circle_model.Calculated{},
	}

	Deformation := &express_circle_model.CalcDeformation{}
	ForcesInBolts := &express_circle_model.CalcForcesInBolts{}
	Bolts := &express_circle_model.CalcBolts{}
	Moment := &express_circle_model.CalcMoment{}

	if d.TypeGasket == express_circle_model.GasketData_Oval {
		// Эффективная ширина прокладки
		Deformation.Width = d.Gasket.Width / 4
		// Расчетный диаметр прокладки
		Deformation.Diameter = data.Gasket.DOut - d.Gasket.Width/2

	} else {
		if d.Gasket.Width <= constants.Bp {
			// Эффективная ширина прокладки
			Deformation.Width = d.Gasket.Width
		} else {
			Deformation.Width = constants.B0 * math.Sqrt(d.Gasket.Width)
		}
		// Расчетный диаметр прокладки
		Deformation.Diameter = data.Gasket.DOut - Deformation.Width
	}

	// Суммарная площадь сечения болтов/шпилек
	ForcesInBolts.Area = float64(d.Bolt.Count) * d.Bolt.Area

	// Усилие необходимое для смятия прокладки при затяжке
	Deformation.Deformation = 0.5 * math.Pi * Deformation.Diameter * Deformation.Width * d.Gasket.Pres

	if data.Pressure >= 0 {
		// Усилие на прокладке в рабочих условиях
		Deformation.Effort = math.Pi * Deformation.Diameter * Deformation.Width * d.Gasket.M * math.Abs(data.Pressure)
	}

	// Равнодействующая нагрузка от давления
	ForcesInBolts.ResultantLoad = 0.785 * math.Pow(Deformation.Diameter, 2) * data.Pressure

	// Минимальное начальное натяжение болтов (шпилек)
	ForcesInBolts.Tension = 0.4 * ForcesInBolts.Area * d.Bolt.SigmaAt20
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	ForcesInBolts.EstimatedLoad2 = math.Max(Deformation.Deformation, ForcesInBolts.Tension)
	// Расчетная нагрузка на болты/шпильки при затяжке необходимая для обеспечения в рабочих условиях давления на прокладку достаточного
	// для герметизации фланцевого соединения
	ForcesInBolts.EstimatedLoad1 = ForcesInBolts.ResultantLoad + Deformation.Effort
	// Расчетная нагрузка на болты/шпильки фланцевых соединений
	ForcesInBolts.DesignLoad = math.Max(ForcesInBolts.EstimatedLoad1, ForcesInBolts.EstimatedLoad2)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений в рабочих условиях
	ForcesInBolts.WorkDesignLoad = ForcesInBolts.DesignLoad

	Kyp := s.Kyp[true]
	Kyz := s.Kyz[data.Condition.String()]
	Kyt := constants.LoadKyt

	// Расчетное напряжение в болтах/шпильках - при затяжке
	Bolts.RatedStress = ForcesInBolts.DesignLoad / ForcesInBolts.Area
	// Условия прочности болтов шпилек - при затяжке
	Bolts.AllowableVoltage = 1.2 * Kyp * Kyz * Kyt * d.Bolt.SigmaAt20
	Bolts.StrengthBolt = &express_circle_model.Condition{X: Bolts.RatedStress, Y: Bolts.AllowableVoltage}

	if d.TypeGasket == express_circle_model.GasketData_Soft {
		// Условие прочности прокладки
		Bolts.StrengthGasket = &express_circle_model.Condition{
			X: math.Max(ForcesInBolts.DesignLoad, ForcesInBolts.WorkDesignLoad) / (math.Pi * Deformation.Diameter * d.Gasket.Width),
			Y: float64(d.Gasket.PermissiblePres),
		}
	}

	ok := Bolts.StrengthBolt.X <= Bolts.StrengthBolt.Y && (d.TypeGasket != express_circle_model.GasketData_Soft ||
		(d.TypeGasket == express_circle_model.GasketData_Soft && Bolts.StrengthGasket.X <= Bolts.StrengthGasket.Y))

	if ok {
		if Bolts.RatedStress > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter {
			Moment.Mkp = s.graphic.CalculateMkp(d.Bolt.Diameter, Bolts.RatedStress)
			Moment.UseGraphic = true
		} else {
			Moment.Mkp = (data.Friction * ForcesInBolts.DesignLoad * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
			// Moment.Mkp = (0.3 * ForcesInBolts.DesignLoad * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		}
		Moment.Friction = data.Friction
		Moment.Mkp1 = 0.75 * Moment.Mkp

		Prek := 0.8 * ForcesInBolts.Area * d.Bolt.SigmaAt20
		Moment.Qrek = Prek / (math.Pi * Deformation.Diameter * d.Gasket.Width)
		Moment.Mrek = (data.Friction * Prek * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		// Moment.Mrek = (0.3 * Prek * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000

		Pmax := Bolts.AllowableVoltage * ForcesInBolts.Area
		Moment.Qmax = Pmax / (math.Pi * Deformation.Diameter * d.Gasket.Width)

		// if d.TypeGasket == express_circle_model.GasketData_Soft && Moment.Qmax > d.Gasket.PermissiblePres {
		// 	Pmax = d.Gasket.PermissiblePres * (math.Pi * Deformation.Diameter * d.Gasket.Width)
		// 	Moment.Qmax = d.Gasket.PermissiblePres
		// }
		if Moment.Qmax > d.Gasket.PermissiblePres {
			Pmax = d.Gasket.PermissiblePres * (math.Pi * Deformation.Diameter * d.Gasket.Width)
			Moment.Qmax = d.Gasket.PermissiblePres
		}

		Moment.Mmax = (data.Friction * Pmax * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		// Moment.Mmax = (0.3 * Pmax * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000

		if Moment.Mrek > Moment.Mmax {
			Moment.Mrek = Moment.Mmax
		}
		if Moment.Qrek > Moment.Qmax {
			Moment.Qrek = Moment.Qmax
		}
	}

	result.Calc.Deformation = Deformation
	result.Calc.ForcesInBolts = ForcesInBolts
	result.Calc.Bolt = Bolts
	result.Calc.Moment = Moment

	if data.IsNeedFormulas {
		result.Formulas = s.formulas.GetFormulas(*data, d, result)
	}

	return &result, nil
}
