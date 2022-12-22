package data

import (
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *DataService) FormatInitData(data *calc_api.FloatRequest) *float_model.DataResult {
	work := map[bool]string{
		true:  "Рабочие условия",
		false: "Условия испытаний",
	}
	typeD := map[string]string{
		"pin":  "Шпилька",
		"bolt": "Болт",
	}
	condition := map[string]string{
		"uncontrollable":  "Неконтролируемая затяжка",
		"controllable":    "Контроль по крутящему моменту",
		"controllablePin": "Контроль по вытяжке шпилек",
	}

	res := &float_model.DataResult{
		Pressure:  data.Pressure,
		HasThorn:  data.HasThorn,
		Work:      work[data.IsWork],
		Type:      typeD[data.Type.String()],
		Condition: condition[data.Condition.String()],
	}
	return res
}
