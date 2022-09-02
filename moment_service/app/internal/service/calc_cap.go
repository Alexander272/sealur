package service

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/formulas"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type CalcCapService struct {
	graphic  *GraphicService
	data     *DataService
	formulas formulas.Cap
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewCalcCapService(graphic *GraphicService, data *DataService, formulas formulas.Cap) *CalcCapService {
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

	return &CalcCapService{
		graphic:  graphic,
		data:     data,
		formulas: formulas,
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

func (s *CalcCapService) Calculation(ctx context.Context, data *moment_api.CalcCapRequest) (*moment_api.CapResponse, error) {
	d, err := s.data.GetDataForCap(ctx, data)
	if err != nil {
		return nil, err
	}

	result := moment_api.CapResponse{
		Data:   s.formatInitData(data),
		Bolt:   d.Bolt,
		Calc:   &moment_api.CalculatedCap{},
		Flange: d.Flange,
		Cap:    d.Cap,
		Embed:  d.Embed,
		Gasket: d.Gasket,
	}

	if data.IsUseWasher {
		result.Washers = append(result.Washers, d.Washer1, d.Washer2)
	}

	if data.Calculation == moment_api.CalcCapRequest_basis {
		result.Calc.Basis = &moment_api.CalcMomentBasis{}
	} else {
		result.Calc.Strength = &moment_api.CalcMomentStrength{}
	}

	Lb0 := d.Gasket.Thickness + d.Flange.H + d.Cap.H

	if d.FType == moment_api.FlangeData_free {
		Lb0 += d.Flange.Hk
	}

	result.Calc.B0 = d.B0
	result.Calc.Dsp = d.Dcp
	result.Bolt.Lenght = Lb0

	var yp float64 = 0
	if d.TypeGasket == "Soft" {
		yp = (d.Gasket.Thickness * d.Gasket.Compression) / (d.Gasket.Epsilon * math.Pi * d.Dcp * d.Gasket.Width)
	}

	// приложение К пояснение к формуле К.2
	Lb := Lb0 + s.typeBolt[data.Type.String()]*d.Bolt.Diameter
	// формула К.2
	yb := Lb / (d.Bolt.EpsilonAt20 * d.Bolt.Area * float64(d.Bolt.Count))
	// фомула 8
	Ab := float64(d.Bolt.Count) * d.Bolt.Area
	result.Calc.A = Ab

	var alpha float64
	if d.TypeGasket == "Oval" || d.FType == moment_api.FlangeData_free {
		// Для фланцев с овальными и восьмигранными прокладками и для свободных фланцев коэффициенты жесткости фланцевого соединения принимают равными 1.
		alpha = 1
	} else {
		// формула (Е.11)
		alpha = 1 - (yp-(d.Flange.Yf*d.Flange.E+d.Cap.Y*d.Flange.B)*d.Flange.B)/
			(yp+yb+(d.Flange.Yf+d.Cap.Y)*math.Pow(d.Flange.B, 2))
	}
	result.Calc.Alpha = alpha

	// формула 6
	Pobg := 0.5 * math.Pi * d.Dcp * d.B0 * d.Gasket.Pres

	var Rp float64 = 0
	if data.Pressure >= 0 {
		// формула 7
		Rp = math.Pi * d.Dcp * d.B0 * d.Gasket.M * math.Abs(data.Pressure)
	}

	// формула 9
	Qd := 0.785 * math.Pow(d.Dcp, 2) * float64(data.Pressure)

	// формула 10
	QFM := float64(data.AxialForce)

	result.Calc.Po = Pobg
	result.Calc.Rp = Rp
	result.Calc.Qd = Qd
	result.Calc.Qfm = QFM

	minB := 0.4 * Ab * d.Bolt.SigmaAt20
	Pb2 := math.Max(Pobg, minB)
	Pb1 := alpha*(Qd+float64(data.AxialForce)) + Rp

	if data.Calculation != moment_api.CalcCapRequest_basis {
		result.Calc.Strength.MinB = minB
		result.Calc.Strength.FPb1 = Pb1
		result.Calc.Strength.FPb2 = Pb2
		result.Calc.Strength.Yp = yp
		result.Calc.Strength.Yb = yb
		result.Calc.Strength.Lb = Lb

		Pbm := math.Max(Pb1, Pb2)
		Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce))
		result.Calc.Strength.FPb = Pbm
		result.Calc.Strength.FPbr = Pbr

		result.Calc.Strength.FSigmaB1 = Pbm / Ab
		result.Calc.Strength.FSigmaB2 = Pbr / Ab

		Kyp := s.Kyp[data.IsWork]
		Kyz := s.Kyz[data.Condition.String()]
		Kyt := constants.NoLoadKyt

		result.Calc.Strength.FDSigmaM = 1.2 * Kyp * Kyz * Kyt * d.Bolt.SigmaAt20
		result.Calc.Strength.FDSigmaR = Kyp * Kyz * Kyt * d.Bolt.Sigma

		if result.Calc.Strength.FSigmaB1 > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter {
			result.Calc.Strength.FMkp = s.graphic.CalculateMkp(d.Bolt.Diameter, result.Calc.Strength.FSigmaB2)
		} else {
			result.Calc.Strength.FMkp = (0.3 * Pbm * float64(d.Bolt.Diameter) / float64(d.Bolt.Count)) / 1000.0
		}

		result.Calc.Strength.FMkp1 = 0.75 * result.Calc.Strength.FMkp

		var qmax float64
		if d.TypeGasket == "Soft" {
			qmax = math.Max(Pbm, Pbr) / (math.Pi * d.Dcp * d.Gasket.Width)
		}
		result.Calc.Strength.FQ = qmax

		strength1 := s.getCalculatedStrength(
			d.Flange,
			d.Bolt,
			d.FType,
			d.Gasket.M,
			data.Pressure,
			Qd,
			d.Dcp,
			result.Calc.Strength.FSigmaB1,
			Pbm,
			Pbr,
			QFM,
			data.AxialForce,
			0,
			data.IsWork,
			false,
		)
		result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, strength1)
	}

	var divider, temp1, temp2 float64
	divider = yp + yb*d.Bolt.EpsilonAt20/d.Bolt.Epsilon + (d.Flange.Yf*d.Flange.EpsilonAt20/d.Flange.Epsilon)*math.Pow(d.Flange.B, 2) +
		+(d.Cap.Y*d.Cap.EpsilonAt20/d.Cap.Epsilon)*math.Pow(d.Flange.B, 2)

	if d.FType == moment_api.FlangeData_free {
		divider += (d.Flange.Yk * d.Flange.EpsilonKAt20 / d.Flange.EpsilonK) * math.Pow(d.Flange.A, 2)
	}

	// формула (Е.8)
	gamma := 1 / divider

	if data.IsUseWasher {
		temp1 = (d.Flange.AlphaF*d.Flange.H+d.Washer1.Alpha*data.Washer[0].Thickness)*(d.Flange.Tf-20) +
			+(d.Cap.Alpha*d.Cap.H+d.Washer2.Alpha*data.Washer[0].Thickness)*(d.Cap.T-20)
	} else {
		temp1 = d.Flange.AlphaF*d.Flange.H*(d.Flange.Tf-20) + d.Cap.Alpha*d.Cap.H*(d.Cap.T-20)
	}
	temp2 = d.Flange.H + d.Flange.H

	if d.FType == moment_api.FlangeData_free {
		temp1 += d.Flange.AlphaK * d.Flange.Hk * (d.Flange.Tk - 20)
		temp2 += d.Flange.Hk
	}
	if data.IsEmbedded {
		temp1 += d.Embed.Alpha * data.Embed.Thickness * (data.Temp - 20)
		temp2 += data.Embed.Thickness
	}

	//? должно быть два варианта формулы с шайбой и без нее
	//формула 11 (в старом 13)
	Qt := gamma * (temp1 - d.Bolt.Alpha*temp2*(d.Bolt.Temp-20))
	result.Calc.Qt = Qt

	Pb1 = math.Max(Pb1, Pb1-Qt)
	Pbm := math.Max(Pb1, Pb2)
	Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + Qt

	SigmaB1 := Pbm / Ab
	SigmaB2 := Pbr / Ab

	Kyp := s.Kyp[data.IsWork]
	Kyz := s.Kyz[data.Condition.String()]
	Kyt := constants.LoadKyt
	// формула Г.3
	DSigmaM := 1.2 * Kyp * Kyz * Kyt * d.Bolt.SigmaAt20
	// формула Г.4
	DSigmaR := Kyp * Kyz * Kyt * d.Bolt.Sigma

	var qmax float64

	if d.TypeGasket == "Soft" {
		qmax = math.Max(Pbm, Pbr) / (math.Pi * d.Dcp * d.Gasket.Width)
	}

	var v_sigmab1, v_sigmab2 bool
	if SigmaB1 <= DSigmaM {
		v_sigmab1 = true
	}
	if SigmaB2 <= DSigmaR {
		v_sigmab2 = true
	}

	if data.Calculation == moment_api.CalcCapRequest_basis {
		result.Calc.Basis.MinB = minB
		result.Calc.Basis.Pb1 = Pb1
		result.Calc.Basis.Pb2 = Pb2
		result.Calc.Basis.Pbr = Pbr
		result.Calc.Basis.Pb = Pbm
		result.Calc.Basis.Q = qmax
		result.Calc.Basis.SigmaB1 = SigmaB1
		result.Calc.Basis.SigmaB2 = SigmaB2
		result.Calc.Basis.DSigmaM = DSigmaM
		result.Calc.Basis.DSigmaR = DSigmaR
		result.Calc.Basis.VSigmaB1 = v_sigmab1
		result.Calc.Basis.VSigmaB2 = v_sigmab2

		if (v_sigmab1 && v_sigmab2 && d.TypeGasket != "Soft") ||
			(v_sigmab1 && v_sigmab2 && qmax <= float64(d.Gasket.PermissiblePres) && d.TypeGasket == "Soft") {
			if result.Calc.Basis.SigmaB1 > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter {
				result.Calc.Basis.Mkp = s.graphic.CalculateMkp(d.Bolt.Diameter, result.Calc.Basis.SigmaB1)
			} else {
				//? вроде как формула изменилась, но почему-то использовалась новая формула
				result.Calc.Basis.Mkp = (0.3 * Pbm * float64(d.Bolt.Diameter) / float64(d.Bolt.Count)) / 1000
			}

			result.Calc.Basis.Mkp1 = 0.75 * result.Calc.Basis.Mkp

			Prek := 0.8 * Ab * d.Bolt.SigmaAt20
			result.Calc.Basis.Qrek = Prek / (math.Pi * d.Dcp * d.Gasket.Width)
			result.Calc.Basis.Mrek = (0.3 * Prek * float64(d.Bolt.Diameter) / float64(d.Bolt.Count)) / 1000

			Pmax := result.Calc.Basis.DSigmaM * Ab
			result.Calc.Basis.Qmax = Pmax / (math.Pi * d.Dcp * d.Gasket.Width)

			if d.TypeGasket == "Soft" && result.Calc.Basis.Qmax > d.Gasket.PermissiblePres {
				Pmax = float64(d.Gasket.PermissiblePres) * (math.Pi * d.Dcp * d.Gasket.Width)
				result.Calc.Basis.Qmax = float64(d.Gasket.PermissiblePres)
			}

			result.Calc.Basis.Mmax = (0.3 * Pmax * float64(d.Bolt.Diameter) / float64(d.Bolt.Count)) / 1000
		}
	} else {
		result.Calc.Strength.Gamma = gamma
		result.Calc.Strength.SPb1 = Pb1
		result.Calc.Strength.SPb2 = Pb2
		result.Calc.Strength.SPbr = Pbr
		result.Calc.Strength.SPb = Pbm
		result.Calc.Strength.SQ = qmax
		result.Calc.Strength.SSigmaB1 = SigmaB1
		result.Calc.Strength.SSigmaB2 = SigmaB2
		result.Calc.Strength.SDSigmaM = DSigmaM
		result.Calc.Strength.SDSigmaR = DSigmaR
		result.Calc.Strength.VSigmaB1 = v_sigmab1
		result.Calc.Strength.VSigmaB2 = v_sigmab2

		strength1 := s.getCalculatedStrength(
			d.Flange,
			d.Bolt,
			d.FType,
			d.Gasket.M,
			data.Pressure,
			Qd,
			d.Dcp,
			result.Calc.Strength.SSigmaB1,
			Pbm,
			Pbr,
			QFM,
			data.AxialForce,
			0,
			data.IsWork,
			true,
		)
		result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, strength1)

		if d.TypeGasket == "Soft" && qmax <= d.Gasket.PermissiblePres {
			result.Calc.Strength.VQmax = true
		}

		if strength1.Teta <= strength1.DTeta {
			result.Calc.Strength.VTeta1 = true
		}

		if d.FType == moment_api.FlangeData_free && strength1.TetaK <= strength1.DTetaK {
			result.Calc.Strength.VTetaK1 = true
		}

		if result.Calc.Strength.SSigmaB1 > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter {
			result.Calc.Strength.SMkp = s.graphic.CalculateMkp(d.Bolt.Diameter, result.Calc.Strength.SSigmaB1)
		} else {
			//? вроде как формула изменилась, но почему-то использовалась новая формула
			result.Calc.Strength.SMkp = (0.3 * Pbm * float64(d.Bolt.Diameter) / float64(d.Bolt.Count)) / 1000
		}
		result.Calc.Strength.SMkp1 = 0.75 * result.Calc.Strength.SMkp

		if (v_sigmab1 && v_sigmab2 && d.TypeGasket != "Soft") ||
			(v_sigmab1 && v_sigmab2 && qmax <= float64(d.Gasket.PermissiblePres) && d.TypeGasket == "Soft") {

			if (result.Calc.Strength.VTeta1 && d.FType != moment_api.FlangeData_free) ||
				(result.Calc.Strength.VTeta1 && d.FType == moment_api.FlangeData_free && result.Calc.Strength.VTetaK1) {

				Prek := 0.8 * Ab * d.Bolt.SigmaAt20
				result.Calc.Strength.Qrek = Prek / (math.Pi * d.Dcp * d.Gasket.Width)
				result.Calc.Strength.Mrek = (0.3 * Prek * float64(d.Bolt.Diameter) / float64(d.Bolt.Count)) / 1000

				Pmax := result.Calc.Strength.SDSigmaM * Ab
				result.Calc.Strength.Qmax = Pmax / (math.Pi * d.Dcp * d.Gasket.Width)

				if d.TypeGasket == "Soft" && result.Calc.Strength.Qmax > d.Gasket.PermissiblePres {
					Pmax = float64(d.Gasket.PermissiblePres) * (math.Pi * d.Dcp * d.Gasket.Width)
					result.Calc.Strength.Qmax = float64(d.Gasket.PermissiblePres)
				}

				result.Calc.Strength.Mmax = (0.3 * Pmax * float64(d.Bolt.Diameter) / float64(d.Bolt.Count)) / 1000
			}
		}
	}

	if data.IsNeedFormulas {
		result.Formulas = s.formulas.GetFormulasForCap(
			d.TypeGasket,
			data.Condition.String(),
			data.IsWork, data.IsUseWasher, data.IsEmbedded,
			d,
			result,
			data.Calculation,
			gamma, yb, yp,
		)
	}

	return &result, nil
}

