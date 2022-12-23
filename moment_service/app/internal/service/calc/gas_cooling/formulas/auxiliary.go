package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
)

func (s *FormulasService) getAuxiliaryFormulas(req *calc_api.GasCoolingRequest, data models.DataGasCooling, result *calc_api.GasCoolingResponse,
) *gas_cooling_model.AuxiliaryFormulas {
	Auxiliary := &gas_cooling_model.AuxiliaryFormulas{}

	// перевод чисел в строки
	width := strconv.FormatFloat(data.Gasket.Width, 'G', 3, 64)
	sizeTrans := strconv.FormatFloat(data.Gasket.SizeTrans, 'G', 3, 64)
	sizeLong := strconv.FormatFloat(data.Gasket.SizeLong, 'G', 3, 64)

	EstimatedGasketWidth := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.EstimatedGasketWidth, 'G', 3, 64), "E", "*10^")

	// расчетная ширина плоской прокладки
	Auxiliary.EstimatedGasketWidth = fmt.Sprintf("min(%s; 3.87 * sqrt(%s))", width, width)
	// Bp - расчетный размер решетки в поперечном направлении
	Auxiliary.SizeTrans = fmt.Sprintf("%s - %s", sizeTrans, EstimatedGasketWidth)
	// Lp - Расчетный размер решетки в продольном направлении
	Auxiliary.SizeLong = fmt.Sprintf("%s - %s", sizeLong, EstimatedGasketWidth)

	return Auxiliary
}
