package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
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
		Calc:         &moment_proto.CalculatedFlange{},
		Gasket:       &moment_proto.GasketResult{},
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
		C:            float64(data.FlangesData[0].Corrosion),
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
			C:            float64(data.FlangesData[1].Corrosion),
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
	} else {
		if bp <= constants.Bp {
			b0 = bp
		} else {
			b0 = constants.B0 * math.Sqrt(bp)
		}
		Dcp = data.Gasket.DOut - b0
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
	}
	result.Calc.Alpha = alpha

	dividend = yb + res1.Yfn*res1.B*(res1.B+res1.E-math.Pow(res1.E, 2)/Dcp) + res2.Yfn*res2.B*(res2.B+res2.E-math.Pow(res2.E, 2)/Dcp)
	divider = yb + yp*math.Pow(flange1.D6/Dcp, 2) + res1.Yfn*math.Pow(res1.B, 2) + res2.Yfn*math.Pow(res2.B, 2)

	if type1 == moment_proto.FlangeData_free {
		dividend += res1.Yfc * math.Pow(res1.A, 2)
		divider += res1.Yfc * math.Pow(res1.A, 2)
	}
	if type2 == moment_proto.FlangeData_free {
		dividend += res2.Yfc * math.Pow(res2.A, 2)
		divider += res2.Yfc * math.Pow(res2.A, 2)
	}

	alphaM := dividend / divider
	result.Calc.AlphaM = alphaM

	Pobg := 0.5 * math.Pi * Dcp * b0 * float64(gasket.SpecificPres)

	var Rp float64 = 0
	if data.Pressure >= 0 {
		//TODO формула изменилась
		Rp = math.Pi * Dcp * b0 * float64(gasket.M) * float64(data.Pressure)
	}

	Qd := 0.785 * math.Pow(Dcp, 2) * float64(data.Pressure)

	temp1 := float64(data.AxialForce) + 4*math.Abs(float64(data.BendingMoment))/Dcp
	temp2 := float64(data.AxialForce) - 4*math.Abs(float64(data.BendingMoment))/Dcp

	QFM := math.Max(temp1, temp2)

	result.Calc.Po = Pobg
	result.Calc.Rp = Rp
	result.Calc.Qd = Qd
	result.Calc.Qfm = QFM

	//* Похоже эти значения используются только при if ($Moment != 1) в остольных случаях они переписываются
	// хз почему это считается здесь
	//? может надо вынести это в одельную функцию?
	// Pb2 := math.Max(Pobg, 0.4*Ab*boltMat.SigmaAt20)
	// Pb1 := alpha*(Qd+float64(data.AxialForce)) + Rp + 4*alphaM*math.Abs(float64(data.BendingMoment))/Dcp

	// PbmFirst := math.Max(Pb1, Pb2)
	// Pbr1 := PbmFirst + (1-alpha)*(Qd+float64(data.AxialForce)) + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/Dcp
	// result.Calc.Pb = PbmFirst

	// sigmaB1 := PbmFirst / Ab
	// sigmaB2 := Pbr1 / Ab

	// var Kyp, Kyz float64
	// if data.IsWork {
	// 	Kyp = 1
	// } else {
	// 	Kyp = 1.35
	// }
	// if data.Condition == moment_proto.FlangeRequest_uncontrollable {
	// 	Kyz = 1
	// } else if data.Condition == moment_proto.FlangeRequest_controllable {
	// 	Kyz = 1.1
	// } else {
	// 	Kyz = 1.3
	// }

	// temp1 = 1
	// d_sigmaM := 1.2 * Kyp * Kyz * temp1 * boltMat.SigmaAt20
	// d_sigmaR := Kyp * Kyz * temp1 * boltMat.Sigma
	// logger.Debug(d_sigmaM, d_sigmaR)

	// var qmax float64
	// //TODO заменить на тип из proto
	// if gasket.Type == "Soft" {
	// 	qmax = math.Max(PbmFirst, Pbr1) / (math.Pi * Dcp * bp)
	// }
	//* end

	// if ($Moment != 1)
	if data.Calculation != moment_proto.FlangeRequest_basis {
		logger.Debug("here  data.Calculation != moment_proto.FlangeRequest_basis")
	} else {
		calc, err := s.getCalculatedBasis(flange1, flange2, type1, type2, data, detMat, boltMat, gasket, gamma, alpha, alphaM, Ab, Pobg, Tb, Qd, Rp, Dcp, bp)
		if err != nil {
			return nil, err
		}
		result.Calc.Qt = calc.Qt
		result.Calc.Pb1 = calc.Pb1
		result.Calc.Pb2 = calc.Pb2
		result.Calc.Pb = calc.Pbm
		result.Calc.Pbr = calc.Pbr
		result.Calc.SigmaB1 = calc.SigmaB1
		result.Calc.SigmaB2 = calc.SigmaB2
		result.Calc.DSigmaM = calc.DSigmaM
		result.Calc.DSigmaR = calc.DSigmaR
		result.Calc.Q = calc.Q
		result.Calc.Mkp = calc.Mkp
		result.Calc.Mkp1 = calc.Mkp1
		result.Calc.Mrek = calc.Mrek
		result.Calc.Qrek = calc.Qrek
		result.Calc.Mmax = calc.Mmax
		result.Calc.Qmax = calc.Qmax
	}

	// temp1 = flange1.AlphaF*flange1.H*(flange1.Tf-20) + flange2.AlphaF*flange2.H*(flange2.Tf-20)
	// temp2 = flange1.H + flange2.H

	// if type1 == moment_proto.FlangeData_free {
	// 	temp1 += flange1.AlphaK * flange1.Hk * (flange1.Tk - 20)
	// 	temp2 += flange1.Hk
	// }
	// if type2 == moment_proto.FlangeData_free {
	// 	temp1 += flange2.AlphaK * flange2.Hk * (flange2.Tk - 20)
	// 	temp2 += flange2.Hk
	// }
	// if data.IsEmbedded {
	// 	temp1 += detMat.AlphaF * data.Embed.Thickness * (data.Temp - 20)
	// 	temp2 += data.Embed.Thickness
	// }

	// //TODO здесь должна быть новая формула (Qt)
	// Qt := gamma * (temp1 - boltMat.AlphaF*temp2*(Tb-20))
	// result.Calc.Qt = Qt

	// Pb1_117 := math.Max(Pb1, Pb1-Qt)
	// Pbm := math.Max(Pb1_117, Pb2)
	// Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + Qt + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/Dcp

	// result.Calc.Pb1 = Pb1_117
	// result.Calc.Pb2 = Pb2
	// result.Calc.Pbr = Pbr
	// result.Calc.SigmaB1 = Pbm / Ab
	// result.Calc.SigmaB2 = Pbr / Ab

	// Kyt := 1.3
	// // формула Г.3
	// result.Calc.DSigmaM = 1.2 * Kyp * Kyz * Kyt * boltMat.SigmaAt20
	// // формула Г.4
	// result.Calc.DSigmaR = Kyp * Kyz * Kyt * boltMat.Sigma

	// //TODO заменить на тип из proto
	// if gasket.Type == "Soft" {
	// 	qmax = math.Max(Pbm, Pbr) / (math.Pi * Dcp * bp)
	// }

	// var v_sigmab1, v_sigmab2 bool
	// if sigmaB1 <= result.Calc.DSigmaM {
	// 	v_sigmab1 = true
	// }
	// if sigmaB2 <= result.Calc.DSigmaR {
	// 	v_sigmab2 = true
	// }

	// // if ($Moment == 1) {

	// //TODO заменить на тип из proto
	// if (v_sigmab1 && v_sigmab2 && gasket.Type != "Soft") || (v_sigmab1 && v_sigmab2 && qmax <= float64(gasket.PermissiblePres) && gasket.Type == "Soft") {
	// 	// var Mkp float64
	// 	if sigmaB1 > constants.MaxSigmaB && flange1.Diameter >= constants.MinDiameter && flange1.Diameter <= constants.MaxDiameter {
	// 		result.Calc.Mkp = s.graphic.CalculateMkp(flange1.Diameter, sigmaB1)
	// 	} else {
	// 		//TODO вроде как формула изменилась
	// 		// зачем-то делится на 1000
	// 		result.Calc.Mkp = (0.3 * Pbm * float64(flange1.Diameter) / float64(flange1.Count)) / 1000
	// 		result.Calc.Mkp1 = 0.75 * result.Calc.Mkp

	// 		Prek := 0.8 * Ab * boltMat.SigmaAt20
	// 		result.Calc.Qrek = Prek / (math.Pi * Dcp * bp)
	// 		result.Calc.Mrek = (0.3 * Prek * float64(flange1.Diameter) / float64(flange1.Count)) / 1000

	// 		Pmax := result.Calc.DSigmaM * Ab
	// 		result.Calc.Qmax = Pmax / (math.Pi * Dcp * bp)

	// 		//TODO заменить на тип из proto
	// 		if gasket.Type == "Soft" && result.Calc.Qmax > float64(gasket.PermissiblePres) {
	// 			Pmax = float64(gasket.PermissiblePres) * (math.Pi * Dcp * bp)
	// 			result.Calc.Qmax = float64(gasket.PermissiblePres)
	// 		}

	// 		result.Calc.Mmax = (0.3 * Pmax * float64(flange1.Diameter) / float64(flange1.Count)) / 1000
	// 	}
	// }

	// else {

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

