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

func (s *FormulasService) getMomentFormulas(req *calc_api.GasCoolingRequest, data models.DataGasCooling, result *calc_api.GasCoolingResponse,
) *gas_cooling_model.MomentFormulas {
	Moment := &gas_cooling_model.MomentFormulas{}

	// перевод чисел в строки
	friction := strconv.FormatFloat(req.Data.Friction, 'G', 5, 64)
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 5, 64)
	count := data.Bolt.Count
	sigmaAt20 := strconv.FormatFloat(data.Bolt.SigmaAt20, 'G', 5, 64)

	width := strconv.FormatFloat(data.Gasket.Width, 'G', 5, 64)
	permissiblePres := strconv.FormatFloat(data.Gasket.PermissiblePres, 'G', 5, 64)

	SizeLong := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeLong, 'G', 5, 64), "E", "*10^")
	SizeTrans := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeTrans, 'G', 5, 64), "E", "*10^")

	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Effort, 'G', 5, 64), "E", "*10^")
	Area := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Area, 'G', 5, 64), "E", "*10^")

	AllowableVoltage := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.AllowableVoltage, 'G', 5, 64), "E", "*10^")

	Mkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Moment.Mkp, 'G', 5, 64), "E", "*10^")

	if !(result.Calc.Bolt.RatedStress > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter) {
		Moment.Mkp = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Effort, diameter, count)
	}
	Moment.Mkp1 = fmt.Sprintf("0.75 * %s", Mkp)

	Prek := fmt.Sprintf("0.8 * %s * %s", Area, sigmaAt20)
	Moment.Qrek = fmt.Sprintf("%s / (2 * (%s + %s) * %s)", Prek, SizeLong, SizeTrans, width)
	Moment.Mrek = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Prek, diameter, count)

	Pmax := fmt.Sprintf("%s * %s", AllowableVoltage, Area)
	Moment.Qmax = fmt.Sprintf("%s / (2 * (%s + %s) * %s)", Pmax, SizeLong, SizeTrans, width)

	if data.TypeGasket == gas_cooling_model.GasketData_Soft && result.Calc.Moment.Qmax > data.Gasket.PermissiblePres {
		Pmax = fmt.Sprintf("%s * (2 * (%s + %s) * %s)", permissiblePres, SizeLong, SizeTrans, width)
	}

	if result.Calc.Moment.Mrek > result.Calc.Moment.Mmax {
		Moment.Mrek = ""
	}
	if result.Calc.Moment.Qrek > result.Calc.Moment.Qmax {
		Moment.Qrek = ""
	}

	Moment.Mmax = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Pmax, diameter, count)

	return Moment
}
