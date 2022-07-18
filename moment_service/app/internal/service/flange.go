package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type FlangeService struct {
	repo          repository.Flange
	materials     *MaterialsService
	gasket        *GasketService
	graphic       *GraphicService
	typeFlangesTF map[string]float64
	typeFlangesTB map[string]float64
	typeFlangesTK map[string]float64
	typeBolt      map[string]float64
	Kyp           map[bool]float64
	Kyz           map[string]float64
}

func NewFlangeService(repo repository.Flange, materials *MaterialsService, gasket *GasketService, graphic *GraphicService) *FlangeService {
	flangesTF := map[string]float64{
		"isolated":    constants.IsolatedFlatTf,
		"nonIsolated": constants.NonIsolatedFlatTf,
	}

	flangesTB := map[string]float64{
		"isolated":         constants.IsolatedFlatTb,
		"nonIsolated":      constants.NonIsolatedFlatTb,
		"isolated-free":    constants.IsolatedFreeTb,
		"nonIsolated-free": constants.NonIsolatedFlatTb,
	}

	flangeTK := map[string]float64{
		"isolated":    constants.IsolatedFreeTk,
		"nonIsolated": constants.NonIsolatedFreeTk,
	}

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

	return &FlangeService{
		repo:          repo,
		materials:     materials,
		gasket:        gasket,
		graphic:       graphic,
		typeFlangesTF: flangesTF,
		typeFlangesTB: flangesTB,
		typeFlangesTK: flangeTK,
		typeBolt:      bolt,
		Kyp:           kp,
		Kyz:           kz,
	}
}

