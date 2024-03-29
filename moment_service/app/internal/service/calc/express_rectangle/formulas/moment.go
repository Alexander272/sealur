package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_rectangle_model"
)

func (s *FormulasService) getMomentFormulas(req calc_api.ExpressRectangleRequest, d models.DataExRect, result calc_api.ExpressRectangleResponse,
) *express_rectangle_model.MomentFormulas {
	Moment := &express_rectangle_model.MomentFormulas{}

	// перевод чисел в строки
	diameter := strconv.FormatFloat(d.Bolt.Diameter, 'G', 3, 64)
	count := d.Bolt.Count
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 3, 64)
	permissiblePres := strconv.FormatFloat(d.Gasket.PermissiblePres, 'G', 3, 64)

	SizeLong := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeLong, 'G', 3, 64), "E", "*10^")
	SizeTrans := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeTrans, 'G', 3, 64), "E", "*10^")

	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Effort, 'G', 3, 64), "E", "*10^")
	Area := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Area, 'G', 3, 64), "E", "*10^")

	AllowableVoltage := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.AllowableVoltage, 'G', 3, 64), "E", "*10^")

	Mkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Moment.Mkp, 'G', 3, 64), "E", "*10^")

	if !(result.Calc.Bolt.RatedStress > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter) {
		Moment.Mkp = fmt.Sprintf("(0.3 * %s * %s / %d) / 1000", Effort, diameter, count)
	}
	Moment.Mkp1 = fmt.Sprintf("0.75 * %s", Mkp)

	Prek := fmt.Sprintf("0.8 * %s * %s", Area, sigmaAt20)
	Moment.Qrek = fmt.Sprintf("%s / (2 * (%s + %s) * %s)", Prek, SizeLong, SizeTrans, width)
	Moment.Mrek = fmt.Sprintf("(0.3 * %s * %s / %d) / 1000", Prek, diameter, count)

	Pmax := fmt.Sprintf("%s * %s", AllowableVoltage, Area)
	Moment.Qmax = fmt.Sprintf("%s / (2 * (%s + %s) * %s)", Pmax, SizeLong, SizeTrans, width)

	if d.TypeGasket == express_rectangle_model.GasketData_Soft && result.Calc.Moment.Qmax > d.Gasket.PermissiblePres {
		Pmax = fmt.Sprintf("%s * (2 * (%s + %s) * %s)", permissiblePres, SizeLong, SizeTrans, width)
	}

	Moment.Mmax = fmt.Sprintf("(0.3 * %s * %s / %d) / 1000", Pmax, diameter, count)

	return Moment
}
