package flange

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/flange/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/flange/formulas"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type FlangeService struct {
	graphic  *graphic.GraphicService
	data     *data.DataService
	formulas *formulas.FormulasService
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewFlangeService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService) *FlangeService {
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

	data := data.NewDataService(flange, materials, gasket, graphic)
	formulas := formulas.NewFormulasService()

	return &FlangeService{
		graphic:  graphic,
		data:     data,
		formulas: formulas,
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

// расчет по ГОСТ 34233.4 - 2017
func (s *FlangeService) Calculation(ctx context.Context, data *moment_api.CalcFlangeRequest) (*moment_api.FlangeResponse, error) {
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := moment_api.FlangeResponse{
		Data:         s.data.FormatInitData(data),
		IsSameFlange: data.IsSameFlange,
		Bolt:         d.Bolt,
		Calc:         &moment_api.CalculatedFlange{},
		Flanges:      []*moment_api.FlangeResult{d.Flange1},
		Embed:        d.Embed,
		Gasket:       d.Gasket,
	}

	if data.IsUseWasher {
		result.Washers = append(result.Washers, d.Washer1)
		if !data.IsSameFlange {
			result.Washers = append(result.Washers, d.Washer2)
		}
	}

	if !data.IsSameFlange {
		result.Flanges = append(result.Flanges, d.Flange2)
	}

	if data.Calculation == moment_api.CalcFlangeRequest_basis {
		result.Calc.Basis = &moment_api.CalcMomentBasis{}
	} else {
		result.Calc.Strength = &moment_api.CalcMomentStrength{}
	}

	Lb0 := d.Gasket.Thickness + d.Flange1.H + d.Flange2.H

	if d.Type1 == moment_api.FlangeData_free {
		Lb0 += d.Flange1.Hk
	}
	if d.Type2 == moment_api.FlangeData_free {
		Lb0 += d.Flange2.Hk
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

	var alpha, dividend, divider float64

	if d.TypeGasket == "Oval" || d.Type1 == moment_api.FlangeData_free || d.Type2 == moment_api.FlangeData_free {
		// Для фланцев с овальными и восьмигранными прокладками и для свободных фланцев коэффициенты жесткости фланцевого соединения принимают равными 1.
		alpha = 1
	} else {
		// формула (Е.11)
		alpha = 1 - (yp-(d.Flange1.Yf*d.Flange1.E*d.Flange1.B+d.Flange2.Yf*d.Flange2.E*d.Flange2.B))/
			(yp+yb+(d.Flange1.Yf*math.Pow(d.Flange1.B, 2)+d.Flange2.Yf*math.Pow(d.Flange2.B, 2)))
	}
	result.Calc.Alpha = alpha

	dividend = yb + d.Flange1.Yfn*d.Flange1.B*(d.Flange1.B+d.Flange1.E-math.Pow(d.Flange1.E, 2)/d.Dcp) +
		+d.Flange2.Yfn*d.Flange2.B*(d.Flange2.B+d.Flange2.E-math.Pow(d.Flange2.E, 2)/d.Dcp)
	divider = yb + yp*math.Pow(d.Flange1.D6/d.Dcp, 2) + d.Flange1.Yfn*math.Pow(d.Flange1.B, 2) + d.Flange2.Yfn*math.Pow(d.Flange2.B, 2)

	if d.Type1 == moment_api.FlangeData_free {
		dividend += d.Flange1.Yfc * math.Pow(d.Flange1.A, 2)
		divider += d.Flange1.Yfc * math.Pow(d.Flange1.A, 2)
	}
	if d.Type2 == moment_api.FlangeData_free {
		dividend += d.Flange2.Yfc * math.Pow(d.Flange2.A, 2)
		divider += d.Flange2.Yfc * math.Pow(d.Flange2.A, 2)
	}

	//формула (Е.13)
	alphaM := dividend / divider
	result.Calc.AlphaM = alphaM
	// формула 6
	Pobg := 0.5 * math.Pi * d.Dcp * d.B0 * d.Gasket.Pres

	var Rp float64 = 0
	if data.Pressure >= 0 {
		// формула 7
		Rp = math.Pi * d.Dcp * d.B0 * d.Gasket.M * math.Abs(data.Pressure)
	}

	// формула 9
	Qd := 0.785 * math.Pow(d.Dcp, 2) * float64(data.Pressure)

	temp1 := float64(data.AxialForce) + 4*math.Abs(float64(data.BendingMoment))/d.Dcp
	temp2 := float64(data.AxialForce) - 4*math.Abs(float64(data.BendingMoment))/d.Dcp

	// формула 10
	QFM := math.Max(temp1, temp2)

	result.Calc.Po = Pobg
	result.Calc.Rp = Rp
	result.Calc.Qd = Qd
	result.Calc.Qfm = QFM

	minB := 0.4 * Ab * d.Bolt.SigmaAt20
	Pb2 := math.Max(Pobg, minB)
	Pb1 := alpha*(Qd+float64(data.AxialForce)) + Rp + 4*alphaM*math.Abs(float64(data.BendingMoment))/d.Dcp

	if data.Calculation != moment_api.CalcFlangeRequest_basis {
		result.Calc.Strength.MinB = minB
		result.Calc.Strength.FPb1 = Pb1
		result.Calc.Strength.FPb2 = Pb2
		result.Calc.Strength.Yp = yp
		result.Calc.Strength.Yb = yb
		result.Calc.Strength.Lb = Lb

		Pbm := math.Max(Pb1, Pb2)
		Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/d.Dcp
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
			d.Flange1,
			d.Bolt,
			d.Type1,
			d.Gasket.M,
			data.Pressure,
			Qd,
			d.Dcp,
			result.Calc.Strength.FSigmaB1,
			Pbm,
			Pbr,
			QFM,
			data.AxialForce,
			data.BendingMoment,
			data.IsWork,
			false,
		)
		result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, strength1)

		if len(data.FlangesData) > 1 {
			strength2 := s.getCalculatedStrength(
				d.Flange2,
				d.Bolt,
				d.Type2,
				d.Gasket.M,
				data.Pressure,
				Qd,
				d.Dcp,
				result.Calc.Strength.FSigmaB2,
				Pbm,
				Pbr,
				QFM,
				data.AxialForce,
				data.BendingMoment,
				data.IsWork,
				false,
			)
			result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, strength2)
		}
	}

	divider = yp + yb*d.Bolt.EpsilonAt20/d.Bolt.Epsilon + (d.Flange1.Yf*d.Flange1.EpsilonAt20/d.Flange1.Epsilon)*math.Pow(d.Flange1.B, 2) +
		+(d.Flange2.Yf*d.Flange2.EpsilonAt20/d.Flange2.Epsilon)*math.Pow(d.Flange2.B, 2)

	if d.Type1 == moment_api.FlangeData_free {
		divider += (d.Flange1.Yk * d.Flange1.EpsilonKAt20 / d.Flange1.EpsilonK) * math.Pow(d.Flange1.A, 2)
	}
	if d.Type2 == moment_api.FlangeData_free {
		divider += (d.Flange2.Yk * d.Flange2.EpsilonKAt20 / d.Flange2.EpsilonK) * math.Pow(d.Flange2.A, 2)
	}

	// формула (Е.8)
	gamma := 1 / divider

	if data.IsUseWasher {
		temp1 = (d.Flange1.AlphaF*d.Flange1.H+d.Washer1.Alpha*data.Washer[0].Thickness)*(d.Flange1.Tf-20) +
			+(d.Flange2.AlphaF*d.Flange2.H+d.Washer2.Alpha*data.Washer[0].Thickness)*(d.Flange2.Tf-20)
	} else {
		temp1 = d.Flange1.AlphaF*d.Flange1.H*(d.Flange1.Tf-20) + d.Flange2.AlphaF*d.Flange2.H*(d.Flange2.Tf-20)
	}
	temp2 = d.Flange1.H + d.Flange2.H

	if d.Type1 == moment_api.FlangeData_free {
		temp1 += d.Flange1.AlphaK * d.Flange1.Hk * (d.Flange1.Tk - 20)
		temp2 += d.Flange1.Hk
	}
	if d.Type2 == moment_api.FlangeData_free {
		temp1 += d.Flange2.AlphaK * d.Flange2.Hk * (d.Flange2.Tk - 20)
		temp2 += d.Flange2.Hk
	}
	if data.IsEmbedded {
		temp1 += d.Embed.Alpha * data.Embed.Thickness * (data.Temp - 20)
		temp2 += data.Embed.Thickness
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	Qt := gamma * (temp1 - d.Bolt.Alpha*temp2*(d.Bolt.Temp-20))
	result.Calc.Qt = Qt

	Pb1 = math.Max(Pb1, Pb1-Qt)
	Pbm := math.Max(Pb1, Pb2)
	Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + Qt + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/d.Dcp

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

	if data.Calculation == moment_api.CalcFlangeRequest_basis {
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
			d.Flange1,
			d.Bolt,
			d.Type1,
			d.Gasket.M,
			data.Pressure,
			Qd,
			d.Dcp,
			result.Calc.Strength.SSigmaB1,
			Pbm,
			Pbr,
			QFM,
			data.AxialForce,
			data.BendingMoment,
			data.IsWork,
			true,
		)
		result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, strength1)

		strength2 := &moment_api.StrengthResult{}
		if len(data.FlangesData) > 1 {
			strength2 = s.getCalculatedStrength(
				d.Flange2,
				d.Bolt,
				d.Type2,
				d.Gasket.M,
				data.Pressure,
				Qd,
				d.Dcp,
				result.Calc.Strength.SSigmaB2,
				Pbm,
				Pbr,
				QFM,
				data.AxialForce,
				data.BendingMoment,
				data.IsWork,
				true,
			)
			result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, strength2)
		}

		if d.TypeGasket == "Soft" && qmax <= d.Gasket.PermissiblePres {
			result.Calc.Strength.VQmax = true
		}

		if strength1.Teta <= strength1.DTeta {
			result.Calc.Strength.VTeta1 = true
		}

		if d.Type1 == moment_api.FlangeData_free && strength1.TetaK <= strength1.DTetaK {
			result.Calc.Strength.VTetaK1 = true
		}

		if !data.IsSameFlange {
			if strength2.Teta <= strength2.DTeta {
				result.Calc.Strength.VTeta2 = true
			}

			if d.Type2 == moment_api.FlangeData_free && strength2.TetaK <= strength2.DTetaK {
				result.Calc.Strength.VTetaK2 = true
			}
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
			ok := false

			if !data.IsSameFlange {
				commonCond := result.Calc.Strength.VTeta1 && result.Calc.Strength.VTeta2
				cond1 := commonCond && d.Type1 != moment_api.FlangeData_free && d.Type2 != moment_api.FlangeData_free
				cond2 := commonCond && d.Type1 == moment_api.FlangeData_free && d.Type2 != moment_api.FlangeData_free && result.Calc.Strength.VTetaK1
				cond3 := commonCond && d.Type1 != moment_api.FlangeData_free && d.Type2 == moment_api.FlangeData_free && result.Calc.Strength.VTetaK2
				cond4 := commonCond && d.Type1 == moment_api.FlangeData_free && d.Type2 == moment_api.FlangeData_free &&
					result.Calc.Strength.VTetaK1 && result.Calc.Strength.VTetaK2

				if cond1 || cond2 || cond3 || cond4 {
					ok = true
				}
			} else {
				if (result.Calc.Strength.VTeta1 && d.Type1 != moment_api.FlangeData_free) ||
					(result.Calc.Strength.VTeta1 && d.Type1 == moment_api.FlangeData_free && result.Calc.Strength.VTetaK1) {
					ok = true
				}
			}

			if ok {
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
		result.Formulas = s.formulas.GetFormulas(
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
