package cap

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/cap/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/cap/formulas"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
	"github.com/goccy/go-json"
)

type CapService struct {
	graphic  *graphic.GraphicService
	data     *data.DataService
	formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
	Kyt      map[bool]float64
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

	kt := map[bool]float64{
		true:  constants.LoadKyt,
		false: constants.NoLoadKyt,
	}

	data := data.NewDataService(flange, materials, gasket, graphic)
	formulas := formulas.NewFormulasService()

	return &CapService{
		graphic:  graphic,
		data:     data,
		formulas: formulas,
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
		Kyt:      kt,
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
		Data:   s.data.FormatInitData(data.Data),
		Bolt:   d.Bolt,
		Flange: d.Flange,
		Cap:    d.Cap,
		Embed:  d.Embed,
		Gasket: d.Gasket,
		Calc:   &cap_model.Calculated{},
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
		result.Formulas = s.formulas.GetFormulas(data, d, &result, aux)
	}

	_, err = json.Marshal(result.Calc)
	if err != nil {
		//? если не можем преобразовать в json обнуляем все значения
		if data.Data.Calculation == cap_model.MainData_basis {
			// расчет основных величин
			result.Calc.Basis = &cap_model.Calculated_Basis{
				Deformation:   &cap_model.CalcDeformation{},
				ForcesInBolts: &cap_model.CalcForcesInBolts{},
				BoltStrength:  &cap_model.CalcBoltStrength{},
				Moment:        &cap_model.CalcMoment{},
			}
		} else {
			// прочностной расчет
			result.Calc.Strength = &cap_model.Calculated_Strength{
				Auxiliary:              &cap_model.CalcAuxiliary{},
				Tightness:              &cap_model.CalcTightness{},
				BoltStrength1:          &cap_model.CalcBoltStrength{},
				Moment1:                &cap_model.CalcMoment{},
				StaticResistance1:      &cap_model.CalcStaticResistance{},
				ConditionsForStrength1: &cap_model.CalcConditionsForStrength{},
				TightnessLoad:          &cap_model.CalcTightnessLoad{},
				BoltStrength2:          &cap_model.CalcBoltStrength{},
				Moment2:                &cap_model.CalcMoment{},
				StaticResistance2:      &cap_model.CalcStaticResistance{},
				ConditionsForStrength2: &cap_model.CalcConditionsForStrength{},
				Deformation:            &cap_model.CalcDeformation{},
				ForcesInBolts:          &cap_model.CalcForcesInBolts{},
				FinalMoment:            &cap_model.CalcMoment{},
			}
		}
		logger.Error("failed to marshal json. error: " + err.Error())
	}

	return &result, nil
}
