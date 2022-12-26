package data

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
)

func (s *DataService) getData(ctx context.Context, data *gas_cooling_model.MainData) *gas_cooling_model.DataResult {
	testPressure := data.TestPressure
	if testPressure == 0 {
		testPressure = 1.5 * data.Pressure.Value
	}

	res := &gas_cooling_model.DataResult{
		Device:        data.Device.Title,
		Factor:        data.Factor.Value,
		Pressure:      data.Pressure.Value,
		Section:       data.Section.Value,
		NumberOfMoves: data.NumberOfMoves.Value,
		TubeLength:    data.TubeLength.Value,
		TestPressure:  testPressure,
		Type:          data.TypeBolt.String(),
		Condition:     data.Condition.String(),
	}
	return res
}

func (s *DataService) FormatInitData(data *gas_cooling_model.MainData) *gas_cooling_model.DataResult {
	typeBolt := map[string]string{
		"pin":  "Шпилька",
		"bolt": "Болт",
	}
	condition := map[string]string{
		"uncontrollable":  "Неконтролируемая затяжка",
		"controllable":    "Контроль по крутящему моменту",
		"controllablePin": "Контроль по вытяжке шпилек",
	}
	testPressure := data.TestPressure
	if testPressure == 0 {
		testPressure = 1.5 * data.Pressure.Value
	}

	res := &gas_cooling_model.DataResult{
		Device:        data.Device.Title,
		Factor:        data.Factor.Value,
		Pressure:      data.Pressure.Value,
		Section:       data.Section.Value,
		TubeCount:     fmt.Sprintf("%d", data.TubeCount.Value),
		NumberOfMoves: data.NumberOfMoves.Value,
		TubeLength:    data.TubeLength.Value,
		TestPressure:  testPressure,
		Type:          typeBolt[data.TypeBolt.String()],
		Condition:     condition[data.Condition.String()],
	}
	return res
}
