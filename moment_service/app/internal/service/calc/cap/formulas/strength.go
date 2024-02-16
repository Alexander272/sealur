package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *FormulasService) strengthFormulas(
	req *calc_api.CapRequest,
	d models.DataCap,
	result *calc_api.CapResponse,
) *cap_model.Formulas_Strength {
	auxiliary := s.auxiliaryFormulas(d, req, result)
	tightness := s.tightnessFormulas(d, req, result)
	bolt1 := s.boltStrengthFormulas(
		req, d, result,
		result.Calc.Strength.Tightness.Pb,
		result.Calc.Strength.Tightness.Pbr,
		result.Calc.Strength.Auxiliary.A,
		result.Calc.Strength.Auxiliary.Dcp,
		false,
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

	static1 := s.staticResistanceFormulas(
		result.Flange,
		result.Calc.Strength.Auxiliary.Flange,
		req.FlangeData.Type,
		result.Calc.Strength.StaticResistance1,
		d, req,
		result.Calc.Strength.Tightness.Pb,
		result.Calc.Strength.Tightness.Pbr,
		result.Calc.Strength.Tightness.Qd,
		result.Calc.Strength.Tightness.Qfm,
	)
	conditions1 := s.conditionsForStrengthFormulas(
		req.FlangeData.Type,
		result.Flange,
		result.Calc.Strength.Auxiliary.Flange,
		result.Calc.Strength.StaticResistance1,
		result.Calc.Strength.ConditionsForStrength1,
		req.Data.IsWork, false,
	)

	tigLoad := s.tightnessLoadFormulas(d, req, result)
	bolt2 := s.boltStrengthFormulas(
		req, d, result,
		result.Calc.Strength.TightnessLoad.Pb,
		result.Calc.Strength.TightnessLoad.Pbr,
		result.Calc.Strength.Auxiliary.A,
		result.Calc.Strength.Auxiliary.Dcp,
		true,
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

	static2 := s.staticResistanceFormulas(
		result.Flange,
		result.Calc.Strength.Auxiliary.Flange,
		req.FlangeData.Type,
		result.Calc.Strength.StaticResistance2,
		d, req,
		result.Calc.Strength.TightnessLoad.Pb,
		result.Calc.Strength.TightnessLoad.Pbr,
		result.Calc.Strength.Tightness.Qd,
		result.Calc.Strength.Tightness.Qfm,
	)
	conditions2 := s.conditionsForStrengthFormulas(
		req.FlangeData.Type,
		result.Flange,
		result.Calc.Strength.Auxiliary.Flange,
		result.Calc.Strength.StaticResistance2,
		result.Calc.Strength.ConditionsForStrength2,
		req.Data.IsWork, true,
	)

	deformation := &cap_model.DeformationFormulas{
		B0:  auxiliary.B0,
		Dcp: auxiliary.Dcp,
		Po:  tightness.Po,
		Rp:  tightness.Rp,
	}
	forces := &cap_model.ForcesInBoltsFormulas{
		A:     auxiliary.A,
		Qd:    tightness.Qd,
		Qfm:   tightness.Qfm,
		Qt:    tigLoad.Qt,
		Pb:    tigLoad.Pb,
		Alpha: auxiliary.Alpha,
		Pb1:   tigLoad.Pb1,
		Pb2:   tightness.Pb2,
		Pbr:   tigLoad.Pbr,
	}
	finalMoment := &cap_model.MomentFormulas{}
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

	formulas := &cap_model.Formulas_Strength{
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

func (s *FormulasService) auxiliaryFormulas(data models.DataCap, req *calc_api.CapRequest, result *calc_api.CapResponse) *cap_model.AuxiliaryFormulas {
	auxiliary := &cap_model.AuxiliaryFormulas{}

	// перевод чисел в строки
	width := strconv.FormatFloat(data.Gasket.Width, 'G', 5, 64)
	thickness := strconv.FormatFloat(data.Gasket.Thickness, 'G', 5, 64)
	dOut := strconv.FormatFloat(data.Gasket.DOut, 'G', 5, 64)
	epsilon := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Epsilon, 'G', 5, 64), "E", "*10^")
	compression := strconv.FormatFloat(data.Gasket.Compression, 'G', 5, 64)

	length := strconv.FormatFloat(data.Bolt.Length, 'G', 5, 64)
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 5, 64)
	area := strconv.FormatFloat(data.Bolt.Area, 'G', 5, 64)
	epsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.EpsilonAt20, 'G', 5, 64), "E", "*10^")
	bEpsilon := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Epsilon, 'G', 5, 64), "E", "*10^")
	count := data.Bolt.Count

	fEpAt201 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.EpsilonAt20, 'G', 5, 64), "E", "*10^")
	fEpAt202 := strings.ReplaceAll(strconv.FormatFloat(data.Cap.EpsilonAt20, 'G', 5, 64), "E", "*10^")
	fEp1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Epsilon, 'G', 5, 64), "E", "*10^")
	fEp2 := strings.ReplaceAll(strconv.FormatFloat(data.Cap.Epsilon, 'G', 5, 64), "E", "*10^")

	typeBolt := strconv.FormatFloat(s.typeBolt[req.Data.Type.String()], 'G', 5, 64)

	B0 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.B0, 'G', 5, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Dcp, 'G', 5, 64), "E", "*10^")
	Lb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Lb, 'G', 5, 64), "E", "*10^")

	Yp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Yp, 'G', 5, 64), "E", "*10^")
	Yb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Yb, 'G', 5, 64), "E", "*10^")
	Yf1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange.Yf, 'G', 5, 64), "E", "*10^")
	Yk1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange.Yk, 'G', 5, 64), "E", "*10^")
	A1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange.A, 'G', 5, 64), "E", "*10^")
	E1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange.E, 'G', 5, 64), "E", "*10^")
	B1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange.B, 'G', 5, 64), "E", "*10^")
	Y := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Cap.Y, 'G', 5, 64), "E", "*10^")

	if data.TypeGasket == cap_model.GasketData_Oval {
		// формула 4
		auxiliary.B0 = fmt.Sprintf("%s / 4", width)
		// формула ?
		auxiliary.Dcp = fmt.Sprintf("%s - %s/2", dOut, width)

	} else {
		if !(data.Gasket.Width <= constants.Bp) {
			// формула 3
			auxiliary.B0 = fmt.Sprintf("%.1f * sqrt(%s)", constants.B0, width)
		}
		// формула 5
		auxiliary.Dcp = fmt.Sprintf("%s - %s", dOut, B0)
	}

	if data.TypeGasket == cap_model.GasketData_Soft {
		// Податливость прокладки
		auxiliary.Yp = fmt.Sprintf("(%s * %s) / (%s * %f * %s * %s)", thickness, compression, epsilon, math.Pi, Dcp, width)
	}
	// приложение К пояснение к формуле К.2
	auxiliary.Lb = fmt.Sprintf("%s + %s * %s", length, typeBolt, diameter)
	// формула К.2
	// Податливость болтов/шпилек
	auxiliary.Yb = fmt.Sprintf("%s / (%s * %s * %d)", Lb, epsilonAt20, area, count)

	// формула 8
	// Суммарная площадь сечения болтов/шпилек
	auxiliary.A = fmt.Sprintf("%d * %s", count, area)

	flange := s.auxFlangeFormulas(req.FlangeData.Type, data.Flange, result.Calc.Strength.Auxiliary.Flange, result.Calc.Strength.Auxiliary.Dcp)
	auxiliary.Flange = flange

	cap := s.auxCapFormulas(req.CapData.Type, data.Cap, data.Flange, result.Calc.Strength.Auxiliary.Cap, result.Calc.Strength.Auxiliary.Dcp)
	auxiliary.Cap = cap

	if !(data.TypeGasket == cap_model.GasketData_Oval || data.FlangeType == cap_model.FlangeData_free) {
		// формула (Е.11)
		// Коэффициент жесткости
		auxiliary.Alpha = fmt.Sprintf("1 - (%s - (%s * %s + %s * %s) * %s)/(%s + %s + (%s + %s) * (%s)^2)",
			Yp, Yf1, E1, Y, B1, B1, Yp, Yb, Yf1, Y, B1)
	}

	divider := fmt.Sprintf("%s + %s * %s/%s + (%s * %s/%s) * (%s)^2 + (%s * %s/%s) * (%s)^2",
		Yp, Yb, epsilonAt20, bEpsilon, Yf1, fEpAt201, fEp1, B1, Y, fEpAt202, fEp2, B1)

	if data.FlangeType == cap_model.FlangeData_free {
		fEpsilonKAt201 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Ring.EpsilonAt20, 'G', 5, 64), "E", "*10^")
		fEpsilonK1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Ring.Epsilon, 'G', 5, 64), "E", "*10^")

		divider += fmt.Sprintf(" + (%s * %s / %s) * (%s)^2", Yk1, fEpsilonKAt201, fEpsilonK1, A1)
	}

	// формула (Е.8)
	auxiliary.Gamma = fmt.Sprintf("1 / %s", divider)

	return auxiliary
}

