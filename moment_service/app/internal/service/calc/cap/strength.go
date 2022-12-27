package cap

import (
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *CapService) strengthCalculate(data models.DataCap, req *calc_api.CapRequest) *cap_model.Calculated_Strength {
	auxiliary := s.auxiliaryCalculate(data, req)
	tightness := s.tightnessCalculate(auxiliary, data, req)
	bolt1 := s.boltStrengthCalculate(data, req, tightness.Pb, tightness.Pbr, auxiliary.A, auxiliary.Dcp, false)
	moment1 := s.momentCalculate(data, bolt1.SigmaB1, bolt1.DSigmaM, tightness.Pb, auxiliary.A, auxiliary.Dcp, false)

	static1 := s.staticResistanceCalculate(data.Flange, auxiliary.Flange, data.FlangeType, data, req, tightness.Pb, tightness.Pbr, tightness.Qd, tightness.Qfm)
	conditions1 := s.conditionsForStrengthCalculate(data.FlangeType, data.Flange, auxiliary.Flange, static1, req.Data.IsWork, false)

	tigLoad := s.tightnessLoadCalculate(auxiliary, tightness, data, req)
	bolt2 := s.boltStrengthCalculate(data, req, tigLoad.Pb, tigLoad.Pbr, auxiliary.A, auxiliary.Dcp, true)
	moment2 := s.momentCalculate(data, bolt2.SigmaB1, bolt2.DSigmaM, tigLoad.Pb, auxiliary.A, auxiliary.Dcp, false)

	static2 := s.staticResistanceCalculate(data.Flange, auxiliary.Flange, data.FlangeType, data, req, tigLoad.Pb, tigLoad.Pbr, tightness.Qd, tightness.Qfm)
	conditions2 := s.conditionsForStrengthCalculate(data.FlangeType, data.Flange, auxiliary.Flange, static2, req.Data.IsWork, true)

	ok := (bolt2.VSigmaB1 && bolt2.VSigmaB2 && data.TypeGasket != cap_model.GasketData_Soft) ||
		(bolt2.VSigmaB1 && bolt2.VSigmaB2 && bolt2.Q <= float64(data.Gasket.PermissiblePres) && data.TypeGasket == cap_model.GasketData_Soft)

	if ok {
		ok = false
		var VTeta1, VTetaK1 bool
		if conditions2.CondTeta.X <= conditions2.CondTeta.Y {
			VTeta1 = true
		}

		if data.FlangeType == cap_model.FlangeData_free && conditions2.CondTeta.X <= conditions2.CondTeta.Y {
			VTetaK1 = true
		}

		if (VTeta1 && data.FlangeType != cap_model.FlangeData_free) ||
			(VTeta1 && data.FlangeType == cap_model.FlangeData_free && VTetaK1) {
			ok = true
		}

	}

	finalMoment := &cap_model.CalcMoment{}
	if ok {
		finalMoment = s.momentCalculate(data, bolt2.SigmaB1, bolt2.DSigmaM, tigLoad.Pb, auxiliary.A, auxiliary.Dcp, true)
	}

	deformation := &cap_model.CalcDeformation{
		B0:  auxiliary.B0,
		Dcp: auxiliary.Dcp,
		Po:  tightness.Po,
		Rp:  tightness.Rp,
	}
	forces := &cap_model.CalcForcesInBolts{
		A:     auxiliary.A,
		Qd:    tightness.Qd,
		Qfm:   tightness.Qfm,
		Qt:    tigLoad.Qt,
		Pb:    tigLoad.Pb,
		Alpha: auxiliary.Alpha,
		Pb1:   tigLoad.Pb1,
		MinB:  0.4 * auxiliary.A * data.Bolt.SigmaAt20,
		Pb2:   tightness.Pb2,
		Pbr:   tigLoad.Pbr,
	}

	res := &cap_model.Calculated_Strength{
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

	return res
}

func (s *CapService) auxiliaryCalculate(data models.DataCap, req *calc_api.CapRequest) *cap_model.CalcAuxiliary {
	auxiliary := &cap_model.CalcAuxiliary{}

	if data.TypeGasket == cap_model.GasketData_Oval {
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

	if data.TypeGasket == cap_model.GasketData_Soft {
		// Податливость прокладки
		auxiliary.Yp = (data.Gasket.Thickness * data.Gasket.Compression) / (data.Gasket.Epsilon * math.Pi * auxiliary.Dcp * data.Gasket.Width)
	}
	// приложение К пояснение к формуле К.2
	auxiliary.Lb = data.Bolt.Length + s.typeBolt[req.Data.Type.String()]*data.Bolt.Diameter
	// формула К.2
	// Податливость болтов/шпилек
	auxiliary.Yb = auxiliary.Lb / (data.Bolt.EpsilonAt20 * data.Bolt.Area * float64(data.Bolt.Count))

	// формула 8
	// Суммарная площадь сечения болтов/шпилек
	auxiliary.A = float64(data.Bolt.Count) * data.Bolt.Area

	flange := s.auxFlangeCalculate(req.FlangeData.Type, data.Flange, auxiliary.Dcp)
	cap := s.auxCapCalculate(req.CapData.Type, data.Cap, data.Flange, auxiliary.Dcp)

	auxiliary.Flange = flange
	auxiliary.Cap = cap

	if data.TypeGasket == cap_model.GasketData_Oval || data.FlangeType == cap_model.FlangeData_free {
		// Для фланцев с овальными и восьмигранными прокладками и для свободных фланцев коэффициенты жесткости фланцевого соединения принимают равными 1.
		auxiliary.Alpha = 1
	} else {
		// формула (Е.11)
		// Коэффициент жесткости
		auxiliary.Alpha = 1 - (auxiliary.Yp-(flange.Yf*flange.E+cap.Y*flange.B)*flange.B)/
			(auxiliary.Yp+auxiliary.Yb+(flange.Yf+cap.Y)*math.Pow(flange.B, 2))
	}

	divider := auxiliary.Yp + auxiliary.Yb*data.Bolt.EpsilonAt20/data.Bolt.Epsilon + (flange.Yf*data.Flange.EpsilonAt20/data.Flange.Epsilon)*
		math.Pow(flange.B, 2) + (cap.Y*data.Cap.EpsilonAt20/data.Cap.Epsilon)*math.Pow(flange.B, 2)

	if data.FlangeType == cap_model.FlangeData_free {
		divider += (flange.Yk * data.Flange.Ring.EpsilonAt20 / data.Flange.Ring.Epsilon) * math.Pow(flange.A, 2)
	}

	// формула (Е.8)
	auxiliary.Gamma = 1 / divider

	return auxiliary
}

func (s *CapService) auxFlangeCalculate(
	flangeType cap_model.FlangeData_Type,
	data *cap_model.FlangeResult,
	Dcp float64,
) *cap_model.CalcAuxiliary_Flange {
	flange := &cap_model.CalcAuxiliary_Flange{}

	if flangeType != cap_model.FlangeData_free {
		// Плечи действия усилий в болтах/шпильках
		flange.B = 0.5 * (data.D6 - Dcp)
	} else {
		flange.A = 0.5 * (data.D6 - data.Ring.Ds)
		flange.B = 0.5 * (data.Ring.Ds - Dcp)
	}

	if flangeType != cap_model.FlangeData_welded {
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

	if flangeType == cap_model.FlangeData_welded && data.S0 != data.S1 {
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

	if flangeType == cap_model.FlangeData_free {
		flange.Psik = 1.28 * (math.Log(data.Ring.Dnk/data.Ring.Dk) / math.Log(10))
		flange.Yk = 1 / (data.Ring.EpsilonAt20 * math.Pow(data.Ring.Hk, 3) * flange.Psik)
	}

	if flangeType != cap_model.FlangeData_free {
		// Угловая податливость фланца нагруженного внешним изгибающим моментом
		flange.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
	} else {
		flange.Yfn = math.Pow(math.Pi/4, 3) * (data.Ring.Ds / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
		flange.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / (data.Ring.EpsilonAt20 * data.Ring.Dnk * math.Pow(data.Ring.Hk, 3)))
	}

	return flange
}

func (s *CapService) auxCapCalculate(
	capType cap_model.CapData_Type,
	data *cap_model.CapResult,
	flange *cap_model.FlangeResult,
	Dcp float64,
) *cap_model.CalcAuxiliary_Cap {
	cap := &cap_model.CalcAuxiliary_Cap{}

	if capType == cap_model.CapData_flat {
		cap.K = flange.DOut / Dcp
		cap.X = (0.67*math.Pow(cap.K, 2)*(1+8.55*math.Log10(cap.K)) - 1) / ((cap.K - 1) *
			(math.Pow(cap.K, 2) - 1 + (1.857*math.Pow(cap.K, 2)+1)*(math.Pow(data.H, 3)/math.Pow(data.Delta, 3))))
		cap.Y = cap.X / (math.Pow(data.Delta, 3) * data.EpsilonAt20)
	} else {
		cap.Lambda = (flange.H / flange.D) * math.Sqrt(data.Radius/flange.S0)
		cap.Omega = 1 / (1 + 1.285*cap.Lambda + 1.63*cap.Lambda*math.Pow((flange.H/flange.S0), 2)*math.Log10(flange.DOut/flange.D))
		cap.Y = ((1 - cap.Omega*(1+1.285*cap.Lambda)) / (data.EpsilonAt20 * math.Pow(flange.H, 3))) *
			((flange.DOut + flange.D) / (flange.DOut - flange.D))
	}

	return cap
}

// Расчет фланцевого соединения на прочность и герметичность без учета нагрузки вызванной стесненностью температурных деформаций
func (s *CapService) tightnessCalculate(aux *cap_model.CalcAuxiliary, data models.DataCap, req *calc_api.CapRequest) *cap_model.CalcTightness {
	tightness := &cap_model.CalcTightness{}

	// формула 6
	// Усилие необходимое для смятия прокладки при затяжке
	tightness.Po = 0.5 * math.Pi * aux.Dcp * aux.B0 * data.Gasket.Pres

	if req.Data.Pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях
		tightness.Rp = math.Pi * aux.Dcp * aux.B0 * data.Gasket.M * math.Abs(req.Data.Pressure)
	}

	// формула 9
	// Равнодействующая нагрузка от давления
	tightness.Qd = 0.785 * math.Pow(aux.Dcp, 2) * req.Data.Pressure

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	tightness.Qfm = float64(req.Data.AxialForce)

	minB := 0.4 * aux.A * data.Bolt.SigmaAt20
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	tightness.Pb2 = math.Max(tightness.Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	tightness.Pb1 = aux.Alpha*(tightness.Qd+float64(req.Data.AxialForce)) + tightness.Rp

	// Расчетная нагрузка на болты/шпильки фланцевых соединений
	tightness.Pb = math.Max(tightness.Pb1, tightness.Pb2)
	// Расчетная нагрузка на болты/шпильки фланцевых соединений в рабочих условиях
	tightness.Pbr = tightness.Pb + (1-aux.Alpha)*(tightness.Qd+float64(req.Data.AxialForce))

	return tightness
}

// Расчет фланца на статическую прочность
func (s *CapService) staticResistanceCalculate(
	flange *cap_model.FlangeResult,
	calcFlange *cap_model.CalcAuxiliary_Flange,
	typeFlange cap_model.FlangeData_Type,
	data models.DataCap,
	req *calc_api.CapRequest,
	Pb, Pbr, Qd, Qfm float64,
) *cap_model.CalcStaticResistance {
	static := &cap_model.CalcStaticResistance{}

	temp1 := math.Pi * flange.D6 / float64(data.Bolt.Count)
	temp2 := 2*float64(data.Bolt.Diameter) + 6*flange.H/(data.Gasket.M+0.5)

	// Коэффициент учитывающий изгиб тарелки фланца между болтами шпильками
	static.Cf = math.Max(1, math.Sqrt(temp1/temp2))

	// Приведенный диаметр приварного встык фланца с конической или прямой втулкой
	var Dzv float64
	if typeFlange == cap_model.FlangeData_welded && flange.D <= 20*flange.S1 {
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

	if typeFlange == cap_model.FlangeData_free {
		static.MMk = static.Cf * Pb * calcFlange.A
		static.Mpk = static.Cf * Pbr * calcFlange.A
	}

	// Меридиональное изгибное напряжение во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке бурта свободного фланца
	if typeFlange == cap_model.FlangeData_welded && flange.S1 != flange.S0 {
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

	if typeFlange == cap_model.FlangeData_free {
		static.SigmaK = calcFlange.BetaY * static.MMk / (math.Pow(flange.Ring.Hk, 2) * flange.Ring.Dk)
	}

	// Меридиональные изгибные напряжения во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке
	// трубе бурта свободного фланца в рабочих условиях
	if typeFlange == cap_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		static.SigmaP1 = static.Mp / (calcFlange.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		static.SigmaP0 = calcFlange.F * static.SigmaP1
	} else {
		static.IsEqualSigma = true
		static.SigmaP1 = static.Mp / (calcFlange.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv)
		static.SigmaP0 = static.SigmaP1
	}

	if typeFlange == cap_model.FlangeData_welded {
		temp := math.Pi * (flange.D + flange.S1) * (flange.S1 - flange.C)
		// формула (ф. 37)
		static.SigmaMp = (0.785*math.Pow(flange.D, 2)*req.Data.Pressure + float64(req.Data.AxialForce) +
			4*math.Abs(float64(0)/(flange.D+flange.S1))) / temp
		static.SigmaMpm = (0.785*math.Pow(flange.D, 2)*req.Data.Pressure + float64(req.Data.AxialForce) -
			4*math.Abs(float64(0)/(flange.D+flange.S1))) / temp
	}

	temp := math.Pi * (flange.D + flange.S0) * (flange.S0 - flange.C)
	// Меридиональные мембранные напряжения во втулке приварного встык фланца обечайке трубе
	// плоского фланца или обечайке трубе бурта свободного фланца в рабочих условиях
	// формула (ф. 37)
	// - для приварных встык фланцев с конической втулкой в сечении S1
	static.SigmaMp0 = (0.785*math.Pow(flange.D, 2)*req.Data.Pressure + float64(req.Data.AxialForce) +
		4*math.Abs(float64(0)/(flange.D+flange.S0))) / temp
	// - для приварных встык фланцев с конической втулкой в сечении S0 приварных фланцев с прямой втулкой плоских фланцев и свободных фланцев
	static.SigmaMpm0 = (0.785*math.Pow(flange.D, 2)*req.Data.Pressure + float64(req.Data.AxialForce) -
		4*math.Abs(float64(0)/(flange.D+flange.S0))) / temp

	// Окружные мембранные напряжения от действия давления во втулке приварного встык фланца обечайке
	// трубе плоского фланца или обечайке трубе бурта свободного фланца в сечении S0
	static.SigmaMop = req.Data.Pressure * flange.D / (2.0 * (flange.S0 - flange.C))

	// Напряжения в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в рабочих условиях
	// - радиальные напряжения
	static.SigmaRp = ((1.33*calcFlange.BetaF*flange.H + calcFlange.L0) / (calcFlange.Lymda * math.Pow(flange.H, 2) * calcFlange.L0 * flange.D)) * static.Mp
	// - окружное напряжения
	static.SigmaTp = calcFlange.BetaY*static.Mp/(math.Pow(flange.H, 2)*flange.D) - calcFlange.BetaZ*static.SigmaRp

	if typeFlange == cap_model.FlangeData_free {
		static.SigmaKp = calcFlange.BetaY * static.Mp / (math.Pow(flange.Ring.Hk, 2) * flange.Ring.Dk)
	}

	return static
}

// Условия статической прочности фланцев
func (s *CapService) conditionsForStrengthCalculate(
	flangeType cap_model.FlangeData_Type,
	flange *cap_model.FlangeResult,
	calcFlange *cap_model.CalcAuxiliary_Flange,
	static *cap_model.CalcStaticResistance,
	isWork, isTemp bool,
) *cap_model.CalcConditionsForStrength {
	conditions := &cap_model.CalcConditionsForStrength{}

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
	DTeta = teta[isWork] * DTeta

	conditions.Teta = static.Mp * calcFlange.Yf * flange.EpsilonAt20 / flange.Epsilon
	conditions.CondTeta = &cap_model.Condition{
		X: conditions.Teta,
		Y: DTeta,
	}

	if flangeType == cap_model.FlangeData_free {
		//strength.DTetaK = 0.002
		DTetaK = 0.02
		DTetaK = teta[isWork] * DTetaK
		conditions.TetaK = static.Mpk * calcFlange.Yk * flange.Ring.EpsilonAt20 / flange.Ring.Epsilon
		conditions.CondTetaK = &cap_model.Condition{
			X: conditions.TetaK,
			Y: DTetaK,
		}
	}

	//* Условия статической прочности фланцев
	if flangeType == cap_model.FlangeData_welded && flange.S1 != flange.S0 {
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

		conditions.Max1 = &cap_model.Condition{
			X: Max1,
			Y: Ks * Kt[isTemp] * flange.SigmaMAt20,
		}
		conditions.Max2 = &cap_model.Condition{
			X: Max2,
			Y: Ks * Kt[isTemp] * flange.SigmaM,
		}
		conditions.Max3 = &cap_model.Condition{
			X: Max3,
			Y: 1.3 * flange.SigmaRAt20,
		}
		conditions.Max4 = &cap_model.Condition{
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

		conditions.Max5 = &cap_model.Condition{
			X: Max5,
			Y: flange.SigmaAt20,
		}
		conditions.Max6 = &cap_model.Condition{
			X: Max6,
			Y: flange.Sigma,
		}
	}

	max7 := math.Max(math.Abs(static.SigmaMp0), math.Abs(static.SigmaMpm0))
	Max7 := math.Max(max7, math.Abs(static.SigmaMop))
	Max8 := math.Max(math.Abs(static.SigmaR), math.Abs(static.SigmaT))
	Max9 := math.Max(math.Abs(static.SigmaRp), math.Abs(static.SigmaTp))

	conditions.Max7 = &cap_model.Condition{X: Max7, Y: flange.Sigma}
	conditions.Max8 = &cap_model.Condition{X: Max8, Y: Kt[isTemp] * flange.SigmaAt20}
	conditions.Max9 = &cap_model.Condition{X: Max9, Y: Kt[isTemp] * flange.Sigma}

	if flangeType == cap_model.FlangeData_free {
		Max10 := static.SigmaK
		Max11 := static.SigmaKp

		conditions.Max10 = &cap_model.Condition{X: Max10, Y: Kt[isTemp] * flange.Ring.SigmaAt20}
		conditions.Max11 = &cap_model.Condition{X: Max11, Y: Kt[isTemp] * flange.Ring.Sigma}
	}

	return conditions
}

func (s *CapService) tightnessLoadCalculate(
	aux *cap_model.CalcAuxiliary,
	tig *cap_model.CalcTightness,
	data models.DataCap,
	req *calc_api.CapRequest,
) *cap_model.CalcTightnessLoad {
	tightness := &cap_model.CalcTightnessLoad{}

	flange := aux.Flange
	cap := aux.Cap

	divider := aux.Yp + aux.Yb*data.Bolt.EpsilonAt20/data.Bolt.Epsilon + (flange.Yf*data.Flange.EpsilonAt20/data.Flange.Epsilon)*math.Pow(flange.B, 2) +
		(cap.Y*data.Cap.EpsilonAt20/data.Cap.Epsilon)*math.Pow(flange.B, 2)

	if data.FlangeType == cap_model.FlangeData_free {
		divider += (flange.Yk * data.Flange.Ring.EpsilonAt20 / data.Flange.Ring.Epsilon) * math.Pow(flange.A, 2)
	}

	// формула (Е.8)
	gamma := 1 / divider

	var temp1, temp2 float64
	if req.IsUseWasher {
		temp1 = (data.Flange.Alpha*data.Flange.H+data.Washer1.Alpha*data.Washer1.Thickness)*(data.Flange.T-20) +
			(data.Cap.Alpha*data.Cap.H+data.Washer2.Alpha*data.Washer2.Thickness)*(data.Cap.T-20)
	} else {
		temp1 = data.Flange.Alpha*data.Flange.H*(data.Flange.T-20) + data.Cap.Alpha*data.Cap.H*(data.Cap.T-20)
	}
	temp2 = data.Flange.H + data.Flange.H

	if data.FlangeType == cap_model.FlangeData_free {
		temp1 += data.Flange.Ring.Alpha * data.Flange.Ring.Hk * (data.Flange.Ring.T - 20)
		temp2 += data.Flange.Ring.Hk
	}
	if req.Data.IsEmbedded {
		temp1 += data.Embed.Alpha * data.Embed.Thickness * (req.Data.Temp - 20)
		temp2 += data.Embed.Thickness
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	Qt := gamma * (temp1 - data.Bolt.Alpha*temp2*(data.Bolt.Temp-20))
	tightness.Qt = Qt

	tightness.Pb1 = math.Max(tig.Pb1, tig.Pb1-Qt)
	tightness.Pb = math.Max(tightness.Pb1, tig.Pb2)
	tightness.Pbr = tig.Pb + (1-aux.Alpha)*(tig.Qd+float64(req.Data.AxialForce)) + Qt

	return tightness
}
