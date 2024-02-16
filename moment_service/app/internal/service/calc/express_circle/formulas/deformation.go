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

func (s *FormulasService) getDeformationFormulas(req *calc_api.ExpressCircleRequest, d models.DataExCircle, result *calc_api.ExpressCircleResponse,
) *express_circle_model.DeformationFormulas {
	Deformation := &express_circle_model.DeformationFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Pressure, 'G', 5, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 5, 64)
	pres := strconv.FormatFloat(d.Gasket.Pres, 'G', 5, 64)
	dOut := strconv.FormatFloat(d.Gasket.DOut, 'G', 5, 64)
	m := strconv.FormatFloat(d.Gasket.M, 'G', 5, 64)

	dWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Deformation.Width, 'G', 5, 64), "E", "*10^")
	dDiameter := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Deformation.Diameter, 'G', 5, 64), "E", "*10^")

	if d.TypeGasket == express_circle_model.GasketData_Oval {
		// Эффективная ширина прокладки
		Deformation.Width = fmt.Sprintf("%s / 4", width)
		// Расчетный диаметр прокладки
		Deformation.Diameter = fmt.Sprintf("%s - %s / 2", dOut, width)

	} else {
		if d.Gasket.Width > constants.Bp {
			Deformation.Width = fmt.Sprintf("%.1f * sqrt(%s)", constants.B0, width)
		}
		// Расчетный диаметр прокладки
		Deformation.Diameter = fmt.Sprintf("%s - %s", dOut, dWidth)
	}

	// Усилие необходимое для смятия прокладки при затяжке
	Deformation.Deformation = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, dDiameter, dWidth, pres)

	if req.Pressure >= 0 {
		// Усилие на прокладке в рабочих условиях
		Deformation.Effort = fmt.Sprintf("%f * %s * %s * %s * |%s|", math.Pi, dDiameter, dWidth, m, pressure)
	}

	return Deformation
}
