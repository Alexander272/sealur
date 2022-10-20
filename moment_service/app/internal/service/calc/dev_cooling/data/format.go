package data

import (
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *DataService) FormatInitData(data *calc_api.DevCoolingRequest) *dev_cooling_model.DataResult {
	//? наверное это надо куда-нибудь вынести и можно будет передавать на клиент
	method := map[string]string{
		"AllThickness":  "На всю толщину решетки",
		"PartThickness": "В части толщины решетки",
		"SteelSheet":    "Стальная решетка с трубами из цветных металлов",
	}
	mounting := map[string]string{
		"flaring": "Развальцовка",
		"welding": "Приварка",
		"rolling": "Приварка с подвальцовкой",
	}
	typeBolt := map[string]string{
		"pin":  "Шпилька",
		"bolt": "Болт",
	}
	typeMounting := map[string]string{
		"flat":   "Гладкое соединение",
		"groove": "Вальцовка в канавку",
	}
	// cameraDiagram := map[string]string{
	// 	"schema1": "Черт. 1. ГОСТ 25822-83",
	// 	"schema2": "Черт. 2. ГОСТ 25822-83",
	// 	"schema3": "Черт. 3. ГОСТ 25822-83",
	// 	"schema4": "Черт. 4. ГОСТ 25822-83",
	// 	"schema5": "Черт. 5. ГОСТ 25822-83",
	// }
	// layout := map[string]string{
	// 	"lSchema1": "Черт. 11. ГОСТ 25822-83",
	// 	"lSchema2": "Черт. 12. ГОСТ 25822-83",
	// }

	res := &dev_cooling_model.DataResult{
		Pressure:      data.Pressure,
		Temp:          data.Temp,
		Method:        method[data.Method.String()],
		TypeBolt:      typeBolt[data.TypeBolt.String()],
		Mounting:      mounting[data.Mounting.String()],
		TypeMounting:  typeMounting[data.TypeMounting.String()],
		CameraDiagram: data.CameraDiagram.String(),
		Layout:        data.Layout.String(),
	}
	return res
}
