package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type CapService struct {
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewCapService() *CapService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}
	kp := map[bool]float64{
		true:  constants.WorkKyp,
		false: constants.TestKyp,
	}
	kz := map[string]float64{
		"uncontrollable":  constants.UncontrollableKyz,
		"controllable":    constants.ControllableKyz,
		"controllablePin": constants.ControllablePinKyz,
	}

	return &CapService{
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

func (s *CapService) GetFormulasForCap(
	TypeGasket, Condition string, IsWork, IsUseWasher, IsEmbedded bool,
	data models.DataCap,
	result moment_api.CapResponse,
	calculation moment_api.CalcCapRequest_Calcutation,
	gamma_, yb_, yp_ float64,
) *moment_api.CalcCapFormulas {
	formulas := &moment_api.CalcCapFormulas{}

	width := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Width, 'G', 3, 64), "E", "*10^")
	DOut := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.DOut, 'G', 3, 64), "E", "*10^")
	b0 := strings.ReplaceAll(strconv.FormatFloat(data.B0, 'G', 3, 64), "E", "*10^")
	th := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Thickness, 'G', 3, 64), "E", "*10^")
	compression := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Compression, 'G', 3, 64), "E", "*10^")
	epsilon := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Epsilon, 'G', 3, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(data.Dcp, 'G', 3, 64), "E", "*10^")
	Lb0 := strings.ReplaceAll(strconv.FormatFloat(result.Bolt.Lenght, 'G', 3, 64), "E", "*10^")
	typeBolt := strings.ReplaceAll(strconv.FormatFloat(s.typeBolt[result.Data.Type], 'G', 3, 64), "E", "*10^")
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 3, 64)
	bEpsilon := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Epsilon, 'G', 3, 64), "E", "*10^")
	bEpsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	area := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Area, 'G', 3, 64), "E", "*10^")
	count := data.Bolt.Count

	yf1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Yf, 'G', 3, 64), "E", "*10^")
	e1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.E, 'G', 3, 64), "E", "*10^")
	b1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.B, 'G', 3, 64), "E", "*10^")
	d6 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.D6, 'G', 3, 64), "E", "*10^")
	a1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.A, 'G', 3, 64), "E", "*10^")
	ep1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Epsilon, 'G', 3, 64), "E", "*10^")
	epAt201 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	y := strings.ReplaceAll(strconv.FormatFloat(data.Cap.Y, 'G', 3, 64), "E", "*10^")

	pres := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Pres, 'G', 3, 64), "E", "*10^")
	permissiblePres := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.PermissiblePres, 'G', 3, 64), "E", "*10^")
	pressure := strings.ReplaceAll(strconv.FormatFloat(math.Abs(result.Data.Pressure), 'G', 3, 64), "E", "*10^")
	m := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.M, 'G', 3, 64), "E", "*10^")
	axialForce := result.Data.AxialForce
	bendingMoment := result.Data.BendingMoment
	po := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Po, 'G', 3, 64), "E", "*10^")
	ab := strings.ReplaceAll(strconv.FormatFloat(result.Calc.A, 'G', 3, 64), "E", "*10^")
	bSigmaAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.SigmaAt20, 'G', 3, 64), "E", "*10^")
	bSigma := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Sigma, 'G', 3, 64), "E", "*10^")
	alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Alpha, 'G', 3, 64), "E", "*10^")
	qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Qd, 'G', 3, 64), "E", "*10^")
	rp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Rp, 'G', 3, 64), "E", "*10^")
	qt := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Qt, 'G', 3, 64), "E", "*10^")
	kyp := strconv.FormatFloat(s.Kyp[IsWork], 'G', 3, 64)
	kyz := strconv.FormatFloat(s.Kyz[Condition], 'G', 3, 64)

	lb := result.Bolt.Lenght + s.typeBolt[result.Data.Type]*float64(data.Bolt.Diameter)
	Lb := strings.ReplaceAll(strconv.FormatFloat(lb, 'G', 3, 64), "E", "*10^")

	yb := strings.ReplaceAll(strconv.FormatFloat(yb_, 'G', 3, 64), "E", "*10^")
	yp := strings.ReplaceAll(strconv.FormatFloat(yp_, 'G', 3, 64), "E", "*10^")

	if TypeGasket == "Oval" {
		formulas.B0 = fmt.Sprintf("%s/4", width)
		formulas.Dcp = fmt.Sprintf("%s - %s/2", DOut, width)
	} else {
		if data.Gasket.Width > constants.Bp {
			formulas.B0 = fmt.Sprintf("%.1f * sqrt(%s)", constants.B0, width)
		}
		formulas.Dcp = fmt.Sprintf("%s - %s", DOut, b0)
	}

	formulas.A = fmt.Sprintf("%d * %s", count, area)

	if !(TypeGasket == "Oval" || data.FType == moment_api.FlangeData_free) {
		formulas.Alpha = fmt.Sprintf("1 - (%s - (%s * %s + %s * %s) * %s)/(%s + %s + (%s + %s) * %s^2)",
			yp, yf1, e1, y, b1, b1, yp, yb, yf1, y, b1)
	}

	formulas.Po = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, Dcp, b0, pres)

	if result.Data.Pressure >= 0 {
		formulas.Rp = fmt.Sprintf("%f * %s * %s * %s * |%s|", math.Pi, Dcp, b0, m, pressure)
	}

	formulas.Qd = fmt.Sprintf("0.785 * %s^2 * %s", Dcp, pressure)
	formulas.Qfm = fmt.Sprintf("max((%d + 4*|%d|/%s);(%d - 4*|%d|/%s))", axialForce, bendingMoment, Dcp, axialForce, bendingMoment, Dcp)

	var tF1, tF2 string
	af1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.AlphaF, 'G', 3, 64), "E", "*10^")
	h1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.H, 'G', 3, 64), "E", "*10^")
	tf1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Tf, 'G', 3, 64), "E", "*10^")
	a := strings.ReplaceAll(strconv.FormatFloat(data.Cap.Alpha, 'G', 3, 64), "E", "*10^")
	h := strings.ReplaceAll(strconv.FormatFloat(data.Cap.H, 'G', 3, 64), "E", "*10^")
	t := strings.ReplaceAll(strconv.FormatFloat(data.Cap.T, 'G', 3, 64), "E", "*10^")

	if IsUseWasher {
		w1 := strings.ReplaceAll(strconv.FormatFloat(data.Washer1.Alpha, 'G', 3, 64), "E", "*10^")
		th := strings.ReplaceAll(strconv.FormatFloat(data.Washer1.Thickness, 'G', 3, 64), "E", "*10^")
		w2 := strings.ReplaceAll(strconv.FormatFloat(data.Washer2.Alpha, 'G', 3, 64), "E", "*10^")

		tF1 = fmt.Sprintf("(%s*%s + %s*%s) * (%s-20) + (%s*%s + %s*%s) * (%s-20)",
			af1, h1, w1, th, tf1, a, h, w2, th, t)
	} else {
		tF1 = fmt.Sprintf("%s * %s * (%s-20) + %s * %s * (%s-20)", af1, h1, tf1, a, h, t)
	}
	tF2 = fmt.Sprintf("%s + %s", h1, h)

	if data.FType == moment_api.FlangeData_free {
		ak := strings.ReplaceAll(strconv.FormatFloat(data.Flange.AlphaK, 'G', 3, 64), "E", "*10^")
		h := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Hk, 'G', 3, 64), "E", "*10^")
		tk := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Tk, 'G', 3, 64), "E", "*10^")

		tF1 += fmt.Sprintf(" + %s * %s * (%s-20)", ak, h, tk)
		tF2 += " + " + h
	}
	if IsEmbedded {
		af := strings.ReplaceAll(strconv.FormatFloat(data.Embed.Alpha, 'G', 3, 64), "E", "*10^")
		h := strings.ReplaceAll(strconv.FormatFloat(data.Embed.Thickness, 'G', 3, 64), "E", "*10^")
		t := strings.ReplaceAll(strconv.FormatFloat(result.Data.Temp, 'G', 3, 64), "E", "*10^")

		tF1 += fmt.Sprintf(" + %s * %s * (%s-20)", af, h, t)
		tF2 += " + " + h
	}

	gamma := strings.ReplaceAll(strconv.FormatFloat(gamma_, 'G', 3, 64), "E", "*10^")
	bAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Alpha, 'G', 3, 64), "E", "*10^")
	bTemp := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Temp, 'G', 3, 64), "E", "*10^")

	formulas.Qt = fmt.Sprintf("%s * (%s - %s * (%s) * (%s-20))", gamma, tF1, bAlpha, tF2, bTemp)

	pb1F := fmt.Sprintf("%s * (%s + %d) + %s", alpha, qd, axialForce, rp)

	if calculation == moment_api.CalcCapRequest_basis {
		formulas.Basis = &moment_api.BasisFormulas{}

		pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pb2, 'G', 3, 64), "E", "*10^")
		pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pb1, 'G', 3, 64), "E", "*10^")
		pbm := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pb, 'G', 3, 64), "E", "*10^")
		pbr := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pbr, 'G', 3, 64), "E", "*10^")
		kyt := strings.ReplaceAll(strconv.FormatFloat(constants.LoadKyt, 'G', 3, 64), "E", "*10^")
		mkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Mkp, 'G', 3, 64), "E", "*10^")
		dSigmaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.DSigmaM, 'G', 3, 64), "E", "*10^")

		formulas.Basis.Pb2 = fmt.Sprintf("max(%s;0.4 * %s * %s)", po, ab, bSigmaAt20)
		formulas.Basis.Pb1 = fmt.Sprintf("max(%s; %s-%s)", pb1F, pb1F, qt)
		formulas.Basis.Pb = fmt.Sprintf("max(%s;%s)", pb1, pb2)
		formulas.Basis.Pbr = fmt.Sprintf("%s + (1-%s) * (%s + %d) + %s", pbm, alpha, qd, axialForce, qt)
		formulas.Basis.SigmaB1 = fmt.Sprintf("%s / %s", pbm, ab)
		formulas.Basis.SigmaB2 = fmt.Sprintf("%s / %s", pbr, ab)
		formulas.Basis.DSigmaM = fmt.Sprintf("1.2 * %s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Basis.DSigmaR = fmt.Sprintf("%s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Basis.Q = fmt.Sprintf("max(%s; %s) / %f * %s *%s", pbm, pbr, math.Pi, Dcp, width)

		if !(result.Calc.Basis.SigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter) {
			formulas.Basis.Mkp = fmt.Sprintf("(0.3 * %s * %s/%d) / 1000", pbm, diameter, count)
		}
		formulas.Basis.Mkp1 = fmt.Sprintf("0.75 * %s", mkp)

		Prek := fmt.Sprintf("0.8 * %s * %s", ab, bSigmaAt20)
		formulas.Basis.Qrek = fmt.Sprintf("(%s) / (%f * %s * %s)", Prek, math.Pi, Dcp, width)
		formulas.Basis.Mrek = fmt.Sprintf("(0.3 * %s * %s/%d) / 1000", Prek, diameter, count)

		Pmax := fmt.Sprintf("%s * %s", dSigmaM, ab)
		formulas.Basis.Qmax = fmt.Sprintf("(%s) / (%f * %s *%s)", Pmax, math.Pi, Dcp, width)

		if TypeGasket == "Soft" && result.Calc.Basis.Qmax > data.Gasket.PermissiblePres {
			Pmax = fmt.Sprintf("%s * %f * %s *%s", permissiblePres, math.Pi, Dcp, width)
		}
		formulas.Basis.Mmax = fmt.Sprintf("(0.3 * %s *%s / %d) / 1000", Pmax, diameter, count)
	} else {
		formulas.Strength = &moment_api.CalcCapFormulas_StrengthFormulas{}
		eAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Cap.EpsilonAt20, 'G', 3, 64), "E", "*10^")
		e := strings.ReplaceAll(strconv.FormatFloat(data.Cap.Epsilon, 'G', 3, 64), "E", "*10^")
		d := strings.ReplaceAll(strconv.FormatFloat(data.Flange.D, 'G', 3, 64), "E", "*10^")
		s0 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.S0, 'G', 3, 64), "E", "*10^")

		if TypeGasket == "Soft" {
			formulas.Strength.Yp = fmt.Sprintf("(%s * %s) / (%s * %f * %s *%s)", th, compression, epsilon, math.Pi, Dcp, width)
		}

		formulas.Strength.Lb = fmt.Sprintf("%s + %s * %s", Lb0, typeBolt, diameter)
		formulas.Strength.Yb = fmt.Sprintf("%s / (%s * %s * %d)", Lb, bEpsilonAt20, area, count)

		divider := fmt.Sprintf("%s + %s * %s/%s + (%s * %s/%s) * %s^2 + (%s * %s/%s) * %s^2",
			yp, yb, bEpsilonAt20, bEpsilon, yf1, epAt201, ep1, b1, y, eAt20, e, b1)

		if data.FType == moment_api.FlangeData_free {
			yk := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Yk, 'G', 3, 64), "E", "*10^")
			ek := strings.ReplaceAll(strconv.FormatFloat(data.Flange.EpsilonK, 'G', 3, 64), "E", "*10^")
			ekAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.EpsilonKAt20, 'G', 3, 64), "E", "*10^")

			divider += fmt.Sprintf("(%s * %s/%s) *%s^2", yk, ekAt20, ek, a1)
		}

		formulas.Strength.Gamma = fmt.Sprintf("1 / %s", divider)

		formulas.Strength.Flange = s.getFlangeFormulas(data.FType, result.Flange, d6, DOut, Dcp)
		formulas.Strength.Cap = s.getCapFormulas(data.CType, data.Cap, h1, d, s0, DOut, Dcp)

		pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.FPb2, 'G', 3, 64), "E", "*10^")
		pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.FPb1, 'G', 3, 64), "E", "*10^")
		fpbm := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.FPb, 'G', 3, 64), "E", "*10^")
		fpbr := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.FPbr, 'G', 3, 64), "E", "*10^")
		fmkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.FMkp, 'G', 3, 64), "E", "*10^")
		qfm := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Qfm, 'G', 3, 64), "E", "*10^")
		spbm := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.SPb, 'G', 3, 64), "E", "*10^")
		spbr := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.SPbr, 'G', 3, 64), "E", "*10^")
		smkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.SMkp, 'G', 3, 64), "E", "*10^")
		kyt := strconv.FormatFloat(constants.NoLoadKyt, 'G', 3, 64)

		formulas.Strength.FPb2 = fmt.Sprintf("max(%s;0.4 * %s * %s)", po, ab, bSigmaAt20)
		formulas.Strength.FPb1 = pb1F
		formulas.Strength.FPb = fmt.Sprintf("max(%s;%s)", pb1, pb2)
		formulas.Strength.FPbr = fmt.Sprintf("%s + (1-%s) * (%s + %d)", fpbm, alpha, qd, axialForce)
		formulas.Strength.FSigmaB1 = fmt.Sprintf("%s / %s", fpbm, ab)
		formulas.Strength.FSigmaB2 = fmt.Sprintf("%s / %s", fpbr, ab)
		formulas.Strength.FDSigmaM = fmt.Sprintf("1.2 * %s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.FDSigmaR = fmt.Sprintf("%s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.FQ = fmt.Sprintf("max(%s; %s) / %f * %s *%s", fpbm, fpbr, math.Pi, Dcp, width)
		if !(result.Calc.Strength.FSigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter) {
			formulas.Strength.FMkp = fmt.Sprintf("(0.3 * %s * %s/%d) / 1000", fpbm, diameter, count)
		}
		formulas.Strength.FMkp1 = fmt.Sprintf("0.75 * %s", fmkp)

		formulas.Strength.Strength = append(formulas.Strength.Strength, s.getStrengthFormulas(data.FType, axialForce, bendingMoment, data, result.Flange, result.Calc.Strength.Strength[0], d6, Dcp, fpbm, fpbr, qd, qfm, pressure, IsWork, false))

		formulas.Strength.Strength = append(formulas.Strength.Strength, s.getStrengthFormulas(data.FType, axialForce, bendingMoment, data, result.Flange, result.Calc.Strength.Strength[1], d6, Dcp, spbm, spbr, qd, qfm, pressure, IsWork, true))

		kyt = strconv.FormatFloat(constants.LoadKyt, 'G', 3, 64)
		dSigmaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.SDSigmaM, 'G', 3, 64), "E", "*10^")

		formulas.Strength.SPb2 = fmt.Sprintf("max(%s;0.4 * %s * %s)", po, ab, bSigmaAt20)
		formulas.Strength.SPb1 = fmt.Sprintf("max(%s; %s-%s)", pb1F, pb1F, qt)
		formulas.Strength.SPb = fmt.Sprintf("max(%s;%s)", pb1, pb2)
		formulas.Strength.SPbr = fmt.Sprintf("%s + (1-%s) * (%s + %d) + %s", spbm, alpha, qd, axialForce, qt)
		formulas.Strength.SSigmaB1 = fmt.Sprintf("%s / %s", spbm, ab)
		formulas.Strength.SSigmaB2 = fmt.Sprintf("%s / %s", spbr, ab)
		formulas.Strength.SDSigmaM = fmt.Sprintf("1.2 * %s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.SDSigmaR = fmt.Sprintf("%s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.SQ = fmt.Sprintf("max(%s; %s) / %f * %s *%s", spbm, spbr, math.Pi, Dcp, width)

		if !(result.Calc.Strength.FSigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter) {
			formulas.Strength.SMkp = fmt.Sprintf("(0.3 * %s * %s/%d) / 1000", spbm, diameter, count)
		}
		formulas.Strength.SMkp1 = fmt.Sprintf("0.75 * %s", smkp)

		Prek := fmt.Sprintf("0.8 * %s * %s", ab, bSigmaAt20)
		formulas.Strength.Qrek = fmt.Sprintf("(%s) / (%f * %s * %s)", Prek, math.Pi, Dcp, width)
		formulas.Strength.Mrek = fmt.Sprintf("(0.3 * %s * %s/%d) / 1000", Prek, diameter, count)

		Pmax := fmt.Sprintf("%s * %s", dSigmaM, ab)
		formulas.Strength.Qmax = fmt.Sprintf("(%s) / (%f * %s *%s)", Pmax, math.Pi, Dcp, width)

		if TypeGasket == "Soft" && result.Calc.Basis.Qmax > data.Gasket.PermissiblePres {
			Pmax = fmt.Sprintf("%s * %f * %s *%s", permissiblePres, math.Pi, Dcp, width)
		}
		formulas.Strength.Mmax = fmt.Sprintf("(0.3 * %s *%s / %d) / 1000", Pmax, diameter, count)
	}

	return formulas
}