func (s *FormulasService) auxFlangeFormulas(
	flangeType cap_model.FlangeData_Type,
	data *cap_model.FlangeResult,
	flangeRes *cap_model.CalcAuxiliary_Flange,
	Dcp float64,
) *cap_model.AuxiliaryFormulas_Flange {
	flange := &cap_model.AuxiliaryFormulas_Flange{}

	// перевод чисел в строки
	d6 := strconv.FormatFloat(data.D6, 'G', 5, 64)
	d := strconv.FormatFloat(data.D, 'G', 5, 64)
	dOut := strconv.FormatFloat(data.DOut, 'G', 5, 64)
	s1 := strconv.FormatFloat(data.S1, 'G', 5, 64)
	s0 := strconv.FormatFloat(data.S0, 'G', 5, 64)
	l := strconv.FormatFloat(data.L, 'G', 5, 64)
	h := strconv.FormatFloat(data.H, 'G', 5, 64)
	epsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.EpsilonAt20, 'G', 5, 64), "E", "*10^")
	var epsilonKAt20, dk, dnk, hk, ds string
	if data.Ring != nil {
		epsilonAt20 = strings.ReplaceAll(strconv.FormatFloat(data.Ring.EpsilonAt20, 'G', 5, 64), "E", "*10^")
		dk = strconv.FormatFloat(data.Ring.Dk, 'G', 5, 64)
		dnk = strconv.FormatFloat(data.Ring.Dnk, 'G', 5, 64)
		hk = strconv.FormatFloat(data.Ring.Hk, 'G', 5, 64)
		ds = strconv.FormatFloat(data.Ring.Ds, 'G', 5, 64)
	}

	Dcp_ := strings.ReplaceAll(strconv.FormatFloat(Dcp, 'G', 5, 64), "E", "*10^")

	Beta := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Beta, 'G', 5, 64), "E", "*10^")
	X := strings.ReplaceAll(strconv.FormatFloat(flangeRes.X, 'G', 5, 64), "E", "*10^")
	Xi := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Xi, 'G', 5, 64), "E", "*10^")
	Se := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Se, 'G', 5, 64), "E", "*10^")
	L0 := strings.ReplaceAll(strconv.FormatFloat(flangeRes.L0, 'G', 5, 64), "E", "*10^")
	Lymda := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Lymda, 'G', 5, 64), "E", "*10^")
	BetaF := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaF, 'G', 5, 64), "E", "*10^")
	BetaT := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaT, 'G', 5, 64), "E", "*10^")
	BetaV := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaV, 'G', 5, 64), "E", "*10^")
	BetaU := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaU, 'G', 5, 64), "E", "*10^")
	Psik := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Psik, 'G', 5, 64), "E", "*10^")

	if flangeType != cap_model.FlangeData_free {
		// Плечи действия усилий в болтах/шпильках
		flange.B = fmt.Sprintf("0.5 * (%s - %s)", d6, Dcp_)
	} else {
		flange.A = fmt.Sprintf("0.5 * (%s - %s)", d6, ds)
		flange.B = fmt.Sprintf("0.5 * (%s - %s)", ds, Dcp_)
	}

	if flangeType == cap_model.FlangeData_welded {
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

	if flangeType == cap_model.FlangeData_free {
		flange.Psik = fmt.Sprintf("1.28 * lg(%s / %s)", dnk, dk)
		flange.Yk = fmt.Sprintf("1 / (%s * (%s)^3 * %s)", epsilonKAt20, hk, Psik)
	}

	if flangeType != cap_model.FlangeData_free {
		// Угловая податливость фланца нагруженного внешним изгибающим моментом
		flange.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * (%s)^3))", math.Pi, d6, epsilonAt20, dOut, h)
	} else {
		flange.Yfn = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * (%s)^3))", math.Pi, ds, epsilonAt20, dOut, h)
		flange.Yfc = fmt.Sprintf("(%f/4)^3 * (%s / (%s * %s * (%s)^3))", math.Pi, d6, epsilonKAt20, dnk, hk)
	}

	return flange
}

