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

	formulas := &flange_model.Formulas_Strength{
		Auxiliary:     s.auxiliaryFormulas(d, req, result),
		Tightness:     s.tightnessFormulas(d, req, result),
		BoltStrength1: bolt1,
		Moment1:       moment1,
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
	epsilon := strconv.FormatFloat(data.Gasket.Epsilon, 'G', 3, 64)
	compression := strconv.FormatFloat(data.Gasket.Compression, 'G', 3, 64)

	lenght := strconv.FormatFloat(data.Bolt.Lenght, 'G', 3, 64)
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 3, 64)
	area := strconv.FormatFloat(data.Bolt.Area, 'G', 3, 64)
	epsilonAt20 := strconv.FormatFloat(data.Bolt.EpsilonAt20, 'G', 3, 64)
	bEpsilon := strconv.FormatFloat(data.Bolt.Epsilon, 'G', 3, 64)
	count := data.Bolt.Count

	db1 := strconv.FormatFloat(data.Flange1.D6, 'G', 3, 64)
	fEpAt201 := strconv.FormatFloat(data.Flange1.EpsilonAt20, 'G', 3, 64)
	fEpAt202 := strconv.FormatFloat(data.Flange2.EpsilonAt20, 'G', 3, 64)
	fEp1 := strconv.FormatFloat(data.Flange1.Epsilon, 'G', 3, 64)
	fEp2 := strconv.FormatFloat(data.Flange2.Epsilon, 'G', 3, 64)
	fEpsilonKAt201 := strconv.FormatFloat(data.Flange1.Ring.EpsilonKAt20, 'G', 3, 64)
	fEpsilonKAt202 := strconv.FormatFloat(data.Flange2.Ring.EpsilonKAt20, 'G', 3, 64)
	fEpsilonK1 := strconv.FormatFloat(data.Flange1.Ring.EpsilonK, 'G', 3, 64)
	fEpsilonK2 := strconv.FormatFloat(data.Flange2.Ring.EpsilonK, 'G', 3, 64)

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
	Yf2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yf, 'G', 3, 64), "E", "*10^")
	Yfc2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yfc, 'G', 3, 64), "E", "*10^")
	Yfn2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yfn, 'G', 3, 64), "E", "*10^")
	A2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.A, 'G', 3, 64), "E", "*10^")
	E2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.E, 'G', 3, 64), "E", "*10^")
	B2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.B, 'G', 3, 64), "E", "*10^")
	Yk2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Strength.Auxiliary.Flange2.Yk, 'G', 3, 64), "E", "*10^")

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
		divider += fmt.Sprintf(" + (%s * %s / %s) * (%s)^2", Yk1, fEpsilonKAt201, fEpsilonK1, A1)
	}
	if data.Type2 == flange_model.FlangeData_free {
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
	epsilonAt20 := strconv.FormatFloat(data.EpsilonAt20, 'G', 3, 64)
	epsilonKAt20 := strconv.FormatFloat(data.Ring.EpsilonKAt20, 'G', 3, 64)

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
	tightness.Pbr = fmt.Sprintf("%s + (1 - %s) * (%s + %d) + 4 * (1 - %s) * |%d| / %s", Pb, Alpha, Qd, axialForce, AlphaM, bendingMoment, Dcp)

	return tightness
}

