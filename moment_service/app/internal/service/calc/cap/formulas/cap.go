package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (s *FormulasService) getCapFormulas(
	capType moment_api.CapData_Type,
	data *moment_api.CapResult,
	h, D, S0, DOut, Dcp string,
) *moment_api.CapFormulas {
	f := &moment_api.CapFormulas{}

	k := strings.ReplaceAll(strconv.FormatFloat(data.K, 'G', 3, 64), "E", "*10^")
	x := strings.ReplaceAll(strconv.FormatFloat(data.X, 'G', 3, 64), "E", "*10^")
	eAt20 := strings.ReplaceAll(strconv.FormatFloat(data.EpsilonAt20, 'G', 3, 64), "E", "*10^")

	if capType == moment_api.CapData_flat {
		H := strings.ReplaceAll(strconv.FormatFloat(data.H, 'G', 3, 64), "E", "*10^")
		delta := strings.ReplaceAll(strconv.FormatFloat(data.Delta, 'G', 3, 64), "E", "*10^")

		f.K = fmt.Sprintf("%s / %s", DOut, Dcp)
		f.X = fmt.Sprintf("0.67 * (%s^2 + (1 + 8.55 * lg(%s) - 1)) / ((%s - 1) * %s^2 - 1 + (1.857 * %s^2 + 1) * %s^3/%s^3)",
			k, k, k, k, k, H, delta)
		f.Y = fmt.Sprintf("%s / (%s * %s)", x, delta, eAt20)
	} else {
		radius := strings.ReplaceAll(strconv.FormatFloat(data.Radius, 'G', 3, 64), "E", "*10^")
		lambda := strings.ReplaceAll(strconv.FormatFloat(data.Lambda, 'G', 3, 64), "E", "*10^")
		omega := strings.ReplaceAll(strconv.FormatFloat(data.Omega, 'G', 3, 64), "E", "*10^")

		f.Lambda = fmt.Sprintf("(%s / %s) * Sqrt(%s / %s)", h, D, radius, S0)
		f.Omega = fmt.Sprintf("1 / (1 + 1.285*%s + 1.63*%s * (%s/%s)^2 * lg(%s/%s)", lambda, lambda, h, S0, DOut, D)
		f.Y = fmt.Sprintf("((1 - %s * (1 + 1.285*%s)) / (%s * %s^3)) * ((%s + %s) / (%s - %s))", omega, lambda, eAt20, h, DOut, D, DOut, D)
	}

	return f
}