func (s *FormulasService) auxCapFormulas(
	capType cap_model.CapData_Type,
	data *cap_model.CapResult,
	flange *cap_model.FlangeResult,
	cap *cap_model.CalcAuxiliary_Cap,
	Dcp float64,
) *cap_model.AuxiliaryFormulas_Cap {
	formulas := &cap_model.AuxiliaryFormulas_Cap{}

	Dcp_ := strings.ReplaceAll(strconv.FormatFloat(Dcp, 'G', 5, 64), "E", "*10^")
	DOut_ := strings.ReplaceAll(strconv.FormatFloat(flange.DOut, 'G', 5, 64), "E", "*10^")
	S0_ := strings.ReplaceAll(strconv.FormatFloat(flange.S0, 'G', 5, 64), "E", "*10^")
	D_ := strings.ReplaceAll(strconv.FormatFloat(flange.D, 'G', 5, 64), "E", "*10^")
	h_ := strings.ReplaceAll(strconv.FormatFloat(flange.H, 'G', 5, 64), "E", "*10^")

	k := strings.ReplaceAll(strconv.FormatFloat(cap.K, 'G', 5, 64), "E", "*10^")
	x := strings.ReplaceAll(strconv.FormatFloat(cap.X, 'G', 5, 64), "E", "*10^")
	eAt20 := strings.ReplaceAll(strconv.FormatFloat(data.EpsilonAt20, 'G', 5, 64), "E", "*10^")

	if capType == cap_model.CapData_flat {
		H := strings.ReplaceAll(strconv.FormatFloat(data.H, 'G', 5, 64), "E", "*10^")
		delta := strings.ReplaceAll(strconv.FormatFloat(data.Delta, 'G', 5, 64), "E", "*10^")

		formulas.K = fmt.Sprintf("%s / %s", DOut_, Dcp_)
		formulas.X = fmt.Sprintf("0.67 * (%s^2 + (1 + 8.55 * lg(%s) - 1)) / ((%s - 1) * %s^2 - 1 + (1.857 * %s^2 + 1) * %s^3/%s^3)",
			k, k, k, k, k, H, delta)
		formulas.Y = fmt.Sprintf("%s / (%s * %s)", x, delta, eAt20)
	} else {
		radius := strings.ReplaceAll(strconv.FormatFloat(data.Radius, 'G', 5, 64), "E", "*10^")
		lambda := strings.ReplaceAll(strconv.FormatFloat(cap.Lambda, 'G', 5, 64), "E", "*10^")
		omega := strings.ReplaceAll(strconv.FormatFloat(cap.Omega, 'G', 5, 64), "E", "*10^")

		formulas.Lambda = fmt.Sprintf("(%s / %s) * Sqrt(%s / %s)", h_, D_, radius, S0_)
		formulas.Omega = fmt.Sprintf("1 / (1 + 1.285*%s + 1.63*%s * (%s/%s)^2 * lg(%s/%s)", lambda, lambda, h_, S0_, DOut_, D_)
		formulas.Y = fmt.Sprintf("((1 - %s * (1 + 1.285*%s)) / (%s * %s^3)) * ((%s + %s) / (%s - %s))", omega, lambda, eAt20, h_, DOut_, D_, DOut_, D_)
	}

	return formulas
}

