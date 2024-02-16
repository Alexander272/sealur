package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_rectangle_model"
)

func (s *FormulasService) getForcesFormulas(req *calc_api.ExpressRectangleRequest, d models.DataExRect, result *calc_api.ExpressRectangleResponse,
) *express_rectangle_model.ForcesInBoltsFormulas {
	ForcesInBolts := &express_rectangle_model.ForcesInBoltsFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Pressure, 'G', 5, 64)
	testPressure := strconv.FormatFloat(result.Data.TestPressure, 'G', 5, 64)

	area := strconv.FormatFloat(d.Bolt.Area, 'G', 5, 64)
	count := d.Bolt.Count

	m := strconv.FormatFloat(d.Gasket.M, 'G', 5, 64)

	EstimatedGasketWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.EstimatedGasketWidth, 'G', 5, 64), "E", "*10^")
	SizeLong := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeLong, 'G', 5, 64), "E", "*10^")
	SizeTrans := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.SizeTrans, 'G', 5, 64), "E", "*10^")

	Area := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Area, 'G', 5, 64), "E", "*10^")

	// Суммарная площадь сечения болтов/шпилек
	ForcesInBolts.Area = fmt.Sprintf("%d * %s", count, area)

	// Fв - Расчетное усилие в болтах (шпильках) в условиях эксплуатации
	ForcesInBolts.WorkEffort = fmt.Sprintf("%s * (%s * %s + 2 * %s * %s * (%s + %s))",
		pressure, SizeLong, SizeTrans, EstimatedGasketWidth, m, SizeLong, SizeTrans)

	tmp1 := fmt.Sprintf("(%s / %s) * %s", testPressure, pressure, Area)
	tmp2 := fmt.Sprintf("%s * (2 * %s * %s + 2 * %s * %s * (%s + %s))", testPressure, SizeLong, SizeTrans, EstimatedGasketWidth, m, SizeLong, SizeTrans)
	// F0 - Расчетное усилие в болтах (шпильках) в условиях испытаний или монтажа
	ForcesInBolts.Effort = fmt.Sprintf("max(%s; %s)", tmp1, tmp2)

	return ForcesInBolts
}
