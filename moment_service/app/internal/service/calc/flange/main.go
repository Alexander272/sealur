package flange

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/flange/data"

	"github.com/Alexander272/sealur/moment_service/internal/service/calc/flange/formulas"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

type FlangeService struct {
	graphic  *graphic.GraphicService
	data     *data.DataService
	formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
	Kyt      map[bool]float64
}

func NewFlangeService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService) *FlangeService {
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

	kt := map[bool]float64{
		true:  constants.LoadKyt,
		false: constants.NoLoadKyt,
	}

	data := data.NewDataService(flange, materials, gasket, graphic)
	formulas := formulas.NewFormulasService()

	return &FlangeService{
		graphic:  graphic,
		data:     data,
		formulas: formulas,
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
		Kyt:      kt,
	}
}

// расчет момента затяжки фланец-фланец по ГОСТ 34233.4 - 2017
func (s *FlangeService) CalculationFlange(ctx context.Context, data *calc_api.FlangeRequest) (*calc_api.FlangeResponse, error) {
	// получение данных (либо из бд, либо либо их передают с клиента) для расчетов (+ там пару формул записано)
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := calc_api.FlangeResponse{
		Data:    s.data.FormatInitData(data),
		Bolt:    d.Bolt,
		Flanges: []*flange_model.FlangeResult{d.Flange1},
		Embed:   d.Embed,
		Gasket:  d.Gasket,
		Calc:    &flange_model.Calculated{},
	}

	if data.IsUseWasher {
		result.Washers = append(result.Washers, d.Washer1)
		if !data.IsSameFlange {
			result.Washers = append(result.Washers, d.Washer2)
		}
	}

	if !data.IsSameFlange {
		result.Flanges = append(result.Flanges, d.Flange2)
	}

	aux := &flange_model.CalcAuxiliary{}
	if data.Calculation == calc_api.FlangeRequest_basis {
		// расчет основных величин
		result.Calc.Basis, aux = s.basisCalculate(d, data)
	} else {
		// прочностной расчет
		result.Calc.Strength = s.strengthCalculate(d, data)
	}

	if data.IsNeedFormulas {
		// получение формул с подставленными значениями переменных
		result.Formulas = s.formulas.GetFormulas(data, d, &result, aux)
	}

	return &result, nil
}