func (s *FormulasService) tightnessFormulas(data models.DataCap, req *calc_api.CapRequest, result *calc_api.CapResponse,
) *cap_model.TightnessFormulas {
	tightness := &cap_model.TightnessFormulas{}

	// перевод чисел в строки
	axialForce := req.Data.AxialForce
	pressure := strconv.FormatFloat(req.Data.Pressure, 'f', -1, 64)

	pres := strconv.FormatFloat(data.Gasket.Pres, 'f', -1, 64)
	m := strconv.FormatFloat(data.Gasket.M, 'f', -1, 64)

	sigmaAt20 := strconv.FormatFloat(data.Bolt.SigmaAt20, 'f', -1, 64)

	B0 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.B0, 'G', 5, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Dcp, 'G', 5, 64), "E", "*10^")
	Ab := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.A, 'G', 5, 64), "E", "*10^")
	Alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Alpha, 'G', 5, 64), "E", "*10^")

	Po := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Po, 'G', 5, 64), "E", "*10^")
	Qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Qd, 'G', 5, 64), "E", "*10^")
	Rp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Rp, 'G', 5, 64), "E", "*10^")
	Pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb1, 'G', 5, 64), "E", "*10^")
	Pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb2, 'G', 5, 64), "E", "*10^")
	Pb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb, 'G', 5, 64), "E", "*10^")

	// формула 6
	// Усилие необходимое для смятия прокладки при затяжке
	tightness.Po = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, Dcp, B0, pres)

	if req.Data.Pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях
		tightness.Rp = fmt.Sprintf("%f * %s * %s * %s * |%s|", math.Pi, Dcp, B0, m, pressure)
	}

	// формула 9
	// Равнодействующая нагрузка от давления
	tightness.Qd = fmt.Sprintf("0.785 * (%s)^2 * %s", Dcp, pressure)

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	// tightness.Qfm = fmt.Sprintf("%d", axialForce)

	// minB := 0.4 * aux.A * data.Bolt.SigmaAt20
	minB := fmt.Sprintf("0.4 * %s * %s", Ab, sigmaAt20)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	tightness.Pb2 = fmt.Sprintf("max(%s; %s)", Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	tightness.Pb1 = fmt.Sprintf("%s * (%s + %d) + %s", Alpha, Qd, axialForce, Rp)

	// Расчетная нагрузка на болты/шпильки фланцевых соединений
	// tightness.Pb = math.Max(tightness.Pb1, tightness.Pb2)
	tightness.Pb = fmt.Sprintf("max(%s; %s)", Pb1, Pb2)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений в рабочих условиях
	tightness.Pbr = fmt.Sprintf("%s + (1 - %s) * (%s + %d)", Pb, Alpha, Qd, axialForce)

	return tightness
}

