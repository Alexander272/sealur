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

func (s *FormulasService) strengthFormulas(req *calc_api.FlangeRequest, d models.DataFlange, result *calc_api.FlangeResponse,
) *flange_model.Formulas_Strength {
	auxiliary := s.auxiliaryFormulas(d, req, result)
	tightness := s.tightnessFormulas(d, req, result)
	bolt1 := s.boltStrengthFormulas(
		req, d, result,
		result.Calc.Strength.Tightness.Pb,
		result.Calc.Strength.Tightness.Pbr,
		result.Calc.Strength.Auxiliary.A,
		result.Calc.Strength.Auxiliary.Dcp,
	)
	moment1 := s.momentFormulas(
		req, d, result,
		result.Calc.Strength.BoltStrength1.SigmaB1,
		result.Calc.Strength.BoltStrength1.DSigmaM,
		result.Calc.Strength.Tightness.Pb,
		result.Calc.Strength.Auxiliary.A,
		result.Calc.Strength.Auxiliary.Dcp,
		result.Calc.Strength.Moment1,
		false,
	)

	static1 := []*flange_model.StaticResistanceFormulas{}
	for i, csr := range result.Calc.Strength.StaticResistance1 {
		flange := result.Calc.Strength.Auxiliary.Flange1
		if i == 1 {
			flange = result.Calc.Strength.Auxiliary.Flange2
		}
		static := s.staticResistanceCalculate(
			result.Flanges[i],
			flange,
			req.FlangesData[i].Type,
			csr, d, req,
			result.Calc.Strength.Tightness.Pb,
			result.Calc.Strength.Tightness.Pbr,
			result.Calc.Strength.Tightness.Qd,
			result.Calc.Strength.Tightness.Qfm,
		)
		static1 = append(static1, static)
	}
	conditions1 := []*flange_model.ConditionsForStrengthFormulas{}
	for i, ccfs := range result.Calc.Strength.ConditionsForStrength1 {
		flange := result.Calc.Strength.Auxiliary.Flange1
		if i == 1 {
			flange = result.Calc.Strength.Auxiliary.Flange2
		}
		cond := s.conditionsForStrengthCalculate(
			req.FlangesData[i].Type,
			result.Flanges[i],
			flange,
			result.Calc.Strength.StaticResistance1[i],
			ccfs,
			req.IsWork, false,
		)
		conditions1 = append(conditions1, cond)
	}

	tigLoad := s.tightnessLoadCalculate(d, req, result)
	bolt2 := s.boltStrengthFormulas(
		req, d, result,
		result.Calc.Strength.TightnessLoad.Pb,
		result.Calc.Strength.TightnessLoad.Pbr,
		result.Calc.Strength.Auxiliary.A,
		result.Calc.Strength.Auxiliary.Dcp,
	)
	moment2 := s.momentFormulas(
		req, d, result,
		result.Calc.Strength.BoltStrength2.SigmaB1,
		result.Calc.Strength.BoltStrength2.DSigmaM,
		result.Calc.Strength.TightnessLoad.Pb,
		result.Calc.Strength.Auxiliary.A,
		result.Calc.Strength.Auxiliary.Dcp,
		result.Calc.Strength.Moment2,
		false,
	)

	static2 := []*flange_model.StaticResistanceFormulas{}
	for i, csr := range result.Calc.Strength.StaticResistance1 {
		flange := result.Calc.Strength.Auxiliary.Flange1
		if i == 1 {
			flange = result.Calc.Strength.Auxiliary.Flange2
		}
		static := s.staticResistanceCalculate(
			result.Flanges[i],
			flange,
			req.FlangesData[i].Type,
			csr, d, req,
			result.Calc.Strength.TightnessLoad.Pb,
			result.Calc.Strength.TightnessLoad.Pbr,
			result.Calc.Strength.Tightness.Qd,
			result.Calc.Strength.Tightness.Qfm,
		)
		static2 = append(static2, static)
	}
	conditions2 := []*flange_model.ConditionsForStrengthFormulas{}
	for i, ccfs := range result.Calc.Strength.ConditionsForStrength1 {
		flange := result.Calc.Strength.Auxiliary.Flange1
		if i == 1 {
			flange = result.Calc.Strength.Auxiliary.Flange2
		}
		cond := s.conditionsForStrengthCalculate(
			req.FlangesData[i].Type,
			result.Flanges[i],
			flange,
			result.Calc.Strength.StaticResistance2[i],
			ccfs,
			req.IsWork, false,
		)
		conditions2 = append(conditions2, cond)
	}

	deformation := &flange_model.DeformationFormulas{
		B0:  auxiliary.B0,
		Dcp: auxiliary.Dcp,
		Po:  tightness.Po,
		Rp:  tightness.Rp,
	}
	forces := &flange_model.ForcesInBoltsFormulas{
		A:      auxiliary.A,
		Qd:     tightness.Qd,
		Qfm:    tightness.Qfm,
		Qt:     tigLoad.Qt,
		Pb:     tigLoad.Pb,
		Alpha:  auxiliary.Alpha,
		AlphaM: auxiliary.AlphaM,
		Pb1:    tigLoad.Pb1,
		Pb2:    tightness.Pb2,
		Pbr:    tigLoad.Pbr,
	}
	finalMoment := &flange_model.MomentFormulas{}
	if result.Calc.Strength.FinalMoment != nil {
		finalMoment = s.momentFormulas(
			req, d, result,
			result.Calc.Strength.BoltStrength2.SigmaB1,
			result.Calc.Strength.BoltStrength2.DSigmaM,
			result.Calc.Strength.TightnessLoad.Pb,
			result.Calc.Strength.Auxiliary.A,
			result.Calc.Strength.Auxiliary.Dcp,
			result.Calc.Strength.FinalMoment,
			true,
		)
	}

	formulas := &flange_model.Formulas_Strength{
		Auxiliary:              auxiliary,
		Tightness:              tightness,
		BoltStrength1:          bolt1,
		Moment1:                moment1,
		StaticResistance1:      static1,
		ConditionsForStrength1: conditions1,
		TightnessLoad:          tigLoad,
		BoltStrength2:          bolt2,
		Moment2:                moment2,
		StaticResistance2:      static2,
		ConditionsForStrength2: conditions2,
		Deformation:            deformation,
		ForcesInBolts:          forces,
		FinalMoment:            finalMoment,
	}

	return formulas
}

