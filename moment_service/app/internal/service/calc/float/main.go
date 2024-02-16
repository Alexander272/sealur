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
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
	"github.com/goccy/go-json"
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

// расчет плавающей головки теплообменного аппарата (хз по какому госту)
func (s *FloatService) CalculationFloat(ctx context.Context, data *calc_api.FloatRequest) (*calc_api.FloatResponse, error) {
	// получение данных (либо из бд, либо либо их передают с клиента) для расчетов (+ там пара формул записана)
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
		Calc:   &float_model.Calculated{},
	}

	// Эффективная ширина прокладки
	result.Calc.B0 = d.B0
	// Расчетный диаметр прокладки
	result.Calc.Dcp = d.Dcp

	var yp, alpha float64 = 0, 1
	if d.TypeGasket == "Soft" {
		// Податливость прокладки
		yp = (d.Gasket.Thickness * d.Gasket.Compression) / (d.Gasket.Epsilon * math.Pi * d.Dcp * d.Gasket.Width)
	}

	result.Calc.Lb = result.Bolt.Lenght + s.typeBolt[data.Type.String()]*d.Bolt.Diameter
	// Податливость болтов/шпилек
	result.Calc.Yb = result.Calc.Lb / (d.Bolt.EpsilonAt20 * d.Bolt.Area * float64(d.Bolt.Count))
	result.Calc.Yp = yp
	// Суммарная площадь сечения болтов/шпилек
	result.Calc.A = float64(d.Bolt.Count) * d.Bolt.Area

	if d.TypeGasket == "Soft" {
		alpha = 1 - (yp-(d.Cap.Y-d.Flange.B)*d.Flange.B)/(yp+result.Calc.Yb)
	}
	// Коэффициент жесткости
	result.Calc.Alpha = alpha

	// Усилие необходимое для смятия прокладки при затяжке
	result.Calc.Po = 0.5 * math.Pi * d.Dcp * d.B0 * d.Gasket.Pres

	if data.Pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях необходимое для обеспечения герметичности фланцевого соединения
		result.Calc.Rp = math.Pi * d.Dcp * d.B0 * d.Gasket.M * math.Abs(data.Pressure)
	}

	// Равнодействующая нагрузка от давления
	result.Calc.Qd = 0.785 * math.Pow(d.Dcp, 2) * float64(data.Pressure)

	// Минимальное начальное натяжение болтов (шпилек)
	result.Calc.MinB = 0.4 * result.Calc.A * d.Bolt.SigmaAt20
	// Расчетная нагрузка на болты/шпильки при затяжке необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	result.Calc.Pb2 = math.Max(result.Calc.Po, result.Calc.MinB)
	// Расчетная нагрузка на болты/шпильки при затяжке необходимая для обеспечения в рабочих условиях давления на прокладку
	// достаточного для герметизации фланцевого соединения
	result.Calc.Pb1 = alpha*result.Calc.Qd + result.Calc.Rp

	// Расчетная нагрузка на болты/шпильки фланцевых соединений
	result.Calc.Pb = math.Max(result.Calc.Pb1, result.Calc.Pb2)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений в рабочих условиях
	result.Calc.Pbr = result.Calc.Pb + (1-alpha)*result.Calc.Qd

	// Расчетное напряжение в болтах/шпильках - при затяжке
	result.Calc.SigmaB1 = result.Calc.Pb / result.Calc.A
	// Расчетное напряжение в болтах/шпильках - в рабочих условиях
	result.Calc.SigmaB2 = result.Calc.Pbr / result.Calc.A

	Kyp := s.Kyp[data.IsWork]
	Kyz := s.Kyz[data.Condition.String()]
	Kyt := constants.NoLoadKyt

	// формула Г.3
	// Допускаемое напряжение для болтов шпилек - при затяжке
	result.Calc.DSigmaM = 1.2 * Kyp * Kyz * Kyt * d.Bolt.SigmaAt20
	// формула Г.4
	// Допускаемое напряжение для болтов шпилек - в рабочих условиях
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
			// Крутящий момент при затяжке болтов/шпилек
			result.Calc.Mkp = s.graphic.CalculateMkp(d.Bolt.Diameter, result.Calc.SigmaB1)
			result.Calc.UseGraphic = true
		} else {
			result.Calc.Mkp = (data.Friction * result.Calc.Pb * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
			// result.Calc.Mkp = (0.3 * result.Calc.Pb * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		}
		result.Calc.Friction = data.Friction

		// Крутящий момент при затяжке болтов/шпилек со смазкой снижается на 25%
		if data.Friction == constants.DefaultFriction {
			result.Calc.Mkp1 = 0.75 * result.Calc.Mkp
		}

		Prek := 0.8 * result.Calc.A * d.Bolt.SigmaAt20
		// Напряжение на прокладке
		result.Calc.Qrek = Prek / (math.Pi * d.Dcp * d.Gasket.Width)
		// Момент затяжки при применении уплотнения на старых (изношенных) фланцах, имеющих перекосы
		result.Calc.Mrek = (data.Friction * Prek * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		// result.Calc.Mrek = (0.3 * Prek * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000

		Pmax := result.Calc.DSigmaM * result.Calc.A
		// Максимальное напряжение на прокладке
		result.Calc.Qmax = Pmax / (math.Pi * d.Dcp * d.Gasket.Width)

		if d.TypeGasket == "Soft" && result.Calc.Qmax > d.Gasket.PermissiblePres {
			Pmax = float64(d.Gasket.PermissiblePres) * (math.Pi * d.Dcp * d.Gasket.Width)
			result.Calc.Qmax = d.Gasket.PermissiblePres
		}
		// if result.Calc.Qmax > d.Gasket.PermissiblePres {
		// 	Pmax = float64(d.Gasket.PermissiblePres) * (math.Pi * d.Dcp * d.Gasket.Width)
		// 	result.Calc.Qmax = d.Gasket.PermissiblePres
		// }

		if result.Calc.Mrek > result.Calc.Mmax {
			result.Calc.Mrek = result.Calc.Mmax
		}
		if result.Calc.Qrek > result.Calc.Qmax {
			result.Calc.Qrek = result.Calc.Qmax
		}

		// Максимальный крутящий момент при затяжке болтов/шпилек
		result.Calc.Mmax = (data.Friction * Pmax * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
		// result.Calc.Mmax = (0.3 * Pmax * d.Bolt.Diameter / float64(d.Bolt.Count)) / 1000
	}

	if data.IsNeedFormulas {
		result.Formulas = s.formulas.GetFormulas(
			data,
			data.Condition.String(), data.Type.String(),
			data.IsWork,
			d,
			result,
		)
	}

	_, err = json.Marshal(result.Calc)
	if err != nil {
		result.Calc = &float_model.Calculated{}
		logger.Error("failed to marshal json. error: " + err.Error())
	}

	return result, nil
}
