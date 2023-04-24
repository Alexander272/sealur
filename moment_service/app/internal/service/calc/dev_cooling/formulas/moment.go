package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *FormulasService) getMomentFormulas(
	data calc_api.DevCoolingRequest,
	d models.DataDevCooling,
	result calc_api.DevCoolingResponse,
) *dev_cooling_model.MomentFormulas {
	Moment := &dev_cooling_model.MomentFormulas{}

	Ab := d.Bolt.Area * float64(d.Bolt.Count)

	// перевод чисел в строки
	friction := strconv.FormatFloat(data.Friction, 'G', 3, 64)
	sAb := strings.ReplaceAll(strconv.FormatFloat(Ab, 'G', 3, 64), "E", "*10^")

	diameter := strconv.FormatFloat(d.Bolt.Diameter, 'G', 3, 64)
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)
	count := d.Bolt.Count

	gWidth := strconv.FormatFloat(d.Gasket.Width, 'G', 3, 64)
	permissiblePres := strconv.FormatFloat(d.Gasket.PermissiblePres, 'G', 3, 64)

	Bp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Bp, 'G', 3, 64), "E", "*10^")

	Lp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Lp, 'G', 3, 64), "E", "*10^")
	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Effort, 'G', 3, 64), "E", "*10^")
	WorkCondY := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.WorkCond.Y, 'G', 3, 64), "E", "*10^")

	Mkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Moment.Mkp, 'G', 3, 64), "E", "*10^")

	if !(result.Calc.Bolt.WorkCond.X > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter) {
		Moment.Mkp = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Effort, diameter, count)
	}

	// Крутящий момент при затяжке болтов/шпилек со смазкой снижается на 25%
	Moment.Mkp1 = fmt.Sprintf("0.75 * %s", Mkp)

	Prek := fmt.Sprintf("0.8 * %s * %s", sAb, sigmaAt20)
	// Напряжение на прокладке
	Moment.Qrek = fmt.Sprintf("%s / (2 * (%s + %s) * %s)", Prek, Lp, Bp, gWidth)
	// Момент затяжки при применении уплотнения на старых (изношенных) фланцах, имеющих перекосы
	Moment.Mrek = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Prek, diameter, count)

	Pmax := fmt.Sprintf("%s * %s", WorkCondY, sAb)
	// Максимальное напряжение на прокладке
	Moment.Qmax = fmt.Sprintf("%s / (2 * (%s + %s) * %s)", Pmax, Lp, Bp, gWidth)
	if result.Calc.Moment.Qmax > d.Gasket.PermissiblePres {
		Pmax = fmt.Sprintf("%s * (2 * (%s + %s) * %s)", permissiblePres, Lp, Bp, gWidth)
		Moment.Qmax = ""
	}

	if result.Calc.Moment.Mrek > result.Calc.Moment.Mmax {
		Moment.Mrek = ""
	}
	if result.Calc.Moment.Qrek > result.Calc.Moment.Qmax {
		Moment.Qrek = ""
	}

	// Максимальный крутящий момент при затяжке болтов/шпилек
	Moment.Mmax = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Pmax, diameter, count)

	return Moment
}