func (s *FormulasService) auxiliaryFormulas(data models.DataFlange, req *calc_api.FlangeRequest, result *calc_api.FlangeResponse,
) *flange_model.AuxiliaryFormulas {
	auxiliary := &flange_model.AuxiliaryFormulas{}

	// перевод чисел в строки
	width := strconv.FormatFloat(data.Gasket.Width, 'G', 3, 64)
	thickness := strconv.FormatFloat(data.Gasket.Thickness, 'G', 3, 64)
	dOut := strconv.FormatFloat(data.Gasket.DOut, 'G', 3, 64)
	epsilon := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Epsilon, 'G', 3, 64), "E", "*10^")
	compression := strconv.FormatFloat(data.Gasket.Compression, 'G', 3, 64)

	lenght := strconv.FormatFloat(data.Bolt.Length, 'G', 3, 64)
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 3, 64)
	area := strconv.FormatFloat(data.Bolt.Area, 'G', 3, 64)
	epsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	bEpsilon := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Epsilon, 'G', 3, 64), "E", "*10^")
	count := data.Bolt.Count

	db1 := strconv.FormatFloat(data.Flange1.D6, 'G', 3, 64)
	fEpAt201 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	fEpAt202 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	fEp1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Epsilon, 'G', 3, 64), "E", "*10^")
	fEp2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Epsilon, 'G', 3, 64), "E", "*10^")

	typeBolt := strconv.FormatFloat(s.typeBolt[req.Type.String()], 'G', 3, 64)

	B0 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.B0, 'G', 3, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Dcp, 'G', 3, 64), "E", "*10^")
	Lb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Lb, 'G', 3, 64), "E", "*10^")

	Yp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Yp, 'G', 3, 64), "E", "*10^")
	Yb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Yb, 'G', 3, 64), "E", "*10^")
	Yf1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange1.Yf, 'G', 3, 64), "E", "*10^")
	Yfc1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange1.Yfc, 'G', 3, 64), "E", "*10^")
	Yfn1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange1.Yfn, 'G', 3, 64), "E", "*10^")
	Yk1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange1.Yk, 'G', 3, 64), "E", "*10^")
	A1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange1.A, 'G', 3, 64), "E", "*10^")
	E1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange1.E, 'G', 3, 64), "E", "*10^")
	B1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange1.B, 'G', 3, 64), "E", "*10^")
	Yf2 := Yf1
	Yfc2 := Yfc1
	Yfn2 := Yfn1
	A2 := A1
	E2 := E1
	B2 := B1
	Yk2 := Yk1
	if result.Calc.Strength.Auxiliary.Flange2 != nil {
		Yf2 = strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yf, 'G', 3, 64), "E", "*10^")
		Yfc2 = strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yfc, 'G', 3, 64), "E", "*10^")
		Yfn2 = strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yfn, 'G', 3, 64), "E", "*10^")
		A2 = strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.A, 'G', 3, 64), "E", "*10^")
		E2 = strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.E, 'G', 3, 64), "E", "*10^")
		B2 = strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.B, 'G', 3, 64), "E", "*10^")
		Yk2 = strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yk, 'G', 3, 64), "E", "*10^")
	}

	if data.TypeGasket == flange_model.GasketData_Oval {
		// фомула 4
		auxiliary.B0 = fmt.Sprintf("%s / 4", width)
		// фомула ?
		auxiliary.Dcp = fmt.Sprintf("%s - %s/2", dOut, width)

	} else {
		if !(data.Gasket.Width <= constants.Bp) {
			// фомула 3
			auxiliary.B0 = fmt.Sprintf("%.1f * sqrt(%s)", constants.B0, width)
		}
		// фомула 5
		auxiliary.Dcp = fmt.Sprintf("%s - %s", dOut, B0)
	}

	if data.TypeGasket == flange_model.GasketData_Soft {
		// Податливость прокладки
		auxiliary.Yp = fmt.Sprintf("(%s * %s) / (%s * %f * %s * %s)", thickness, compression, epsilon, math.Pi, Dcp, width)
	}
	// приложение К пояснение к формуле К.2
	auxiliary.Lb = fmt.Sprintf("%s + %s * %s", lenght, typeBolt, diameter)
	// формула К.2
	// Податливость болтов/шпилек
	auxiliary.Yb = fmt.Sprintf("%s / (%s * %s * %d)", Lb, epsilonAt20, area, count)

	// фомула 8
	// Суммарная площадь сечения болтов/шпилек
	auxiliary.A = fmt.Sprintf("%d * %s", count, area)

	flange1 := s.auxFlangeFormulas(req.FlangesData[0].Type, data.Flange1, result.Calc.Strength.Auxiliary.Flange1, result.Calc.Strength.Auxiliary.Dcp)
	flange2 := flange1
	auxiliary.Flange1 = flange1
	if len(req.FlangesData) > 1 {
		flange2 = s.auxFlangeFormulas(req.FlangesData[1].Type, data.Flange2, result.Calc.Strength.Auxiliary.Flange2, result.Calc.Strength.Auxiliary.Dcp)
		auxiliary.Flange2 = flange2
	}

	if !(data.TypeGasket == flange_model.GasketData_Oval || data.Type1 == flange_model.FlangeData_free || data.Type2 == flange_model.FlangeData_free) {
		// формула (Е.11)
		// Коэффициент жесткости
		auxiliary.Alpha = fmt.Sprintf("1 - (%s - (%s * %s * %s + %s * %s * %s)) / (%s + %s + (%s * (%s)^2 + %s * (%s)^2))",
			Yp, Yf1, E1, B1, Yf2, E2, B2, Yp, Yb, Yf1, B1, Yf2, B2)
	}

	dividend := fmt.Sprintf("%s + %s * %s * (%s + %s - (%s)^2 / %s) + %s * %s * (%s + %s - (%s)^2 / %s)",
		Yb, Yfn1, B1, B1, E1, E1, Dcp, Yfn2, B2, B2, E2, E2, Dcp)
	divider := fmt.Sprintf("%s + %s * (%s / %s)^2 + %s * (%s)^2 + %s * (%s)^2", Yb, Yp, db1, Dcp, Yfn1, B1, Yfn2, B2)

	if data.Type1 == flange_model.FlangeData_free {
		dividend += fmt.Sprintf(" + %s * (%s)^2", Yfc1, A1)
		divider += fmt.Sprintf(" + %s * (%s)^2", Yfc1, A1)
	}
	if data.Type2 == flange_model.FlangeData_free {
		dividend += fmt.Sprintf(" + %s * (%s)^2", Yfc2, A2)
		divider += fmt.Sprintf(" + %s * (%s)^2", Yfc2, A2)
	}

	// формула (Е.13)
	// Коэффициент жесткости фланцевого соединения нагруженного внешним изгибающим моментом
	auxiliary.AlphaM = fmt.Sprintf("%s / %s", dividend, divider)

	divider = fmt.Sprintf("%s + %s * %s / %s + (%s * %s / %s) * (%s)^2 + (%s * %s / %s) * (%s)^2",
		Yp, Yb, epsilonAt20, bEpsilon, Yf1, fEpAt201, fEp1, B1, Yf2, fEpAt202, fEp2, B2)

	if data.Type1 == flange_model.FlangeData_free {
		fEpsilonKAt201 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Ring.EpsilonKAt20, 'G', 3, 64), "E", "*10^")
		fEpsilonK1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Ring.EpsilonK, 'G', 3, 64), "E", "*10^")

		divider += fmt.Sprintf(" + (%s * %s / %s) * (%s)^2", Yk1, fEpsilonKAt201, fEpsilonK1, A1)
	}
	if data.Type2 == flange_model.FlangeData_free {
		fEpsilonKAt202 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Ring.EpsilonKAt20, 'G', 3, 64), "E", "*10^")
		fEpsilonK2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Ring.EpsilonK, 'G', 3, 64), "E", "*10^")

		divider += fmt.Sprintf(" + (%s * %s / %s) * (%s)^2", Yk2, fEpsilonKAt202, fEpsilonK2, A2)
	}

	// формула (Е.8)
	auxiliary.Gamma = fmt.Sprintf("1 / %s", divider)

	return auxiliary
}