//? можно расчет по основным формулам вынести в отдельный пакет, а потом просто использовать (должно сделать код более понятным)
func (s *FlangeService) Calculation(ctx context.Context, data *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error) {
	result := moment_proto.FlangeResponse{
		IsSameFlange: data.IsSameFlange,
		Bolt:         &moment_proto.BoltResult{},
		Calc: &moment_proto.CalculatedFlange{
			Strength: &moment_proto.CalcMoment{},
			Basis:    &moment_proto.CalcMoment{},
		},
		Gasket:   &moment_proto.GasketResult{},
		Formulas: &moment_proto.CalcFormulas{},
	}

	flange1, err := s.getDataFlange(ctx, data.FlangesData[0], data.Flanges.String(), data.Temp)
	if err != nil {
		return nil, err
	}

	result.Flanges = append(result.Flanges, &moment_proto.FlangeResult{
		DOut:         flange1.DOut,
		D:            flange1.D,
		Dk:           flange1.Dk,
		Dnk:          flange1.Dnk,
		Ds:           flange1.Ds,
		H:            flange1.H,
		Hk:           flange1.Hk,
		S0:           flange1.S0,
		S1:           flange1.S1,
		L:            flange1.L,
		D6:           flange1.D6,
		C:            flange1.C,
		Tf:           flange1.Tf,
		Tk:           flange1.Tk,
		AlphaK:       flange1.AlphaK,
		EpsilonKAt20: flange1.EpsilonKAt20,
		EpsilonK:     flange1.EpsilonK,
		SigmaKAt20:   flange1.SigmaKAt20,
		SigmaK:       flange1.SigmaK,
		AlphaF:       flange1.AlphaF,
		EpsilonAt20:  flange1.EpsilonAt20,
		Epsilon:      flange1.Epsilon,
		Sigma:        flange1.Sigma,
		SigmaAt20:    flange1.SigmaAt20,
		SigmaM:       flange1.SigmaM,
		SigmaMAt20:   flange1.SigmaMAt20,
		SigmaR:       flange1.SigmaR,
		SigmaRAt20:   flange1.SigmaRAt20,
		Material:     flange1.Material,
	})

	type1 := data.FlangesData[0].Type
	var type2 moment_proto.FlangeData_Type

	var flange2 models.InitialDataFlange
	if len(data.FlangesData) > 1 {
		flange2, err = s.getDataFlange(ctx, data.FlangesData[1], data.Flanges.String(), data.Temp)
		if err != nil {
			return nil, err
		}

		// res := moment_proto.FlangeResult(flange2)
		result.Flanges = append(result.Flanges, &moment_proto.FlangeResult{
			DOut:         flange2.DOut,
			D:            flange2.D,
			Dk:           flange2.Dk,
			Dnk:          flange2.Dnk,
			Ds:           flange2.Ds,
			H:            flange2.H,
			Hk:           flange2.Hk,
			S0:           flange2.S0,
			S1:           flange2.S1,
			L:            flange2.L,
			D6:           flange2.D6,
			C:            flange2.C,
			Tf:           flange2.Tf,
			Tk:           flange2.Tk,
			AlphaK:       flange2.AlphaK,
			EpsilonKAt20: flange2.EpsilonKAt20,
			EpsilonK:     flange2.EpsilonK,
			SigmaKAt20:   flange2.SigmaKAt20,
			SigmaK:       flange2.SigmaK,
			AlphaF:       flange2.AlphaF,
			EpsilonAt20:  flange2.EpsilonAt20,
			Epsilon:      flange2.Epsilon,
			Sigma:        flange2.Sigma,
			SigmaAt20:    flange2.SigmaAt20,
			SigmaM:       flange2.SigmaM,
			SigmaMAt20:   flange2.SigmaMAt20,
			SigmaR:       flange2.SigmaR,
			SigmaRAt20:   flange2.SigmaRAt20,
			Material:     flange2.Material,
		})
		type2 = data.FlangesData[1].Type
	} else {
		flange2 = flange1
		type2 = type1
	}

	Tb := s.typeFlangesTB[data.Flanges.String()] * data.Temp
	if data.FlangesData[0].Type == moment_proto.FlangeData_free {
		Tb = s.typeFlangesTB[data.Flanges.String()+"-free"] * data.Temp
	}
	//TODO учитывать возможность ввода вручную

	boltMat, err := s.materials.GetMatFotCalculate(ctx, data.Bolts.MarkId, Tb)
	if err != nil {
		return nil, err
	}
	result.Bolt = &moment_proto.BoltResult{
		Diameter:    flange1.Diameter,
		Area:        flange1.Area,
		Count:       flange1.Count,
		Lenght:      flange1.L,
		Temp:        Tb,
		Alpha:       boltMat.AlphaF,
		EpsilonAt20: boltMat.EpsilonAt20,
		Epsilon:     boltMat.Epsilon,
		SigmaAt20:   boltMat.SigmaAt20,
		Sigma:       boltMat.Sigma,
	}

	//TODO учесть ввод данных для прокладки (все значения заносятся ручками)
	gasket, err := s.gasket.Get(ctx, models.GetGasket{GasketId: data.Gasket.GasketId, EnvId: data.Gasket.EnvId, Thickness: data.Gasket.Thickness})
	if err != nil {
		return nil, err
	}

	Lb0 := gasket.Thickness
	Lb0 += flange1.H + flange2.H

	if type1 == moment_proto.FlangeData_free {
		Lb0 += flange1.Hk
	}
	if type2 == moment_proto.FlangeData_free {
		Lb0 += flange2.Hk
	}

	var detMat models.MaterialsResult
	if data.IsEmbedded {
		//* тут было получиние прокладки еще раз, но входные данные не менялись, так что я это убрал
		Lb0 += gasket.Thickness
		Lb0 += data.Embed.Thickness

		detMat, err = s.materials.GetMatFotCalculate(ctx, data.Embed.MarkId, data.Temp)
		if err != nil {
			return nil, err
		}

		result.Embed = &moment_proto.EmbedResult{
			MarkId:    data.Embed.MarkId,
			Thickness: data.Embed.Thickness,
			Alpfa:     detMat.AlphaF,
			Temp:      data.Temp,
		}
	}

	bp := (data.Gasket.DOut - data.Gasket.DIn) / 2

	result.Gasket = &moment_proto.GasketResult{
		GasketId:        data.Gasket.GasketId,
		EnvId:           data.Gasket.EnvId,
		Thickness:       data.Gasket.Thickness,
		DOut:            data.Gasket.DOut,
		Width:           bp,
		M:               gasket.M,
		Pres:            gasket.SpecificPres,
		PermissiblePres: gasket.PermissiblePres,
		Compression:     gasket.Compression,
		Epsilon:         gasket.Epsilon,
	}

	var b0, Dcp float64

	//TODO заменить на тип из proto
	if gasket.Type == "Oval" {
		b0 = bp / 4
		Dcp = data.Gasket.DOut - bp/2

		if data.IsNeedFormulas {
			result.Formulas.B0 = fmt.Sprintf("%.3f/4", bp)
			result.Formulas.Dsp = fmt.Sprintf("%.0f - %.3f/2", data.Gasket.DOut, bp)
		}
	} else {
		if bp <= constants.Bp {
			b0 = bp
		} else {
			b0 = constants.B0 * math.Sqrt(bp)

			if data.IsNeedFormulas {
				result.Formulas.B0 = fmt.Sprintf("%.1f * sqrt(%.3f)", constants.B0, bp)
			}
		}
		Dcp = data.Gasket.DOut - b0

		if data.IsNeedFormulas {
			result.Formulas.Dsp = fmt.Sprintf("%.0f - %.3f", data.Gasket.DOut, b0)
		}
	}
	result.Calc.B0 = b0
	result.Calc.Dsp = Dcp
	result.Bolt.Lenght = Lb0

	var yp float64 = 0
	//TODO заменить на тип из proto
	if gasket.Type == "Soft" {
		yp = float64(gasket.Thickness*gasket.Compression) / (gasket.Epsilon * math.Pi * Dcp * bp)
	}

	Lb := Lb0 + s.typeBolt[data.Type.String()]*float64(flange1.Diameter)

	yb := Lb / (boltMat.EpsilonAt20 * flange1.Area * float64(flange1.Count))
	Ab := float64(flange1.Count) * flange1.Area
	result.Calc.A = Ab

	if data.IsNeedFormulas {
		result.Formulas.A = fmt.Sprintf("%d * %.3f", flange1.Count, flange1.Area)
	}

	res1, err := s.getCalculatedData(ctx, data.FlangesData[0], flange1, Dcp)
	if err != nil {
		return nil, err
	}

	var res2 models.CalculatedData
	if len(data.FlangesData) > 1 {
		res2, err = s.getCalculatedData(ctx, data.FlangesData[1], flange2, Dcp)
		if err != nil {
			return nil, err
		}
	} else {
		res2 = res1
	}

	var alpha, dividend, divider float64

	divider = yp + yb*boltMat.EpsilonAt20/boltMat.Epsilon + (res1.Yf*flange1.EpsilonAt20/flange1.Epsilon)*math.Pow(res1.B, 2) + (res2.Yf*flange2.EpsilonAt20/flange2.Epsilon)*math.Pow(res2.B, 2)

	if type1 == moment_proto.FlangeData_free {
		divider += (res1.Yk * flange1.EpsilonKAt20 / flange1.EpsilonK) * math.Pow(res1.A, 2)
	}
	if type2 == moment_proto.FlangeData_free {
		divider += (res2.Yk * flange2.EpsilonKAt20 / flange2.EpsilonK) * math.Pow(res2.A, 2)
	}

	gamma := 1 / divider

	//TODO заменить на тип из proto
	if gasket.Type == "Oval" || type1 == moment_proto.FlangeData_free || type2 == moment_proto.FlangeData_free {
		alpha = 1
	} else {
		alpha = 1 - (yp-(res1.Yf*res1.E*res1.B+res2.Yf*res2.E*res2.B))/(yp+yb+(res1.Yf*math.Pow(res1.B, 2)+res2.Yf*math.Pow(res2.B, 2)))

		if data.IsNeedFormulas {
			result.Formulas.Alpha = fmt.Sprintf("1 - (%f - (%f*%f*%f + %f*%f*%f)/(%f + %f + (%f*%f^2 + %f*%f^2)))", yp, res1.Yf, res1.E, res1.B, res2.Yf, res2.E, res2.B, yp, yb, res1.Yf, res1.B, res2.Yf, res2.B)
		}
	}
	result.Calc.Alpha = alpha

	dividend = yb + res1.Yfn*res1.B*(res1.B+res1.E-math.Pow(res1.E, 2)/Dcp) + res2.Yfn*res2.B*(res2.B+res2.E-math.Pow(res2.E, 2)/Dcp)
	divider = yb + yp*math.Pow(flange1.D6/Dcp, 2) + res1.Yfn*math.Pow(res1.B, 2) + res2.Yfn*math.Pow(res2.B, 2)
	dividendF := fmt.Sprintf("(%f + %f*%f*(%f + %f - %f^2/%f) + %f*%f*(%f + %f - %f^2/%f)", yb, res1.Yfn, res1.B, res1.B, res1.E, res1.E, Dcp, res2.Yfn, res2.B, res2.B, res2.E, res2.E, Dcp)
	dividerF := fmt.Sprintf("(%f + %f*(%f/%f)^2 + %f*%f^2 + %f*%f^2", yb, yp, flange1.D6, Dcp, res1.Yfn, res1.B, res2.Yfn, res2.B)

	if type1 == moment_proto.FlangeData_free {
		dividend += res1.Yfc * math.Pow(res1.A, 2)
		divider += res1.Yfc * math.Pow(res1.A, 2)

		dividendF += fmt.Sprintf("%f*%f^2", res1.Yfc, res1.A)
		dividerF += fmt.Sprintf("%f*%f^2", res1.Yfc, res1.A)
	}
	if type2 == moment_proto.FlangeData_free {
		dividend += res2.Yfc * math.Pow(res2.A, 2)
		divider += res2.Yfc * math.Pow(res2.A, 2)

		dividendF += fmt.Sprintf("%f*%f^2", res2.Yfc, res2.A)
		dividerF += fmt.Sprintf("%f*%f^2", res2.Yfc, res2.A)
	}

	alphaM := dividend / divider
	result.Calc.AlphaM = alphaM
	if data.IsNeedFormulas {
		result.Formulas.AlphaM = dividendF + ") / " + dividerF + ")"
	}

	Pobg := 0.5 * math.Pi * Dcp * b0 * gasket.SpecificPres
	if data.IsNeedFormulas {
		result.Formulas.Po = fmt.Sprintf("0.5 * %f * %.3f * %f * %f", math.Pi, Dcp, b0, gasket.SpecificPres)
	}

	var Rp float64 = 0
	if data.Pressure >= 0 {
		//TODO формула изменилась
		Rp = math.Pi * Dcp * b0 * gasket.M * data.Pressure

		if data.IsNeedFormulas {
			result.Formulas.Rp = fmt.Sprintf("%f * %f * %f * %f * %.0f", math.Pi, Dcp, b0, gasket.M, data.Pressure)
		}
	}

	Qd := 0.785 * math.Pow(Dcp, 2) * float64(data.Pressure)

	temp1 := float64(data.AxialForce) + 4*math.Abs(float64(data.BendingMoment))/Dcp
	temp2 := float64(data.AxialForce) - 4*math.Abs(float64(data.BendingMoment))/Dcp

	QFM := math.Max(temp1, temp2)

	if data.IsNeedFormulas {
		result.Formulas.Qd = fmt.Sprintf("0.785 * %.3f^2 * %.0f", Dcp, data.Pressure)
		result.Formulas.Qfm = fmt.Sprintf("max((%d + 4*|%d/%.3f|);(%d - 4*|%d/%.3f|))", data.AxialForce, data.BendingMoment, Dcp, data.AxialForce, data.BendingMoment, Dcp)
	}

	result.Calc.Po = Pobg
	result.Calc.Rp = Rp
	result.Calc.Qd = Qd
	result.Calc.Qfm = QFM

	Pb2 := math.Max(Pobg, 0.4*Ab*boltMat.SigmaAt20)
	Pb1 := alpha*(Qd+float64(data.AxialForce)) + Rp + 4*alphaM*math.Abs(float64(data.BendingMoment))/Dcp

	// if ($Moment != 1)
	if data.Calculation != moment_proto.FlangeRequest_basis {
		result.Calc.Strength.Pb1 = Pb1
		result.Calc.Strength.Pb2 = Pb2

		Pbm := math.Max(Pb1, Pb2)
		Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/Dcp
		result.Calc.Strength.Pb = Pbm
		result.Calc.Strength.Pbr = Pbr

		result.Calc.Strength.SigmaB1 = Pbm / Ab
		result.Calc.Strength.SigmaB2 = Pbr / Ab

		Kyp := s.Kyp[data.IsWork]
		Kyz := s.Kyz[data.Condition.String()]
		Kyt := constants.NoLoadKyt

		result.Calc.Strength.DSigmaM = 1.2 * Kyp * Kyz * Kyt * boltMat.SigmaAt20
		result.Calc.Strength.DSigmaR = Kyp * Kyz * Kyt * boltMat.Sigma

		var qmax float64
		//TODO заменить на тип из proto
		if gasket.Type == "Soft" {
			qmax = math.Max(Pbm, Pbr) / (math.Pi * Dcp * bp)
		}
		result.Calc.Strength.Q = qmax

		strength1 := s.getCalculatedStrength(flange1, res1, type1, gasket.M, data.Pressure, Qd, Dcp, result.Calc.Strength.SigmaB1, Pbm, Pbr, QFM, data.AxialForce, data.BendingMoment)
		st := moment_proto.StrengthResult(strength1)
		result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, &st)

		if len(data.FlangesData) > 1 {
			strength2 := s.getCalculatedStrength(flange2, res2, type2, gasket.M, data.Pressure, Qd, Dcp, result.Calc.Strength.SigmaB2, Pbm, Pbr, QFM, data.AxialForce, data.BendingMoment)
			st := moment_proto.StrengthResult(strength2)
			result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, &st)
		}
	}

	temp1 = flange1.AlphaF*flange1.H*(flange1.Tf-20) + flange2.AlphaF*flange2.H*(flange2.Tf-20)
	temp2 = flange1.H + flange2.H

	if type1 == moment_proto.FlangeData_free {
		temp1 += flange1.AlphaK * flange1.Hk * (flange1.Tk - 20)
		temp2 += flange1.Hk
	}
	if type2 == moment_proto.FlangeData_free {
		temp1 += flange2.AlphaK * flange2.Hk * (flange2.Tk - 20)
		temp2 += flange2.Hk
	}
	if data.IsEmbedded {
		temp1 += detMat.AlphaF * data.Embed.Thickness * (data.Temp - 20)
		temp2 += data.Embed.Thickness
	}

	//TODO здесь должна быть новая формула (Qt)
	Qt := gamma * (temp1 - boltMat.AlphaF*temp2*(Tb-20))
	result.Calc.Qt = Qt

	Pb1 = math.Max(Pb1, Pb1-Qt)
	Pbm := math.Max(Pb1, Pb2)
	Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + Qt + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/Dcp

	SigmaB1 := Pbm / Ab
	SigmaB2 := Pbr / Ab

	Kyp := s.Kyp[data.IsWork]
	Kyz := s.Kyz[data.Condition.String()]
	Kyt := constants.LoadKyt
	// формула Г.3
	DSigmaM := 1.2 * Kyp * Kyz * Kyt * boltMat.SigmaAt20
	// формула Г.4
	DSigmaR := Kyp * Kyz * Kyt * boltMat.Sigma

	var qmax float64
	//TODO заменить на тип из proto
	if gasket.Type == "Soft" {
		qmax = math.Max(Pbm, Pbr) / (math.Pi * Dcp * bp)
	}

	var v_sigmab1, v_sigmab2 bool
	if SigmaB1 <= DSigmaM {
		v_sigmab1 = true
	}
	if SigmaB2 <= DSigmaR {
		v_sigmab2 = true
	}

	if data.Calculation == moment_proto.FlangeRequest_basis {
		result.Calc.Basis.Pb1 = Pb1
		result.Calc.Basis.Pb2 = Pb2
		result.Calc.Basis.Pbr = Pbr
		result.Calc.Basis.Pb = Pbm
		result.Calc.Basis.Q = qmax
		result.Calc.Basis.SigmaB1 = SigmaB1
		result.Calc.Basis.SigmaB2 = SigmaB2
		result.Calc.Basis.DSigmaM = DSigmaM
		result.Calc.Basis.DSigmaR = DSigmaR

		//TODO заменить на тип из proto
		if (v_sigmab1 && v_sigmab2 && gasket.Type != "Soft") || (v_sigmab1 && v_sigmab2 && qmax <= float64(gasket.PermissiblePres) && gasket.Type == "Soft") {
			// var Mkp float64
			if result.Calc.Basis.SigmaB1 > constants.MaxSigmaB && flange1.Diameter >= constants.MinDiameter && flange1.Diameter <= constants.MaxDiameter {
				result.Calc.Basis.Mkp = s.graphic.CalculateMkp(flange1.Diameter, result.Calc.Basis.SigmaB1)
			} else {
				//TODO вроде как формула изменилась
				// зачем-то делится на 1000
				result.Calc.Basis.Mkp = (0.3 * Pbm * float64(flange1.Diameter) / float64(flange1.Count)) / 1000
			}

			result.Calc.Basis.Mkp1 = 0.75 * result.Calc.Basis.Mkp

			Prek := 0.8 * Ab * boltMat.SigmaAt20
			result.Calc.Basis.Qrek = Prek / (math.Pi * Dcp * bp)
			result.Calc.Basis.Mrek = (0.3 * Prek * float64(flange1.Diameter) / float64(flange1.Count)) / 1000

			Pmax := result.Calc.Basis.DSigmaM * Ab
			result.Calc.Basis.Qmax = Pmax / (math.Pi * Dcp * bp)

			//TODO заменить на тип из proto
			if gasket.Type == "Soft" && result.Calc.Basis.Qmax > float64(gasket.PermissiblePres) {
				Pmax = float64(gasket.PermissiblePres) * (math.Pi * Dcp * bp)
				result.Calc.Basis.Qmax = float64(gasket.PermissiblePres)
			}

			result.Calc.Basis.Mmax = (0.3 * Pmax * float64(flange1.Diameter) / float64(flange1.Count)) / 1000
		}
	} else {
		strength1 := s.getCalculatedStrength(flange1, res1, type1, gasket.M, data.Pressure, Qd, Dcp, result.Calc.Strength.SigmaB1, Pbm, Pbr, QFM, data.AxialForce, data.BendingMoment)
		st := moment_proto.StrengthResult(strength1)
		result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, &st)

		if len(data.FlangesData) > 1 {
			strength2 := s.getCalculatedStrength(flange2, res2, type2, gasket.M, data.Pressure, Qd, Dcp, result.Calc.Strength.SigmaB2, Pbm, Pbr, QFM, data.AxialForce, data.BendingMoment)
			st := moment_proto.StrengthResult(strength2)
			result.Calc.Strength.Strength = append(result.Calc.Strength.Strength, &st)
		}

		//TODO
	}

	return &result, nil
}