// расчеты если мы считаем основные величины
func (s *FlangeService) getCalculatedBasis(
	flange1, flange2 models.InitialDataFlange,
	type1, type2 moment_proto.FlangeData_Type,
	data *moment_proto.FlangeRequest,
	detMat, boltMat models.MaterialsResult,
	gasket models.Gasket,
	gamma, alpha, alphaM, Ab, Pobg, Tb, Qd, Rp, Dcp, bp float64,
) (models.CalculatedBasis, error) {
	//! как-то дофига переменных приходится передавать
	var basis models.CalculatedBasis

	temp1 := flange1.AlphaF*flange1.H*(flange1.Tf-20) + flange2.AlphaF*flange2.H*(flange2.Tf-20)
	temp2 := flange1.H + flange2.H

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
	basis.Qt = gamma * (temp1 - boltMat.AlphaF*temp2*(Tb-20))

	basis.Pb2 = math.Max(Pobg, 0.4*Ab*boltMat.SigmaAt20)
	Pb1Temp := alpha*(Qd+float64(data.AxialForce)) + Rp + 4*alphaM*math.Abs(float64(data.BendingMoment))/Dcp

	basis.Pb1 = math.Max(Pb1Temp, Pb1Temp-basis.Qt)
	basis.Pbm = math.Max(basis.Pb1, basis.Pb2)
	basis.Pbr = basis.Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + basis.Qt + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/Dcp

	basis.SigmaB1 = basis.Pbm / Ab
	basis.SigmaB2 = basis.Pbr / Ab

	Kyp := s.Kyp[data.IsWork]
	Kyz := s.Kyz[data.Condition.String()]

	// формула Г.3
	basis.DSigmaM = 1.2 * Kyp * Kyz * constants.LoadKyt * boltMat.SigmaAt20
	// формула Г.4
	basis.DSigmaR = Kyp * Kyz * constants.LoadKyt * boltMat.Sigma

	//TODO заменить на тип из proto
	if gasket.Type == "Soft" {
		basis.Q = math.Max(basis.Pbm, basis.Pbr) / (math.Pi * Dcp * bp)
	}

	var v_sigmab1, v_sigmab2 bool
	if basis.SigmaB1 <= basis.DSigmaM {
		v_sigmab1 = true
	}
	if basis.SigmaB2 <= basis.DSigmaR {
		v_sigmab2 = true
	}

	//TODO заменить на тип из proto
	if (v_sigmab1 && v_sigmab2 && gasket.Type != "Soft") || (v_sigmab1 && v_sigmab2 && basis.Q <= float64(gasket.PermissiblePres) && gasket.Type == "Soft") {
		if basis.SigmaB1 > constants.MaxSigmaB && flange1.Diameter >= constants.MinDiameter && flange1.Diameter <= constants.MaxDiameter {
			basis.Mkp = s.graphic.CalculateMkp(flange1.Diameter, basis.SigmaB1)
		} else {
			//TODO вроде как формула изменилась
			// зачем-то делится на 1000
			basis.Mkp = (0.3 * basis.Pbm * float64(flange1.Diameter) / float64(flange1.Count)) / 1000
		}

		basis.Mkp1 = 0.75 * basis.Mkp
		Prek := 0.8 * Ab * boltMat.SigmaAt20
		basis.Qrek = Prek / (math.Pi * Dcp * bp)
		basis.Mrek = (0.3 * Prek * float64(flange1.Diameter) / float64(flange1.Count)) / 1000

		Pmax := basis.DSigmaM * Ab
		basis.Qmax = Pmax / (math.Pi * Dcp * bp)

		//TODO заменить на тип из proto
		if gasket.Type == "Soft" && basis.Qmax > float64(gasket.PermissiblePres) {
			Pmax = float64(gasket.PermissiblePres) * (math.Pi * Dcp * bp)
			basis.Qmax = float64(gasket.PermissiblePres)
		}

		basis.Mmax = (0.3 * Pmax * float64(flange1.Diameter) / float64(flange1.Count)) / 1000
	}

	return basis, nil
}