func (s *FormulasService) auxFlangeFormulas(
	flangeType flange_model.FlangeData_Type,
	data *flange_model.FlangeResult,
	flangeRes *flange_model.CalcAuxiliary_Flange,
	Dcp float64,
) *flange_model.AuxiliaryFormulas_Flange {
	flange := &flange_model.AuxiliaryFormulas_Flange{}

	// перевод чисел в строки
	d6 := strconv.FormatFloat(data.D6, 'G', 3, 64)
	ds := strconv.FormatFloat(data.Ds, 'G', 3, 64)
	d := strconv.FormatFloat(data.D, 'G', 3, 64)
	dk := strconv.FormatFloat(data.Dk, 'G', 3, 64)
	dnk := strconv.FormatFloat(data.Dnk, 'G', 3, 64)
	dOut := strconv.FormatFloat(data.DOut, 'G', 3, 64)
	s1 := strconv.FormatFloat(data.S1, 'G', 3, 64)
	s0 := strconv.FormatFloat(data.S0, 'G', 3, 64)
	l := strconv.FormatFloat(data.L, 'G', 3, 64)
	h := strconv.FormatFloat(data.H, 'G', 3, 64)
	hk := strconv.FormatFloat(data.Hk, 'G', 3, 64)
	epsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	epsilonKAt20 := ""
	if data.Ring != nil {
		epsilonAt20 = strings.ReplaceAll(strconv.FormatFloat(data.Ring.EpsilonKAt20, 'G', 3, 64), "E", "*10^")
	}

	Dcp_ := strings.ReplaceAll(strconv.FormatFloat(Dcp, 'G', 3, 64), "E", "*10^")

	Beta := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Beta, 'G', 3, 64), "E", "*10^")
	X := strings.ReplaceAll(strconv.FormatFloat(flangeRes.X, 'G', 3, 64), "E", "*10^")
	Xi := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Xi, 'G', 3, 64), "E", "*10^")
	Se := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Se, 'G', 3, 64), "E", "*10^")
	L0 := strings.ReplaceAll(strconv.FormatFloat(flangeRes.L0, 'G', 3, 64), "E", "*10^")
	Lymda := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Lymda, 'G', 3, 64), "E", "*10^")
	BetaF := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaF, 'G', 3, 64), "E", "*10^")
	BetaT := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaT, 'G', 3, 64), "E", "*10^")
	BetaV := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaV, 'G', 3, 64), "E", "*10^")
	BetaU := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaU, 'G', 3, 64), "E", "*10^")
	Psik := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Psik, 'G', 3, 64), "E", "*10^")

	if flangeType != flange_model.FlangeData_free {
		// Плечи действия усилий в болтах/шпильках
		flange.B = fmt.Sprintf("0.5 * (%s - %s)", d6, Dcp_)
	} else {
		flange.A = fmt.Sprintf("0.5 * (%s - %s)", d6, ds)
		flange.B = fmt.Sprintf("0.5 * (%s - %s)", ds, Dcp_)
	}

	if flangeType == flange_model.FlangeData_welded {
		flange.X = fmt.Sprintf("%s / (sqrt(%s * %s))", l, d, s0)
		flange.Beta = fmt.Sprintf("%s / %s", s1, s0)
		// Коэффициент зависящий от соотношения размеров конической втулки фланца
		flange.Xi = fmt.Sprintf("1 + (%s - 1) * %s / (%s + (1 + %s)/4)", Beta, X, X, Beta)
		flange.Se = fmt.Sprintf("%s * %s", Xi, s0)
	}

	// Плечо усилия от действия давления на фланец
	flange.E = fmt.Sprintf("0.5 * (%s - %s - %s)", Dcp_, d, Se)
	// Параметр длины обечайки
	flange.L0 = fmt.Sprintf("sqrt(%s * %s)", d, s0)
	// Отношение наружного диаметра тарелки фланца к внутреннему диаметру
	flange.K = fmt.Sprintf("%s / %s", dOut, d)

	flange.Lymda = fmt.Sprintf("(%s * %s + %s)/(%s * %s) + (%s * (%s)^3) / (%s * %s * (%s)^2)", BetaF, h, L0, BetaT, L0, BetaV, h, BetaU, L0, s0)

	// Угловая податливость фланца при затяжке
	flange.Yf = fmt.Sprintf("(0.91 * %s) / (%s * %s * (%s)^2 * %s)", BetaV, epsilonAt20, Lymda, s0, L0)

	if flangeType == flange_model.FlangeData_free {
		flange.Psik = fmt.Sprintf("1.28 * lg(%s / %s)", dnk, dk)
		flange.Yk = fmt.Sprintf("1 / (%s * (%s)^3 * %s)", epsilonKAt20, hk, Psik)
	}

	if flangeType != flange_model.FlangeData_free {
		// Угловая податливость фланца нагруженного внешним изгибающим моментом
		flange.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * (%s)^3))", math.Pi, d6, epsilonAt20, dOut, h)
	} else {
		flange.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * (%s)^3))", math.Pi, ds, epsilonAt20, dOut, h)
		flange.Yfc = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * (%s)^3))", math.Pi, d6, epsilonKAt20, dnk, hk)
	}

	return flange
}

