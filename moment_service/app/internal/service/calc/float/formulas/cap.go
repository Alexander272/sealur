package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *FormulasService) getCapFormulas(flange *float_model.FlangeResult, cap *float_model.CapResult) *float_model.CapFormulas {
	formulas := &float_model.CapFormulas{}

	h := strings.ReplaceAll(strconv.FormatFloat(cap.H, 'G', 3, 64), "E", "*10^")
	cs := strings.ReplaceAll(strconv.FormatFloat(cap.S, 'G', 3, 64), "E", "*10^")
	radius := strings.ReplaceAll(strconv.FormatFloat(cap.Radius, 'G', 3, 64), "E", "*10^")
	labmda := strings.ReplaceAll(strconv.FormatFloat(cap.Lambda, 'G', 3, 64), "E", "*10^")
	omega := strings.ReplaceAll(strconv.FormatFloat(cap.Omega, 'G', 3, 64), "E", "*10^")
	eAt20 := strings.ReplaceAll(strconv.FormatFloat(cap.EpsilonAt20, 'G', 3, 64), "E", "*10^")

	d := strings.ReplaceAll(strconv.FormatFloat(flange.D, 'G', 3, 64), "E", "*10^")
	dOut := strings.ReplaceAll(strconv.FormatFloat(flange.DOut, 'G', 3, 64), "E", "*10^")

	formulas.Lambda = fmt.Sprintf("(%s / %s) * sqrt(%s / %s)", h, d, radius, cs)
	formulas.Omega = fmt.Sprintf("1 / (1 + 1.285*%s + 1.63*%s*((%s/%s)^2)*lg(%s/%s))", labmda, labmda, h, d, dOut, d)
	formulas.Y = fmt.Sprintf("((1 - %s*(1+1.285*%s)) / (%s * %s^3)) * ((%s + %s) / (%s - %s))", omega, labmda, eAt20, h, dOut, d, dOut, d)

	return formulas
}
