package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *FormulasService) GetBoltFormulas(
	Lambda1, Lambda2, Alpha1, Alpha2 float64,
	data calc_api.DevCoolingRequest,
	d models.DataDevCooling,
	result calc_api.DevCoolingResponse,
) *dev_cooling_model.BoltFormulas {
	Bolt := &dev_cooling_model.BoltFormulas{}

	// формулы (чтобы не передевать всю эту кучу фи)
	// Phi для Угловые податливости крышки
	var Phi1, Phi2, Phi3, Phi4, Phi5, Phi6 float64
	if data.CameraDiagram == calc_api.DevCoolingRequest_schema1 || data.CameraDiagram == calc_api.DevCoolingRequest_schema4 {
		Phi1 = 1
		Phi2 = 8 * math.Pow(d.Cap.Depth/d.Cap.InnerSize, 3)
		Phi4 = 1
		Phi5 = 2 * (d.Cap.Depth / d.Cap.InnerSize)
	} else {
		Phi1 = 1 + 0.85*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 2) - 12.55*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 3) +
			13.7*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 2)*(d.Cap.Depth/d.Cap.InnerSize)
		Phi2 = 8*math.Pow(d.Cap.Depth/d.Cap.InnerSize, 3) - 12*(d.Cap.Depth/d.Cap.InnerSize)*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 2) +
			4*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 3)
		Phi4 = 1 - 1.14*(d.Cap.Radius/d.Cap.InnerSize)
		Phi5 = 2*(d.Cap.Depth/d.Cap.InnerSize) - 2*(d.Cap.Radius/d.Cap.InnerSize)
	}
	Phi3 = 12*math.Pow(d.Cap.Depth/d.Cap.InnerSize, 2)*(d.Cap.FlangeThick/d.Cap.InnerSize) - 4*math.Pow(d.Cap.FlangeThick/d.Cap.InnerSize, 3)
	Phi6 = 2 * (d.Cap.FlangeThick / d.Cap.InnerSize)

	Lb := d.Bolt.Lenght + s.typeBolt[data.TypeBolt.String()]*d.Bolt.Diameter

	// перевод чисел в строки
	sPhi1 := strings.ReplaceAll(strconv.FormatFloat(Phi1, 'G', 3, 64), "E", "*10^")
	sPhi2 := strings.ReplaceAll(strconv.FormatFloat(Phi2, 'G', 3, 64), "E", "*10^")
	sPhi3 := strings.ReplaceAll(strconv.FormatFloat(Phi3, 'G', 3, 64), "E", "*10^")
	sPhi4 := strings.ReplaceAll(strconv.FormatFloat(Phi4, 'G', 3, 64), "E", "*10^")
	sPhi5 := strings.ReplaceAll(strconv.FormatFloat(Phi5, 'G', 3, 64), "E", "*10^")
	sPhi6 := strings.ReplaceAll(strconv.FormatFloat(Phi6, 'G', 3, 64), "E", "*10^")
	sLambda1 := strings.ReplaceAll(strconv.FormatFloat(Lambda1, 'G', 3, 64), "E", "*10^")
	sLambda2 := strings.ReplaceAll(strconv.FormatFloat(Lambda2, 'G', 3, 64), "E", "*10^")
	sAlpha1 := strings.ReplaceAll(strconv.FormatFloat(Alpha1, 'G', 3, 64), "E", "*10^")
	sAlpha2 := strings.ReplaceAll(strconv.FormatFloat(Alpha2, 'G', 3, 64), "E", "*10^")
	sLb := strings.ReplaceAll(strconv.FormatFloat(Lb, 'G', 3, 64), "E", "*10^")

	pressure := strconv.FormatFloat(data.Pressure, 'G', 3, 64)

	innerSize := strconv.FormatFloat(d.Cap.InnerSize, 'G', 3, 64)
	bottomThick := strconv.FormatFloat(d.Cap.BottomThick, 'G', 3, 64)
	wallThick := strconv.FormatFloat(d.Cap.WallThick, 'G', 3, 64)
	capEpsilon := strconv.FormatFloat(d.Cap.Epsilon, 'G', 3, 64)

	sizeLong := strconv.FormatFloat(d.Gasket.SizeLong, 'G', 3, 64)
	gThickness := strconv.FormatFloat(d.Gasket.Thickness, 'G', 3, 64)
	gWidth := strconv.FormatFloat(d.Gasket.Width, 'G', 3, 64)
	m := strconv.FormatFloat(d.Gasket.M, 'G', 3, 64)

	bEpsilon := strconv.FormatFloat(d.Bolt.Epsilon, 'G', 3, 64)
	area := strconv.FormatFloat(d.Bolt.Area, 'G', 3, 64)
	count := d.Bolt.Count

	zoneThick := strconv.FormatFloat(d.TubeSheet.ZoneThick, 'G', 3, 64)
	outZoneThick := strconv.FormatFloat(d.TubeSheet.OutZoneThick, 'G', 3, 64)
	tsEpsilon := strconv.FormatFloat(d.TubeSheet.Epsilon, 'G', 3, 64)

	CapPsi := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Cap.Psi, 'G', 3, 64), "E", "*10^")

	TsPsi := strings.ReplaceAll(strconv.FormatFloat(result.Calc.TubeSheet.Psi, 'G', 3, 64), "E", "*10^")

	EstimatedGasketWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.EstimatedGasketWidth, 'G', 3, 64), "E", "*10^")
	EstimatedZoneWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.EstimatedZoneWidth, 'G', 3, 64), "E", "*10^")
	RelativeWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.RelativeWidth, 'G', 3, 64), "E", "*10^")
	Bp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Bp, 'G', 3, 64), "E", "*10^")
	Arm1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Arm1, 'G', 3, 64), "E", "*10^")

	Lp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Lp, 'G', 3, 64), "E", "*10^")
	UpsilonB := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.UpsilonB, 'G', 3, 64), "E", "*10^")
	UpsilonP := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.UpsilonP, 'G', 3, 64), "E", "*10^")
	CapUpsilonM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.CapUpsilonM, 'G', 3, 64), "E", "*10^")
	SheetUpsilonM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.SheetUpsilonM, 'G', 3, 64), "E", "*10^")
	CapUpsilonP := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.CapUpsilonP, 'G', 3, 64), "E", "*10^")
	SheetUpsilonP := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.SheetUpsilonP, 'G', 3, 64), "E", "*10^")

	tmp1 := fmt.Sprintf("(%s)^3 / (%s * (%s)^3)", innerSize, capEpsilon, bottomThick)
	var tmp2, tmp3, tmpPhi string
	if data.CameraDiagram == calc_api.DevCoolingRequest_schema4 {
		tmpPhi = sPhi4
		tmp2 = fmt.Sprintf("%s * %s", sPhi1, sLambda1)
	} else {
		tmp2 = fmt.Sprintf("(%s + (%s - %s) * (%s / %s)^3) * %s", sPhi1, sPhi2, sPhi3, bottomThick, wallThick, sLambda1)
		tmpPhi = fmt.Sprintf("%s + (%s - %s) * (%s / %s)^3", sPhi4, sPhi5, sPhi6, bottomThick, wallThick)
	}
	tmp3 = fmt.Sprintf("1/8 * %s * %s * %s", tmpPhi, CapPsi, sLambda2)
	// Угловые податливости крышки
	Bolt.CapUpsilonP = fmt.Sprintf("10.9 * %s * (%s + %s)", tmp1, tmp2, tmp3)

	Bolt.Lp = fmt.Sprintf("%s - %s", sizeLong, EstimatedGasketWidth)
	// Угловые податливости крышки
	Bolt.CapUpsilonM = fmt.Sprintf("10.9 * (%s / (2 * %s * (%s)^3 * (%s + %s))) * %s * %s", innerSize, capEpsilon, bottomThick, Lp, Bp, tmpPhi, sLambda2)

	tmp1 = fmt.Sprintf("0.23 * (%s)^3 / (%s * (%s)^3)", EstimatedZoneWidth, tsEpsilon, zoneThick)
	tmp2 = fmt.Sprintf("%s * (2 * %s - %s) * (%s / %s)^3", RelativeWidth, TsPsi, RelativeWidth, zoneThick, outZoneThick)
	tmp3 = fmt.Sprintf("1.7 * (%s * %s + 4 * %s)", TsPsi, sAlpha1, sAlpha2)
	// Угловые податливости решетки
	Bolt.SheetUpsilonP = fmt.Sprintf("%s * (%s + %s)", tmp1, tmp2, tmp3)

	tmp1 = fmt.Sprintf("%s / (2 * %s * (%s)^3 * (%s + %s))", EstimatedZoneWidth, tsEpsilon, zoneThick, Lp, Bp)
	tmp2 = fmt.Sprintf("2 * %s * (%s / %s)^3 + 1.1 * %s", RelativeWidth, zoneThick, outZoneThick, sAlpha1)
	// Угловые податливости решетки
	Bolt.SheetUpsilonM = fmt.Sprintf("2.7 * %s * %s", tmp1, tmp2)

	// Yb Линейная податливость болта (шпильки)
	Bolt.UpsilonB = fmt.Sprintf("%s / (%s * %s * %d)", sLb, bEpsilon, area, count)
	// Yp Линейная податливость прокладки
	Bolt.UpsilonP = fmt.Sprintf("%s / (2 * %s * (%s + %s) * %s)", gThickness, bEpsilon, Lp, Bp, gWidth)

	tmp1 = fmt.Sprintf("%s + (%s + %s) * (%s)^2", UpsilonB, CapUpsilonM, SheetUpsilonM, Arm1)
	tmp2 = fmt.Sprintf("((%s + %s) / (%s * %s)) * %s", CapUpsilonP, SheetUpsilonP, Lp, Bp, Arm1)
	tmp3 = fmt.Sprintf("%s + %s + (%s + %s) * (%s)^2", UpsilonB, UpsilonP, CapUpsilonM, SheetUpsilonM, Arm1)
	// Коэффициент податливости фланцевого соединения крышки и решетки
	Bolt.Eta = fmt.Sprintf("(%s + %s) / %s", tmp1, tmp2, tmp3)

	// Fв - Расчетное усилие в болтах (шпильках) в условиях эксплуатации
	Bolt.WorkEffort = fmt.Sprintf("%s * (%s * %s + 2 * %s * %s * (%s + %s))", pressure, Lp, Bp, EstimatedGasketWidth, m, Lp, Bp)

	//TODO в оригинале почему-то тут не WorkEffort, а площадь и количество болтов
	// tmp1 = (result.Calc.Pressure / data.Pressure) * Bolt.WorkEffort
	// tmp2 = result.Calc.Pressure * (Bolt.Eta*Bolt.Lp*Auxiliary.Bp + 2*Auxiliary.EstimatedGasketWidth*d.Gasket.M*(Bolt.Lp+Auxiliary.Bp))
	// F0 - Расчетное усилие в болтах (шпильках) в условиях испытаний или монтажа
	// Bolt.Effort = math.Max(tmp1, tmp2)

	// Ab := d.Bolt.Area * float64(d.Bolt.Count)
	// Условия прочности болтов/шпилек - в условиях испытания или монтажа
	// Bolt.TestCond = &dev_cooling_model.Condition{
	// 	X: Bolt.Effort / Ab,
	// 	Y: d.Bolt.SigmaAt20,
	// }
	// // Условия прочности болтов/шпилек - в условиях эксплуатации
	// Bolt.WorkCond = &dev_cooling_model.Condition{
	// 	X: Bolt.WorkEffort / Ab,
	// 	Y: d.Bolt.Sigma,
	// }

	return Bolt
}