func (s *FormulasService) tightnessFormulas(data models.DataFlange, req *calc_api.FlangeRequest, result *calc_api.FlangeResponse,
) *flange_model.TightnessFormulas {
	tightness := &flange_model.TightnessFormulas{}

	// перевод чисел в строки
	axialForce := req.AxialForce
	bendingMoment := req.BendingMoment
	pressure := strconv.FormatFloat(req.Pressure, 'G', 3, 64)

	pres := strconv.FormatFloat(data.Gasket.Pres, 'G', 3, 64)
	m := strconv.FormatFloat(data.Gasket.M, 'G', 3, 64)

	sigmaAt20 := strconv.FormatFloat(data.Bolt.SigmaAt20, 'G', 3, 64)

	B0 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.B0, 'G', 3, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Dcp, 'G', 3, 64), "E", "*10^")
	Ab := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.A, 'G', 3, 64), "E", "*10^")
	Alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Alpha, 'G', 3, 64), "E", "*10^")
	AlphaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.AlphaM, 'G', 3, 64), "E", "*10^")

	Po := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Po, 'G', 3, 64), "E", "*10^")
	Qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Qd, 'G', 3, 64), "E", "*10^")
	Rp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Rp, 'G', 3, 64), "E", "*10^")
	Pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb1, 'G', 3, 64), "E", "*10^")
	Pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb2, 'G', 3, 64), "E", "*10^")
	Pb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb, 'G', 3, 64), "E", "*10^")

	// формула 6
	// Усилие необходимое для смятия прокладки при затяжке
	tightness.Po = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, Dcp, B0, pres)

	if req.Pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях
		tightness.Rp = fmt.Sprintf("%f * %s * %s * %s *|%s|", math.Pi, Dcp, B0, m, pressure)
	}

	// формула 9
	// Равнодействующая нагрузка от давления
	tightness.Qd = fmt.Sprintf("0.785 * (%s)^2 * %s", Dcp, pressure)

	// temp1 := float64(req.AxialForce) + 4*math.Abs(float64(req.BendingMoment))/aux.Dcp
	// temp2 := float64(req.AxialForce) - 4*math.Abs(float64(req.BendingMoment))/aux.Dcp
	temp1 := fmt.Sprintf("%d + 4 * |%d| / %s", axialForce, bendingMoment, Dcp)
	temp2 := fmt.Sprintf("%d - 4 * |%d| / %s", axialForce, bendingMoment, Dcp)

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	tightness.Qfm = fmt.Sprintf("max(%s; %s)", temp1, temp2)

	// minB := 0.4 * aux.A * data.Bolt.SigmaAt20
	minB := fmt.Sprintf("0.4 * %s * %s", Ab, sigmaAt20)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	tightness.Pb2 = fmt.Sprintf("max(%s; %s)", Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	tightness.Pb1 = fmt.Sprintf("%s * (%s + %d) + %s + 4 * %s * |%d| / %s", Alpha, Qd, axialForce, Rp, AlphaM, bendingMoment, Dcp)

	// Расчетная нагрузка на болты/шпильки фланцевых соединений
	// tightness.Pb = math.Max(tightness.Pb1, tightness.Pb2)
	tightness.Pb = fmt.Sprintf("max(%s; %s)", Pb1, Pb2)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений в рабочих условиях
	tightness.Pbr = fmt.Sprintf("%s + (1 - %s) * (%s + %d) + 4 * (1 - %s * |%d|) / %s", Pb, Alpha, Qd, axialForce, AlphaM, bendingMoment, Dcp)

	return tightness
}

