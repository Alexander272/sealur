package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
)

func (s *FormulasService) getMomentFormulas(req calc_api.ExpressCircleRequest, d models.DataExCircle, result calc_api.ExpressCircleResponse,
) *express_circle_model.MomentFormulas {
	Moment := &express_circle_model.MomentFormulas{}

	// перевод чисел в строки
	friction := strconv.FormatFloat(req.Friction, 'G', 3, 64)
	diameter := strconv.FormatFloat(d.Bolt.Diameter, 'G', 3, 64)
	count := d.Bolt.Count
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 3, 64)
	permissiblePres := strconv.FormatFloat(d.Gasket.PermissiblePres, 'G', 3, 64)

	dDiameter := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Deformation.Diameter, 'G', 3, 64), "E", "*10^")

	DesignLoad := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.DesignLoad, 'G', 3, 64), "E", "*10^")
	Area := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Area, 'G', 3, 64), "E", "*10^")

	AllowableVoltage := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.AllowableVoltage, 'G', 3, 64), "E", "*10^")

	Mkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Moment.Mkp, 'G', 3, 64), "E", "*10^")

	if !(result.Calc.Bolt.RatedStress > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter) {
		Moment.Mkp = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, DesignLoad, diameter, count)
	}
	Moment.Mkp1 = fmt.Sprintf("0.75 * %s", Mkp)

	Prek := fmt.Sprintf("0.8 * %s * %s", Area, sigmaAt20)
	Moment.Qrek = fmt.Sprintf("%s / (%f * %s * %s)", Prek, math.Pi, dDiameter, width)
	Moment.Mrek = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Prek, diameter, count)

	Pmax := fmt.Sprintf("%s * %s", AllowableVoltage, Area)
	Moment.Qmax = fmt.Sprintf("%s / (%f * %s * %s)", Pmax, math.Pi, dDiameter, width)

	if d.TypeGasket == express_circle_model.GasketData_Soft && result.Calc.Moment.Qmax > d.Gasket.PermissiblePres {
		Pmax = fmt.Sprintf("%s * (%f * %s * %s)", permissiblePres, math.Pi, dDiameter, width)
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