func (s *FormulasService) staticResistanceCalculate(
	flange *flange_model.FlangeResult,
	calcFlange *flange_model.CalcAuxiliary_Flange,
	typeFlange flange_model.FlangeData_Type,
	data models.DataFlange,
	req *calc_api.FlangeRequest,
	Pb, Pbr, Qd, Qfm float64,
) *flange_model.CalcStaticResistance {
	static := &flange_model.CalcStaticResistance{}

	temp1 := math.Pi * flange.D6 / float64(data.Bolt.Count)
	temp2 := 2*float64(data.Bolt.Diameter) + 6*flange.H/(data.Gasket.M+0.5)

	// Коэффициент учитывающий изгиб тарелки фланца между болтами шпильками
	static.Cf = math.Max(1, math.Sqrt(temp1/temp2))

	// Приведенный диаметр приварного встык фланца с конической или прямой втулкой
	var Dzv float64
	if typeFlange == flange_model.FlangeData_welded && flange.D <= 20*flange.S1 {
		if calcFlange.F > 1 {
			Dzv = flange.D + flange.S0
		} else {
			Dzv = flange.D + flange.S1
		}
	} else {
		Dzv = flange.D
	}
	static.Dzv = Dzv

	// Расчетный изгибающий момент действующий на фланец при затяжке
	static.MM = static.Cf * Pb * calcFlange.B
	// Расчетный изгибающий момент действующий на фланец в рабочих условиях
	static.Mp = static.Cf * math.Max(Pbr*calcFlange.B+(Qd+Qfm)*calcFlange.E, math.Abs(Qd+Qfm)*calcFlange.E)

	if typeFlange == flange_model.FlangeData_free {
		static.MMk = static.Cf * Pb * calcFlange.A
		static.Mpk = static.Cf * Pbr * calcFlange.A
	}

	// Меридиональное изгибное напряжение во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке бурта свободного фланца
	if typeFlange == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		static.SigmaM1 = static.MM / (calcFlange.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		static.SigmaM0 = calcFlange.F * static.SigmaM1
	} else {
		static.SigmaM1 = static.MM / (calcFlange.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv)
		static.SigmaM0 = static.SigmaM1
	}

	// Радиальное напряжение в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в условиях затяжки
	static.SigmaR = ((1.33*calcFlange.BetaF*flange.H + calcFlange.L0) / (calcFlange.Lymda * math.Pow(flange.H, 2) * calcFlange.L0 * flange.D)) * static.MM
	// Окружное напряжение в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в условиях затяжки
	static.SigmaT = calcFlange.BetaY*static.MM/(math.Pow(flange.H, 2)*flange.D) - calcFlange.BetaZ*static.SigmaR

	if typeFlange == flange_model.FlangeData_free {
		static.SigmaK = calcFlange.BetaY * static.MMk / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	// Меридиональные изгибные напряжения во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке
	// трубе бурта свободного фланца в рабочих условиях
	if typeFlange == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		static.SigmaP1 = static.Mp / (calcFlange.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		static.SigmaP0 = calcFlange.F * static.SigmaP1
	} else {
		static.IsEqualSigma = true
		static.SigmaP1 = static.Mp / (calcFlange.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv)
		static.SigmaP0 = static.SigmaP1
	}

	if typeFlange == flange_model.FlangeData_welded {
		temp := math.Pi * (flange.D + flange.S1) * (flange.S1 - flange.C)
		// формула (ф. 37)
		static.SigmaMp = (0.785*math.Pow(flange.D, 2)*req.Pressure + float64(req.AxialForce) +
			4*math.Abs(float64(req.BendingMoment)/(flange.D+flange.S1))) / temp
		static.SigmaMpm = (0.785*math.Pow(flange.D, 2)*req.Pressure + float64(req.AxialForce) -
			4*math.Abs(float64(req.BendingMoment)/(flange.D+flange.S1))) / temp
	}

	temp := math.Pi * (flange.D + flange.S0) * (flange.S0 - flange.C)
	// Меридиональные мембранные напряжения во втулке приварного встык фланца обечайке трубе
	// плоского фланца или обечайке трубе бурта свободного фланца в рабочих условиях
	// формула (ф. 37)
	// - для приварных встык фланцев с конической втулкой в сечении S1
	static.SigmaMp0 = (0.785*math.Pow(flange.D, 2)*req.Pressure + float64(req.AxialForce) +
		4*math.Abs(float64(req.BendingMoment)/(flange.D+flange.S0))) / temp
	// - для приварных встык фланцев с конической втулкой в сечении S0 приварных фланцев с прямой втулкой плоских фланцев и свободных фланцев
	static.SigmaMpm0 = (0.785*math.Pow(flange.D, 2)*req.Pressure + float64(req.AxialForce) -
		4*math.Abs(float64(req.BendingMoment)/(flange.D+flange.S0))) / temp

	// Окружные мембранные напряжения от действия давления во втулке приварного встык фланца обечайке
	// трубе плоского фланца или обечайке трубе бурта свободного фланца в сечениии S0
	static.SigmaMop = req.Pressure * flange.D / (2.0 * (flange.S0 - flange.C))

	// Напряжения в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в рабочих условиях
	// - радиальные напряжения
	static.SigmaRp = ((1.33*calcFlange.BetaF*flange.H + calcFlange.L0) / (calcFlange.Lymda * math.Pow(flange.H, 2) * calcFlange.L0 * flange.D)) * static.Mp
	// - окружное напряжения
	static.SigmaTp = calcFlange.BetaY*static.Mp/(math.Pow(flange.H, 2)*flange.D) - calcFlange.BetaZ*static.SigmaRp

	if typeFlange == flange_model.FlangeData_free {
		static.SigmaKp = calcFlange.BetaY * static.Mp / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	return static
}

func (s *FormulasService) conditionsForStrengthCalculate(
	flangeType flange_model.FlangeData_Type,
	flange *flange_model.FlangeResult,
	calcFlange *flange_model.CalcAuxiliary_Flange,
	static *flange_model.CalcStaticResistance,
	isWork, isTemp bool,
) *flange_model.CalcConditionsForStrength {
	conditions := &flange_model.CalcConditionsForStrength{}

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

	var DTeta, DTetaK float64
	if flangeType == flange_model.FlangeData_welded {
		if flange.D <= constants.MinD {
			DTeta = constants.MinDTetta
		} else if flange.D > constants.MaxD {
			DTeta = constants.MaxDTetta
		} else {
			DTeta = ((flange.D-constants.MinD)/(constants.MaxD-constants.MinD))*
				(constants.MaxDTetta-constants.MinDTetta) + constants.MinDTetta
		}
	} else {
		DTeta = constants.MaxDTetta
	}
	DTeta = teta[isWork] * DTeta

	conditions.Teta = static.Mp * calcFlange.Yf * flange.EpsilonAt20 / flange.Epsilon
	conditions.CondTeta = &flange_model.Condition{
		X: conditions.Teta,
		Y: DTeta,
	}

	if flangeType == flange_model.FlangeData_free {
		//strength.DTetaK = 0.002
		DTetaK = 0.02
		DTetaK = teta[isWork] * DTetaK
		conditions.TetaK = static.Mpk * calcFlange.Yk * flange.Ring.EpsilonKAt20 / flange.Ring.EpsilonK
		conditions.CondTetaK = &flange_model.Condition{
			X: conditions.TetaK,
			Y: DTetaK,
		}
	}

	//* Условия статической прочности фланцев
	if flangeType == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		Max1 := math.Max(math.Abs(static.SigmaM1+static.SigmaR), math.Abs(static.SigmaM1+static.SigmaT))

		t1 := math.Max(math.Abs(static.SigmaP1-static.SigmaMp+static.SigmaRp), math.Abs(static.SigmaP1-static.SigmaMpm+static.SigmaRp))
		t2 := math.Max(math.Abs(static.SigmaP1-static.SigmaMp+static.SigmaTp), math.Abs(static.SigmaP1-static.SigmaMpm+static.SigmaTp))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(static.SigmaP1+static.SigmaMp), math.Abs(static.SigmaP1+static.SigmaMpm))

		Max2 := math.Max(t1, t2)
		Max3 := static.SigmaM0

		t1 = math.Max(math.Abs(static.SigmaP0+static.SigmaMp0), math.Abs(static.SigmaP0-static.SigmaMp0))
		t2 = math.Max(math.Abs(static.SigmaP0+static.SigmaMpm0), math.Abs(static.SigmaP0-static.SigmaMpm0))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.3*static.SigmaP0+static.SigmaMop), math.Abs(0.3*static.SigmaP0-static.SigmaMop))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.7*static.SigmaP0+(static.SigmaMp0-static.SigmaMop)), math.Abs(0.7*static.SigmaP0-(static.SigmaMp0-static.SigmaMop)))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.7*static.SigmaP0+(static.SigmaMpm0-static.SigmaMop)), math.Abs(0.7*static.SigmaP0-(static.SigmaMpm0-static.SigmaMop)))

		Max4 := math.Max(t1, t2)

		conditions.Max1 = &flange_model.Condition{
			X: Max1,
			Y: Ks * Kt[isTemp] * flange.SigmaMAt20,
		}
		conditions.Max2 = &flange_model.Condition{
			X: Max2,
			Y: Ks * Kt[isTemp] * flange.SigmaM,
		}
		conditions.Max3 = &flange_model.Condition{
			X: Max3,
			Y: 1.3 * flange.SigmaRAt20,
		}
		conditions.Max4 = &flange_model.Condition{
			X: Max4,
			Y: 1.3 * flange.SigmaR,
		}
	} else {
		Max5 := math.Max(math.Abs(static.SigmaM0+static.SigmaR), math.Abs(static.SigmaM0+static.SigmaT))

		t1 := math.Max(math.Abs(static.SigmaP0-static.SigmaMp0+static.SigmaTp), math.Abs(static.SigmaP0-static.SigmaMpm0+static.SigmaTp))
		t2 := math.Max(math.Abs(static.SigmaP0-static.SigmaMp0+static.SigmaRp), math.Abs(static.SigmaP0-static.SigmaMpm0+static.SigmaRp))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(static.SigmaP0+static.SigmaMp0), math.Abs(static.SigmaP0+static.SigmaMpm0))

		Max6 := math.Max(t1, t2)

		conditions.Max5 = &flange_model.Condition{
			X: Max5,
			Y: flange.SigmaAt20,
		}
		conditions.Max6 = &flange_model.Condition{
			X: Max6,
			Y: flange.Sigma,
		}
	}

	max7 := math.Max(math.Abs(static.SigmaMp0), math.Abs(static.SigmaMpm0))
	Max7 := math.Max(max7, math.Abs(static.SigmaMop))
	Max8 := math.Max(math.Abs(static.SigmaR), math.Abs(static.SigmaT))
	Max9 := math.Max(math.Abs(static.SigmaRp), math.Abs(static.SigmaTp))

	conditions.Max7 = &flange_model.Condition{X: Max7, Y: flange.Sigma}
	conditions.Max8 = &flange_model.Condition{X: Max8, Y: Kt[isTemp] * flange.SigmaAt20}
	conditions.Max9 = &flange_model.Condition{X: Max9, Y: Kt[isTemp] * flange.Sigma}

	if flangeType == flange_model.FlangeData_free {
		Max10 := static.SigmaK
		Max11 := static.SigmaKp

		conditions.Max10 = &flange_model.Condition{X: Max10, Y: Kt[isTemp] * flange.Ring.SigmaKAt20}
		conditions.Max11 = &flange_model.Condition{X: Max11, Y: Kt[isTemp] * flange.Ring.SigmaK}
	}

	return conditions
}