func (s *FormulasService) staticResistanceCalculate(
	flange *flange_model.FlangeResult,
	flangeRes *flange_model.CalcAuxiliary_Flange,
	typeFlange flange_model.FlangeData_Type,
	staticRes *flange_model.CalcStaticResistance,
	data models.DataFlange,
	req *calc_api.FlangeRequest,
	Pb_, Pbr_, Qd_, Qfm_ float64,
) *flange_model.StaticResistanceFormulas {
	static := &flange_model.StaticResistanceFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Pressure, 'G', 3, 64)
	bendingMoment := req.BendingMoment
	axialForce := req.AxialForce

	fD6 := strconv.FormatFloat(flange.D6, 'G', 3, 64)
	fH := strconv.FormatFloat(flange.H, 'G', 3, 64)
	fD := strconv.FormatFloat(flange.D, 'G', 3, 64)
	fS0 := strconv.FormatFloat(flange.S0, 'G', 3, 64)
	fS1 := strconv.FormatFloat(flange.S1, 'G', 3, 64)
	fC := strconv.FormatFloat(flange.C, 'G', 3, 64)

	gM := strconv.FormatFloat(data.Gasket.M, 'G', 3, 64)

	count := data.Bolt.Count
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 3, 64)

	Pb := strings.ReplaceAll(strconv.FormatFloat(Pb_, 'G', 3, 64), "E", "*10^")
	Pbr := strings.ReplaceAll(strconv.FormatFloat(Pbr_, 'G', 3, 64), "E", "*10^")
	Qd := strings.ReplaceAll(strconv.FormatFloat(Qd_, 'G', 3, 64), "E", "*10^")
	Qfm := strings.ReplaceAll(strconv.FormatFloat(Qfm_, 'G', 3, 64), "E", "*10^")

	B := strings.ReplaceAll(strconv.FormatFloat(flangeRes.B, 'G', 3, 64), "E", "*10^")
	E := strings.ReplaceAll(strconv.FormatFloat(flangeRes.E, 'G', 3, 64), "E", "*10^")
	A := strings.ReplaceAll(strconv.FormatFloat(flangeRes.A, 'G', 3, 64), "E", "*10^")
	F := strings.ReplaceAll(strconv.FormatFloat(flangeRes.F, 'G', 3, 64), "E", "*10^")
	Lymda := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Lymda, 'G', 3, 64), "E", "*10^")
	L0 := strings.ReplaceAll(strconv.FormatFloat(flangeRes.L0, 'G', 3, 64), "E", "*10^")
	BetaF := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaF, 'G', 3, 64), "E", "*10^")
	BetaY := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaY, 'G', 3, 64), "E", "*10^")
	BetaZ := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaZ, 'G', 3, 64), "E", "*10^")

	Cf := strings.ReplaceAll(strconv.FormatFloat(staticRes.Cf, 'G', 3, 64), "E", "*10^")
	MM := strings.ReplaceAll(strconv.FormatFloat(staticRes.MM, 'G', 3, 64), "E", "*10^")
	Dzv := strings.ReplaceAll(strconv.FormatFloat(staticRes.Dzv, 'G', 3, 64), "E", "*10^")
	Mp := strings.ReplaceAll(strconv.FormatFloat(staticRes.Mp, 'G', 3, 64), "E", "*10^")
	SigmaM1 := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaM1, 'G', 3, 64), "E", "*10^")
	SigmaR := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaR, 'G', 3, 64), "E", "*10^")
	SigmaRp := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaRp, 'G', 3, 64), "E", "*10^")
	SigmaP1 := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaP1, 'G', 3, 64), "E", "*10^")

	temp1 := fmt.Sprintf("%f * %s / %d", math.Pi, fD6, count)
	temp2 := fmt.Sprintf("2 * %s + 6 * %s / (%s + 0.5)", diameter, fH, gM)

	// Коэффициент учитывающий изгиб тарелки фланца между болтами шпильками
	static.Cf = fmt.Sprintf("max(1; sqrt((%s) / (%s)))", temp1, temp2)

	// Приведенный диаметр приварного встык фланца с конической или прямой втулкой
	if typeFlange == flange_model.FlangeData_welded && flange.D <= 20*flange.S1 {
		if flangeRes.F > 1 {
			static.Dzv = fmt.Sprintf("%s + %s", fD, fS0)
		} else {
			static.Dzv = fmt.Sprintf("%s + %s", fD, fS1)
		}
	}

	// Расчетный изгибающий момент действующий на фланец при затяжке
	static.MM = fmt.Sprintf("%s * %s * %s", Cf, Pb, B)
	// Расчетный изгибающий момент действующий на фланец в рабочих условиях
	static.Mp = fmt.Sprintf("%s * max(%s * %s + (%s + %s) * %s; |%s + %s| * %s)", Cf, Pbr, B, Qd, Qfm, E, Qd, Qfm, E)

	if typeFlange == flange_model.FlangeData_free {
		static.MMk = fmt.Sprintf("%s * %s * %s", Cf, Pb, A)
		static.Mpk = fmt.Sprintf("%s * %s * %s", Cf, Pbr, A)
	}

	// Меридиональное изгибное напряжение во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке бурта свободного фланца
	if typeFlange == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		static.SigmaM1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", MM, Lymda, fS1, fC, Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		static.SigmaM0 = fmt.Sprintf("%s * %s", F, SigmaM1)
	} else {
		static.SigmaM1 = fmt.Sprintf("%s/ (%s * (%s - %s)^2 * %s)", MM, Lymda, fS0, fC, Dzv)
		static.SigmaM0 = SigmaM1
	}

	// Радиальное напряжение в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в условиях затяжки
	static.SigmaR = fmt.Sprintf("((1.33 * %s * %s + %s) / (%s * (%s)^2 * %s * %s)) * %s", BetaF, fH, L0, Lymda, fH, L0, fD, MM)
	// Окружное напряжение в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в условиях затяжки
	static.SigmaT = fmt.Sprintf("%s * %s / ((%s)^2 * %s) - %s * %s", BetaY, MM, fH, fD, BetaZ, SigmaR)

	if typeFlange == flange_model.FlangeData_free {
		Hk := strconv.FormatFloat(flange.Hk, 'G', 3, 64)
		Dk := strconv.FormatFloat(flange.Dk, 'G', 3, 64)
		MMk := strings.ReplaceAll(strconv.FormatFloat(staticRes.MMk, 'G', 3, 64), "E", "*10^")

		static.SigmaK = fmt.Sprintf("%s * %s / ((%s)^2 * %s)", BetaY, MMk, Hk, Dk)
	}

	// Меридиональные изгибные напряжения во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке
	// трубе бурта свободного фланца в рабочих условиях
	if typeFlange == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		static.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", Mp, Lymda, fS1, fC, Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		static.SigmaP0 = fmt.Sprintf("%s * %s", F, SigmaP1)
	} else {
		static.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", Mp, Lymda, fS0, fC, Dzv)
		static.SigmaP0 = SigmaP1
	}

	if typeFlange == flange_model.FlangeData_welded {
		temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, fD, fS1, fS1, fC)
		// формула (ф. 37)
		static.SigmaMp = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d + 4 * |%d / (%s + %s)|) / %s", fD, pressure, axialForce, bendingMoment, fD, fS1, temp)
	}

	temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, fD, fS0, fS0, fC)
	// Меридиональные мембранные напряжения во втулке приварного встык фланца обечайке трубе
	// плоского фланца или обечайке трубе бурта свободного фланца в рабочих условиях
	// формула (ф. 37)
	// - для приварных встык фланцев с конической втулкой в сечении S1
	static.SigmaMp0 = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d + 4 * |%d / (%s + %s)|) / %s", fD, pressure, axialForce, bendingMoment, fD, fS0, temp)

	// Окружные мембранные напряжения от действия давления во втулке приварного встык фланца обечайке
	// трубе плоского фланца или обечайке трубе бурта свободного фланца в сечениии S0
	static.SigmaMop = fmt.Sprintf("%s * %s / (2 * (%s - %s))", pressure, fD, fS0, fC)

	// Напряжения в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в рабочих условиях
	// - радиальные напряжения
	static.SigmaRp = fmt.Sprintf("((1.33 * %s * %s + %s) / (%s * (%s)^2 * %s * %s)) * %s", BetaF, fH, L0, Lymda, fH, L0, fD, Mp)
	// - окружное напряжения
	static.SigmaTp = fmt.Sprintf("%s * %s / ((%s)^2 * %s) - %s * %s", BetaY, Mp, fH, fD, BetaZ, SigmaRp)

	if typeFlange == flange_model.FlangeData_free {
		Hk := strconv.FormatFloat(flange.Hk, 'G', 3, 64)
		Dk := strconv.FormatFloat(flange.Dk, 'G', 3, 64)

		static.SigmaKp = fmt.Sprintf("%s * %s / ((%s)^2 * %s)", BetaY, Mp, Hk, Dk)
	}

	return static
}