func (s *FormulasService) staticResistanceFormulas(
	flange *cap_model.FlangeResult,
	flangeRes *cap_model.CalcAuxiliary_Flange,
	typeFlange cap_model.FlangeData_Type,
	staticRes *cap_model.CalcStaticResistance,
	data models.DataCap,
	req *calc_api.CapRequest,
	Pb_, Pbr_, Qd_, Qfm_ float64,
) *cap_model.StaticResistanceFormulas {
	static := &cap_model.StaticResistanceFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Data.Pressure, 'G', 5, 64)
	axialForce := req.Data.AxialForce

	fD6 := strconv.FormatFloat(flange.D6, 'G', 5, 64)
	fH := strconv.FormatFloat(flange.H, 'G', 5, 64)
	fD := strconv.FormatFloat(flange.D, 'G', 5, 64)
	fS0 := strconv.FormatFloat(flange.S0, 'G', 5, 64)
	fS1 := strconv.FormatFloat(flange.S1, 'G', 5, 64)
	fC := strconv.FormatFloat(flange.C, 'G', 5, 64)

	gM := strconv.FormatFloat(data.Gasket.M, 'G', 5, 64)

	count := data.Bolt.Count
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 5, 64)

	Pb := strings.ReplaceAll(strconv.FormatFloat(Pb_, 'G', 5, 64), "E", "*10^")
	Pbr := strings.ReplaceAll(strconv.FormatFloat(Pbr_, 'G', 5, 64), "E", "*10^")
	Qd := strings.ReplaceAll(strconv.FormatFloat(Qd_, 'G', 5, 64), "E", "*10^")
	Qfm := strings.ReplaceAll(strconv.FormatFloat(Qfm_, 'G', 5, 64), "E", "*10^")

	B := strings.ReplaceAll(strconv.FormatFloat(flangeRes.B, 'G', 5, 64), "E", "*10^")
	E := strings.ReplaceAll(strconv.FormatFloat(flangeRes.E, 'G', 5, 64), "E", "*10^")
	A := strings.ReplaceAll(strconv.FormatFloat(flangeRes.A, 'G', 5, 64), "E", "*10^")
	F := strings.ReplaceAll(strconv.FormatFloat(flangeRes.F, 'G', 5, 64), "E", "*10^")
	Lymda := strings.ReplaceAll(strconv.FormatFloat(flangeRes.Lymda, 'G', 5, 64), "E", "*10^")
	L0 := strings.ReplaceAll(strconv.FormatFloat(flangeRes.L0, 'G', 5, 64), "E", "*10^")
	BetaF := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaF, 'G', 5, 64), "E", "*10^")
	BetaY := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaY, 'G', 5, 64), "E", "*10^")
	BetaZ := strings.ReplaceAll(strconv.FormatFloat(flangeRes.BetaZ, 'G', 5, 64), "E", "*10^")

	Cf := strings.ReplaceAll(strconv.FormatFloat(staticRes.Cf, 'G', 5, 64), "E", "*10^")
	MM := strings.ReplaceAll(strconv.FormatFloat(staticRes.MM, 'G', 5, 64), "E", "*10^")
	Dzv := strings.ReplaceAll(strconv.FormatFloat(staticRes.Dzv, 'G', 5, 64), "E", "*10^")
	Mp := strings.ReplaceAll(strconv.FormatFloat(staticRes.Mp, 'G', 5, 64), "E", "*10^")
	SigmaM1 := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaM1, 'G', 5, 64), "E", "*10^")
	SigmaR := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaR, 'G', 5, 64), "E", "*10^")
	SigmaRp := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaRp, 'G', 5, 64), "E", "*10^")
	SigmaP1 := strings.ReplaceAll(strconv.FormatFloat(staticRes.SigmaP1, 'G', 5, 64), "E", "*10^")

	temp1 := fmt.Sprintf("%f * %s / %d", math.Pi, fD6, count)
	temp2 := fmt.Sprintf("2 * %s + 6 * %s / (%s + 0.5)", diameter, fH, gM)

	// Коэффициент учитывающий изгиб тарелки фланца между болтами шпильками
	static.Cf = fmt.Sprintf("max(1; sqrt((%s) / (%s)))", temp1, temp2)

	// Приведенный диаметр приварного встык фланца с конической или прямой втулкой
	if typeFlange == cap_model.FlangeData_welded && flange.D <= 20*flange.S1 {
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

	if typeFlange == cap_model.FlangeData_free {
		static.MMk = fmt.Sprintf("%s * %s * %s", Cf, Pb, A)
		static.Mpk = fmt.Sprintf("%s * %s * %s", Cf, Pbr, A)
	}

	// Меридиональное изгибное напряжение во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке бурта свободного фланца
	if typeFlange == cap_model.FlangeData_welded && flange.S1 != flange.S0 {
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
	static.SigmaT = fmt.Sprintf("(%s * %s) / ((%s)^2 * %s) - %s * %s", BetaY, MM, fH, fD, BetaZ, SigmaR)

	if typeFlange == cap_model.FlangeData_free {
		Hk := strconv.FormatFloat(flange.Ring.Hk, 'G', 5, 64)
		Dk := strconv.FormatFloat(flange.Ring.Dk, 'G', 5, 64)
		MMk := strings.ReplaceAll(strconv.FormatFloat(staticRes.MMk, 'G', 5, 64), "E", "*10^")

		static.SigmaK = fmt.Sprintf("%s * %s / ((%s)^2 * %s)", BetaY, MMk, Hk, Dk)
	}

	// Меридиональные изгибные напряжения во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке
	// трубе бурта свободного фланца в рабочих условиях
	if typeFlange == cap_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		static.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", Mp, Lymda, fS1, fC, Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		static.SigmaP0 = fmt.Sprintf("%s * %s", F, SigmaP1)
	} else {
		static.SigmaP1 = fmt.Sprintf("%s / (%s * (%s - %s)^2 * %s)", Mp, Lymda, fS0, fC, Dzv)
		static.SigmaP0 = SigmaP1
	}

	if typeFlange == cap_model.FlangeData_welded {
		temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, fD, fS1, fS1, fC)
		// формула (ф. 37)
		static.SigmaMp = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d + (4 * |%d|) / (%s + %s)) / %s", fD, pressure, axialForce, 0, fD, fS1, temp)
	}

	temp := fmt.Sprintf("%f * (%s + %s) * (%s - %s)", math.Pi, fD, fS0, fS0, fC)
	// Меридиональные мембранные напряжения во втулке приварного встык фланца обечайке трубе
	// плоского фланца или обечайке трубе бурта свободного фланца в рабочих условиях
	// формула (ф. 37)
	// - для приварных встык фланцев с конической втулкой в сечении S1
	static.SigmaMp0 = fmt.Sprintf("(0.785 * (%s)^2 * %s + %d + (4 * |%d|) / (%s + %s)) / %s", fD, pressure, axialForce, 0, fD, fS0, temp)

	// Окружные мембранные напряжения от действия давления во втулке приварного встык фланца обечайке
	// трубе плоского фланца или обечайке трубе бурта свободного фланца в сечении S0
	static.SigmaMop = fmt.Sprintf("%s * %s / (2 * (%s - %s))", pressure, fD, fS0, fC)

	// Напряжения в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в рабочих условиях
	// - радиальные напряжения
	static.SigmaRp = fmt.Sprintf("((1.33 * %s * %s + %s) / (%s * (%s)^2 * %s * %s)) * %s", BetaF, fH, L0, Lymda, fH, L0, fD, Mp)
	// - окружное напряжения
	static.SigmaTp = fmt.Sprintf("%s * %s / ((%s)^2 * %s) - %s * %s", BetaY, Mp, fH, fD, BetaZ, SigmaRp)

	if typeFlange == cap_model.FlangeData_free {
		Hk := strconv.FormatFloat(flange.Ring.Hk, 'G', 5, 64)
		Dk := strconv.FormatFloat(flange.Ring.Dk, 'G', 5, 64)

		static.SigmaKp = fmt.Sprintf("%s * %s / ((%s)^2 * %s)", BetaY, Mp, Hk, Dk)
	}

	return static
}