func (s *CalcCapService) formatInitData(data *moment_api.CalcCapRequest) *moment_api.DataResult {
	work := map[bool]string{
		true:  "Рабочие условия",
		false: "Условия испытаний",
	}
	flanges := map[string]string{
		"isolated":    "Изолированные фланцы",
		"nonIsolated": "Неизолированные фланцы",
		"manually":    "Задается вручную",
	}
	sameFlange := map[bool]string{
		true:  "Да",
		false: "Нет",
	}
	typeD := map[string]string{
		"pin":  "Шпилька",
		"bolt": "Болт",
	}
	condition := map[string]string{
		"uncontrollable":  "Неконтролируемая затяжка",
		"controllable":    "Контроль по крутящему моменту",
		"controllablePin": "Контроль по вытяжке шпилек",
	}

	res := &moment_api.DataResult{
		Pressure:   data.Pressure,
		AxialForce: data.AxialForce,
		Temp:       data.Temp,
		Work:       work[data.IsWork],
		Flanges:    flanges[data.Flanges.String()],
		Embedded:   sameFlange[data.IsEmbedded],
		Type:       typeD[data.Type.String()],
		Condition:  condition[data.Condition.String()],
	}
	return res
}

// расчеты если выполняется прочностной расчет
func (s *CalcCapService) getCalculatedStrength(
	flange *moment_api.FlangeResult,
	bolt *moment_api.BoltResult,
	typeF moment_api.FlangeData_Type,
	M, Pressure, Qd, Dcp, SigmaB, Pbm, Pbr, QFM float64,
	AxialForce, BendingMoment int32,
	isWork, isTemp bool,
) *moment_api.StrengthResult {
	//* большинтсво переменный называются +- так же как и в оригинале

	strength := &moment_api.StrengthResult{}
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

	strength.Cf = math.Max(1, math.Sqrt(temp1/temp2))

	var Dzv float64
	if typeF == moment_api.FlangeData_welded && flange.D <= 20*flange.S1 {
		if flange.F > 1 {
			Dzv = flange.D + flange.S0
		} else {
			Dzv = flange.D + flange.S1
		}
	} else {
		Dzv = flange.D
	}
	strength.Dzv = Dzv

	strength.MM = strength.Cf * Pbm * flange.B
	strength.Mp = strength.Cf * math.Max(Pbr*flange.B+(Qd+QFM)*flange.E, math.Abs(Qd+QFM)*flange.E)

	var sigmaM1, sigmaM0 float64
	if typeF == moment_api.FlangeData_free {
		strength.MMk = strength.Cf * Pbm * flange.A
		strength.Mpk = strength.Cf * Pbr * flange.A
	}

	if typeF == moment_api.FlangeData_welded && flange.S1 != flange.S0 {
		sigmaM1 = strength.MM / (flange.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv)
		sigmaM0 = flange.F * sigmaM1
	} else {
		sigmaM1 = strength.MM / (flange.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv)
		sigmaM0 = sigmaM1
	}

	sigmaR := ((1.33*flange.BetaF*flange.H + flange.L0) / (flange.Lymda * math.Pow(flange.H, 2) * flange.L0 * flange.D)) * strength.MM
	sigmaT := flange.BetaY*strength.MM/(math.Pow(flange.H, 2)*flange.D) - flange.BetaZ*sigmaR

	strength.SigmaR = sigmaR
	strength.SigmaT = sigmaT

	var sigmaK, sigmaP1, sigmaP0, sigmaMp, sigmaMpm float64
	if typeF == moment_api.FlangeData_free {
		sigmaK = flange.BetaY * strength.MMk / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	if typeF == moment_api.FlangeData_welded && flange.S1 != flange.S0 {
		sigmaP1 = strength.Mp / (flange.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv)
		sigmaP0 = flange.F * sigmaP1
	} else {
		strength.IsSameSigma = true
		sigmaP1 = strength.Mp / (flange.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv)
		sigmaP0 = sigmaP1
	}

	if typeF == moment_api.FlangeData_welded {
		temp := math.Pi * (flange.D + flange.S1) * (flange.S1 - flange.C)
		// формула (ф. 37)
		sigmaMp = (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) +
			4*math.Abs(float64(BendingMoment)/(flange.D+flange.S1))) / temp
		sigmaMpm = (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) -
			4*math.Abs(float64(BendingMoment)/(flange.D+flange.S1))) / temp
	}

	temp := math.Pi * (flange.D + flange.S0) * (flange.S0 - flange.C)
	// формула (ф. 37)
	sigmaMp0 := (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) +
		4*math.Abs(float64(BendingMoment)/(flange.D+flange.S0))) / temp
	sigmaMpm0 := (0.785*math.Pow(flange.D, 2)*Pressure + float64(AxialForce) -
		4*math.Abs(float64(BendingMoment)/(flange.D+flange.S0))) / temp
	sigmaMop := Pressure * flange.D / (2.0 * (flange.S0 - flange.C))

	sigmaRp := ((1.33*flange.BetaF*flange.H + flange.L0) / (flange.Lymda * math.Pow(flange.H, 2) * flange.L0 * flange.D)) * strength.Mp
	sigmaTp := flange.BetaY*strength.Mp/(math.Pow(flange.H, 2)*flange.D) - flange.BetaZ*sigmaRp

	var sigmaKp float64
	if typeF == moment_api.FlangeData_free {
		sigmaKp = flange.BetaY * strength.Mp / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	if typeF == moment_api.FlangeData_welded {
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

	if typeF == moment_api.FlangeData_free {
		//strength.DTetaK = 0.002
		strength.DTetaK = 0.02
		strength.DTetaK = teta[isWork] * strength.DTetaK
		strength.TetaK = strength.Mpk * flange.Yk * flange.EpsilonKAt20 / flange.EpsilonK
	}

	if typeF == moment_api.FlangeData_welded && flange.S1 != flange.S0 {
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

	if typeF == moment_api.FlangeData_free {
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