func (s *CapService) getFlangeFormulas(
	typeF moment_api.FlangeData_Type,
	data *moment_api.FlangeResult,
	D6, DOut, Dcp string,
) *moment_api.FlangeFormulas {
	f := &moment_api.FlangeFormulas{}

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
	// k := strconv.FormatFloat(result.K, 'G', 3, 64)

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

	if typeF != moment_api.FlangeData_free {
		f.B = fmt.Sprintf("0.5 * (%s - %s)", D6, Dcp)
	} else {
		Ds := fmt.Sprintf("0.5 * (%s + %s + 2*%s)", DOut, dk, h0)
		f.A = fmt.Sprintf("0.5 * (%s - %s)", D6, Ds)
		f.B = fmt.Sprintf("0.5 * (%s - %s)", Ds, Dcp)
	}

	if typeF == moment_api.FlangeData_welded {
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

	if typeF == moment_api.FlangeData_free {
		f.Psik = fmt.Sprintf("1.28 * (log(%s/%s) / log(10))", dnk, dk)
		f.Yk = fmt.Sprintf("1 / (%s * %s^3 * %s)", epsilonKAt20, hk, psik)
	}

	if typeF != moment_api.FlangeData_free {
		f.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * %s^3))", math.Pi, D6, epsilonAt20, DOut, h)
	} else {
		f.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * %s^3))", math.Pi, ds, epsilonAt20, DOut, h)
		f.Yfc = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * %s^3))", math.Pi, D6, epsilonKAt20, dnk, hk)
	}

	return f
}

