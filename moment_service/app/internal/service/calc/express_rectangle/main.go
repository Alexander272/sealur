package express_rectangle

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_rectangle/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/express_rectangle/formulas"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_rectangle_model"
)

type ExRectService struct {
	graphic  *graphic.GraphicService
	data     *data.DataService
	formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewExRectServiceService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService) *ExRectService {
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

	return &ExRectService{
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
		graphic:  graphic,
		data:     data,
		formulas: formulas,
	}
}

// экспресс оценка момента затяжки
func (s *ExRectService) CalculateExRect(ctx context.Context, data *calc_api.ExpressRectangleRequest) (*calc_api.ExpressRectangleResponse, error) {
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := calc_api.ExpressRectangleResponse{
		Data:     s.data.FormatInitData(data),
		Bolts:    d.Bolt,
		Gasket:   d.Gasket,
		Calc:     &express_rectangle_model.Calculated{},
		Formulas: &express_rectangle_model.Formulas{},
	}

	Auxiliary := &express_rectangle_model.CalcAuxiliary{}
	ForcesInBolts := &express_rectangle_model.CalcForcesInBolts{}
	Bolts := &express_rectangle_model.CalcBolts{}
	Moment := &express_rectangle_model.CalcMoment{}

	if result.Data.TestPressure == 0 {
		result.Data.TestPressure = 1.5 * result.Data.Pressure
	}

	// расчетная ширина плоской прокладки
	Auxiliary.EstimatedGasketWidth = math.Min(d.Gasket.Width, 3.87*math.Sqrt(d.Gasket.Width))
	// Bp - расчетный размер решетки в поперечном направлении
	Auxiliary.SizeTrans = d.Gasket.SizeTrans - Auxiliary.EstimatedGasketWidth
	// Lp - Расчетный размер решетки в продольном направлении
	Auxiliary.SizeLong = d.Gasket.SizeLong - Auxiliary.EstimatedGasketWidth

	// Суммарная площадь сечения болтов/шпилек
	ForcesInBolts.Area = float64(d.Bolt.Count) * d.Bolt.Area
	// Fв - Расчетное усилие в болтах (шпильках) в условиях эксплуатации
	ForcesInBolts.WorkEffort = data.Pressure * (Auxiliary.SizeLong*Auxiliary.SizeTrans + 2*Auxiliary.EstimatedGasketWidth*
		d.Gasket.M*(Auxiliary.SizeLong+Auxiliary.SizeTrans))

	tmp1 := (result.Data.TestPressure / data.Pressure) * ForcesInBolts.Area
	tmp2 := result.Data.TestPressure * (2*Auxiliary.SizeLong*Auxiliary.SizeTrans +
		2*Auxiliary.EstimatedGasketWidth*d.Gasket.M*(Auxiliary.SizeLong+Auxiliary.SizeTrans))
	// F0 - Расчетное усилие в болтах (шпильках) в условиях испытаний или монтажа
	ForcesInBolts.Effort = math.Max(tmp1, tmp2)

	Kyp := s.Kyp[true]
	Kyz := s.Kyz[data.Condition.String()]
	Kyt := constants.LoadKyt

	// Расчетное напряжение в болтах/шпильках - при затяжке
	Bolts.RatedStress = ForcesInBolts.Effort / ForcesInBolts.Area
	// Условия прочности болтов шпилек - при затяжке
	Bolts.AllowableVoltage = 1.2 * Kyp * Kyz * Kyt * d.Bolt.SigmaAt20
	Bolts.StrengthBolt = &express_rectangle_model.Condition{X: Bolts.RatedStress, Y: Bolts.AllowableVoltage}

	// Условие прочности прокладки
	Bolts.StrengthGasket = &express_rectangle_model.Condition{
		X: math.Max(ForcesInBolts.WorkEffort, ForcesInBolts.Effort) / (2 * (Auxiliary.SizeLong + Auxiliary.SizeTrans) * d.Gasket.Width),
		Y: float64(d.Gasket.PermissiblePres),
	}

	ok := Bolts.StrengthBolt.X <= Bolts.StrengthBolt.Y && (d.TypeGasket != express_rectangle_model.GasketData_Soft ||
		(d.TypeGasket == express_rectangle_model.GasketData_Soft && Bolts.StrengthGasket.X <= Bolts.StrengthGasket.Y))

	if ok {
		if Bolts.RatedStress > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter {
			Moment.Mkp = s.graphic.CalculateMkp(d.Bolt.Diameter, Bolts.RatedStress)
		} else {
			Moment.Mkp = (0.3 * ForcesInBolts.Effort * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		}
		Moment.Mkp1 = 0.75 * Moment.Mkp

		Prek := 0.8 * ForcesInBolts.Area * d.Bolt.SigmaAt20
		Moment.Qrek = Prek / (2 * (Auxiliary.SizeLong + Auxiliary.SizeTrans) * d.Gasket.Width)
		Moment.Mrek = (0.3 * Prek * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000

		Pmax := Bolts.AllowableVoltage * ForcesInBolts.Area
		Moment.Qmax = Pmax / (2 * (Auxiliary.SizeLong + Auxiliary.SizeTrans) * d.Gasket.Width)

		if d.TypeGasket == express_rectangle_model.GasketData_Soft && Moment.Qmax > d.Gasket.PermissiblePres {
			Pmax = d.Gasket.PermissiblePres * (2 * (Auxiliary.SizeLong + Auxiliary.SizeTrans) * d.Gasket.Width)
			Moment.Qmax = d.Gasket.PermissiblePres
		}

		Moment.Mmax = (0.3 * Pmax * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
	}

	result.Calc.Auxiliary = Auxiliary
	result.Calc.ForcesInBolts = ForcesInBolts
	result.Calc.Bolt = Bolts
	result.Calc.Moment = Moment

	if data.IsNeedFormulas {
		result.Formulas = s.formulas.GetFormulas(*data, d, result)
	}

	return &result, nil
}
