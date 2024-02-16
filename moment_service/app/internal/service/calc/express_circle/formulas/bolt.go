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

func (s *FormulasService) getBoltFormulas(req *calc_api.ExpressCircleRequest, d models.DataExCircle, result *calc_api.ExpressCircleResponse,
) *express_circle_model.BoltsFormulas {
	Bolts := &express_circle_model.BoltsFormulas{}

	// перевод чисел в строки
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 5, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 5, 64)

	Diameter := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Deformation.Diameter, 'G', 5, 64), "E", "*10^")

	Area := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.Area, 'G', 5, 64), "E", "*10^")
	DesignLoad := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.DesignLoad, 'G', 5, 64), "E", "*10^")
	WorkDesignLoad := strings.ReplaceAll(strconv.FormatFloat(result.Calc.ForcesInBolts.WorkDesignLoad, 'G', 5, 64), "E", "*10^")

	Kyp := s.Kyp[true]
	Kyz := s.Kyz[req.Condition.String()]
	Kyt := constants.LoadKyt

	// Расчетное напряжение в болтах/шпильках - при затяжке
	Bolts.RatedStress = fmt.Sprintf("%s / %s", DesignLoad, Area)
	// Условия прочности болтов шпилек - при затяжке
	Bolts.AllowableVoltage = fmt.Sprintf("1.2 * %.f * %.1f * %.1f * %s", Kyp, Kyz, Kyt, sigmaAt20)

	if d.TypeGasket == express_circle_model.GasketData_Soft {
		// Условие прочности прокладки
		Bolts.StrengthGasket = fmt.Sprintf("max(%s; %s) / (%f * %s * %s)", DesignLoad, WorkDesignLoad, math.Pi, Diameter, width)
	}

	return Bolts
}