// Функция для получения данных необходимых для расчетов
func (s *FlangeService) getDataFlange(
	ctx context.Context,
	flange *moment_proto.FlangeData,
	typeFlange string,
	temp float64,
) (models.InitialDataFlange, error) {
	size, err := s.repo.GetSize(ctx, float64(flange.Dy), flange.Py, flange.StandartId)
	if err != nil {
		return models.InitialDataFlange{}, fmt.Errorf("failed to get size. error: %w", err)
	}

	dataFlange := models.InitialDataFlange{
		DOut:     size.D1,
		D:        size.D,
		H:        size.B,
		S0:       size.S0,
		S1:       size.S1,
		L:        size.Lenght,
		D6:       size.D2,
		Diameter: size.Diameter,
		Count:    size.Count,
		Area:     size.Area,
		C:        flange.Corrosion,
	}

	dataFlange.Tf = s.typeFlangesTF[typeFlange] * temp

	if flange.Type == moment_proto.FlangeData_free {
		dataFlange.Tk = s.typeFlangesTK[typeFlange] * temp

		//TODO тут неправильная марка указана
		mat, err := s.materials.GetMatFotCalculate(ctx, flange.MarkId, dataFlange.Tk)
		if err != nil {
			return models.InitialDataFlange{}, err
		}

		dataFlange.AlphaK = mat.AlphaF
		dataFlange.EpsilonKAt20 = mat.EpsilonAt20
		dataFlange.EpsilonK = mat.Epsilon
		dataFlange.SigmaKAt20 = mat.SigmaAt20
		dataFlange.SigmaK = mat.Sigma
	}

	mat, err := s.materials.GetMatFotCalculate(ctx, flange.MarkId, dataFlange.Tf)
	if err != nil {
		return models.InitialDataFlange{}, err
	}

	dataFlange.AlphaF = mat.AlphaF
	dataFlange.EpsilonAt20 = mat.EpsilonAt20
	dataFlange.Epsilon = mat.Epsilon
	dataFlange.SigmaAt20 = mat.SigmaAt20
	dataFlange.Sigma = mat.Sigma

	dataFlange.SigmaM = constants.SigmaM * mat.Sigma
	dataFlange.SigmaMAt20 = constants.SigmaM * mat.SigmaAt20
	dataFlange.SigmaR = constants.SigmaR * mat.Sigma
	dataFlange.SigmaRAt20 = constants.SigmaR * mat.SigmaAt20

	return dataFlange, nil
}