func (s *FormulasService) conditionsForStrengthCalculate(
	flangeType flange_model.FlangeData_Type,
	flange *flange_model.FlangeResult,
	calcFlange *flange_model.CalcAuxiliary_Flange,
	static *flange_model.CalcStaticResistance,
	cond *flange_model.CalcConditionsForStrength,
	isWork, isTemp bool,
) *flange_model.ConditionsForStrengthFormulas {
	conditions := &flange_model.ConditionsForStrengthFormulas{}

	// перевод чисел в строки
	fEpsilon := strings.ReplaceAll(strconv.FormatFloat(flange.Epsilon, 'G', 3, 64), "E", "*10^")
	fEpsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	fSigmaAt20 := strconv.FormatFloat(flange.SigmaAt20, 'G', 3, 64)
	fSigma := strconv.FormatFloat(flange.Sigma, 'G', 3, 64)
	fSigmaMAt20 := strconv.FormatFloat(flange.SigmaMAt20, 'G', 3, 64)
	fSigmaM := strconv.FormatFloat(flange.SigmaM, 'G', 3, 64)
	fSigmaRAt20 := strconv.FormatFloat(flange.SigmaRAt20, 'G', 3, 64)
	fSigmaR := strconv.FormatFloat(flange.SigmaR, 'G', 3, 64)

	Yf := strings.ReplaceAll(strconv.FormatFloat(calcFlange.Yf, 'G', 3, 64), "E", "*10^")

	Mp := strings.ReplaceAll(strconv.FormatFloat(static.Mp, 'G', 3, 64), "E", "*10^")
	SigmaM1 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaM1, 'G', 3, 64), "E", "*10^")
	SigmaM0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaM0, 'G', 3, 64), "E", "*10^")
	SigmaR := strings.ReplaceAll(strconv.FormatFloat(static.SigmaR, 'G', 3, 64), "E", "*10^")
	SigmaT := strings.ReplaceAll(strconv.FormatFloat(static.SigmaT, 'G', 3, 64), "E", "*10^")
	SigmaP1 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaP1, 'G', 3, 64), "E", "*10^")
	SigmaP0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaP0, 'G', 3, 64), "E", "*10^")
	SigmaMp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMp, 'G', 3, 64), "E", "*10^")
	SigmaMp0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMp0, 'G', 3, 64), "E", "*10^")
	SigmaMpm := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMpm, 'G', 3, 64), "E", "*10^")
	SigmaMpm0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMpm0, 'G', 3, 64), "E", "*10^")
	SigmaRp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaRp, 'G', 3, 64), "E", "*10^")
	SigmaTp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaTp, 'G', 3, 64), "E", "*10^")
	SigmaMop := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMop, 'G', 3, 64), "E", "*10^")

	teta := map[bool]float64{
		true:  constants.WorkTeta,
		false: constants.TestTeta,
	}
	var Ks float64
	if calcFlange.K <= constants.MinK {
		Ks = constants.MinKs
	} else if calcFlange.K >= constants.MaxK {
		Ks = constants.MaxKs
	} else {
		Ks = ((calcFlange.K-constants.MinK)/(constants.MaxK-constants.MinK))*(constants.MaxKs-constants.MinKs) + constants.MinKs
	}
	Kt := map[bool]float64{
		true:  constants.TempKt,
		false: constants.Kt,
	}

	var DTeta float64
	if flangeType == flange_model.FlangeData_welded {
		if flange.D <= constants.MinD {
			DTeta = constants.MinDTeta
		} else if flange.D > constants.MaxD {
			DTeta = constants.MaxDTeta
		} else {
			DTeta = ((flange.D-constants.MinD)/(constants.MaxD-constants.MinD))*
				(constants.MaxDTeta-constants.MinDTeta) + constants.MinDTeta
		}
	} else {
		DTeta = constants.MaxDTeta
	}
	tmp := strconv.FormatFloat(DTeta, 'G', 3, 64)
	DTeta_ := fmt.Sprintf("%.1f * %s", teta[isWork], tmp)

	conditions.Teta = fmt.Sprintf("%s * %s * %s / %s", Mp, Yf, fEpsilonAt20, fEpsilon)
	conditions.CondTeta = &flange_model.ConditionFormulas{
		X: conditions.Teta,
		Y: DTeta_,
	}

	if flangeType == flange_model.FlangeData_free {
		fEpsilonK := strconv.FormatFloat(flange.Ring.EpsilonK, 'G', 3, 64)
		fEpsilonKAt20 := strconv.FormatFloat(flange.Ring.EpsilonKAt20, 'G', 3, 64)
		Yk := strings.ReplaceAll(strconv.FormatFloat(calcFlange.Yk, 'G', 3, 64), "E", "*10^")
		Mpk := strings.ReplaceAll(strconv.FormatFloat(static.Mpk, 'G', 3, 64), "E", "*10^")

		DTetaK := fmt.Sprintf("%.1f * 0.02", teta[isWork])
		conditions.TetaK = fmt.Sprintf("%s * %s * %s / %s", Mpk, Yk, fEpsilonKAt20, fEpsilonK)
		conditions.CondTetaK = &flange_model.ConditionFormulas{
			X: conditions.TetaK,
			Y: DTetaK,
		}
	}

	//* Условия статической прочности фланцев
	if flangeType == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		Max1 := fmt.Sprintf("max(|%s + %s|; |%s + %s|)", SigmaM1, SigmaR, SigmaM1, SigmaT)
		max1Y := fmt.Sprintf("%f * %.1f * %s", Ks, Kt[isTemp], fSigmaMAt20)

		t1 := fmt.Sprintf("max(|%s - %s + %s|; |%s - %s + %s|)", SigmaP1, SigmaMp, SigmaRp, SigmaP1, SigmaMpm, SigmaRp)
		t2 := fmt.Sprintf("max(|%s - %s + %s|; |%s - %s + %s|)", SigmaP1, SigmaMp, SigmaTp, SigmaP1, SigmaMpm, SigmaTp)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("max(|%s + %s|; |%s + %s|)", SigmaP1, SigmaMp, SigmaP1, SigmaMpm)

		Max2 := fmt.Sprintf("max(%s; %s)", t1, t2)
		max2Y := fmt.Sprintf("%f * %.1f * %s", Ks, Kt[isTemp], fSigmaM)

		max3Y := fmt.Sprintf("1.3 * %s", fSigmaRAt20)

		t1 = fmt.Sprintf("max(|%s + %s|; |%s - %s|)", SigmaP0, SigmaMp0, SigmaP0, SigmaMp0)
		t2 = fmt.Sprintf("max(|%s + %s|; |%s - %s|)", SigmaP0, SigmaMpm0, SigmaP0, SigmaMpm0)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("max(|0.3 * %s + %s); |0.3 * %s - %s|)", SigmaP0, SigmaMop, SigmaP0, SigmaMop)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("max(|0.7 * %s + (%s - %s)|; |0.7 * %s - (%s - %s)|)", SigmaP0, SigmaMp0, SigmaMop, SigmaP0, SigmaMp0, SigmaMop)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("max(|0.7 * %s + (%s - %s)|; |0.7 * %s - (%s - %s)|)", SigmaP0, SigmaMpm0, SigmaMop, SigmaP0, SigmaMpm0, SigmaMop)

		Max4 := fmt.Sprintf("max(%s; %s)", t1, t2)
		max4Y := fmt.Sprintf("1.3 * %s", fSigmaR)

		conditions.Max1 = &flange_model.ConditionFormulas{
			X: Max1,
			Y: max1Y,
		}
		conditions.Max2 = &flange_model.ConditionFormulas{
			X: Max2,
			Y: max2Y,
		}
		conditions.Max3 = &flange_model.ConditionFormulas{
			X: SigmaM0,
			Y: max3Y,
		}
		conditions.Max4 = &flange_model.ConditionFormulas{
			X: Max4,
			Y: max4Y,
		}
	} else {
		Max5 := fmt.Sprintf("max(|%s + %s|; |%s + %s|)", SigmaM0, SigmaR, SigmaM0, SigmaT)

		t1 := fmt.Sprintf("max(|%s - %s + %s); |%s - %s + %s|)", SigmaP0, SigmaMp0, SigmaTp, SigmaP0, SigmaMpm0, SigmaTp)
		t2 := fmt.Sprintf("max(|%s - %s + %s|; |%s - %s + %s|)", SigmaP0, SigmaMp0, SigmaRp, SigmaP0, SigmaMpm0, SigmaRp)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("max(|%s + %s|; |%s + %s|)", SigmaP0, SigmaMp0, SigmaP0, SigmaMpm0)

		Max6 := fmt.Sprintf("max(%s; %s)", t1, t2)

		conditions.Max5 = &flange_model.ConditionFormulas{
			X: Max5,
			Y: fSigmaAt20,
		}
		conditions.Max6 = &flange_model.ConditionFormulas{
			X: Max6,
			Y: fSigma,
		}
	}

	max7 := fmt.Sprintf("|%s|; |%s|", SigmaMp0, SigmaMpm0)
	Max7 := fmt.Sprintf("max(%s; |%s|)", max7, SigmaMop)
	Max8 := fmt.Sprintf("max(|%s|; |%s|)", SigmaR, SigmaT)
	max8Y := fmt.Sprintf("%.1f * %s", Kt[isTemp], fSigmaAt20)
	Max9 := fmt.Sprintf("max(|%s|; |%s|)", SigmaRp, SigmaTp)
	max9Y := fmt.Sprintf("%.1f * %s", Kt[isTemp], fSigma)

	conditions.Max7 = &flange_model.ConditionFormulas{X: Max7, Y: fSigma}
	conditions.Max8 = &flange_model.ConditionFormulas{X: Max8, Y: max8Y}
	conditions.Max9 = &flange_model.ConditionFormulas{X: Max9, Y: max9Y}

	if flangeType == flange_model.FlangeData_free {
		SigmaK := strings.ReplaceAll(strconv.FormatFloat(static.SigmaK, 'G', 3, 64), "E", "*10^")
		SigmaKp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaKp, 'G', 3, 64), "E", "*10^")
		fSigmaKAt20 := strconv.FormatFloat(flange.Ring.SigmaKAt20, 'G', 3, 64)
		fSigmaK := strconv.FormatFloat(flange.Ring.SigmaK, 'G', 3, 64)

		max10Y := fmt.Sprintf("%.1f * %s", Kt[isTemp], fSigmaKAt20)
		max11Y := fmt.Sprintf("%.1f * %s", Kt[isTemp], fSigmaK)

		conditions.Max10 = &flange_model.ConditionFormulas{X: SigmaK, Y: max10Y}
		conditions.Max11 = &flange_model.ConditionFormulas{X: SigmaKp, Y: max11Y}
	}

	return conditions
}