func (s *FormulasService) conditionsForStrengthFormulas(
	flangeType cap_model.FlangeData_Type,
	flange *cap_model.FlangeResult,
	calcFlange *cap_model.CalcAuxiliary_Flange,
	static *cap_model.CalcStaticResistance,
	cond *cap_model.CalcConditionsForStrength,
	isWork, isTemp bool,
) *cap_model.ConditionsForStrengthFormulas {
	conditions := &cap_model.ConditionsForStrengthFormulas{}

	// перевод чисел в строки
	fEpsilon := strings.ReplaceAll(strconv.FormatFloat(flange.Epsilon, 'G', 5, 64), "E", "*10^")
	fEpsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(flange.EpsilonAt20, 'G', 5, 64), "E", "*10^")
	fSigmaAt20 := strconv.FormatFloat(flange.SigmaAt20, 'G', 5, 64)
	fSigma := strconv.FormatFloat(flange.Sigma, 'G', 5, 64)
	fSigmaMAt20 := strconv.FormatFloat(flange.SigmaMAt20, 'G', 5, 64)
	fSigmaM := strconv.FormatFloat(flange.SigmaM, 'G', 5, 64)
	fSigmaRAt20 := strconv.FormatFloat(flange.SigmaRAt20, 'G', 5, 64)
	fSigmaR := strconv.FormatFloat(flange.SigmaR, 'G', 5, 64)

	Yf := strings.ReplaceAll(strconv.FormatFloat(calcFlange.Yf, 'G', 5, 64), "E", "*10^")

	Mp := strings.ReplaceAll(strconv.FormatFloat(static.Mp, 'G', 5, 64), "E", "*10^")
	SigmaM1 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaM1, 'G', 5, 64), "E", "*10^")
	SigmaM0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaM0, 'G', 5, 64), "E", "*10^")
	SigmaR := strings.ReplaceAll(strconv.FormatFloat(static.SigmaR, 'G', 5, 64), "E", "*10^")
	SigmaT := strings.ReplaceAll(strconv.FormatFloat(static.SigmaT, 'G', 5, 64), "E", "*10^")
	SigmaP1 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaP1, 'G', 5, 64), "E", "*10^")
	SigmaP0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaP0, 'G', 5, 64), "E", "*10^")
	SigmaMp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMp, 'G', 5, 64), "E", "*10^")
	SigmaMp0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMp0, 'G', 5, 64), "E", "*10^")
	SigmaMpm := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMpm, 'G', 5, 64), "E", "*10^")
	SigmaMpm0 := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMpm0, 'G', 5, 64), "E", "*10^")
	SigmaRp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaRp, 'G', 5, 64), "E", "*10^")
	SigmaTp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaTp, 'G', 5, 64), "E", "*10^")
	SigmaMop := strings.ReplaceAll(strconv.FormatFloat(static.SigmaMop, 'G', 5, 64), "E", "*10^")

	chooseTeta := map[bool]float64{
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

	kt := strconv.FormatFloat(Kt[isTemp], 'G', 5, 64)
	ks := strconv.FormatFloat(Ks, 'G', 5, 64)
	teta := strconv.FormatFloat(chooseTeta[isWork], 'G', 5, 64)

	var DTeta float64
	if flangeType == cap_model.FlangeData_welded {
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
	tmp := strconv.FormatFloat(DTeta, 'G', 5, 64)
	DTeta_ := fmt.Sprintf("%s * %s", teta, tmp)

	conditions.Teta = fmt.Sprintf("%s * %s * %s / %s", Mp, Yf, fEpsilonAt20, fEpsilon)
	conditions.CondTeta = &cap_model.ConditionFormulas{
		X: conditions.Teta,
		Y: DTeta_,
	}

	if flangeType == cap_model.FlangeData_free {
		fEpsilonK := strconv.FormatFloat(flange.Ring.Epsilon, 'G', 5, 64)
		fEpsilonKAt20 := strconv.FormatFloat(flange.Ring.EpsilonAt20, 'G', 5, 64)
		Yk := strings.ReplaceAll(strconv.FormatFloat(calcFlange.Yk, 'G', 5, 64), "E", "*10^")
		Mpk := strings.ReplaceAll(strconv.FormatFloat(static.Mpk, 'G', 5, 64), "E", "*10^")

		DTetaK := fmt.Sprintf("%s * 0.02", teta)
		conditions.TetaK = fmt.Sprintf("%s * %s * %s / %s", Mpk, Yk, fEpsilonKAt20, fEpsilonK)
		conditions.CondTetaK = &cap_model.ConditionFormulas{
			X: conditions.TetaK,
			Y: DTetaK,
		}
	}

	//* Условия статической прочности фланцев
	if flangeType == cap_model.FlangeData_welded && flange.S1 != flange.S0 {
		Max1 := fmt.Sprintf("max(|%s + %s|; |%s + %s|)", SigmaM1, SigmaR, SigmaM1, SigmaT)
		max1Y := fmt.Sprintf("%s * %s * %s", ks, kt, fSigmaMAt20)

		t1 := fmt.Sprintf("|%s - %s + %s|; |%s - %s + %s|", SigmaP1, SigmaMp, SigmaRp, SigmaP1, SigmaMpm, SigmaRp)
		t2 := fmt.Sprintf("|%s - %s + %s|; |%s - %s + %s|", SigmaP1, SigmaMp, SigmaTp, SigmaP1, SigmaMpm, SigmaTp)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("|%s + %s|; |%s + %s|", SigmaP1, SigmaMp, SigmaP1, SigmaMpm)

		Max2 := fmt.Sprintf("max(%s; %s)", t1, t2)
		max2Y := fmt.Sprintf("%s * %s * %s", ks, kt, fSigmaM)

		max3Y := fmt.Sprintf("1.3 * %s", fSigmaRAt20)

		t1 = fmt.Sprintf("|%s + %s|; |%s - %s|", SigmaP0, SigmaMp0, SigmaP0, SigmaMp0)
		t2 = fmt.Sprintf("|%s + %s|; |%s - %s|", SigmaP0, SigmaMpm0, SigmaP0, SigmaMpm0)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("|0.3 * %s + %s|; |0.3 * %s - %s|", SigmaP0, SigmaMop, SigmaP0, SigmaMop)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("|0.7 * %s + (%s - %s)|; |0.7 * %s - (%s - %s)|", SigmaP0, SigmaMp0, SigmaMop, SigmaP0, SigmaMp0, SigmaMop)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("|0.7 * %s + (%s - %s)|; |0.7 * %s - (%s - %s)|", SigmaP0, SigmaMpm0, SigmaMop, SigmaP0, SigmaMpm0, SigmaMop)

		Max4 := fmt.Sprintf("max(%s; %s)", t1, t2)
		max4Y := fmt.Sprintf("1.3 * %s", fSigmaR)

		conditions.Max1 = &cap_model.ConditionFormulas{
			X: Max1,
			Y: max1Y,
		}
		conditions.Max2 = &cap_model.ConditionFormulas{
			X: Max2,
			Y: max2Y,
		}
		conditions.Max3 = &cap_model.ConditionFormulas{
			X: SigmaM0,
			Y: max3Y,
		}
		conditions.Max4 = &cap_model.ConditionFormulas{
			X: Max4,
			Y: max4Y,
		}
	} else {
		Max5 := fmt.Sprintf("max(|%s + %s|; |%s + %s|)", SigmaM0, SigmaR, SigmaM0, SigmaT)

		t1 := fmt.Sprintf("|%s - %s + %s); |%s - %s + %s|", SigmaP0, SigmaMp0, SigmaTp, SigmaP0, SigmaMpm0, SigmaTp)
		t2 := fmt.Sprintf("|%s - %s + %s|; |%s - %s + %s|", SigmaP0, SigmaMp0, SigmaRp, SigmaP0, SigmaMpm0, SigmaRp)
		t1 = fmt.Sprintf("%s; %s", t1, t2)
		t2 = fmt.Sprintf("|%s + %s|; |%s + %s|", SigmaP0, SigmaMp0, SigmaP0, SigmaMpm0)

		Max6 := fmt.Sprintf("max(%s; %s)", t1, t2)

		conditions.Max5 = &cap_model.ConditionFormulas{
			X: Max5,
			Y: fSigmaAt20,
		}
		conditions.Max6 = &cap_model.ConditionFormulas{
			X: Max6,
			Y: fSigma,
		}
	}

	max7 := fmt.Sprintf("|%s|; |%s|", SigmaMp0, SigmaMpm0)
	Max7 := fmt.Sprintf("max(%s; |%s|)", max7, SigmaMop)
	Max8 := fmt.Sprintf("max(|%s|; |%s|)", SigmaR, SigmaT)
	max8Y := fmt.Sprintf("%s * %s", kt, fSigmaAt20)
	Max9 := fmt.Sprintf("max(|%s|; |%s|)", SigmaRp, SigmaTp)
	max9Y := fmt.Sprintf("%s * %s", kt, fSigma)

	conditions.Max7 = &cap_model.ConditionFormulas{X: Max7, Y: fSigma}
	conditions.Max8 = &cap_model.ConditionFormulas{X: Max8, Y: max8Y}
	conditions.Max9 = &cap_model.ConditionFormulas{X: Max9, Y: max9Y}

	if flangeType == cap_model.FlangeData_free {
		SigmaK := strings.ReplaceAll(strconv.FormatFloat(static.SigmaK, 'G', 5, 64), "E", "*10^")
		SigmaKp := strings.ReplaceAll(strconv.FormatFloat(static.SigmaKp, 'G', 5, 64), "E", "*10^")
		fSigmaKAt20 := strconv.FormatFloat(flange.Ring.SigmaAt20, 'G', 5, 64)
		fSigmaK := strconv.FormatFloat(flange.Ring.Sigma, 'G', 5, 64)

		max10Y := fmt.Sprintf("%s * %s", kt, fSigmaKAt20)
		max11Y := fmt.Sprintf("%s * %s", kt, fSigmaK)

		conditions.Max10 = &cap_model.ConditionFormulas{X: SigmaK, Y: max10Y}
		conditions.Max11 = &cap_model.ConditionFormulas{X: SigmaKp, Y: max11Y}
	}

	return conditions
}

func (s *FormulasService) tightnessLoadFormulas(
	data models.DataCap,
	req *calc_api.CapRequest,
	result *calc_api.CapResponse,
) *cap_model.TightnessLoadFormulas {
	tightness := &cap_model.TightnessLoadFormulas{}

	// перевод чисел в строки
	fAlpha1 := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Alpha, 'G', 5, 64), "E", "*10^")
	fH1 := strconv.FormatFloat(data.Flange.H, 'G', 5, 64)
	fT1 := strconv.FormatFloat(data.Flange.T, 'G', 5, 64)
	fAlpha2 := strings.ReplaceAll(strconv.FormatFloat(data.Cap.Alpha, 'G', 5, 64), "E", "*10^")
	fH2 := strconv.FormatFloat(data.Cap.H, 'G', 5, 64)
	fT2 := strconv.FormatFloat(data.Cap.T, 'G', 5, 64)

	bAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Alpha, 'G', 5, 64), "E", "*10^")
	bTemp := strconv.FormatFloat(data.Bolt.Temp, 'G', 5, 64)

	Gamma := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Gamma, 'G', 5, 64), "E", "*10^")
	Alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Alpha, 'G', 5, 64), "E", "*10^")

	Rp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Rp, 'G', 5, 64), "E", "*10^")
	Pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb1, 'G', 5, 64), "E", "*10^")
	Pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb2, 'G', 5, 64), "E", "*10^")
	Pb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Pb, 'G', 5, 64), "E", "*10^")
	Qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Tightness.Qd, 'G', 5, 64), "E", "*10^")

	Qt := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.TightnessLoad.Qt, 'G', 5, 64), "E", "*10^")

	var temp1, temp2 string
	if req.IsUseWasher {
		wAlpha1 := strings.ReplaceAll(strconv.FormatFloat(data.Washer1.Alpha, 'G', 5, 64), "E", "*10^")
		wThick1 := strconv.FormatFloat(data.Washer1.Thickness, 'G', 5, 64)
		wAlpha2 := strings.ReplaceAll(strconv.FormatFloat(data.Washer2.Alpha, 'G', 5, 64), "E", "*10^")
		wThick2 := strconv.FormatFloat(data.Washer2.Thickness, 'G', 5, 64)

		temp1 = fmt.Sprintf("(%s * %s + %s * %s) * (%s - 20) + (%s * %s + %s * %s)*(%s - 20)",
			fAlpha1, fH1, wAlpha1, wThick1, fT1, fAlpha2, fH2, wAlpha2, wThick2, fT2)
	} else {
		temp1 = fmt.Sprintf("%s * %s * (%s - 20) + %s * %s * (%s - 20)", fAlpha1, fH1, fT1, fAlpha2, fH2, fT2)
	}
	temp2 = fmt.Sprintf("%s + %s", fH1, fH2)

	if data.FlangeType == cap_model.FlangeData_free {
		fAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Flange.Ring.Alpha, 'G', 5, 64), "E", "*10^")
		fH := strconv.FormatFloat(data.Flange.Ring.Hk, 'G', 5, 64)
		fT := strconv.FormatFloat(data.Flange.Ring.T, 'G', 5, 64)

		temp1 += fmt.Sprintf(" + (%s * %s) * (%s - 20)", fAlpha, fH, fT)
		temp2 += " + " + fH
	}

	if req.Data.IsEmbedded {
		eAlpha := strings.ReplaceAll(strconv.FormatFloat(data.Embed.Alpha, 'G', 5, 64), "E", "*10^")
		eThick := strconv.FormatFloat(data.Embed.Thickness, 'G', 5, 64)
		temp := strconv.FormatFloat(req.Data.Temp, 'G', 5, 64)

		temp1 += fmt.Sprintf(" + (%s * %s) * (%s - 20)", eAlpha, eThick, temp)
		temp2 += " + " + eThick
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	tightness.Qt = fmt.Sprintf("%s * (%s - %s * (%s) * (%s - 20))", Gamma, temp1, bAlpha, temp2, bTemp)

	pb1 := fmt.Sprintf("%s * (%s + %d) + %s", Alpha, Qd, req.Data.AxialForce, Rp)

	tightness.Pb1 = fmt.Sprintf("max(%s; %s - %s)", pb1, pb1, Qt)
	tightness.Pb = fmt.Sprintf("max(%s; %s)", Pb1, Pb2)
	tightness.Pbr = fmt.Sprintf("%s + (1 - %s) * (%s + %d) + %s", Pb, Alpha, Qd, req.Data.AxialForce, Qt)

	return tightness
}
