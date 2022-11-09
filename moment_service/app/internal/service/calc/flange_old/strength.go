package flange_old

import (
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

// расчеты если выполняется прочностной расчет
func (s *FlangeService) getCalculatedStrength(
	flange *flange_model.OldFlangeResult,
	bolt *flange_model.BoltResult,
	typeF flange_model.FlangeData_Type,
	M, Pressure, Qd, Dcp, SigmaB, Pbm, Pbr, QFM float64,
	AxialForce, BendingMoment int32,
	isWork, isTemp bool,
) *flange_model.StrengthResult {
	//* большинтсво переменный называются +- так же как и в оригинале

	strength := &flange_model.StrengthResult{}
	teta := map[bool]float64{
		true:  constants.WorkTeta,
		false: constants.TestTeta,
	}
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

	temp1 := math.Pi * flange.D6 / float64(bolt.Count)
	temp2 := 2*float64(bolt.Diameter) + 6*flange.H/(M+0.5)

	// Коэффициент учитывающий изгиб тарелки фланца между болтами шпильками
	strength.Cf = math.Max(1, math.Sqrt(temp1/temp2))

	// Приведенный диаметр приварного встык фланца с конической или прямой втулкой
	var Dzv float64
	if typeF == flange_model.FlangeData_welded && flange.D <= 20*flange.S1 {
		if flange.F > 1 {
			Dzv = flange.D + flange.S0
		} else {
			Dzv = flange.D + flange.S1
		}
	} else {
		Dzv = flange.D
	}
	strength.Dzv = Dzv

	// Расчетный изгибающий момент действующий на фланец при затяжке
	strength.MM = strength.Cf * Pbm * flange.B
	// Расчетный изгибающий момент действующий на фланец в рабочих условиях
	strength.Mp = strength.Cf * math.Max(Pbr*flange.B+(Qd+QFM)*flange.E, math.Abs(Qd+QFM)*flange.E)

	if typeF == flange_model.FlangeData_free {
		strength.MMk = strength.Cf * Pbm * flange.A
		strength.Mpk = strength.Cf * Pbr * flange.A
	}

	// Меридиональное изгибное напряжение во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке бурта свободного фланца
	var sigmaM1, sigmaM0 float64
	if typeF == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		sigmaM1 = strength.MM / (flange.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		sigmaM0 = flange.F * sigmaM1
	} else {
		sigmaM1 = strength.MM / (flange.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv)
		sigmaM0 = sigmaM1
	}

	// Радиальное напряжение в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в условиях затяжки
	sigmaR := ((1.33*flange.BetaF*flange.H + flange.L0) / (flange.Lymda * math.Pow(flange.H, 2) * flange.L0 * flange.D)) * strength.MM
	// Окружное напряжение в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в условиях затяжки
	sigmaT := flange.BetaY*strength.MM/(math.Pow(flange.H, 2)*flange.D) - flange.BetaZ*sigmaR

	strength.SigmaR = sigmaR
	strength.SigmaT = sigmaT

	var sigmaK, sigmaP1, sigmaP0, sigmaMp, sigmaMpm float64
	if typeF == flange_model.FlangeData_free {
		sigmaK = flange.BetaY * strength.MMk / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	// Меридиональные изгибные напряжения во втулке приварного встык фланца обечайке трубе плоского фланца или обечайке
	// трубе бурта свободного фланца в рабочих условиях
	if typeF == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		// - для приварных встык фланцев с конической втулкой в сечении S1
		sigmaP1 = strength.Mp / (flange.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv)
		// - для приварных встык фланцев с конической втулкой в сечении S0
		sigmaP0 = flange.F * sigmaP1
	} else {
		strength.IsSameSigma = true
		sigmaP1 = strength.Mp / (flange.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv)
		sigmaP0 = sigmaP1
	}

	if typeF == flange_model.FlangeData_welded {
		temp := math.Pi * (flange.D + flange.S1) * (flange.S1 - flange.C)
		// формула (ф. 37)
		sigmaMp = (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) +
			4*math.Abs(float64(BendingMoment)/(flange.D+flange.S1))) / temp
		sigmaMpm = (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) -
			4*math.Abs(float64(BendingMoment)/(flange.D+flange.S1))) / temp
	}

	temp := math.Pi * (flange.D + flange.S0) * (flange.S0 - flange.C)
	// Меридиональные мембранные напряжения во втулке приварного встык фланца обечайке трубе
	// плоского фланца или обечайке трубе бурта свободного фланца в рабочих условиях
	// формула (ф. 37)
	// - для приварных встык фланцев с конической втулкой в сечении S1
	sigmaMp0 := (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) +
		4*math.Abs(float64(BendingMoment)/(flange.D+flange.S0))) / temp
	// - для приварных встык фланцев с конической втулкой в сечении S0 приварных фланцев с прямой втулкой плоских фланцев и свободных фланцев
	sigmaMpm0 := (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) -
		4*math.Abs(float64(BendingMoment)/(flange.D+flange.S0))) / temp

	// Окружные мембранные напряжения от действия давления во втулке приварного встык фланца обечайке
	// трубе плоского фланца или обечайке трубе бурта свободного фланца в сечениии S0
	sigmaMop := Pressure * flange.D / (2.0 * (flange.S0 - flange.C))

	// Напряжения в тарелке приварного встык фланца плоского фланца и бурте свободного фланца в рабочих условиях
	// - радиальные напряжения
	sigmaRp := ((1.33*flange.BetaF*flange.H + flange.L0) / (flange.Lymda * math.Pow(flange.H, 2) * flange.L0 * flange.D)) * strength.Mp
	// - окружное напряжения
	sigmaTp := flange.BetaY*strength.Mp/(math.Pow(flange.H, 2)*flange.D) - flange.BetaZ*sigmaRp

	var sigmaKp float64
	if typeF == flange_model.FlangeData_free {
		sigmaKp = flange.BetaY * strength.Mp / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	if typeF == flange_model.FlangeData_welded {
		if flange.D <= constants.MinD {
			strength.DTeta = constants.MinDTetta
		} else if flange.D > constants.MaxD {
			strength.DTeta = constants.MaxDTetta
		} else {
			strength.DTeta = ((flange.D-constants.MinD)/(constants.MaxD-constants.MinD))*
				(constants.MaxDTetta-constants.MinDTetta) + constants.MinDTetta
		}
	} else {
		strength.DTeta = constants.MaxDTetta
	}
	strength.DTeta = teta[isWork] * strength.DTeta

	strength.Teta = strength.Mp * flange.Yf * flange.EpsilonAt20 / flange.Epsilon

	if typeF == flange_model.FlangeData_free {
		//strength.DTetaK = 0.002
		strength.DTetaK = 0.02
		strength.DTetaK = teta[isWork] * strength.DTetaK
		strength.TetaK = strength.Mpk * flange.Yk * flange.EpsilonKAt20 / flange.EpsilonK
	}

	//* Условия статической прочности фланцев
	if typeF == flange_model.FlangeData_welded && flange.S1 != flange.S0 {
		strength.Max1 = math.Max(math.Abs(sigmaM1+sigmaR), math.Abs(sigmaM1+sigmaT))

		t1 := math.Max(math.Abs(sigmaP1-sigmaMp+sigmaRp), math.Abs(sigmaP1-sigmaMpm+sigmaRp))
		t2 := math.Max(math.Abs(sigmaP1-sigmaMp+sigmaTp), math.Abs(sigmaP1-sigmaMpm+sigmaTp))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(sigmaP1+sigmaMp), math.Abs(sigmaP1+sigmaMpm))

		strength.Max2 = math.Max(t1, t2)
		strength.Max3 = sigmaM0

		t1 = math.Max(math.Abs(sigmaP0+sigmaMp0), math.Abs(sigmaP0-sigmaMp0))
		t2 = math.Max(math.Abs(sigmaP0+sigmaMpm0), math.Abs(sigmaP0-sigmaMpm0))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.3*sigmaP0+sigmaMop), math.Abs(0.3*sigmaP0-sigmaMop))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.7*sigmaP0+(sigmaMp0-sigmaMop)), math.Abs(0.7*sigmaP0-(sigmaMp0-sigmaMop)))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.7*sigmaP0+(sigmaMpm0-sigmaMop)), math.Abs(0.7*sigmaP0-(sigmaMpm0-sigmaMop)))

		strength.Max4 = math.Max(t1, t2)

		strength.CondMax1 = Ks * Kt[isTemp] * flange.SigmaMAt20
		strength.CondMax2 = Ks * Kt[isTemp] * flange.SigmaM
		strength.CondMax3 = 1.3 * flange.SigmaRAt20
		strength.CondMax4 = 1.3 * flange.SigmaR
	} else {
		strength.Max5 = math.Max(math.Abs(sigmaM0+sigmaR), math.Abs(sigmaM0+sigmaT))

		t1 := math.Max(math.Abs(sigmaP0-sigmaMp0+sigmaTp), math.Abs(sigmaP0-sigmaMpm0+sigmaTp))
		t2 := math.Max(math.Abs(sigmaP0-sigmaMp0+sigmaRp), math.Abs(sigmaP0-sigmaMpm0+sigmaRp))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(sigmaP0+sigmaMp0), math.Abs(sigmaP0+sigmaMpm0))

		strength.Max6 = math.Max(t1, t2)

		strength.CondMax5 = flange.SigmaAt20
		strength.CondMax6 = flange.Sigma
	}

	max7 := math.Max(math.Abs(sigmaMp0), math.Abs(sigmaMpm0))
	strength.Max7 = math.Max(max7, math.Abs(sigmaMop))
	strength.Max8 = math.Max(math.Abs(sigmaR), math.Abs(sigmaT))
	strength.Max9 = math.Max(math.Abs(sigmaRp), math.Abs(sigmaTp))

	strength.CondMax7 = flange.Sigma
	strength.CondMax8 = Kt[isTemp] * flange.SigmaAt20
	strength.CondMax9 = Kt[isTemp] * flange.Sigma

	if typeF == flange_model.FlangeData_free {
		strength.Max10 = sigmaK
		strength.Max11 = sigmaKp

		strength.CondMax10 = Kt[isTemp] * flange.SigmaKAt20
		strength.CondMax11 = Kt[isTemp] * flange.SigmaK
	}

	strength.SigmaM0 = sigmaM0
	strength.SigmaM1 = sigmaM1
	strength.SigmaTp = sigmaTp
	strength.SigmaRp = sigmaRp
	strength.SigmaK = sigmaK
	strength.SigmaP1 = sigmaP1
	strength.SigmaP0 = sigmaP0
	strength.SigmaMp = sigmaMp
	strength.SigmaMpm = sigmaMpm
	strength.SigmaMp0 = sigmaMp0
	strength.SigmaMpm0 = sigmaMpm0
	strength.SigmaMop = sigmaMop
	strength.SigmaKp = sigmaKp

	return strength
}