// расчеты
func (s *FlangeService) getCalculatedData(
	ctx context.Context,
	flange *moment_proto.FlangeData,
	data models.InitialDataFlange,
	Dcp float64,
) (models.CalculatedData, error) {
	var calculated models.CalculatedData
	if flange.Type != moment_proto.FlangeData_free {
		calculated.B = 0.5 * (data.D6 - Dcp)
	} else {
		calculated.Ds = 0.5 * (data.DOut + data.Dk + 2*data.H)
		calculated.A = 0.5 * (data.D6 - data.Ds)
		calculated.B = 0.5 * (data.Ds - Dcp)
	}

	if flange.Type != moment_proto.FlangeData_welded {
		calculated.Se = data.S0
	} else {
		calculated.X = data.L / (math.Sqrt(data.D * data.S0))
		calculated.Betta = data.S1 / data.S0
		calculated.Zak = 1 + (calculated.Betta-1)*calculated.X/(calculated.X+(1+calculated.Betta)/4)
		calculated.Se = calculated.Zak * data.S0
	}

	calculated.E = 0.5 * (Dcp - data.D - calculated.Se)
	calculated.L0 = math.Sqrt(data.D * data.S0)
	calculated.K = data.DOut / data.D

	dividend := math.Pow(calculated.K, 2)*(1+8.55*(math.Log(calculated.K)/math.Log(10))) - 1
	divider := (1.05 + 1.945*math.Pow(calculated.K, 2)) * (calculated.K - 1)
	calculated.BettaT = dividend / divider

	divider = 1.36 * (math.Pow(calculated.K, 2) - 1) * (calculated.K - 1)
	calculated.BettaU = dividend / divider

	dividend = 1 / (calculated.K - 1)
	divider = 0.69 + 5.72*((math.Pow(calculated.K, 2)*(math.Log(calculated.K)/math.Log(10)))/(math.Pow(calculated.K, 2)-1))
	calculated.BettaY = dividend * divider

	dividend = math.Pow(calculated.K, 2) + 1
	divider = math.Pow(calculated.K, 2) - 1
	calculated.BettaZ = dividend / divider

	if flange.Type == moment_proto.FlangeData_welded && data.S0 != data.S1 {
		betta := data.S1 / data.S0
		x := data.L / calculated.L0

		calculated.BettaF = s.graphic.CalculateBettaF(betta, x)
		calculated.BettaV = s.graphic.CalculateBettaV(betta, x)
		calculated.F = s.graphic.CalculateF(betta, x)
	} else {
		calculated.BettaF = constants.InitBettaF
		calculated.BettaV = constants.InitBettaV
		calculated.F = constants.InitF
	}

	calculated.Lymda = (calculated.BettaF*data.H+calculated.L0)/(calculated.BettaT*calculated.L0) + (calculated.BettaV*math.Pow(data.H, 3))/(calculated.BettaU*calculated.L0*math.Pow(data.S0, 2))
	calculated.Yf = (0.91 * calculated.BettaV) / (data.EpsilonAt20 * calculated.Lymda * math.Pow(data.S0, 2) * calculated.L0)

	if flange.Type == moment_proto.FlangeData_free {
		calculated.Psik = 1.28 * (math.Log(data.Dnk/data.Dk) / math.Log(10))
		calculated.Yk = 1 / (data.EpsilonKAt20 * math.Pow(data.Hk, 3) * calculated.Psik)
	}

	if flange.Type != moment_proto.FlangeData_free {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
	} else {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.Ds / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
		calculated.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonKAt20 * data.Dnk * math.Pow(data.Hk, 3)))
	}

	return calculated, nil
}

