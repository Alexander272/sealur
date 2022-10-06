package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

func (s *FormulasService) getFlangeFormulas(
	typeF flange_model.FlangeData_Type,
	data *flange_model.FlangeResult,
	D6, DOut, Dcp string,
) *flange_model.FlangeFormulas {
	f := &flange_model.FlangeFormulas{}

	dk := strings.ReplaceAll(strconv.FormatFloat(data.Dk, 'G', 3, 64), "E", "*10^")
	ds := strings.ReplaceAll(strconv.FormatFloat(data.Ds, 'G', 3, 64), "E", "*10^")
	dnk := strings.ReplaceAll(strconv.FormatFloat(data.Dnk, 'G', 3, 64), "E", "*10^")
	h0 := strings.ReplaceAll(strconv.FormatFloat(data.H0, 'G', 3, 64), "E", "*10^")
	l := strings.ReplaceAll(strconv.FormatFloat(data.L, 'G', 3, 64), "E", "*10^")
	d := strings.ReplaceAll(strconv.FormatFloat(data.D, 'G', 3, 64), "E", "*10^")
	s0 := strings.ReplaceAll(strconv.FormatFloat(data.S0, 'G', 3, 64), "E", "*10^")
	s1 := strings.ReplaceAll(strconv.FormatFloat(data.S1, 'G', 3, 64), "E", "*10^")
	beta := strings.ReplaceAll(strconv.FormatFloat(data.Beta, 'G', 3, 64), "E", "*10^")
	x := strings.ReplaceAll(strconv.FormatFloat(data.X, 'G', 3, 64), "E", "*10^")
	xi := strings.ReplaceAll(strconv.FormatFloat(data.Xi, 'G', 3, 64), "E", "*10^")
	Se := strings.ReplaceAll(strconv.FormatFloat(data.Se, 'G', 3, 64), "E", "*10^")

	h := strings.ReplaceAll(strconv.FormatFloat(data.H, 'G', 3, 64), "E", "*10^")
	l0 := strings.ReplaceAll(strconv.FormatFloat(data.L0, 'G', 3, 64), "E", "*10^")
	betaF := strings.ReplaceAll(strconv.FormatFloat(data.BetaF, 'G', 3, 64), "E", "*10^")
	betaT := strings.ReplaceAll(strconv.FormatFloat(data.BetaT, 'G', 3, 64), "E", "*10^")
	betaU := strings.ReplaceAll(strconv.FormatFloat(data.BetaU, 'G', 3, 64), "E", "*10^")
	betaV := strings.ReplaceAll(strconv.FormatFloat(data.BetaV, 'G', 3, 64), "E", "*10^")
	epsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	epsilonKAt20 := strings.ReplaceAll(strconv.FormatFloat(data.EpsilonKAt20, 'G', 3, 64), "E", "*10^")
	lymda := strings.ReplaceAll(strconv.FormatFloat(data.Lymda, 'G', 3, 64), "E", "*10^")
	hk := strings.ReplaceAll(strconv.FormatFloat(data.Hk, 'G', 3, 64), "E", "*10^")
	psik := strings.ReplaceAll(strconv.FormatFloat(data.Psik, 'G', 3, 64), "E", "*10^")

	if typeF != flange_model.FlangeData_free {
		f.B = fmt.Sprintf("0.5 * (%s - %s)", D6, Dcp)
	} else {
		Ds := fmt.Sprintf("0.5 * (%s + %s + 2*%s)", DOut, dk, h0)
		f.A = fmt.Sprintf("0.5 * (%s - %s)", D6, Ds)
		f.B = fmt.Sprintf("0.5 * (%s - %s)", Ds, Dcp)
	}

	if typeF == flange_model.FlangeData_welded {
		f.X = fmt.Sprintf("%s / Sqrt(%s * %s)", l, d, s0)
		f.Beta = fmt.Sprintf("%s / %s", s1, s0)
		f.Xi = fmt.Sprintf("1 + (%s - 1) * %s / (%s + (1 + %s)/4)", beta, x, x, beta)
		f.Se = fmt.Sprintf("%s * %s", xi, s0)
	}

	f.E = fmt.Sprintf("0.5 * (%s - %s - %s)", Dcp, d, Se)
	f.L0 = fmt.Sprintf("Sqrt(%s * %s)", d, s0)
	f.K = fmt.Sprintf("%s / %s", DOut, d)

	f.Lymda = fmt.Sprintf("(%s * %s + %s)/(%s + %s) + (%s * %s^3)/(%s * %s * %s^2)",
		betaF, h, l0, betaT, l0, betaV, h, betaU, l0, s0)
	f.Yf = fmt.Sprintf("(0.91 * %s)/(%s * %s * %s^2) * %s", betaV, epsilonAt20, lymda, s0, l0)

	if typeF == flange_model.FlangeData_free {
		f.Psik = fmt.Sprintf("1.28 * (log(%s/%s) / log(10))", dnk, dk)
		f.Yk = fmt.Sprintf("1 / (%s * %s^3 * %s)", epsilonKAt20, hk, psik)
	}

	if typeF != flange_model.FlangeData_free {
		f.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * %s^3))", math.Pi, D6, epsilonAt20, DOut, h)
	} else {
		f.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * %s^3))", math.Pi, ds, epsilonAt20, DOut, h)
		f.Yfc = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * %s^3))", math.Pi, D6, epsilonKAt20, dnk, hk)
	}

	return f
}