func (s *CapService) getCapFormulas(
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
		f.X = fmt.Sprintf("0.67 * %s^2 - 1 + (1 + 8.55 * lg(%s) - 1) / ((%s - 1) * %s^2 - 1 + (1.857 * %s^2 + 1) * %s^3/%s^3)",
			k, k, k, k, k, H, delta)
		f.Y = fmt.Sprintf("%s / (%s * %s)", x, delta, eAt20)
	} else {
		radius := strings.ReplaceAll(strconv.FormatFloat(data.Radius, 'G', 3, 64), "E", "*10^")

		f.K = fmt.Sprintf("(%s / %s) * Sqrt(%s / %s)", h, D, radius, S0)
		f.X = fmt.Sprintf("1 / (1 + 1.285*%s + 1.63*%s * (%s/%s)^2 * lg(%s/%s)", k, k, h, S0, DOut, D)
		f.Y = fmt.Sprintf("((1 - %s * (1 + 1.285*%s)) / (%s * %s^3)) * ((%s + %s) / (%s - %s))", x, k, eAt20, h, DOut, D, DOut, D)
	}

	return f
}

func (s *CapService) getStrengthFormulas(
	typeF moment_api.FlangeData_Type,
	AxialForce, BendingMoment int32,
	data models.DataCap,
	flange *moment_api.FlangeResult,
	res *moment_api.StrengthResult,
	D6, Dcp, Pbm, Pbr, Qd, QFM, pressure string,
	isWork, isTemp bool,
) *moment_api.AddStrengthFormulas {
	strength := &moment_api.AddStrengthFormulas{}

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

	if typeF == moment_api.FlangeData_free {
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

	if typeF == moment_api.FlangeData_welded && flange.S1 != flange.S0 {
		strength.SigmaM1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mm, lymda, s1, c, dvz)
		strength.SigmaM0 = fmt.Sprintf("%s * %s", f, sigmaM1)
	} else {
		strength.SigmaM1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mm, lymda, s0, c, dvz)
	}

	strength.SigmaR = fmt.Sprintf("((1.33 * %s * %s + %s)/(%s * %s^2 * %s * %s)) * %s",
		betaF, h, l0, lymda, h, l0, d, mm)
	strength.SigmaT = fmt.Sprintf("(%s * %s)/(%s^2 * %s) - %s * %s", betaY, mm, h, d, betaZ, sigmaR)

	hk := strings.ReplaceAll(strconv.FormatFloat(flange.Hk, 'G', 3, 64), "E", "*10^")
	dk := strings.ReplaceAll(strconv.FormatFloat(flange.Dk, 'G', 3, 64), "E", "*10^")
	mmk := strings.ReplaceAll(strconv.FormatFloat(res.MMk, 'G', 3, 64), "E", "*10^")
	mp := strings.ReplaceAll(strconv.FormatFloat(res.Mp, 'G', 3, 64), "E", "*10^")
	sigmaP1 := strings.ReplaceAll(strconv.FormatFloat(res.SigmaP1, 'G', 3, 64), "E", "*10^")

	if typeF == moment_api.FlangeData_free {
		strength.SigmaK = fmt.Sprintf("%s * %s / %s^2 * %s", betaY, mmk, hk, dk)
	}

	if typeF == moment_api.FlangeData_welded && flange.S1 != flange.S0 {
		strength.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mp, lymda, s1, c, dvz)
		strength.SigmaP0 = fmt.Sprintf("%s * %s", f, sigmaP1)
	} else {
		strength.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", mp, lymda, s0, c, dvz)
	}

	if typeF == moment_api.FlangeData_welded {
		temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, d, s1, s1, c)
		strength.SigmaMp = fmt.Sprintf("(0.785 * %s^2 * %s + %d + (4 * |%d|)/(%s + %s)) / %s",
			d, pressure, AxialForce, BendingMoment, d, s1, temp)
		strength.SigmaMpm = fmt.Sprintf("(0.785 * %s^2 * %s + %d - (4 * |%d|)/(%s + %s)) / %s",
			d, pressure, AxialForce, BendingMoment, d, s1, temp)
	}

	temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, d, s0, s0, c)
	strength.SigmaMp0 = fmt.Sprintf("(0.785 * %s^2 * %s + %d + (4 * |%d|)/(%s + %s)) / %s",
		d, pressure, AxialForce, BendingMoment, d, s0, temp)
	strength.SigmaMpm0 = fmt.Sprintf("(0.785 * %s^2 * %s + %d - (4 * |%d|)/(%s + %s)) / %s",
		d, pressure, AxialForce, BendingMoment, d, s0, temp)

	strength.SigmaMop = fmt.Sprintf("%s * %s / (2 * (%s - %s))", pressure, d, s0, c)

	sigmaRp := strings.ReplaceAll(strconv.FormatFloat(res.SigmaRp, 'G', 3, 64), "E", "*10^")
	sigmaTp := strings.ReplaceAll(strconv.FormatFloat(res.SigmaTp, 'G', 3, 64), "E", "*10^")

	strength.SigmaRp = fmt.Sprintf("((1.33 * %s * %s + %s)/(%s * %s^2 * %s * %s)) * %s", betaF, h, l0, lymda, h, l0, d, mp)
	strength.SigmaTp = fmt.Sprintf("(%s * %s)/(%s^2 * %s) - %s * %s", betaY, mp, h, d, betaZ, sigmaRp)

	if typeF == moment_api.FlangeData_free {
		strength.SigmaKp = fmt.Sprintf("(%s * %s)/(%s^2 * %s)", betaY, mp, hk, dk)
	}

	epsilon := strings.ReplaceAll(strconv.FormatFloat(flange.Epsilon, 'G', 3, 64), "E", "*10^")
	epsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	yf := strings.ReplaceAll(strconv.FormatFloat(flange.Yf, 'G', 3, 64), "E", "*10^")

	strength.Teta = fmt.Sprintf("%s * %s * %s/%s", mp, yf, epsilonAt20, epsilon)

	if typeF == moment_api.FlangeData_free {
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

	if typeF == moment_api.FlangeData_welded && flange.S1 != flange.S0 {
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

	if typeF == moment_api.FlangeData_free {
		sigmaK := strings.ReplaceAll(strconv.FormatFloat(res.SigmaK, 'G', 3, 64), "E", "*10^")
		sigmaKp := strings.ReplaceAll(strconv.FormatFloat(res.SigmaKp, 'G', 3, 64), "E", "*10^")
		sigmaKAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaKAt20, 'G', 3, 64), "E", "*10^")
		sigma := strings.ReplaceAll(strconv.FormatFloat(flange.SigmaK, 'G', 3, 64), "E", "*10^")

		strength.Max10 = fmt.Sprintf("%s ≤ %s * %s", sigmaK, kt, sigmaKAt20)
		strength.Max11 = fmt.Sprintf("%s ≤ %s * %s", sigmaKp, kt, sigma)
	}

	return strength
}
