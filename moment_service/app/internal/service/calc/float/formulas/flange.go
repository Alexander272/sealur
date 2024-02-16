package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *FormulasService) getFlangeFormulas(flange *float_model.FlangeResult, cap *float_model.CapResult, Dcp string,
) *float_model.FlangeFormulas {
	formulas := &float_model.FlangeFormulas{}

	d6 := strings.ReplaceAll(strconv.FormatFloat(flange.D6, 'G', 5, 64), "E", "*10^")
	d := strings.ReplaceAll(strconv.FormatFloat(flange.D, 'G', 5, 64), "E", "*10^")
	cs := strings.ReplaceAll(strconv.FormatFloat(cap.S, 'G', 5, 64), "E", "*10^")

	formulas.B = fmt.Sprintf("0.5 * (%s - %s)", d6, Dcp)
	formulas.L0 = fmt.Sprintf("sqrt(%s * %s)", d, cs)
	// formulas.Y = 0

	return formulas
}
