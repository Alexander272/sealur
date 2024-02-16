package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *FormulasService) getTubeSheetFormulas(
	data *calc_api.DevCoolingRequest,
	d models.DataDevCooling,
	result *calc_api.DevCoolingResponse,
) *dev_cooling_model.TubeSheetFormulas {
	TubeSheet := &dev_cooling_model.TubeSheetFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(data.Pressure, 'G', 5, 64)

	tLength := strconv.FormatFloat(d.Tube.Length, 'G', 5, 64)

	zoneThick := strconv.FormatFloat(d.TubeSheet.ZoneThick, 'G', 5, 64)
	tsSigma := strconv.FormatFloat(d.TubeSheet.Sigma, 'G', 5, 64)
	corrosion := strconv.FormatFloat(d.TubeSheet.Corrosion, 'G', 5, 64)

	CPressure := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Pressure, 'G', 5, 64), "E", "*10^")

	RelativeWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.RelativeWidth, 'G', 5, 64), "E", "*10^")
	EstimatedZoneWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.EstimatedZoneWidth, 'G', 5, 64), "E", "*10^")
	Upsilon := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Upsilon, 'G', 5, 64), "E", "*10^")
	Bp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Bp, 'G', 5, 64), "E", "*10^")
	Arm1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Arm1, 'G', 5, 64), "E", "*10^")
	Arm2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Arm2, 'G', 5, 64), "E", "*10^")
	LoadTube := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.LoadTube, 'G', 5, 64), "E", "*10^")
	PhiT := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.PhiT, 'G', 5, 64), "E", "*10^")
	Phi := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Phi, 'G', 5, 64), "E", "*10^")

	WorkEffort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.WorkEffort, 'G', 5, 64), "E", "*10^")
	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Effort, 'G', 5, 64), "E", "*10^")
	Eta := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Eta, 'G', 5, 64), "E", "*10^")
	Lp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Lp, 'G', 5, 64), "E", "*10^")

	OmegaP := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.OmegaP, 'G', 5, 64), "E", "*10^")
	Psi := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.Psi, 'G', 5, 64), "E", "*10^")
	Lambda := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.Lambda, 'G', 5, 64), "E", "*10^")
	TsEffort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.Effort, 'G', 5, 64), "E", "*10^")
	Omega := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.Omega, 'G', 5, 64), "E", "*10^")
	ZF := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.ZF, 'G', 5, 64), "E", "*10^")
	ZM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.ZM, 'G', 5, 64), "E", "*10^")

	//
	TubeSheet.Psi = fmt.Sprintf("%s * (%s + 2)", RelativeWidth, RelativeWidth)
	TubeSheet.Omega = fmt.Sprintf("1.6 * (%s / %s) * (%s * %s / %s)^(1/4)", EstimatedZoneWidth, zoneThick, Upsilon, zoneThick, tLength)

	// коэффициенты для Толщина трубной решетки в пределах зоны перфорации
	TubeSheet.Lambda = fmt.Sprintf("(4 * %s * %s) / (%s * (%s + %s) * (%s)^2)", WorkEffort, Arm1, pressure, Lp, Bp, EstimatedZoneWidth)
	var tmp1, tmp2, tmp3 string
	if data.Pressure*result.Calc.Auxiliary.Eta <= result.Calc.Auxiliary.PhiT*result.Calc.Auxiliary.LoadTube {
		TubeSheet.OmegaP = fmt.Sprintf("%s / (%s + %s * %s)", pressure, LoadTube, pressure, Eta)
	} else {
		tmp1 = fmt.Sprintf("%s * %s - %s * %s", pressure, Eta, PhiT, LoadTube)
		tmp2 = fmt.Sprintf("%s - %s * (2 - %s)", LoadTube, pressure, Eta)
		tmp3 = fmt.Sprintf("%s * %s * (1 + %s)", pressure, LoadTube, PhiT)
		TubeSheet.OmegaP = fmt.Sprintf("((%s)^2 + (%s) * (%s)) / %s", pressure, tmp1, tmp2, tmp3)
	}

	tmp1 = fmt.Sprintf("sqrt(%s / (%s * %s))", pressure, Phi, tsSigma)
	tmp2 = fmt.Sprintf("sqrt(%s + %s + %s + 1.5 * %s / (%s * %s))", Lambda, Psi, OmegaP, pressure, Phi, tsSigma)
	// s1 (s1min) - Толщина трубной решетки в пределах зоны перфорации
	TubeSheet.ZoneThick = fmt.Sprintf("0.71 * %s * %s * %s + %s", EstimatedZoneWidth, tmp1, tmp2, corrosion)

	// F1 - Расчетное усилие
	TubeSheet.Effort = fmt.Sprintf("(%s / (%s + %s)) * (%s / %s)", Effort, Lp, Bp, pressure, CPressure)

	tmp1 = fmt.Sprintf("sqrt(%s / %s)", TsEffort, tsSigma)
	tmp2 = fmt.Sprintf("sqrt(4 * %s + 1.5 * (%s / %s))", Arm1, TsEffort, tsSigma)
	tmp3 = fmt.Sprintf("sqrt(4 * %s + 1.5 * (%s / %s))", Arm2, TsEffort, tsSigma)
	// s2 (s2min) - Толщина трубной решетки в месте уплотнения
	TubeSheet.PlaceThick = fmt.Sprintf("0.71 * %s * %s + %s", tmp1, tmp2, corrosion)
	// s3 (s3min) - Толщина трубной решетки вне зоны уплотнения
	TubeSheet.OutZoneThick = fmt.Sprintf("0.71 * %s * %s + %s", tmp1, tmp3, corrosion)

	tmp1 = fmt.Sprintf("sh(%s) + sin(%s)", Omega, Omega)
	TubeSheet.ZF = fmt.Sprintf("%s * ((ch(%s) + cos(%s)) / (%s))", Omega, Omega, Omega, tmp1)
	TubeSheet.ZM = fmt.Sprintf("((%s)^2 / 4) * ((sh(%s) - sin(%s)) / (%s))", Omega, Omega, Omega, tmp1)

	//Условие прочности крепления труб в решетке
	TubeSheet.Strength = fmt.Sprintf("%s * (%s - %s + %s * (%s + %s))", pressure, ZF, Eta, ZM, Lambda, Psi)

	return TubeSheet
}