func (s *FlangeService) getCalculatedStrength(
	flange1, flange2 models.InitialDataFlange,
	res1, res2 models.CalculatedData,
	type1, type2 moment_proto.FlangeData_Type,
	data *moment_proto.FlangeRequest,
	detMat, boltMat models.MaterialsResult,
	gasket models.Gasket,
	gamma, alpha, alphaM, Ab, Pobg, Tb, Qd, Rp, Dcp, bp float64,
) models.CalculatedStrength {
	var strength models.CalculatedStrength

	Pb2 := math.Max(Pobg, 0.4*Ab*boltMat.SigmaAt20)
	Pb1 := alpha*(Qd+float64(data.AxialForce)) + Rp + 4*alphaM*math.Abs(float64(data.BendingMoment))/Dcp

	Pbm := math.Max(Pb1, Pb2)
	Pbr1 := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + 4*(1-alphaM)*math.Abs(float64(data.BendingMoment))/Dcp

	SigmaB1 := Pbm / Ab
	SigmaB2 := Pbr1 / Ab

	Kyp := s.Kyp[data.IsWork]
	Kyz := s.Kyz[data.Condition.String()]

	DSigmaM := 1.2 * Kyp * Kyz * constants.NoLoadKyt * boltMat.SigmaAt20
	DSigmaR := Kyp * Kyz * constants.NoLoadKyt * boltMat.Sigma
	logger.Debug(DSigmaM, DSigmaR)

	var Q float64
	//TODO заменить на тип из proto
	if gasket.Type == "Soft" {
		Q = math.Max(Pbm, Pbr1) / (math.Pi * Dcp * bp)
	}

	var Mkp, Mkp1 float64
	if SigmaB1 > constants.MaxSigmaB && flange1.Diameter >= constants.MinDiameter && flange1.Diameter <= constants.MaxDiameter {
		Mkp = s.graphic.CalculateMkp(flange1.Diameter, SigmaB1)
	} else {
		Mkp = (0.3 * Pbm * float64(flange1.Diameter) / float64(flange1.Count)) / 1000.0
	}

	Mkp1 = 0.75 * Mkp

	temp1 := math.Pi * flange1.D6 / float64(flange1.Count)
	temp2 := 2*float64(flange1.Diameter) + 6*flange1.H/(gasket.M+0.5)

	Cf1 := math.Max(1, math.Sqrt(temp1/temp2))
	var Dzv1 float64

	if type1 == moment_proto.FlangeData_welded && flange1.D <= 20*flange1.S1 {
		// if flange1.D > 20 * flange1.S1 {
		// 	Dzv1 = flange1.D
		// } else {
		if res1.F > 1 {
			Dzv1 = flange1.D + flange1.S0
		} else {
			Dzv1 = flange1.D + flange1.S1
		}
	} else {
		Dzv1 = flange1.D
	}

	// $MM1_ = $Cf1 * $Pbm_first * $b1;
	// if ($TipF1 == 2) {
	// 	$MMk1 = $Cf1 * $Pbm_first * $a1;
	// }
	// $Mp1_ = $Cf1 * max($Pbr1 * $b1 + ($Qd + $QFM) * $e1, abs($Qd + $QFM) * $e1);
	// if ($TipF1 == 2) {
	// 	$Mpk1 = $Cf1 * $Pbr1 * $a1;
	// }
	// if ($TipF1 == 0 and $S11 != $S01) {
	// 	$sigma1m1 = $MM1_ / ($lymda1 * (pow(($S11 - $c1), 2)) * $Dzv1);
	// 	$sigma0m1 = $f1 * $sigma1m1;
	// } else {
	// 	$sigma1m1 = $MM1_ / ($lymda1 * (pow(($S01 - $c1), 2)) * $Dzv1);
	// 	$sigma0m1 = $sigma1m1;
	// }

	return strength
}
