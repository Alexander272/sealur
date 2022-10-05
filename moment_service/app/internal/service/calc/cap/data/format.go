package data

import "github.com/Alexander272/sealur_proto/api/moment_api"

func (s *DataService) FormatInitData(data *moment_api.CalcCapRequest) *moment_api.DataResult {
	work := map[bool]string{
		true:  "Рабочие условия",
		false: "Условия испытаний",
	}
	flanges := map[string]string{
		"isolated":    "Изолированные фланцы",
		"nonIsolated": "Неизолированные фланцы",
		"manually":    "Задается вручную",
	}
	sameFlange := map[bool]string{
		true:  "Да",
		false: "Нет",
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

	res := &moment_api.DataResult{
		Pressure:   data.Pressure,
		AxialForce: data.AxialForce,
		Temp:       data.Temp,
		Work:       work[data.IsWork],
		Flanges:    flanges[data.Flanges.String()],
		Embedded:   sameFlange[data.IsEmbedded],
		Type:       typeD[data.Type.String()],
		Condition:  condition[data.Condition.String()],
	}
	return res
}
