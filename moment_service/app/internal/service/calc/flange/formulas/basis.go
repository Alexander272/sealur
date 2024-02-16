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

func (s *FormulasService) basisFormulas(
	req *calc_api.FlangeRequest,
	d models.DataFlange,
	result *calc_api.FlangeResponse,
	aux *flange_model.CalcAuxiliary,
) *flange_model.Formulas_Basis {
	bolt := s.boltStrengthFormulas(
		req, d, result,
		result.Calc.Basis.ForcesInBolts.Pb,
		result.Calc.Basis.ForcesInBolts.Pbr,
		result.Calc.Basis.ForcesInBolts.A,
		result.Calc.Basis.Deformation.Dcp,
		true,
	)

	moment := s.momentFormulas(
		req, d, result,
		result.Calc.Basis.BoltStrength.SigmaB1,
		result.Calc.Basis.BoltStrength.DSigmaM,
		result.Calc.Basis.ForcesInBolts.Pb,
		result.Calc.Basis.ForcesInBolts.A,
		result.Calc.Basis.Deformation.Dcp,
		result.Calc.Basis.Moment,
		true,
	)

	formulas := &flange_model.Formulas_Basis{
		Deformation:   s.deformationFormulas(req, d, result),
		ForcesInBolts: s.forcesInBoltsCalculte(req, d, result, aux),
		BoltStrength:  bolt,
		Moment:        moment,
	}

	return formulas
}

func (s *FormulasService) deformationFormulas(req *calc_api.FlangeRequest, d models.DataFlange, result *calc_api.FlangeResponse,
) *flange_model.DeformationFormulas {
	deformation := &flange_model.DeformationFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Pressure, 'G', 5, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 5, 64)
	dOut := strconv.FormatFloat(d.Gasket.DOut, 'G', 5, 64)
	pres := strconv.FormatFloat(d.Gasket.Pres, 'G', 5, 64)
	m := strconv.FormatFloat(d.Gasket.M, 'G', 5, 64)

	B0 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.B0, 'G', 5, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Dcp, 'G', 5, 64), "E", "*10^")

	if d.TypeGasket == flange_model.GasketData_Oval {
		// фомула 4
		deformation.B0 = fmt.Sprintf("%s / 4", width)
		// фомула ?
		deformation.Dcp = fmt.Sprintf("%s - %s/2", dOut, width)

	} else {
		if !(d.Gasket.Width <= constants.Bp) {
			// фомула 3
			deformation.B0 = fmt.Sprintf("%.1f * sqrt(%s)", constants.B0, width)
		}
		// фомула 5
		deformation.Dcp = fmt.Sprintf("%s - %s", dOut, B0)
	}

	// формула 6
	// Усилие необходимое для смятия прокладки при затяжке
	deformation.Po = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, Dcp, B0, pres)

	if req.Pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях
		deformation.Rp = fmt.Sprintf("%f * %s * %s * %s *|%s|", math.Pi, Dcp, B0, m, pressure)
	}

	return deformation
}

