package flange

import (
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

func (s *FlangeService) strengthCalculate(data models.DataFlange, req *calc_api.FlangeRequest) *flange_model.Calculated_Strength {
	auxiliary := s.auxiliaryCalculate(data, req)
	tightness := s.tightnessCalculate(auxiliary, data, req)
	bolt1 := s.boltStrengthCalculate(data, req, tightness.Pb, tightness.Pbr, auxiliary.A, auxiliary.Dcp, false)
	moment1 := s.momentCalculate(req.Friction, data, bolt1.SigmaB1, bolt1.DSigmaM, tightness.Pb, auxiliary.A, auxiliary.Dcp, false)

	static1 := s.staticResistanceCalculate(data.Flange1, auxiliary.Flange1, data.Type1, data, req, tightness.Pb, tightness.Pbr, tightness.Qd, tightness.Qfm)
	conditions1 := s.conditionsForStrengthCalculate(data.Type1, data.Flange1, auxiliary.Flange1, static1, req.IsWork, false)
	Static1 := []*flange_model.CalcStaticResistance{}
	Conditions1 := []*flange_model.CalcConditionsForStrength{}
	Static1 = append(Static1, static1)
	Conditions1 = append(Conditions1, conditions1)
	if len(req.FlangesData) > 1 {
		static2 := s.staticResistanceCalculate(data.Flange2, auxiliary.Flange2, data.Type2, data, req, tightness.Pb, tightness.Pbr, tightness.Qd, tightness.Qfm)
		conditions2 := s.conditionsForStrengthCalculate(data.Type2, data.Flange2, auxiliary.Flange2, static2, req.IsWork, false)
		Static1 = append(Static1, static2)
		Conditions1 = append(Conditions1, conditions2)
	}

	tigLoad := s.tightnessLoadCalculate(auxiliary, tightness, data, req)
	bolt2 := s.boltStrengthCalculate(data, req, tigLoad.Pb, tigLoad.Pbr, auxiliary.A, auxiliary.Dcp, true)
	moment2 := s.momentCalculate(req.Friction, data, bolt2.SigmaB1, bolt2.DSigmaM, tigLoad.Pb, auxiliary.A, auxiliary.Dcp, false)

	static1 = s.staticResistanceCalculate(data.Flange1, auxiliary.Flange1, data.Type1, data, req, tigLoad.Pb, tigLoad.Pbr, tightness.Qd, tightness.Qfm)
	conditions1 = s.conditionsForStrengthCalculate(data.Type1, data.Flange1, auxiliary.Flange1, static1, req.IsWork, true)
	Static2 := []*flange_model.CalcStaticResistance{}
	Conditions2 := []*flange_model.CalcConditionsForStrength{}
	Static2 = append(Static2, static1)
	Conditions2 = append(Conditions2, conditions1)
	if len(req.FlangesData) > 1 {
		static2 := s.staticResistanceCalculate(data.Flange2, auxiliary.Flange2, data.Type2, data, req, tigLoad.Pb, tigLoad.Pbr, tightness.Qd, tightness.Qfm)
		conditions2 := s.conditionsForStrengthCalculate(data.Type2, data.Flange2, auxiliary.Flange2, static2, req.IsWork, true)
		Static2 = append(Static2, static2)
		Conditions2 = append(Conditions2, conditions2)
	}

	ok := (bolt2.VSigmaB1 && bolt2.VSigmaB2 && data.TypeGasket != flange_model.GasketData_Soft) ||
		(bolt2.VSigmaB1 && bolt2.VSigmaB2 && bolt2.Q <= float64(data.Gasket.PermissiblePres) && data.TypeGasket == flange_model.GasketData_Soft)

	if ok {
		ok = false
		var VTeta1, VTeta2, VTetaK1, VTetaK2 bool
		if Conditions2[0].CondTeta.X <= Conditions2[0].CondTeta.Y {
			VTeta1 = true
		}

		if data.Type1 == flange_model.FlangeData_free && Conditions2[0].CondTeta.X <= Conditions2[0].CondTeta.Y {
			VTetaK1 = true
		}

		if !req.IsSameFlange {
			if len(Conditions2) > 1 && Conditions2[1].CondTeta.X <= Conditions2[1].CondTeta.Y {
				VTeta2 = true
			}

			if data.Type2 == flange_model.FlangeData_free && len(Conditions2) > 1 && Conditions2[1].CondTetaK.X <= Conditions2[1].CondTetaK.Y {
				VTetaK2 = true
			}
		}

		if !req.IsSameFlange {
			commonCond := VTeta1 && VTeta2
			cond1 := commonCond && data.Type1 != flange_model.FlangeData_free && data.Type2 != flange_model.FlangeData_free
			cond2 := commonCond && data.Type1 == flange_model.FlangeData_free && data.Type2 != flange_model.FlangeData_free && VTetaK1
			cond3 := commonCond && data.Type1 != flange_model.FlangeData_free && data.Type2 == flange_model.FlangeData_free && VTetaK2
			cond4 := commonCond && data.Type1 == flange_model.FlangeData_free && data.Type2 == flange_model.FlangeData_free && VTetaK1 && VTetaK2

			if cond1 || cond2 || cond3 || cond4 {
				ok = true
			}
		} else {
			if (VTeta1 && data.Type1 != flange_model.FlangeData_free) ||
				(VTeta1 && data.Type1 == flange_model.FlangeData_free && VTetaK1) {
				ok = true
			}
		}
	}

	finalMoment := &flange_model.CalcMoment{}
	if ok {
		finalMoment = s.momentCalculate(req.Friction, data, bolt2.SigmaB1, bolt2.DSigmaM, tigLoad.Pb, auxiliary.A, auxiliary.Dcp, true)
	}

	deformation := &flange_model.CalcDeformation{
		B0:  auxiliary.B0,
		Dcp: auxiliary.Dcp,
		Po:  tightness.Po,
		Rp:  tightness.Rp,
	}
	forces := &flange_model.CalcForcesInBolts{
		A:      auxiliary.A,
		Qd:     tightness.Qd,
		Qfm:    tightness.Qfm,
		Qt:     tigLoad.Qt,
		Pb:     tigLoad.Pb,
		Alpha:  auxiliary.Alpha,
		AlphaM: auxiliary.AlphaM,
		Pb1:    tigLoad.Pb1,
		MinB:   0.4 * auxiliary.A * data.Bolt.SigmaAt20,
		Pb2:    tightness.Pb2,
		Pbr:    tigLoad.Pbr,
	}

	res := &flange_model.Calculated_Strength{
		Auxiliary:              auxiliary,
		Tightness:              tightness,
		BoltStrength1:          bolt1,
		Moment1:                moment1,
		StaticResistance1:      Static1,
		ConditionsForStrength1: Conditions1,
		TightnessLoad:          tigLoad,
		BoltStrength2:          bolt2,
		Moment2:                moment2,
		StaticResistance2:      Static2,
		ConditionsForStrength2: Conditions2,
		Deformation:            deformation,
		ForcesInBolts:          forces,
		FinalMoment:            finalMoment,
	}

	return res
}

// Расчет вспомогательных величин
func (s *FlangeService) auxiliaryCalculate(data models.DataFlange, req *calc_api.FlangeRequest) *flange_model.CalcAuxiliary {
	auxiliary := &flange_model.CalcAuxiliary{}

	if data.TypeGasket == flange_model.GasketData_Oval {
		// формула 4
		auxiliary.B0 = data.Gasket.Width / 4
		// формула ?
		auxiliary.Dcp = data.Gasket.DOut - data.Gasket.Width/2

	} else {
		if data.Gasket.Width <= constants.Bp {
			// формула 2
			auxiliary.B0 = data.Gasket.Width
		} else {
			// формула 3
			auxiliary.B0 = constants.B0 * math.Sqrt(data.Gasket.Width)
		}
		// формула 5
		auxiliary.Dcp = data.Gasket.DOut - auxiliary.B0
	}

	if data.TypeGasket == flange_model.GasketData_Soft {
		// Податливость прокладки
		auxiliary.Yp = (data.Gasket.Thickness * data.Gasket.Compression) / (data.Gasket.Epsilon * math.Pi * auxiliary.Dcp * data.Gasket.Width)
	}
	// приложение К пояснение к формуле К.2
	auxiliary.Lb = data.Bolt.Length + s.typeBolt[req.Type.String()]*data.Bolt.Diameter
	// формула К.2
	// Податливость болтов/шпилек
	auxiliary.Yb = auxiliary.Lb / (data.Bolt.EpsilonAt20 * data.Bolt.Area * float64(data.Bolt.Count))

	// формула 8
	// Суммарная площадь сечения болтов/шпилек
	auxiliary.A = float64(data.Bolt.Count) * data.Bolt.Area

	flange1 := s.auxFlangeCalculate(req.FlangesData[0].Type, data.Flange1, auxiliary.Dcp)
	flange2 := flange1
	auxiliary.Flange1 = flange1
	if len(req.FlangesData) > 1 {
		flange2 = s.auxFlangeCalculate(req.FlangesData[1].Type, data.Flange2, auxiliary.Dcp)
		auxiliary.Flange2 = flange2
	}

	if data.TypeGasket == flange_model.GasketData_Oval || data.Type1 == flange_model.FlangeData_free || data.Type2 == flange_model.FlangeData_free {
		// Для фланцев с овальными и восьмигранными прокладками и для свободных фланцев коэффициенты жесткости фланцевого соединения принимают равными 1.
		auxiliary.Alpha = 1
	} else {
		// формула (Е.11)
		// Коэффициент жесткости
		auxiliary.Alpha = 1 - (auxiliary.Yp-(flange1.Yf*flange1.E*flange1.B+flange2.Yf*flange2.E*flange2.B))/
			(auxiliary.Yp+auxiliary.Yb+(flange1.Yf*math.Pow(flange1.B, 2)+flange2.Yf*math.Pow(flange2.B, 2)))
	}

	dividend := auxiliary.Yb + flange1.Yfn*flange1.B*(flange1.B+flange1.E-math.Pow(flange1.E, 2)/auxiliary.Dcp) +
		+flange2.Yfn*flange2.B*(flange2.B+flange2.E-math.Pow(flange2.E, 2)/auxiliary.Dcp)
	divider := auxiliary.Yb + auxiliary.Yp*math.Pow(data.Flange1.D6/auxiliary.Dcp, 2) + flange1.Yfn*math.Pow(flange1.B, 2) + flange2.Yfn*math.Pow(flange2.B, 2)

	if data.Type1 == flange_model.FlangeData_free {
		dividend += flange1.Yfc * math.Pow(flange1.A, 2)
		divider += flange1.Yfc * math.Pow(flange1.A, 2)
	}
	if data.Type2 == flange_model.FlangeData_free {
		dividend += flange2.Yfc * math.Pow(flange2.A, 2)
		divider += flange2.Yfc * math.Pow(flange2.A, 2)
	}

	// формула (Е.13)
	// Коэффициент жесткости фланцевого соединения нагруженного внешним изгибающим моментом
	auxiliary.AlphaM = dividend / divider

	divider = auxiliary.Yp + auxiliary.Yb*data.Bolt.EpsilonAt20/data.Bolt.Epsilon + (flange1.Yf*data.Flange1.EpsilonAt20/data.Flange1.Epsilon)*
		math.Pow(flange1.B, 2) + (flange2.Yf*data.Flange2.EpsilonAt20/data.Flange2.Epsilon)*math.Pow(flange2.B, 2)

	if data.Type1 == flange_model.FlangeData_free {
		divider += (flange1.Yk * data.Flange1.Ring.EpsilonKAt20 / data.Flange1.Ring.EpsilonK) * math.Pow(flange1.A, 2)
	}
	if data.Type2 == flange_model.FlangeData_free {
		divider += (flange2.Yk * data.Flange2.Ring.EpsilonKAt20 / data.Flange2.Ring.EpsilonK) * math.Pow(flange2.A, 2)
	}

	// формула (Е.8)
	auxiliary.Gamma = 1 / divider

	return auxiliary
}

// дополнительные величины связанные с фланцем
func (s *FlangeService) auxFlangeCalculate(
	flangeType flange_model.FlangeData_Type,
	data *flange_model.FlangeResult,
	Dcp float64,
) *flange_model.CalcAuxiliary_Flange {
	flange := &flange_model.CalcAuxiliary_Flange{}
	if flangeType != flange_model.FlangeData_free {
		// Плечи действия усилий в болтах/шпильках
		flange.B = 0.5 * (data.D6 - Dcp)
	} else {
		flange.A = 0.5 * (data.D6 - data.Ds)
		flange.B = 0.5 * (data.Ds - Dcp)
	}

	if flangeType != flange_model.FlangeData_welded {
		// Эквивалентная толщина втулки
		flange.Se = data.S0
	} else {
		flange.X = data.L / (math.Sqrt(data.D * data.S0))
		flange.Beta = data.S1 / data.S0
		// Коэффициент зависящий от соотношения размеров конической втулки фланца
		flange.Xi = 1 + (flange.Beta-1)*flange.X/(flange.X+(1+flange.Beta)/4)
		flange.Se = flange.Xi * data.S0
	}

	// Плечо усилия от действия давления на фланец
	flange.E = 0.5 * (Dcp - data.D - flange.Se)
	// Параметр длины обечайки
	flange.L0 = math.Sqrt(data.D * data.S0)
	// Отношение наружного диаметра тарелки фланца к внутреннему диаметру
	flange.K = data.DOut / data.D

	dividend := math.Pow(flange.K, 2)*(1+8.55*(math.Log(flange.K)/math.Log(10))) - 1
	divider := (1.05 + 1.945*math.Pow(flange.K, 2)) * (flange.K - 1)
	flange.BetaT = dividend / divider

	divider = 1.36 * (math.Pow(flange.K, 2) - 1) * (flange.K - 1)
	flange.BetaU = dividend / divider

	dividend = 1 / (flange.K - 1)
	divider = 0.69 + 5.72*((math.Pow(flange.K, 2)*(math.Log(flange.K)/math.Log(10)))/(math.Pow(flange.K, 2)-1))
	flange.BetaY = dividend * divider

	dividend = math.Pow(flange.K, 2) + 1
	divider = math.Pow(flange.K, 2) - 1
	flange.BetaZ = dividend / divider

	if flangeType == flange_model.FlangeData_welded && data.S0 != data.S1 {
		flange.BetaF = s.graphic.CalculateBetaF(flange.Beta, flange.X)
		flange.BetaV = s.graphic.CalculateBetaV(flange.Beta, flange.X)
		flange.F = s.graphic.CalculateF(flange.Beta, flange.X)
	} else {
		flange.BetaF = constants.InitBetaF
		flange.BetaV = constants.InitBetaV
		flange.F = constants.InitF
	}

	flange.Lymda = (flange.BetaF*data.H+flange.L0)/(flange.BetaT*flange.L0) +
		(flange.BetaV*math.Pow(data.H, 3))/(flange.BetaU*flange.L0*math.Pow(data.S0, 2))

	// Угловая податливость фланца при затяжке
	flange.Yf = (0.91 * flange.BetaV) / (data.EpsilonAt20 * flange.Lymda * math.Pow(data.S0, 2) * flange.L0)

	if flangeType == flange_model.FlangeData_free {
		flange.Psik = 1.28 * (math.Log(data.Dnk/data.Dk) / math.Log(10))
		flange.Yk = 1 / (data.Ring.EpsilonKAt20 * math.Pow(data.Hk, 3) * flange.Psik)
	}

	if flangeType != flange_model.FlangeData_free {
		// Угловая податливость фланца нагруженного внешним изгибающим моментом
		flange.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
	} else {
		flange.Yfn = math.Pow(math.Pi/4, 3) * (data.Ds / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
		flange.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / (data.Ring.EpsilonKAt20 * data.Dnk * math.Pow(data.Hk, 3)))
	}

	return flange
}

// Расчет фланцевого соединения на прочность и герметичность без учета нагрузки вызванной стесненностью температурных деформаций
func (s *FlangeService) tightnessCalculate(aux *flange_model.CalcAuxiliary, data models.DataFlange, req *calc_api.FlangeRequest) *flange_model.CalcTightness {
	tightness := &flange_model.CalcTightness{}

	// формула 6
	// Усилие необходимое для смятия прокладки при затяжке
	tightness.Po = 0.5 * math.Pi * aux.Dcp * aux.B0 * data.Gasket.Pres

	if req.Pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях
		tightness.Rp = math.Pi * aux.Dcp * aux.B0 * data.Gasket.M * math.Abs(req.Pressure)
	}

	// формула 9
	// Равнодействующая нагрузка от давления
	tightness.Qd = 0.785 * math.Pow(aux.Dcp, 2) * req.Pressure

	temp1 := float64(req.AxialForce) + 4*math.Abs(float64(req.BendingMoment))/aux.Dcp
	temp2 := float64(req.AxialForce) - 4*math.Abs(float64(req.BendingMoment))/aux.Dcp

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	tightness.Qfm = math.Max(temp1, temp2)

	minB := 0.4 * aux.A * data.Bolt.SigmaAt20
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	tightness.Pb2 = math.Max(tightness.Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	tightness.Pb1 = aux.Alpha*(tightness.Qd+float64(req.AxialForce)) + tightness.Rp + 4*aux.AlphaM*math.Abs(float64(req.BendingMoment))/aux.Dcp

	// Расчетная нагрузка на болты/шпильки фланцевых соединений
	tightness.Pb = math.Max(tightness.Pb1, tightness.Pb2)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений в рабочих условиях
	tightness.Pbr = tightness.Pb + (1-aux.Alpha)*(tightness.Qd+float64(req.AxialForce)) + 4*(1-aux.AlphaM*math.Abs(float64(req.BendingMoment)))/aux.Dcp

	return tightness
}

// Расчет фланца на статическую прочность
func (s *FlangeService) staticResistanceCalculate(
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
			(4*math.Abs(float64(req.BendingMoment)))/(flange.D+flange.S1)) / temp
		static.SigmaMpm = (0.785*math.Pow(flange.D, 2)*req.Pressure + float64(req.AxialForce) -
			(4*math.Abs(float64(req.BendingMoment)))/(flange.D+flange.S1)) / temp
	}

	temp := math.Pi * (flange.D + flange.S0) * (flange.S0 - flange.C)
	// Меридиональные мембранные напряжения во втулке приварного встык фланца обечайке трубе
	// плоского фланца или обечайке трубе бурта свободного фланца в рабочих условиях
	// формула (ф. 37)
	// - для приварных встык фланцев с конической втулкой в сечении S1
	static.SigmaMp0 = (0.785*math.Pow(flange.D, 2)*req.Pressure + float64(req.AxialForce) +
		(4*math.Abs(float64(req.BendingMoment)))/(flange.D+flange.S0)) / temp
	// - для приварных встык фланцев с конической втулкой в сечении S0 приварных фланцев с прямой втулкой плоских фланцев и свободных фланцев
	static.SigmaMpm0 = (0.785*math.Pow(flange.D, 2)*req.Pressure + float64(req.AxialForce) -
		(4*math.Abs(float64(req.BendingMoment)))/(flange.D+flange.S0)) / temp

	// Окружные мембранные напряжения от действия давления во втулке приварного встык фланца обечайке
	// трубе плоского фланца или обечайке трубе бурта свободного фланца в сечении S0
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

// Условия статической прочности фланцев
func (s *FlangeService) conditionsForStrengthCalculate(
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

func (s *FlangeService) tightnessLoadCalculate(
	aux *flange_model.CalcAuxiliary,
	tig *flange_model.CalcTightness,
	data models.DataFlange,
	req *calc_api.FlangeRequest,
) *flange_model.CalcTightnessLoad {
	tightness := &flange_model.CalcTightnessLoad{}

	flange1 := aux.Flange1
	flange2 := flange1
	if aux.Flange2 != nil {
		flange2 = aux.Flange2
	}

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
	tightness.Pbr = tig.Pb + (1-aux.Alpha)*(tig.Qd+float64(req.AxialForce)) + Qt + 4*(1-aux.AlphaM*math.Abs(float64(req.BendingMoment)))/aux.Dcp

	return tightness
}
