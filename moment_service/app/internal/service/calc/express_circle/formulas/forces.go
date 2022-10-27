package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
)

func (s *FormulasService) getForcesFormulas(req calc_api.ExpressCircleRequest, d models.DataExCircle, result calc_api.ExpressCircleResponse,
) *express_circle_model.ForcesInBoltsFormulas {
	ForcesInBolts := &express_circle_model.ForcesInBoltsFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Pressure, 'G', 3, 64)

	area := strconv.FormatFloat(d.Bolt.Area, 'G', 3, 64)
	count := d.Bolt.Count
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)

	dDiameter := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Deformation.Diameter, 'G', 3, 64), "E", "*10^")
	dDeformation := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Deformation.Deformation, 'G', 3, 64), "E", "*10^")
	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Deformation.Effort, 'G', 3, 64), "E", "*10^")

	FArea := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Area, 'G', 3, 64), "E", "*10^")
	ResultantLoad := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.ResultantLoad, 'G', 3, 64), "E", "*10^")
	EstimatedLoad1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.EstimatedLoad1, 'G', 3, 64), "E", "*10^")
	EstimatedLoad2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.EstimatedLoad2, 'G', 3, 64), "E", "*10^")
	DesignLoad := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.DesignLoad, 'G', 3, 64), "E", "*10^")

	// Суммарная площадь сечения болтов/шпилек
	ForcesInBolts.Area = fmt.Sprintf("%d * %s", count, area)

	// Равнодействующая нагрузка от давления
	ForcesInBolts.ResultantLoad = fmt.Sprintf("0.785 * (%s)^2 * %s", dDiameter, pressure)

	// Минимальное начальное натяжение болтов (шпилек)
	Tension := fmt.Sprintf("0.4 * %s * %s", FArea, sigmaAt20)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	ForcesInBolts.EstimatedLoad2 = fmt.Sprintf("max(%s; %s)", dDeformation, Tension)
	// Расчетная нагрузка на болты/шпильки при затяжке необходимая для обеспечения в рабочих условиях давления на прокладку достаточного
	// для герметизации фланцевого соединения
	ForcesInBolts.EstimatedLoad1 = fmt.Sprintf("1 * (%s) + %s", ResultantLoad, Effort)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений
	ForcesInBolts.DesignLoad = fmt.Sprintf("max(%s; %s)", EstimatedLoad1, EstimatedLoad2)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений в рабочих условиях
	ForcesInBolts.WorkDesignLoad = fmt.Sprintf("%s + (1 - 1) * %s", DesignLoad, ResultantLoad)

	return ForcesInBolts
}