func (s *FormulasService) forcesInBoltsCalculte(
	req *calc_api.FlangeRequest, d models.DataFlange, result *calc_api.FlangeResponse, aux *flange_model.CalcAuxiliary,
) *flange_model.ForcesInBoltsFormulas {
	forces := &flange_model.ForcesInBoltsFormulas{}

	// перевод чисел в строки
	axialForce := req.AxialForce
	bendingMoment := req.BendingMoment
	pressure := strconv.FormatFloat(req.Pressure, 'G', 5, 64)
	temp := strconv.FormatFloat(req.Temp, 'G', 5, 64)

	count := d.Bolt.Count
	area := strconv.FormatFloat(d.Bolt.Area, 'G', 5, 64)
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 5, 64)
	bAlpha := strings.ReplaceAll(strconv.FormatFloat(d.Bolt.Alpha, 'G', 5, 64), "E", "*10^")
	bTemp := strconv.FormatFloat(d.Bolt.Temp, 'G', 5, 64)

	db1 := strconv.FormatFloat(d.Flange1.D6, 'G', 5, 64)
	alphaF1 := strings.ReplaceAll(strconv.FormatFloat(d.Flange1.AlphaF, 'G', 5, 64), "E", "*10^")
	alphaF2 := strings.ReplaceAll(strconv.FormatFloat(d.Flange2.AlphaF, 'G', 5, 64), "E", "*10^")
	h1 := strconv.FormatFloat(d.Flange1.H, 'G', 5, 64)
	h2 := strconv.FormatFloat(d.Flange2.H, 'G', 5, 64)
	tf1 := strconv.FormatFloat(d.Flange1.Tf, 'G', 5, 64)
	tf2 := strconv.FormatFloat(d.Flange2.Tf, 'G', 5, 64)

	var wAlpha1, wAlpha2, thick1, thick2 string
	if req.IsUseWasher {
		wAlpha1 = strings.ReplaceAll(strconv.FormatFloat(d.Washer1.Alpha, 'G', 5, 64), "E", "*10^")
		wAlpha2 = strings.ReplaceAll(strconv.FormatFloat(d.Washer2.Alpha, 'G', 5, 64), "E", "*10^")
		thick1 = strconv.FormatFloat(d.Washer1.Thickness, 'G', 5, 64)
		thick2 = strconv.FormatFloat(d.Washer2.Thickness, 'G', 5, 64)
	}

	var eAlpha, eThick string
	if req.IsEmbedded {
		eAlpha = strings.ReplaceAll(strconv.FormatFloat(d.Embed.Alpha, 'G', 5, 64), "E", "*10^")
		eThick = strconv.FormatFloat(d.Embed.Thickness, 'G', 5, 64)
	}

	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Dcp, 'G', 5, 64), "E", "*10^")
	Po := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Po, 'G', 5, 64), "E", "*10^")
	Rp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Rp, 'G', 5, 64), "E", "*10^")

	Qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Qd, 'G', 5, 64), "E", "*10^")
	Ab := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.A, 'G', 5, 64), "E", "*10^")
	Alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Alpha, 'G', 5, 64), "E", "*10^")
	AlphaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.AlphaM, 'G', 5, 64), "E", "*10^")
	Qt := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Qt, 'G', 5, 64), "E", "*10^")
	Pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Pb1, 'G', 5, 64), "E", "*10^")
	Pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Pb2, 'G', 5, 64), "E", "*10^")
	Pb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Pb, 'G', 5, 64), "E", "*10^")

	Yp := strings.ReplaceAll(strconv.FormatFloat(aux.Yp, 'G', 5, 64), "E", "*10^")
	Yb := strings.ReplaceAll(strconv.FormatFloat(aux.Yb, 'G', 5, 64), "E", "*10^")
	Yf1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange1.Yf, 'G', 5, 64), "E", "*10^")
	Yfn1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange1.Yfn, 'G', 5, 64), "E", "*10^")
	E1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange1.E, 'G', 5, 64), "E", "*10^")
	B1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange1.B, 'G', 5, 64), "E", "*10^")
	Yf2 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange2.Yf, 'G', 5, 64), "E", "*10^")
	Yfn2 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange2.Yfn, 'G', 5, 64), "E", "*10^")
	E2 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange2.E, 'G', 5, 64), "E", "*10^")
	B2 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange2.B, 'G', 5, 64), "E", "*10^")
	Gamma := strings.ReplaceAll(strconv.FormatFloat(aux.Gamma, 'G', 5, 64), "E", "*10^")

	// фомула 8
	// Суммарная площадь сечения болтов/шпилек
	forces.A = fmt.Sprintf("%d * %s", count, area)

	// формула 9
	// Равнодействующая нагрузка от давления
	forces.Qd = fmt.Sprintf("0.785 * (%s)^2 * %s", Dcp, pressure)

	temp1 := fmt.Sprintf("%d + 4 * |%d| / %s", axialForce, bendingMoment, Dcp)
	temp2 := fmt.Sprintf("%d - 4 * |%d| / %s", axialForce, bendingMoment, Dcp)

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	forces.Qfm = fmt.Sprintf("max(%s; %s)", temp1, temp2)

	if !(d.TypeGasket == flange_model.GasketData_Oval || d.Type1 == flange_model.FlangeData_free || d.Type2 == flange_model.FlangeData_free) {
		// формула (Е.11)
		// Коэффициент жесткости
		forces.Alpha = fmt.Sprintf("1 - (%s - (%s * %s * %s + %s * %s * %s)) / (%s + %s + (%s * (%s)^2 + %s * (%s)^2))",
			Yp, Yf1, E1, B1, Yf2, E2, B2, Yp, Yb, Yf1, B1, Yf2, B2)
	}

	dividend := fmt.Sprintf("%s + %s * %s * (%s + %s - (%s)^2 / %s) + %s * %s * (%s + %s - (%s)^2 / %s)",
		Yb, Yfn1, B1, B1, E1, E1, Dcp, Yfn2, B2, B2, E2, E2, Dcp)
	divider := fmt.Sprintf("%s + %s * (%s / %s)^2 + %s * (%s)^2 + %s * (%s)^2", Yb, Yp, db1, Dcp, Yfn1, B1, Yfn2, B2)

	if d.Type1 == flange_model.FlangeData_free {
		A1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange1.A, 'G', 5, 64), "E", "*10^")
		Yfc1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange1.Yfc, 'G', 5, 64), "E", "*10^")

		dividend += fmt.Sprintf(" + %s * (%s)^2", Yfc1, A1)
		divider += fmt.Sprintf(" + %s * (%s)^2", Yfc1, A1)
	}
	if d.Type2 == flange_model.FlangeData_free {
		A2 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange2.A, 'G', 5, 64), "E", "*10^")
		Yfc2 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange2.Yfc, 'G', 5, 64), "E", "*10^")

		dividend += fmt.Sprintf(" + %s * (%s)^2", Yfc2, A2)
		divider += fmt.Sprintf(" + %s * (%s)^2", Yfc2, A2)
	}

	// формула (Е.13)
	// Коэффициент жесткости фланцевого соединения нагруженного внешним изгибающим моментом
	forces.AlphaM = fmt.Sprintf("%s / %s", dividend, divider)

	minB := fmt.Sprintf("0.4 * %s * %s", Ab, sigmaAt20)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	forces.Pb2 = fmt.Sprintf("max(%s; %s)", Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	forces.Pb1 = fmt.Sprintf("%s * (%s + %d) + %s + 4 * %s * |%d| / %s", Alpha, Qd, axialForce, Rp, AlphaM, bendingMoment, Dcp)

	if req.IsUseWasher {
		temp1 = fmt.Sprintf("(%s * %s + %s * %s) * (%s - 20) + (%s * %s + %s * %s) * (%s - 20)",
			alphaF1, h1, wAlpha1, thick1, tf1, alphaF2, h2, wAlpha2, thick2, tf2)
	} else {
		temp1 = fmt.Sprintf("%s * %s * (%s - 20) + %s * %s *(%s - 20)", alphaF1, h1, tf1, alphaF2, h2, tf2)
	}
	temp2 = fmt.Sprintf("%s + %s", h1, h2)

	if d.Type1 == flange_model.FlangeData_free {
		alphaK1 := strings.ReplaceAll(strconv.FormatFloat(d.Flange1.Ring.AlphaK, 'G', 5, 64), "E", "*10^")
		hk1 := strconv.FormatFloat(d.Flange1.Hk, 'G', 5, 64)
		tk1 := strconv.FormatFloat(d.Flange1.Ring.Tk, 'G', 5, 64)

		temp1 += fmt.Sprintf(" + %s * %s * (%s - 20)", alphaK1, hk1, tk1)
		temp2 += fmt.Sprintf(" + %s", hk1)
	}
	if d.Type2 == flange_model.FlangeData_free {
		alphaK2 := strings.ReplaceAll(strconv.FormatFloat(d.Flange2.Ring.AlphaK, 'G', 5, 64), "E", "*10^")
		hk2 := strconv.FormatFloat(d.Flange2.Hk, 'G', 5, 64)
		tk2 := strconv.FormatFloat(d.Flange2.Ring.Tk, 'G', 5, 64)

		temp1 += fmt.Sprintf(" + %s * %s * (%s - 20)", alphaK2, hk2, tk2)
		temp2 += fmt.Sprintf(" + %s", hk2)
	}
	if req.IsEmbedded {
		temp1 += fmt.Sprintf(" + %s * %s * (%s - 20)", eAlpha, eThick, temp)
		temp2 += fmt.Sprintf(" + %s", eThick)
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	forces.Qt = fmt.Sprintf("%s * ((%s) - %s * (%s) * (%s - 20))", Gamma, temp1, bAlpha, temp2, bTemp)

	forces.Pb1 = fmt.Sprintf("max(%s; %s-%s)", forces.Pb1, forces.Pb1, Qt)
	forces.Pb = fmt.Sprintf("max(%s; %s)", Pb1, Pb2)
	forces.Pbr = fmt.Sprintf("%s + (1 - %s) * (%s + %d) + %s + 4 * (1 - %s * |%d|) / %s", Pb, Alpha, Qd, axialForce, Qt, AlphaM, bendingMoment, Dcp)
	return forces
}

func (s *FormulasService) boltStrengthFormulas(
	req *calc_api.FlangeRequest, d models.DataFlange, result *calc_api.FlangeResponse,
	Pbm, Pbr, Ab, Dcp float64,
	isLoad bool,
) *flange_model.BoltStrengthFormulas {
	bolt := &flange_model.BoltStrengthFormulas{}

	// перевод чисел в строки
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 5, 64)
	sigma := strconv.FormatFloat(d.Bolt.Sigma, 'G', 5, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 5, 64)

	Ab_ := strings.ReplaceAll(strconv.FormatFloat(Ab, 'G', 5, 64), "E", "*10^")
	Pb_ := strings.ReplaceAll(strconv.FormatFloat(Pbm, 'G', 5, 64), "E", "*10^")
	Pbr_ := strings.ReplaceAll(strconv.FormatFloat(Pbr, 'G', 5, 64), "E", "*10^")

	Dcp_ := strings.ReplaceAll(strconv.FormatFloat(Dcp, 'G', 5, 64), "E", "*10^")

	bolt.SigmaB1 = fmt.Sprintf("%s / %s", Pb_, Ab_)
	bolt.SigmaB2 = fmt.Sprintf("%s / %s", Pbr_, Ab_)

	Kyp := s.Kyp[req.IsWork]
	Kyz := s.Kyz[req.Condition.String()]
	Kyt := s.Kyt[isLoad]
	// формула Г.3
	bolt.DSigmaM = fmt.Sprintf("1.2 * %.2f * %.1f * %.1f * %s", Kyp, Kyz, Kyt, sigmaAt20)
	// формула Г.4
	bolt.DSigmaR = fmt.Sprintf("%.2f * %.1f * %.1f * %s", Kyp, Kyz, Kyt, sigma)

	if d.TypeGasket == flange_model.GasketData_Soft {
		bolt.Q = fmt.Sprintf("max(%s; %s) / (%f * %s * %s)", Pb_, Pbr_, math.Pi, Dcp_, width)
	}

	return bolt
}

