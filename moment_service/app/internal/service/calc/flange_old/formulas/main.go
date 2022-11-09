package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

type FormulasService struct {
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewFormulasService() *FormulasService {
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

	return &FormulasService{
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

func (s *FormulasService) GetFormulas(
	TypeGasket, TypeBolt, Condition string,
	IsWork, IsUseWasher, IsEmbedded bool,
	data models.DataFlangeOld,
	result calc_api.FlangeResponseOld,
	calculation calc_api.FlangeRequest_Calcutation,
	gamma_, yb_, yp_ float64,
) *flange_model.CalcFlangeFormulas {
	formulas := &flange_model.CalcFlangeFormulas{}

	width := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Width, 'G', 3, 64), "E", "*10^")
	DOut := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.DOut, 'G', 3, 64), "E", "*10^")
	b0 := strings.ReplaceAll(strconv.FormatFloat(data.B0, 'G', 3, 64), "E", "*10^")
	th := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Thickness, 'G', 3, 64), "E", "*10^")
	compression := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Compression, 'G', 3, 64), "E", "*10^")
	epsilon := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Epsilon, 'G', 3, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(data.Dcp, 'G', 3, 64), "E", "*10^")
	Lb0 := strings.ReplaceAll(strconv.FormatFloat(result.Bolt.Lenght, 'G', 3, 64), "E", "*10^")
	typeBolt := strings.ReplaceAll(strconv.FormatFloat(s.typeBolt[TypeBolt], 'G', 3, 64), "E", "*10^")
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 3, 64)
	bEpsilon := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Epsilon, 'G', 3, 64), "E", "*10^")
	bEpsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	area := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Area, 'G', 3, 64), "E", "*10^")
	count := data.Bolt.Count

	yf1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Yf, 'G', 3, 64), "E", "*10^")
	e1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.E, 'G', 3, 64), "E", "*10^")
	b1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.B, 'G', 3, 64), "E", "*10^")
	yf2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Yf, 'G', 3, 64), "E", "*10^")
	e2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.E, 'G', 3, 64), "E", "*10^")
	b2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.B, 'G', 3, 64), "E", "*10^")
	d6 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.D6, 'G', 3, 64), "E", "*10^")
	yfn1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Yfn, 'G', 3, 64), "E", "*10^")
	yfn2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Yfn, 'G', 3, 64), "E", "*10^")
	yfc1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Yfc, 'G', 3, 64), "E", "*10^")
	a1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.A, 'G', 3, 64), "E", "*10^")
	yfc2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Yfc, 'G', 3, 64), "E", "*10^")
	a2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.A, 'G', 3, 64), "E", "*10^")
	ep1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Epsilon, 'G', 3, 64), "E", "*10^")
	ep2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Epsilon, 'G', 3, 64), "E", "*10^")
	epAt201 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	epAt202 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.EpsilonAt20, 'G', 3, 64), "E", "*10^")

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
	alphaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.AlphaM, 'G', 3, 64), "E", "*10^")
	qt := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Qt, 'G', 3, 64), "E", "*10^")
	kyp := strings.ReplaceAll(strconv.FormatFloat(s.Kyp[IsWork], 'G', 3, 64), "E", "*10^")
	kyz := strings.ReplaceAll(strconv.FormatFloat(s.Kyz[Condition], 'G', 3, 64), "E", "*10^")

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

	if !(TypeGasket == "Oval" || data.Type1 == flange_model.FlangeData_free || data.Type2 == flange_model.FlangeData_free) {
		formulas.Alpha = fmt.Sprintf("1 - (%s - (%s * %s * %s + %s * %s * %s))/(%s + %s + (%s * (%s)^2 + %s * (%s)^2))",
			yp, yf1, e1, b1, yf2, e2, b2, yp, yb, yf1, b1, yf2, b2)
	}

	dividendF := fmt.Sprintf("(%s + %s * %s * (%s + %s - (%s)^2/%s) + %s * %s * (%s + %s - (%s)^2/%s)",
		yb, yfn1, b1, b1, e1, e1, Dcp, yfn2, b2, b2, e2, e2, Dcp)
	dividerF := fmt.Sprintf("(%s + %s * (%s/%s)^2 + %s * (%s)^2 + %s * (%s)^2",
		yb, yp, d6, Dcp, yfn1, b1, yfn2, b2)

	if data.Type1 == flange_model.FlangeData_free {
		dividendF += fmt.Sprintf("%s * (%s)^2", yfc1, a1)
		dividerF += fmt.Sprintf("%s * (%s)^2", yfc1, a1)
	}
	if data.Type2 == flange_model.FlangeData_free {
		dividendF += fmt.Sprintf("%s * (%s)^2", yfc2, a2)
		dividerF += fmt.Sprintf("%s * (%s)^2", yfc2, a2)
	}
	formulas.AlphaM = dividendF + ") / " + dividerF + ")"

	formulas.Po = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, Dcp, b0, pres)

	if result.Data.Pressure >= 0 {
		formulas.Rp = fmt.Sprintf("%f * %s * %s * %s * |%s|", math.Pi, Dcp, b0, m, pressure)
	}

	formulas.Qd = fmt.Sprintf("0.785 * (%s)^2 * %s", Dcp, pressure)
	formulas.Qfm = fmt.Sprintf("max((%d + 4*|%d|/%s);(%d - 4*|%d|/%s))", axialForce, bendingMoment, Dcp, axialForce, bendingMoment, Dcp)

	var tF1, tF2 string
	af1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.AlphaF, 'G', 3, 64), "E", "*10^")
	h1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.H, 'G', 3, 64), "E", "*10^")
	tf1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Tf, 'G', 3, 64), "E", "*10^")
	af2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.AlphaF, 'G', 3, 64), "E", "*10^")
	h2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.H, 'G', 3, 64), "E", "*10^")
	tf2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Tf, 'G', 3, 64), "E", "*10^")

	if IsUseWasher {
		w1 := strings.ReplaceAll(strconv.FormatFloat(data.Washer1.Alpha, 'G', 3, 64), "E", "*10^")
		th := strings.ReplaceAll(strconv.FormatFloat(data.Washer1.Thickness, 'G', 3, 64), "E", "*10^")
		w2 := strings.ReplaceAll(strconv.FormatFloat(data.Washer2.Alpha, 'G', 3, 64), "E", "*10^")

		tF1 = fmt.Sprintf("(%s*%s + %s*%s) * (%s-20) + (%s*%s + %s*%s) * (%s-20)",
			af1, h1, w1, th, tf1, af2, h2, w2, th, tf2)
	} else {
		tF1 = fmt.Sprintf("%s * %s * (%s-20) + %s * %s * (%s-20)", af1, h1, tf1, af2, h2, tf2)
	}
	tF2 = fmt.Sprintf("%s + %s", h1, h2)

	if data.Type1 == flange_model.FlangeData_free {
		ak := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.AlphaK, 'G', 3, 64), "E", "*10^")
		h := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Hk, 'G', 3, 64), "E", "*10^")
		tk := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Tk, 'G', 3, 64), "E", "*10^")

		tF1 += fmt.Sprintf(" + %s * %s * (%s-20)", ak, h, tk)
		tF2 += " + " + h
	}
	if data.Type2 == flange_model.FlangeData_free {
		ak := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.AlphaK, 'G', 3, 64), "E", "*10^")
		h := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Hk, 'G', 3, 64), "E", "*10^")
		tk := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Tk, 'G', 3, 64), "E", "*10^")

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

	pb1F := fmt.Sprintf("%s * (%s + %d) + %s + 4 * %s * |%d|/%s", alpha, qd, axialForce, rp, alphaM, bendingMoment, Dcp)

	if calculation == calc_api.FlangeRequest_basis {
		formulas.Basis = &flange_model.BasisFormulas{}

		pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pb2, 'G', 3, 64), "E", "*10^")
		pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pb1, 'G', 3, 64), "E", "*10^")
		pbm := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pb, 'G', 3, 64), "E", "*10^")
		pbr := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Pbr, 'G', 3, 64), "E", "*10^")
		kyt := strings.ReplaceAll(strconv.FormatFloat(constants.LoadKyt, 'G', 3, 64), "E", "*10^")
		mkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Mkp, 'G', 3, 64), "E", "*10^")
		dSigmaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.DSigmaM, 'G', 3, 64), "E", "*10^")

		formulas.Basis.Pb2 = fmt.Sprintf("max(%s; 0.4 * %s * %s)", po, ab, bSigmaAt20)
		formulas.Basis.Pb1 = fmt.Sprintf("max(%s; %s-%s)", pb1F, pb1F, qt)
		formulas.Basis.Pb = fmt.Sprintf("max(%s; %s)", pb1, pb2)
		formulas.Basis.Pbr = fmt.Sprintf("%s + (1-%s) * (%s + %d) + %s + 4 * (1-%s) * |%d|/%s",
			pbm, alpha, qd, axialForce, qt, alphaM, bendingMoment, Dcp)
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
		formulas.Strength = &flange_model.StrengthFormulas{}

		if TypeGasket == "Soft" {
			formulas.Strength.Yp = fmt.Sprintf("(%s * %s) / (%s * %f * %s *%s)", th, compression, epsilon, math.Pi, Dcp, width)
		}

		formulas.Strength.Lb = fmt.Sprintf("%s + %s * %s", Lb0, typeBolt, diameter)
		formulas.Strength.Yb = fmt.Sprintf("%s / (%s * %s * %d)", Lb, bEpsilonAt20, area, count)

		divider := fmt.Sprintf("%s + %s * %s/%s + (%s * %s/%s) * (%s)^2 + (%s * %s/%s) * (%s)^2",
			yp, yb, bEpsilonAt20, bEpsilon, yf1, epAt201, ep1, b1, yf2, epAt202, ep2, b2)

		if data.Type1 == flange_model.FlangeData_free {
			yk := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Yk, 'G', 3, 64), "E", "*10^")
			ek := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.EpsilonK, 'G', 3, 64), "E", "*10^")
			ekAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.EpsilonKAt20, 'G', 3, 64), "E", "*10^")

			divider += fmt.Sprintf("(%s * %s/%s) * (%s)^2", yk, ekAt20, ek, a1)
		}
		if data.Type2 == flange_model.FlangeData_free {
			yk := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Yk, 'G', 3, 64), "E", "*10^")
			ek := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.EpsilonK, 'G', 3, 64), "E", "*10^")
			ekAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.EpsilonKAt20, 'G', 3, 64), "E", "*10^")

			divider += fmt.Sprintf("(%s * %s/%s) * (%s)^2", yk, ekAt20, ek, a2)
		}

		formulas.Strength.Gamma = fmt.Sprintf("1 / %s", divider)

		formulas.Strength.Flange = append(formulas.Strength.Flange, s.getFlangeFormulas(data.Type1, result.Flanges[0], d6, DOut, Dcp))
		if !result.IsSameFlange {
			formulas.Strength.Flange = append(formulas.Strength.Flange, s.getFlangeFormulas(data.Type2, result.Flanges[1], d6, DOut, Dcp))
		}

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

		formulas.Strength.FPb2 = fmt.Sprintf("max(%s; 0.4 * %s * %s)", po, ab, bSigmaAt20)
		formulas.Strength.FPb1 = pb1F
		formulas.Strength.FPb = fmt.Sprintf("max(%s; %s)", pb1, pb2)
		formulas.Strength.FPbr = fmt.Sprintf("%s + (1-%s) * (%s + %d) + 4 * (1-%s) * |%d|/%s",
			fpbm, alpha, qd, axialForce, alphaM, bendingMoment, Dcp)
		formulas.Strength.FSigmaB1 = fmt.Sprintf("%s / %s", fpbm, ab)
		formulas.Strength.FSigmaB2 = fmt.Sprintf("%s / %s", fpbr, ab)
		formulas.Strength.FDSigmaM = fmt.Sprintf("1.2 * %s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.FDSigmaR = fmt.Sprintf("%s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.FQ = fmt.Sprintf("max(%s; %s) / %f * %s *%s", fpbm, fpbr, math.Pi, Dcp, width)
		if !(result.Calc.Strength.FSigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter) {
			formulas.Strength.FMkp = fmt.Sprintf("(0.3 * %s * %s/%d) / 1000", fpbm, diameter, count)
		}
		formulas.Strength.FMkp1 = fmt.Sprintf("0.75 * %s", fmkp)

		formulas.Strength.Strength = append(formulas.Strength.Strength, s.getStrengthFormulas(
			data.Type1,
			axialForce, bendingMoment,
			data,
			result.Flanges[0],
			result.Calc.Strength.Strength[0],
			d6, Dcp, fpbm, fpbr, qd, qfm, pressure,
			IsWork, false,
		))
		if !result.IsSameFlange {
			formulas.Strength.Strength = append(formulas.Strength.Strength, s.getStrengthFormulas(
				data.Type1,
				axialForce, bendingMoment,
				data,
				result.Flanges[0],
				result.Calc.Strength.Strength[1],
				d6, Dcp, fpbm, fpbr, qd, qfm, pressure,
				IsWork, false,
			))
		}

		if result.IsSameFlange {
			formulas.Strength.Strength = append(formulas.Strength.Strength, s.getStrengthFormulas(
				data.Type1,
				axialForce,
				bendingMoment, data,
				result.Flanges[0],
				result.Calc.Strength.Strength[1],
				d6, Dcp, spbm, spbr, qd, qfm, pressure,
				IsWork, true,
			))
		} else {
			formulas.Strength.Strength = append(formulas.Strength.Strength, s.getStrengthFormulas(
				data.Type1,
				axialForce, bendingMoment,
				data,
				result.Flanges[0],
				result.Calc.Strength.Strength[2],
				d6, Dcp, spbm, spbr, qd, qfm, pressure,
				IsWork, true,
			))
			formulas.Strength.Strength = append(formulas.Strength.Strength, s.getStrengthFormulas(
				data.Type1,
				axialForce, bendingMoment,
				data,
				result.Flanges[1],
				result.Calc.Strength.Strength[3],
				d6, Dcp, spbm, spbr, qd, qfm, pressure,
				IsWork, true,
			))
		}

		kyt = strconv.FormatFloat(constants.LoadKyt, 'G', 3, 64)
		dSigmaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.SDSigmaM, 'G', 3, 64), "E", "*10^")

		formulas.Strength.SPb2 = fmt.Sprintf("max(%s; 0.4 * %s * %s)", po, ab, bSigmaAt20)
		formulas.Strength.SPb1 = fmt.Sprintf("max(%s; %s-%s)", pb1F, pb1F, qt)
		formulas.Strength.SPb = fmt.Sprintf("max(%s; %s)", pb1, pb2)
		formulas.Strength.SPbr = fmt.Sprintf("%s + (1-%s) * (%s + %d) + %s + 4 * (1-%s) * |%d|/%s",
			spbm, alpha, qd, axialForce, qt, alphaM, bendingMoment, Dcp)
		formulas.Strength.SSigmaB1 = fmt.Sprintf("%s / %s", spbm, ab)
		formulas.Strength.SSigmaB2 = fmt.Sprintf("%s / %s", spbr, ab)
		formulas.Strength.SDSigmaM = fmt.Sprintf("1.2 * %s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.SDSigmaR = fmt.Sprintf("%s * %s * %s * %s", kyp, kyz, kyt, bSigma)
		formulas.Strength.SQ = fmt.Sprintf("max(%s; %s) / %f * %s *%s", spbm, spbr, math.Pi, Dcp, width)

		if !(result.Calc.Strength.FSigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter &&
			data.Bolt.Diameter <= constants.MaxDiameter) {
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