// расчеты если выполняется прочностной расчет
func (s *FlangeService) getCalculatedStrength(
	flange models.InitialDataFlange,
	res models.CalculatedData,
	type1 moment_proto.FlangeData_Type,
	M, Pressure, Qd, Dcp, SigmaB, Pbm, Pbr, QFM float64,
	AxialForce, BendingMoment int32,
) models.CalculatedStrength {
	//* большинтсво переменный называются +- так же как и в оригинале
	var strength models.CalculatedStrength

	if SigmaB > constants.MaxSigmaB && flange.Diameter >= constants.MinDiameter && flange.Diameter <= constants.MaxDiameter {
		strength.Mkp = s.graphic.CalculateMkp(flange.Diameter, SigmaB)
	} else {
		strength.Mkp = (0.3 * Pbm * float64(flange.Diameter) / float64(flange.Count)) / 1000.0
	}

	strength.Mkp1 = 0.75 * strength.Mkp

	temp1 := math.Pi * flange.D6 / float64(flange.Count)
	temp2 := 2*float64(flange.Diameter) + 6*flange.H/(M+0.5)

	Cf1 := math.Max(1, math.Sqrt(temp1/temp2))
	strength.Cf = Cf1
	var Dzv1 float64

	if type1 == moment_proto.FlangeData_welded && flange.D <= 20*flange.S1 {
		// if flange1.D > 20 * flange1.S1 {
		// 	Dzv1 = flange1.D
		// } else {
		if res.F > 1 {
			Dzv1 = flange.D + flange.S0
		} else {
			Dzv1 = flange.D + flange.S1
		}
	} else {
		Dzv1 = flange.D
	}
	strength.Dzv = Dzv1

	MM1 := Cf1 * Pbm * res.B
	Mp1 := Cf1 * math.Max(Pbr*res.B+(Qd+QFM)*res.E, math.Abs(Qd+QFM)*res.E)
	var MMk1, Mpk, sigmaM1, sigmaM0 float64
	if type1 == moment_proto.FlangeData_free {
		MMk1 = Cf1 * Pbm * res.A
		Mpk = Cf1 * Pbr * res.A
	}

	if type1 == moment_proto.FlangeData_welded && flange.S1 != flange.S0 {
		sigmaM1 = MM1 / (res.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv1)
		sigmaM0 = res.F * sigmaM1
	} else {
		sigmaM1 = MM1 / (res.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv1)
		sigmaM0 = sigmaM1
	}

	sigmaR := ((1.33*res.BettaF*flange.H + res.L0) / (res.Lymda * math.Pow(flange.H, 2) * res.L0 * flange.D)) * MM1
	sigmaT := res.BettaY*MM1/(math.Pow(flange.H, 2)*flange.D) - res.BettaZ*sigmaR

	strength.SigmaR = sigmaR
	strength.SigmaT = sigmaT

	var sigmaK, sigmaP1, sigmaP0, sigmaMp, sigmaMpm float64
	if type1 == moment_proto.FlangeData_free {
		sigmaK = res.BettaY * MMk1 / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	if type1 == moment_proto.FlangeData_welded && flange.S1 != flange.S0 {
		sigmaP1 = Mp1 / (res.Lymda * math.Pow(flange.S1-flange.C, 2) * Dzv1)
		sigmaP0 = res.F * sigmaP1
	} else {
		sigmaP1 = Mp1 / (res.Lymda * math.Pow(flange.S0-flange.C, 2) * Dzv1)
		sigmaP0 = sigmaP1
	}

	if type1 == moment_proto.FlangeData_welded {
		temp := math.Pi * (flange.D + flange.S1) * (flange.S1 - flange.C)
		sigmaMp = (Qd + float64(AxialForce) + 4*math.Abs(float64(BendingMoment)/Dcp)) / temp
		sigmaMpm = (Qd + float64(AxialForce) - 4*math.Abs(float64(BendingMoment)/Dcp)) / temp
	}

	temp := math.Pi * (flange.D + flange.S0) * (flange.S0 - flange.C)
	sigmaMp0 := (Qd + float64(AxialForce) + 4*math.Abs(float64(BendingMoment)/Dcp)) / temp
	sigmaMpm0 := (Qd + float64(AxialForce) - 4*math.Abs(float64(BendingMoment)/Dcp)) / temp
	sigmaMop := Pressure * flange.D / (2.0 * (flange.S0 - flange.C))

	sigmaRp := ((1.33*res.BettaF*flange.H + res.L0) / (res.Lymda * math.Pow(flange.H, 2) * res.L0 * flange.D)) * Mp1
	sigmaTp := res.BettaY*Mp1/(math.Pow(flange.H, 2)*flange.D) - res.BettaZ*sigmaRp

	var sigmaKp float64
	if type1 == moment_proto.FlangeData_free {
		sigmaKp = res.BettaY * Mp1 / (math.Pow(flange.Hk, 2) * flange.Dk)
	}

	if type1 == moment_proto.FlangeData_welded {
		if flange.D <= constants.MinD {
			strength.DTeta = constants.MinDTetta
		} else if flange.D > constants.MaxD {
			strength.DTeta = constants.MaxDTetta
		} else {
			strength.DTeta = ((flange.D-constants.MinD)/(constants.MaxD-constants.MinD))*(constants.MaxDTetta-constants.MinDTetta) + constants.MinDTetta
		}
	} else {
		strength.DTeta = constants.MaxDTetta
	}

	strength.Teta = Mp1 * res.Yf * flange.EpsilonAt20 / flange.Epsilon

	if type1 == moment_proto.FlangeData_free {
		strength.DTetaK = 0.002
		strength.TetaK = Mpk * res.Yk * flange.EpsilonKAt20 / flange.EpsilonK
	}

	var max1, max2, max3, max4, max5, max6, max7, max8, max9, max10, max11 float64
	if type1 == moment_proto.FlangeData_welded && flange.S1 != flange.S0 {
		max1 = math.Max(math.Abs(sigmaM1+sigmaR), math.Abs(sigmaM1+sigmaT))

		t1 := math.Max(math.Abs(sigmaP1-sigmaMp+sigmaRp), math.Abs(sigmaP1-sigmaMpm+sigmaRp))
		t2 := math.Max(math.Abs(sigmaP1-sigmaMp+sigmaTp), math.Abs(sigmaP1-sigmaMpm+sigmaTp))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(sigmaP1+sigmaMp), math.Abs(sigmaP1+sigmaMpm))

		max2 = math.Max(t1, t2)
		max3 = sigmaM0

		t1 = math.Max(math.Abs(sigmaP0+sigmaMp0), math.Abs(sigmaP0-sigmaMp0))
		t2 = math.Max(math.Abs(sigmaP0+sigmaMpm0), math.Abs(sigmaP0-sigmaMpm0))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.3*sigmaP0+sigmaMop), math.Abs(0.3*sigmaP0-sigmaMop))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.7*sigmaP0+(sigmaMp0-sigmaMop)), math.Abs(0.7*sigmaP0-(sigmaMp0-sigmaMop)))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(0.7*sigmaP0+(sigmaMpm0-sigmaMop)), math.Abs(0.7*sigmaP0-(sigmaMpm0-sigmaMop)))

		max4 = math.Max(t1, t2)
	} else {
		max5 = math.Max(math.Abs(sigmaM0+sigmaR), math.Abs(sigmaM0+sigmaT))

		t1 := math.Max(math.Abs(sigmaP0-sigmaMp0+sigmaTp), math.Abs(sigmaP0-sigmaMpm0+sigmaTp))
		t2 := math.Max(math.Abs(sigmaP0-sigmaMp0+sigmaRp), math.Abs(sigmaP0-sigmaMpm0+sigmaRp))
		t1 = math.Max(t1, t2)
		t2 = math.Max(math.Abs(sigmaP0+sigmaMp0), math.Abs(sigmaP0+sigmaMpm0))

		max6 = math.Max(t1, t2)
	}

	max7 = math.Max(math.Abs(sigmaMp0), math.Abs(sigmaMpm0))
	max7 = math.Max(max7, math.Abs(sigmaMop))
	max8 = math.Max(math.Abs(sigmaR), math.Abs(sigmaT))
	max9 = math.Max(math.Abs(sigmaRp), math.Abs(sigmaTp))

	if type1 == moment_proto.FlangeData_free {
		max10 = sigmaK
		max11 = sigmaKp
	}

	strength.MM = MM1
	strength.Mp = Mp1
	strength.MMk = MMk1
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

	strength.Max1 = max1
	strength.Max2 = max2
	strength.Max3 = max3
	strength.Max4 = max4
	strength.Max5 = max5
	strength.Max6 = max6
	strength.Max7 = max7
	strength.Max8 = max8
	strength.Max9 = max9
	strength.Max10 = max10
	strength.Max11 = max11

	return strength
}
