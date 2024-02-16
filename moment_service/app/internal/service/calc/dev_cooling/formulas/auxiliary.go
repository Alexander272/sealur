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

func (s *FormulasService) getAuxiliaryFormulas(
	data *calc_api.DevCoolingRequest,
	d models.DataDevCooling,
	result *calc_api.DevCoolingResponse,
) *dev_cooling_model.AuxiliaryFormulas {
	Auxiliary := &dev_cooling_model.AuxiliaryFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(data.Pressure, 'G', 5, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 5, 64)
	sizeTrans := strconv.FormatFloat(d.Gasket.SizeTrans, 'G', 5, 64)

	count := d.TubeSheet.Count
	stepTrans := strconv.FormatFloat(d.TubeSheet.StepTrans, 'G', 5, 64)
	stepLong := strconv.FormatFloat(d.TubeSheet.StepLong, 'G', 5, 64)
	tsDiameter := strconv.FormatFloat(d.TubeSheet.Diameter, 'G', 5, 64)
	tsSigma := strconv.FormatFloat(d.TubeSheet.Sigma, 'G', 5, 64)

	diameter := strconv.FormatFloat(d.Tube.Diameter, 'G', 5, 64)
	thickness := strconv.FormatFloat(d.Tube.Thickness, 'G', 5, 64)
	corrosion := strconv.FormatFloat(d.Tube.Corrosion, 'G', 5, 64)
	depth := strconv.FormatFloat(d.Tube.Depth, 'G', 5, 64)
	size := strconv.FormatFloat(d.Tube.Size, 'G', 5, 64)
	sigma := strconv.FormatFloat(d.Tube.Sigma, 'G', 5, 64)
	epsilon := strconv.FormatFloat(d.Tube.Epsilon, 'G', 5, 64)
	reducedLength := strconv.FormatFloat(d.Tube.ReducedLength, 'G', 5, 64)

	distance := strconv.FormatFloat(d.Bolt.Distance, 'G', 5, 64)

	EstimatedGasketWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.EstimatedGasketWidth, 'G', 5, 64), "E", "*10^")
	EstimatedZoneWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.EstimatedZoneWidth, 'G', 5, 64), "E", "*10^")
	Bp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Bp, 'G', 5, 64), "E", "*10^")
	D := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.D, 'G', 5, 64), "E", "*10^")
	Upsilon := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Upsilon, 'G', 5, 64), "E", "*10^")
	Mu := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Mu, 'G', 5, 64), "E", "*10^")

	// расчетная ширина плоской прокладки
	Auxiliary.EstimatedGasketWidth = fmt.Sprintf("min(%s; 3.87 * sqrt(%s))", width, width)
	// расчетный размер решетки в поперечном направлении
	Auxiliary.Bp = fmt.Sprintf("%s - %s", sizeTrans, EstimatedGasketWidth)
	// расчетная ширина перфорированной зоны решетки
	Auxiliary.EstimatedZoneWidth = fmt.Sprintf("min(%d*%s; %s)", count, stepTrans, Bp)
	// относительная ширина беструбной зоны решетки
	Auxiliary.RelativeWidth = fmt.Sprintf("(%s - %s) / %s", Bp, EstimatedZoneWidth, EstimatedZoneWidth)

	// Вспомогательные коэффициенты
	Auxiliary.Upsilon = fmt.Sprintf("(%f * (%s - %s) * (%s - %s)) / (%s * %s)",
		math.Pi, diameter, thickness, thickness, corrosion, stepLong, stepTrans)
	Auxiliary.Eta = fmt.Sprintf("1 - (%f/4)*((%s - 2 * %s)^2 / (%s * %s))", math.Pi, diameter, thickness, stepLong, stepTrans)

	// эффективный диаметр отверстия решетки или задней стенке
	if data.Method == calc_api.DevCoolingRequest_AllThickness {
		Auxiliary.D = fmt.Sprintf("%s - 2 * %s", tsDiameter, thickness)
	} else if data.Method == calc_api.DevCoolingRequest_PartThickness {
		Auxiliary.D = fmt.Sprintf("%s - %s", tsDiameter, thickness)
	}

	// коэффициент ослабления решетки и задней стенки
	Auxiliary.Phi = fmt.Sprintf("1 - %s /%s", D, stepLong)
	// допускаемая нагрузка из условия прочности труб
	Auxiliary.LoadTube = fmt.Sprintf("%s * (1 - ((%s - %s)/(2 * (%s - %s))) * (%s / %s)) * %s",
		Upsilon, diameter, thickness, thickness, corrosion, pressure, sigma, sigma)

	// допускаемое напряжение из условия прочности крепления трубы в решетке
	if data.Mounting == calc_api.DevCoolingRequest_flaring {
		Auxiliary.Load = fmt.Sprintf("%s * %s * ((2 * %s) / (%s * %s)) * %s", Upsilon, Mu, depth, diameter, thickness, sigma)
	} else if data.Mounting == calc_api.DevCoolingRequest_welding {
		Auxiliary.Load = fmt.Sprintf("0.7 * %s * (%s / %s) * min(%s; %s)", Upsilon, size, thickness, sigma, tsSigma)
	} else {
		Auxiliary.Load = fmt.Sprintf("0.7*%s * (%s / %s) * min(%s; %s) + 0.6*(%s * %s * ((2 * %s)/(%s * %s)) * %s)",
			Upsilon, size, thickness, sigma, tsSigma, Upsilon, Mu, depth, diameter, thickness, sigma)
	}

	// коэффициент уменьшения допускаемых напряжений при продольном изгибе
	Auxiliary.PhiT = fmt.Sprintf("1 / sqrt(1 +[1.8*(%s / %s) * (%s / (%s - %s))^2]^2)", sigma, epsilon, reducedLength, diameter, thickness)

	// l1, l2 - Плечи изгибающих моментов
	Auxiliary.Arm1 = fmt.Sprintf("0.5 * (%s - %s)", distance, Bp)
	Auxiliary.Arm2 = fmt.Sprintf("0.5 * (%s - %s)", distance, sizeTrans)

	return Auxiliary
}
