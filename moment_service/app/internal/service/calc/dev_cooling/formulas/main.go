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

type FormulasService struct {
	typeBolt map[string]float64
}

func NewFormulasService() *FormulasService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	return &FormulasService{
		typeBolt: bolt,
	}
}

func (s *FormulasService) GetFormulas(
	Ab, Lambda1, Lambda2, Alpha1, Alpha2 float64,
	req calc_api.DevCoolingRequest,
	d models.DataDevCooling,
	result calc_api.DevCoolingResponse,
) *dev_cooling_model.Formulas {
	formulas := &dev_cooling_model.Formulas{
		Auxiliary: s.getAuxiliaryFormulas(req, d, result),
		Bolt:      s.getBoltFormulas(Lambda1, Lambda2, Alpha1, Alpha2, req, d, result),
		TubeSheet: s.getTubeSheetFormulas(req, d, result),
		Cap:       s.getCapFormulas(req, d, result),
		Moment:    s.getMomentFormulas(req, d, result),
	}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Pressure, 'G', 3, 64)

	zoneThick := strconv.FormatFloat(d.TubeSheet.ZoneThick, 'G', 3, 64)
	corrosion := strconv.FormatFloat(d.TubeSheet.Corrosion, 'G', 3, 64)
	tsSigmaAt20 := strconv.FormatFloat(d.TubeSheet.SigmaAt20, 'G', 3, 64)
	tsSigma := strconv.FormatFloat(d.TubeSheet.Sigma, 'G', 3, 64)

	tSigmaAt20 := strconv.FormatFloat(d.Tube.SigmaAt20, 'G', 3, 64)
	tSigma := strconv.FormatFloat(d.Tube.Sigma, 'G', 3, 64)

	bottomThick := strconv.FormatFloat(d.Cap.BottomThick, 'G', 3, 64)
	capCorrosion := strconv.FormatFloat(d.Cap.Corrosion, 'G', 3, 64)
	capSigmaAt20 := strconv.FormatFloat(d.Cap.SigmaAt20, 'G', 3, 64)
	capSigma := strconv.FormatFloat(d.Cap.Sigma, 'G', 3, 64)

	bSigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)
	bSigma := strconv.FormatFloat(d.Bolt.Sigma, 'G', 3, 64)

	gWidth := strconv.FormatFloat(d.Gasket.Width, 'G', 3, 64)

	Bp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Auxiliary.Bp, 'G', 3, 64), "E", "*10^")

	Lp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Lp, 'G', 3, 64), "E", "*10^")
	WorkEffort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.WorkEffort, 'G', 3, 64), "E", "*10^")
	Effort := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Bolt.Effort, 'G', 3, 64), "E", "*10^")

	// Условия применения формул
	formulas.Condition1 = fmt.Sprintf("(%s - %s) / %s", zoneThick, corrosion, Bp)
	formulas.Condition2 = fmt.Sprintf("(%s - %s) / %s", bottomThick, capCorrosion, Bp)

	cond1 := fmt.Sprintf("(%s / %s); (%s / %s)", capSigmaAt20, capSigma, tsSigmaAt20, tsSigma)
	cond2 := fmt.Sprintf("(%s / %s); (%s / %s)", tSigmaAt20, tSigma, bSigmaAt20, bSigma)
	// Пробное давление
	formulas.Pressure = fmt.Sprintf("1.25 * %s * min(%s; %s)", pressure, cond1, cond2)

	// Условие прочности прокладки
	formulas.GasketCond = fmt.Sprintf("max(%s; %s) / (2 * (%s + %s) * %s)", WorkEffort, Effort, Lp, Bp, gWidth)

	return formulas
}
