package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

func (s *FormulasService) getStrengthFormulas(
	typeF flange_model.FlangeData_Type,
	AxialForce, BendingMoment int32,
	data models.DataFlange,
	flange *flange_model.FlangeResult,
	res *flange_model.StrengthResult,
	D6, Dcp, Pbm, Pbr, Qd, QFM, pressure string,
	isWork, isTemp bool,
) *flange_model.AddStrengthFormulas {
	strength := &flange_model.AddStrengthFormulas{}

	var Ks float64
	if flange.K <= constants.MinK {
		Ks = constants.MinKs
	} else if flange.K >= constants.MaxK {
		Ks = constants.MaxKs
	} else {
		Ks = ((flange.K-constants.MinK)/(constants.MaxK-constants.MinK))*(constants.MaxKs-constants.MinKs) + constants.MinKs
	}
	Kt := map[bool]float64{
		true:  constants.TempKt,
		false: constants.Kt,
	}

	count := data.Bolt.Count
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 3, 64)
	h := strings.ReplaceAll(strconv.FormatFloat(flange.H, 'G', 3, 64), "E", "*10^")
	m := strconv.FormatFloat(data.Gasket.M, 'G', 3, 64)
	b := strings.ReplaceAll(strconv.FormatFloat(flange.B, 'G', 3, 64), "E", "*10^")
	e := strings.ReplaceAll(strconv.FormatFloat(flange.E, 'G', 3, 64), "E", "*10^")
	a := strings.ReplaceAll(strconv.FormatFloat(flange.A, 'G', 3, 64), "E", "*10^")
	cf := strings.ReplaceAll(strconv.FormatFloat(res.Cf, 'G', 3, 64), "E", "*10^")

	temp1 := fmt.Sprintf("%f * %s/%d", math.Pi, D6, count)
	temp2 := fmt.Sprintf("2*%s + 6*%s/(%s + 0.5)", diameter, h, m)

	strength.Cf = fmt.Sprintf("max(1; Sqrt((%s) / (%s)))", temp1, temp2)

	strength.MM = fmt.Sprintf("%s * %s * %s", cf, Pbm, b)
	strength.Mp = fmt.Sprintf("%s * max(%s * %s + (%s + %s)*%s; |%s + %s|*%s)", cf, Pbr, b, Qd, QFM, e, Qd, QFM, e)

	if typeF == flange_model.FlangeData_free {
		strength.MMk = fmt.Sprintf("%s * %s * %s", cf, Pbm, a)
		strength.Mpk = fmt.Sprintf("%s * %s * %s", cf, Pbr, a)
	}

	mm := strings.ReplaceAll(strconv.FormatFloat(res.MM, 'G', 3, 64), "E", "*10^")
	dvz := strings.ReplaceAll(strconv.FormatFloat(res.Dzv, 'G', 3, 64), "E", "*10^")
	sigmaM1 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaM1, 'G', 3, 64), "E", "*10^")
	s1 := strings.ReplaceAll(strconv.FormatFloat(flange.S1, 'G', 3, 64), "E", "*10^")
	s0 := strings.ReplaceAll(strconv.FormatFloat(flange.S0, 'G', 3, 64), "E", "*10^")
	c := strings.ReplaceAll(strconv.FormatFloat(flange.C, 'G', 3, 64), "E", "*10^")
	f := strings.ReplaceAll(strconv.FormatFloat(flange.F, 'G', 3, 64), "E", "*10^")
	l0 := strings.ReplaceAll(strconv.FormatFloat(flange.L0, 'G', 3, 64), "E", "*10^")
	d := strings.ReplaceAll(strconv.FormatFloat(flange.D, 'G', 3, 64), "E", "*10^")
	lymda := strings.ReplaceAll(strconv.FormatFloat(flange.Lymda, 'G', 3, 64), "E", "*10^")
	betaF := strings.ReplaceAll(strconv.FormatFloat(flange.BetaF, 'G', 3, 64), "E", "*10^")
	betaY := strings.ReplaceAll(strconv.FormatFloat(flange.BetaY, 'G', 3, 64), "E", "*10^")
	betaZ := strings.ReplaceAll(strconv.FormatFloat(flange.BetaZ, 'G', 3, 64), "E", "*10^")
	sigmaR := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaR, 'G', 3, 64), "E", "*10^")

	if typeF == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		strength.SigmaM1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mm, lymda, s1, c, dvz)
		strength.SigmaM0 = fmt.Sprintf("%s * %s", f, sigmaM1)
	} else {
		strength.SigmaM1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mm, lymda, s0, c, dvz)
	}

	strength.SigmaR = fmt.Sprintf("((1.33 * %s * %s + %s)/(%s * (%s)^2 * %s * %s)) * %s",
		betaF, h, l0, lymda, h, l0, d, mm)
	strength.SigmaT = fmt.Sprintf("(%s * %s)/((%s)^2 * %s) - %s * %s", betaY, mm, h, d, betaZ, sigmaR)

	hk := strings.ReplaceAll(strconv.FormatFloat(flange.Hk, 'G', 3, 64), "E", "*10^")
	dk := strings.ReplaceAll(strconv.FormatFloat(flange.Dk, 'G', 3, 64), "E", "*10^")
	mmk := strings.ReplaceAll(strconv.FormatFloat(res.MMk, 'G', 3, 64), "E", "*10^")
	mp := strings.ReplaceAll(strconv.FormatFloat(res.Mp, 'G', 3, 64), "E", "*10^")
	sigmaP1 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaP1, 'G', 3, 64), "E", "*10^")

	if typeF == flange_model.FlangeData_free {
		strength.SigmaK = fmt.Sprintf("%s * %s / (%s)^2 * %s", betaY, mmk, hk, dk)
	}

	if typeF == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		strength.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mp, lymda, s1, c, dvz)
		strength.SigmaP0 = fmt.Sprintf("%s * %s", f, sigmaP1)
	} else {
		strength.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mp, lymda, s0, c, dvz)
	}

	if typeF == flange_model.FlangeData_welded {
		temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, d, s1, s1, c)
		strength.SigmaMp = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d + (4 * |%d|)/(%s + %s)) / %s",
			d, pressure, AxialForce, BendingMoment, d, s1, temp)
		strength.SigmaMpm = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d - (4 * |%d|)/(%s + %s)) / %s",
			d, pressure, AxialForce, BendingMoment, d, s1, temp)
	}

	temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, d, s0, s0, c)
	strength.SigmaMp0 = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d + (4 * |%d|)/(%s + %s)) / %s",
		d, pressure, AxialForce, BendingMoment, d, s0, temp)
	strength.SigmaMpm0 = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d - (4 * |%d|)/(%s + %s)) / %s",
		d, pressure, AxialForce, BendingMoment, d, s0, temp)

	strength.SigmaMop = fmt.Sprintf("%s * %s / (2 * (%s - %s))", pressure, d, s0, c)

	sigmaRp := strings.ReplaceAll(strconv.FormatFloat(res.SigmaRp, 'G', 3, 64), "E", "*10^")
	sigmaTp := strings.ReplaceAll(strconv.FormatFloat(res.SigmaTp, 'G', 3, 64), "E", "*10^")

	strength.SigmaRp = fmt.Sprintf("((1.33 * %s * %s + %s)/(%s * (%s)^2 * %s * %s)) * %s", betaF, h, l0, lymda, h, l0, d, mp)
	strength.SigmaTp = fmt.Sprintf("(%s * %s)/((%s)^2 * %s) - %s * %s", betaY, mp, h, d, betaZ, sigmaRp)

	if typeF == flange_model.FlangeData_free {
		strength.SigmaKp = fmt.Sprintf("(%s * %s)/((%s)^2 * %s)", betaY, mp, hk, dk)
	}

	epsilon := strings.ReplaceAll(strconv.FormatFloat(flange.Epsilon, 'G', 3, 64), "E", "*10^")
	epsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	yf := strings.ReplaceAll(strconv.FormatFloat(flange.Yf, 'G', 3, 64), "E", "*10^")

	strength.Teta = fmt.Sprintf("%s * %s * %s/%s", mp, yf, epsilonAt20, epsilon)

	if typeF == flange_model.FlangeData_free {
		epsilonK := strings.ReplaceAll(strconv.FormatFloat(flange.EpsilonK, 'G', 3, 64), "E", "*10^")
		epsilonKAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.EpsilonKAt20, 'G', 3, 64), "E", "*10^")
		yk := strings.ReplaceAll(strconv.FormatFloat(flange.Yk, 'G', 3, 64), "E", "*10^")
		mpk := strings.ReplaceAll(strconv.FormatFloat(res.Mpk, 'G', 3, 64), "E", "*10^")

		strength.Teta = fmt.Sprintf("%s * %s * %s/%s", mpk, yk, epsilonKAt20, epsilonK)
	}

	kt := strconv.FormatFloat(Kt[isTemp], 'G', 3, 64)
	sigmaT := strings.ReplaceAll(strconv.FormatFloat(res.SigmaT, 'G', 3, 64), "E", "*10^")
	sigmaM0 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaM0, 'G', 3, 64), "E", "*10^")
	sigmaP0 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaP0, 'G', 3, 64), "E", "*10^")
	sigmaMp0 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaMp0, 'G', 3, 64), "E", "*10^")
	sigmaMpm0 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaMpm0, 'G', 3, 64), "E", "*10^")
	sigmaMop := strings.ReplaceAll(strconv.FormatFloat(res.SigmaMop, 'G', 3, 64), "E", "*10^")
	sigmaAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaAt20, 'G', 3, 64), "E", "*10^")
	sigma := strings.ReplaceAll(strconv.FormatFloat(flange.Sigma, 'G', 3, 64), "E", "*10^")

	if typeF == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		sigmaM := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaM, 'G', 3, 64), "E", "*10^")
		sigmaMAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaMAt20, 'G', 3, 64), "E", "*10^")
		sigmaRAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaRAt20, 'G', 3, 64), "E", "*10^")
		sigmaM1 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaM1, 'G', 3, 64), "E", "*10^")
		sigmaMp := strings.ReplaceAll(strconv.FormatFloat(res.SigmaMp, 'G', 3, 64), "E", "*10^")
		sigmaMpm := strings.ReplaceAll(strconv.FormatFloat(res.SigmaMpm, 'G', 3, 64), "E", "*10^")
		sigmaP1 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaP1, 'G', 3, 64), "E", "*10^")
		ks := strconv.FormatFloat(Ks, 'G', 3, 64)

		strength.Max1 = fmt.Sprintf(`max(|%s + %s|; |%s + %s|) ≤ %s * %s * %s`, sigmaM1, sigmaR, sigmaM1, sigmaT, ks, kt, sigmaMAt20)
		strength.Max2 = fmt.Sprintf(`max(|%s-%s+%s|;|%s-%s+%s|;|%s-%s+%s|;|%s-%s+%s|;|%s+%s|;|%s+%s|) ≤ %s*%s*%s`,
			sigmaP1, sigmaMp, sigmaRp, sigmaP1, sigmaMpm, sigmaRp,
			sigmaP1, sigmaMp, sigmaTp, sigmaP1, sigmaMpm, sigmaTp,
			sigmaP1, sigmaMp, sigmaP1, sigmaMpm, ks, kt, sigmaM,
		)

		strength.Max3 = fmt.Sprintf("%s ≤ 1.3 * %s", sigmaM0, sigmaRAt20)
		strength.Max4 = fmt.Sprintf(`max(|%s+%s|;|%s-%s|;|%s+%s|;|%s-%s|;|0.3*%s+%s|;|0.3*%s-%s|;|0.7*%s+(%s-%s)|;|0.7*%s-(%s-%s)|;|0.7*%s+(%s-%s)|;|0.7*%s-(%s-%s)|) ≤ 1.3*%s`,
			sigmaP0, sigmaMp0, sigmaP0, sigmaMp0,
			sigmaP0, sigmaMpm0, sigmaP0, sigmaMpm0,
			sigmaP0, sigmaMop, sigmaP0, sigmaMop,
			sigmaP0, sigmaMp0, sigmaMop, sigmaP0, sigmaMp0, sigmaMop,
			sigmaP0, sigmaMpm0, sigmaMop, sigmaP0, sigmaMpm0, sigmaMop, sigmaR,
		)
	} else {
		strength.Max5 = fmt.Sprintf("max(|%s+%s|;|%s+%s|) ≤ %s", sigmaM0, sigmaR, sigmaM0, sigmaT, sigmaAt20)
		strength.Max6 = fmt.Sprintf("max(|%s-%s+%s|;|%s-%s+%s|;|%s-%s+%s|;|%s-%s+%s|;|%s+%s|;|%s+%s|) ≤ %s",
			sigmaP0, sigmaMp0, sigmaTp, sigmaP0, sigmaMpm0, sigmaTp,
			sigmaP0, sigmaMp0, sigmaRp, sigmaP0, sigmaMpm0, sigmaRp,
			sigmaP0, sigmaMp0, sigmaP0, sigmaMpm0, sigma,
		)
	}

	strength.Max7 = fmt.Sprintf("max(|%s|;|%s|;|%s|) ≤ %s", sigmaMp0, sigmaMpm0, sigmaMop, sigma)
	strength.Max8 = fmt.Sprintf("max(|%s|; |%s|) ≤ %s * %s", sigmaR, sigmaT, kt, sigma)
	strength.Max9 = fmt.Sprintf("max(|%s|; |%s|) ≤ %s * %s", sigmaRp, sigmaTp, kt, sigmaAt20)

	if typeF == flange_model.FlangeData_free {
		sigmaK := strings.ReplaceAll(strconv.FormatFloat(res.SigmaK, 'G', 3, 64), "E", "*10^")
		sigmaKp := strings.ReplaceAll(strconv.FormatFloat(res.SigmaKp, 'G', 3, 64), "E", "*10^")
		sigmaKAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaKAt20, 'G', 3, 64), "E", "*10^")
		sigma := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaK, 'G', 3, 64), "E", "*10^")

		strength.Max10 = fmt.Sprintf("%s ≤ %s * %s", sigmaK, kt, sigmaKAt20)
		strength.Max11 = fmt.Sprintf("%s ≤ %s * %s", sigmaKp, kt, sigma)
	}

	return strength
}