func (s *FormulasService) tightnessLoadCalculate(
	data models.DataFlange,
	req *calc_api.FlangeRequest,
	result *calc_api.FlangeResponse,
) *flange_model.TightnessLoadFormulas {
	tightness := &flange_model.TightnessLoadFormulas{}

	// перевод чисел в строки
	fAlpha1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.AlphaF, 'G', 3, 64), "E", "*10^")
	fH1 := strconv.FormatFloat(data.Flange1.H, 'G', 3, 64)
	fT1 := strconv.FormatFloat(data.Flange1.Tf, 'G', 3, 64)
	fAlpha2 := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.AlphaF, 'G', 3, 64), "E", "*10^")
	fH2 := strconv.FormatFloat(data.Flange2.H, 'G', 3, 64)
	fT2 := strconv.FormatFloat(data.Flange2.Tf, 'G', 3, 64)

	bAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Alpha, 'G', 3, 64), "E", "*10^")
	bTemp := strconv.FormatFloat(data.Bolt.Temp, 'G', 3, 64)

	Gamma := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Gamma, 'G', 3, 64), "E", "*10^")
	Alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Alpha, 'G', 3, 64), "E", "*10^")
	AlphaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.AlphaM, 'G', 3, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Dcp, 'G', 3, 64), "E", "*10^")

	Pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb1, 'G', 3, 64), "E", "*10^")
	Pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb2, 'G', 3, 64), "E", "*10^")
	Pb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb, 'G', 3, 64), "E", "*10^")
	Qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Qd, 'G', 3, 64), "E", "*10^")

	Qt := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.TightnessLoad.Qt, 'G', 3, 64), "E", "*10^")

	var temp1, temp2 string
	if req.IsUseWasher {
		wAlpha1 := strings.ReplaceAll(strconv.FormatFloat(data.Washer1.Alpha, 'G', 3, 64), "E", "*10^")
		wThick1 := strconv.FormatFloat(data.Washer1.Thickness, 'G', 3, 64)
		wAlpha2 := strings.ReplaceAll(strconv.FormatFloat(data.Washer2.Alpha, 'G', 3, 64), "E", "*10^")
		wThick2 := strconv.FormatFloat(data.Washer2.Thickness, 'G', 3, 64)

		temp1 = fmt.Sprintf("(%s * %s + %s * %s) * (%s - 20) + (%s * %s + %s * %s)*(%s - 20)",
			fAlpha1, fH1, wAlpha1, wThick1, fT1, fAlpha2, fH2, wAlpha2, wThick2, fT2)
	} else {
		temp1 = fmt.Sprintf("%s * %s * (%s - 20) + %s * %s * (%s - 20)", fAlpha1, fH1, fT1, fAlpha2, fH2, fT2)
	}
	temp2 = fmt.Sprintf("%s + %s", fH1, fH2)

	if data.Type1 == flange_model.FlangeData_free {
		fAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Flange1.Ring.AlphaK, 'G', 3, 64), "E", "*10^")
		fH := strconv.FormatFloat(data.Flange1.Hk, 'G', 3, 64)
		fT := strconv.FormatFloat(data.Flange1.Ring.Tk, 'G', 3, 64)

		temp1 += fmt.Sprintf("%s * %s * (%s - 20)", fAlpha, fH, fT)
		temp2 += fH
	}
	if data.Type2 == flange_model.FlangeData_free {
		fAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Flange2.Ring.AlphaK, 'G', 3, 64), "E", "*10^")
		fH := strconv.FormatFloat(data.Flange2.Hk, 'G', 3, 64)
		fT := strconv.FormatFloat(data.Flange2.Ring.Tk, 'G', 3, 64)

		temp1 += fmt.Sprintf("%s * %s * (%s - 20)", fAlpha, fH, fT)
		temp2 += fH
	}
	if req.IsEmbedded {
		eAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Embed.Alpha, 'G', 3, 64), "E", "*10^")
		eThick := strconv.FormatFloat(data.Embed.Thickness, 'G', 3, 64)
		temp := strconv.FormatFloat(req.Temp, 'G', 3, 64)

		temp1 += fmt.Sprintf("%s * %s * (%s - 20)", eAlpha, eThick, temp)
		temp2 += eThick
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	tightness.Qt = fmt.Sprintf("%s * (%s - %s * %s * (%s - 20))", Gamma, temp1, bAlpha, temp2, bTemp)

	tightness.Pb1 = fmt.Sprintf("max(%s; %s - %s)", Pb1, Pb1, Qt)
	tightness.Pb = fmt.Sprintf("max(%s; %s)", Pb1, Pb2)
	tightness.Pbr = fmt.Sprintf("%s + (1 - %s) * (%s + %d) + %s + 4 * (1 - %s * |%d|) / %s",
		Pb, Alpha, Qd, req.AxialForce, Qt, AlphaM, req.BendingMoment, Dcp)

	return tightness
}
