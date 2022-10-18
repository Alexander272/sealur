package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *FormulasService) getCapFormulas(
	data calc_api.DevCoolingRequest,
	d models.DataDevCooling,
	result calc_api.DevCoolingResponse,
) *dev_cooling_model.CapFormulas {
	Cap := &dev_cooling_model.CapFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(data.Pressure, 'G', 3, 64)

	innerSize := strconv.FormatFloat(d.Cap.InnerSize, 'G', 3, 64)
	l := strconv.FormatFloat(d.Cap.L, 'G', 3, 64)
	depth := strconv.FormatFloat(d.Cap.Depth, 'G', 3, 64)
	flangeThick := strconv.FormatFloat(d.Cap.FlangeThick, 'G', 3, 64)
	wallThick := strconv.FormatFloat(d.Cap.WallThick, 'G', 3, 64)
	bottomThick := strconv.FormatFloat(d.Cap.BottomThick, 'G', 3, 64)
	sigma := strconv.FormatFloat(d.Cap.Sigma, 'G', 3, 64)
	strength := strconv.FormatFloat(d.Cap.Strength, 'G', 3, 64)
	corrosion := strconv.FormatFloat(d.Cap.Corrosion, 'G', 3, 64)

	distance := strconv.FormatFloat(d.Bolt.Distance, 'G', 3, 64)

	tsSigma := strconv.FormatFloat(d.TubeSheet.Sigma, 'G', 3, 64)

	Bp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Bp, 'G', 3, 64), "E", "*10^")
	Arm1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Arm1, 'G', 3, 64), "E", "*10^")

	Lp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Lp, 'G', 3, 64), "E", "*10^")
	WorkEffort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.WorkEffort, 'G', 3, 64), "E", "*10^")

	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.Effort, 'G', 3, 64), "E", "*10^")

	F1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.F1, 'G', 3, 64), "E", "*10^")
	F2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.F2, 'G', 3, 64), "E", "*10^")
	Lambda := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.Lambda, 'G', 3, 64), "E", "*10^")
	Psi := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.Psi, 'G', 3, 64), "E", "*10^")
	ChiK := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.ChiK, 'G', 3, 64), "E", "*10^")
	Chi := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.Chi, 'G', 3, 64), "E", "*10^")
	WallThick := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.WallThick, 'G', 3, 64), "E", "*10^")

	//
	Cap.Psi = fmt.Sprintf("((%s / %s)^2-1)*(%s / (%s + %s)) - 4 * (%s / %s)^2", Bp, innerSize, l, l, innerSize, depth, innerSize)

	// коэффициенты для Толщина донышка крышки
	Cap.Lambda = fmt.Sprintf("(4 * %s * %s) / (%s * (%s + %s) * (%s)^2)", WorkEffort, Arm1, pressure, Lp, Bp, innerSize)

	var tmp1, tmp2 string
	if data.CameraDiagram != calc_api.DevCoolingRequest_schema4 {
		tmp1 = fmt.Sprintf("1.5 * (%s - %s) - %s", distance, innerSize, flangeThick)
		Cap.Chi = fmt.Sprintf("(0.8 / %s) * (%s) * (%s / %s)^2", l, tmp1, flangeThick, wallThick)
	} else {
		tmp1 = fmt.Sprintf("6 * %s - %s - %s", flangeThick, distance, innerSize)
		tmp2 = fmt.Sprintf("((%s - %s) / %s)^2", distance, innerSize, bottomThick)
		Cap.Chi = fmt.Sprintf("(0.1 / %s) * (%s) * (%s)", l, tmp1, tmp2)
	}

	if data.CameraDiagram != calc_api.DevCoolingRequest_schema5 {
		Cap.Psi = fmt.Sprintf("((%s / %s)^2 - 1) * (%s / (%s + %s)) - 4 * (%s / %s)^2", Bp, innerSize, l, l, innerSize, depth, innerSize)
		Cap.F1 = fmt.Sprintf("1 / (1 + (%s / %s) + (%s / %s)^2)", innerSize, l, innerSize, l)
		Cap.F2 = fmt.Sprintf("0.5 * %s", F1)
		if data.CameraDiagram != calc_api.DevCoolingRequest_schema4 {
			tmp1 = fmt.Sprintf("(1.5 * (%s - %s) - %s) * (%s / %s)^2", distance, innerSize, flangeThick, flangeThick, bottomThick)
			tmp2 = fmt.Sprintf("(3 * (%s - %s) + 2 * %s) * (%s / %s)^2", depth, flangeThick, wallThick, wallThick, bottomThick)
			Cap.ChiK = fmt.Sprintf("(0.8 / %s) * (%s + %s)", l, tmp1, tmp2)
		} else {
			tmp1 = fmt.Sprintf("6 * %s - %s - %s", flangeThick, distance, innerSize)
			tmp2 = fmt.Sprintf("((%s - %s) / %s)^2", distance, innerSize, bottomThick)
			Cap.ChiK = fmt.Sprintf("(0.1 / %s) * (%s) * (%s)", l, tmp1, tmp2)
		}
		tmp1 = fmt.Sprintf("(%s + %s + %s) / (1 + %s)", Lambda, Psi, F1, ChiK)
		// s4 (s4min) - Толщина донышка крышки
		Cap.BottomThick = fmt.Sprintf("0.71 * %s * sqrt(%s / %s) * (sqrt(max((%s); %s) + 1.5 * %s/%s)) + %s",
			innerSize, pressure, sigma, tmp1, F2, pressure, sigma, corrosion)
	} else {
		// s4 (s4min) - Толщина донышка крышки
		Cap.BottomThick = fmt.Sprintf("0.71 * %s * sqrt(%s / %s) * sqrt((%s / (%s + %s)) + 0.5 * (%s / ((%s)^2 * %s))) + %s",
			innerSize, pressure, sigma, Lambda, strength, Chi, pressure, strength, sigma, corrosion)
	}

	tmp1 = fmt.Sprintf("sqrt(%s / %s)", Effort, tsSigma)
	tmp2 = fmt.Sprintf("sqrt(4 * %s / %s + %s)", Arm1, strength, Chi)
	// s5 (s5min) - Толщина стенки крышки в месте присоединения к фланцу
	Cap.WallThick = fmt.Sprintf("0.71 * (%s) * (%s) + %s", tmp1, tmp2, corrosion)

	tmp1 = fmt.Sprintf("sqrt(%s / %s)", Effort, sigma)
	tmp2 = fmt.Sprintf("sqrt(4 * %s + 1.5 * (%s / %s))", Arm1, Effort, sigma)
	// s6 (s6min) - Толщина фланца крышки
	Cap.FlangeThick = fmt.Sprintf("0.71 * (%s) * (%s) + %s", tmp1, tmp2, corrosion)

	if data.CameraDiagram == calc_api.DevCoolingRequest_schema5 {
		Cap.SideWallThick = fmt.Sprintf("max(%s; (0.25 * %s * sqrt(%s / %s) + %s))", WallThick, innerSize, pressure, sigma, corrosion)
	}

	return Cap
}
