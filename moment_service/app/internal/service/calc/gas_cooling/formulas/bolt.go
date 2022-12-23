package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
)

func (s *FormulasService) getBoltFormulas(req *calc_api.GasCoolingRequest, data models.DataGasCooling, result *calc_api.GasCoolingResponse,
) *gas_cooling_model.BoltsFormulas {
	Bolts := &gas_cooling_model.BoltsFormulas{}

	// перевод чисел в строки
	sigmaAt20 := strconv.FormatFloat(data.Bolt.SigmaAt20, 'G', 3, 64)

	width := strconv.FormatFloat(data.Gasket.Width, 'G', 3, 64)

	SizeLong := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeLong, 'G', 3, 64), "E", "*10^")
	SizeTrans := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeTrans, 'G', 3, 64), "E", "*10^")

	Area := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Area, 'G', 3, 64), "E", "*10^")
	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Effort, 'G', 3, 64), "E", "*10^")
	WorkEffort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.WorkEffort, 'G', 3, 64), "E", "*10^")

	Kyp := s.Kyp[true]
	Kyz := s.Kyz[req.Data.Condition.String()]
	Kyt := constants.LoadKyt

	// Расчетное напряжение в болтах/шпильках - при затяжке
	Bolts.RatedStress = fmt.Sprintf("%s / %s", Effort, Area)
	// Условия прочности болтов шпилек - при затяжке
	Bolts.AllowableVoltage = fmt.Sprintf("1.2 * %.f * %.1f * %.1f * %s", Kyp, Kyz, Kyt, sigmaAt20)

	if data.TypeGasket == gas_cooling_model.GasketData_Soft {
		// Условие прочности прокладки
		Bolts.StrengthGasket = fmt.Sprintf("max(%s; %s) / (2 * (%s + %s) * %s)", WorkEffort, Effort, SizeLong, SizeTrans, width)
	}

	return Bolts
}
