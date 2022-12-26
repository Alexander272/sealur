package cap

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/cap/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

type CapService struct {
	graphic *graphic.GraphicService
	data    *data.DataService
	// formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewCapService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService) *CapService {
	//значение зависит от поля "Тип соединения"
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

	data := data.NewDataService(flange, materials, gasket, graphic)
	// formulas := formulas.NewFormulasService()

	return &CapService{
		graphic: graphic,
		data:    data,
		// formulas: formulas,
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

// расчет момента затяжка фланец-крышка по ГОСТ 34233.4 - 2017
func (s *CapService) CalculationCap(ctx context.Context, data *calc_api.CapRequest) (*calc_api.CapResponse, error) {
	// получение данных (либо из бд, либо либо их передают с клиента) для расчетов (+ там пару формул записано)
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := calc_api.CapResponse{
		Data:     s.data.FormatInitData(data.Data),
		Bolt:     d.Bolt,
		Flange:   d.Flange,
		Cap:      d.Cap,
		Embed:    d.Embed,
		Gasket:   d.Gasket,
		Calc:     &cap_model.Calculated{},
		Formulas: &cap_model.Formulas{},
	}

	if data.IsUseWasher {
		result.Washers = append(result.Washers, d.Washer1, d.Washer2)
	}

	aux := &cap_model.CalcAuxiliary{}
	if data.Data.Calculation == cap_model.MainData_basis {
		// расчет основных величин
		result.Calc.Basis, aux = s.basisCalculate(d, data)
	} else {
		// прочностной расчет
		result.Calc.Strength = s.strengthCalculate(d, data)
	}

	if data.IsNeedFormulas {
		// получение формул с подставленными значениями переменных
		// result.Formulas = s.formulas.GetFormulas(data, d, &result, aux)
		//TODO
		logger.Debug(aux.A)
	}

	return &result, nil
}
