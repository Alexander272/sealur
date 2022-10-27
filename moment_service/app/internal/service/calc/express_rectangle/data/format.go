package data

import (
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_rectangle_model"
)

func (s *DataService) FormatInitData(data *calc_api.ExpressRectangleRequest) *express_rectangle_model.DataResult {
	typeBolt := map[string]string{
		"pin":  "Шпилька",
		"bolt": "Болт",
	}
	condition := map[string]string{
		"uncontrollable":  "Неконтролируемая затяжка",
		"controllable":    "Контроль по крутящему моменту",
		"controllablePin": "Контроль по вытяжке шпилек",
	}

	res := &express_rectangle_model.DataResult{
		Pressure:     data.Pressure,
		TestPressure: data.TestPressure,
		Type:         typeBolt[data.TypeBolt.String()],
		Condition:    condition[data.Condition.String()],
	}
	return res
}