func (s *FormulasService) momentFormulas(
	req *calc_api.FlangeRequest, d models.DataFlange, result *calc_api.FlangeResponse,
	SigmaB1, DSigmaM, Pbm, Ab, Dcp float64,
	mom *flange_model.CalcMoment,
	fullCalculate bool,
) *flange_model.MomentFormulas {
	moment := &flange_model.MomentFormulas{}

	// перевод чисел в строки
	friction := strconv.FormatFloat(req.Friction, 'G', 5, 64)
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 5, 64)
	diameter := strconv.FormatFloat(d.Bolt.Diameter, 'G', 5, 64)
	count := d.Bolt.Count

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 5, 64)
	perPres := strconv.FormatFloat(d.Gasket.PermissiblePres, 'G', 5, 64)

	Ab_ := strings.ReplaceAll(strconv.FormatFloat(Ab, 'G', 5, 64), "E", "*10^")
	Pb_ := strings.ReplaceAll(strconv.FormatFloat(Pbm, 'G', 5, 64), "E", "*10^")
	Dcp_ := strings.ReplaceAll(strconv.FormatFloat(Dcp, 'G', 5, 64), "E", "*10^")
	DSigmaM_ := strings.ReplaceAll(strconv.FormatFloat(DSigmaM, 'G', 5, 64), "E", "*10^")
	Mkp := strings.ReplaceAll(strconv.FormatFloat(mom.Mkp, 'G', 5, 64), "E", "*10^")

	if !(SigmaB1 > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter) {
		moment.Mkp = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Pb_, diameter, count)
	}

	moment.Mkp1 = fmt.Sprintf("0.75 * %s", Mkp)

	if fullCalculate {
		Prek := fmt.Sprintf("0.8 * %s * %s", Ab_, sigmaAt20)
		moment.Qrek = fmt.Sprintf("%s / (%f * %s * %s)", Prek, math.Pi, Dcp_, width)
		moment.Mrek = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Prek, diameter, count)

		Pmax := fmt.Sprintf("%s * %s", DSigmaM_, Ab_)
		moment.Qmax = fmt.Sprintf("%s / (%f * %s * %s)", Pmax, math.Pi, Dcp_, width)

		if d.TypeGasket == flange_model.GasketData_Soft && mom.Qmax > d.Gasket.PermissiblePres {
			Pmax = fmt.Sprintf("%s * (%f * %s * %s)", perPres, math.Pi, Dcp_, width)
			moment.Qmax = ""
		}

		if mom.Mrek > mom.Mmax {
			moment.Mrek = ""
		}
		if mom.Qrek > mom.Qmax {
			moment.Qrek = ""
		}

		moment.Mmax = fmt.Sprintf("(%s * %s * %s / %d) / 1000", friction, Pmax, diameter, count)
	}

	return moment
}