func (s *FormulasService) tightnessLoadCalculate(
	aux *flange_model.CalcAuxiliary,
	tig *flange_model.CalcTightness,
	data models.DataFlange,
	req *calc_api.FlangeRequest,
) *flange_model.CalcTightnessLoad {
	tightness := &flange_model.CalcTightnessLoad{}

	flange1 := aux.Flange1
	flange2 := aux.Flange2

	divider := aux.Yp + aux.Yb*data.Bolt.EpsilonAt20/data.Bolt.Epsilon + (flange1.Yf*data.Flange1.EpsilonAt20/data.Flange1.Epsilon)*math.Pow(flange1.B, 2) +
		+(flange2.Yf*data.Flange2.EpsilonAt20/data.Flange2.Epsilon)*math.Pow(flange2.B, 2)

	if data.Type1 == flange_model.FlangeData_free {
		divider += (flange1.Yk * data.Flange1.Ring.EpsilonKAt20 / data.Flange1.Ring.EpsilonK) * math.Pow(flange1.A, 2)
	}
	if data.Type2 == flange_model.FlangeData_free {
		divider += (flange2.Yk * data.Flange2.Ring.EpsilonKAt20 / data.Flange2.Ring.EpsilonK) * math.Pow(flange2.A, 2)
	}

	// формула (Е.8)
	gamma := 1 / divider

	var temp1, temp2 float64
	if req.IsUseWasher {
		temp1 = (data.Flange1.AlphaF*data.Flange1.H+data.Washer1.Alpha*data.Washer1.Thickness)*(data.Flange1.Tf-20) +
			+(data.Flange2.AlphaF*data.Flange2.H+data.Washer2.Alpha*data.Washer2.Thickness)*(data.Flange2.Tf-20)
	} else {
		temp1 = data.Flange1.AlphaF*data.Flange1.H*(data.Flange1.Tf-20) + data.Flange2.AlphaF*data.Flange2.H*(data.Flange2.Tf-20)
	}
	temp2 = data.Flange1.H + data.Flange2.H

	if data.Type1 == flange_model.FlangeData_free {
		temp1 += data.Flange1.Ring.AlphaK * data.Flange1.Hk * (data.Flange1.Ring.Tk - 20)
		temp2 += data.Flange1.Hk
	}
	if data.Type2 == flange_model.FlangeData_free {
		temp1 += data.Flange2.Ring.AlphaK * data.Flange2.Hk * (data.Flange2.Ring.Tk - 20)
		temp2 += data.Flange2.Hk
	}
	if req.IsEmbedded {
		temp1 += data.Embed.Alpha * data.Embed.Thickness * (req.Temp - 20)
		temp2 += data.Embed.Thickness
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	Qt := gamma * (temp1 - data.Bolt.Alpha*temp2*(data.Bolt.Temp-20))
	tightness.Qt = Qt

	tightness.Pb1 = math.Max(tig.Pb1, tig.Pb1-Qt)
	tightness.Pb = math.Max(tightness.Pb1, tig.Pb2)
	tightness.Pbr = tig.Pb + (1-aux.Alpha)*(tig.Qd+float64(req.AxialForce)) + Qt + 4*(1-aux.AlphaM)*math.Abs(float64(req.BendingMoment))/aux.Dcp

	return tightness
}
