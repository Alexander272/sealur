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
	typeFlangesTD map[string]float64
	typeBolt      map[string]float64
}

func NewFlangeService(repo repository.Flange, materials *MaterialsService, gasket *GasketService, graphic *GraphicService) *FlangeService {
	flangesTF := map[string]float64{
		"isolated":    constants.IsolatedFlatTf,
		"nonIsolated": constants.NonIsolatedFlatTf,
	}

	flangesTD := map[string]float64{
		"isolated":         constants.IsolatedFlatTb,
		"nonIsolated":      constants.NonIsolatedFlatTb,
		"isolated-free":    constants.IsolatedFreeTb,
		"nonIsolated-free": constants.NonIsolatedFlatTb,
	}

	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	return &FlangeService{
		repo:          repo,
		materials:     materials,
		gasket:        gasket,
		graphic:       graphic,
		typeFlangesTF: flangesTF,
		typeFlangesTD: flangesTD,
		typeBolt:      bolt,
	}
}

//? можно расчет по основным формулам вынести в отдельный пакет, а потом просто использовать (должно сделать код более понятным)

// TODO в зависимости от госта можно будет вызывать отдельные функции (возможно придется делать все в одной функции)
func (s *FlangeService) Calculation(ctx context.Context, data *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error) {
	dataFlangeFirst, err := s.getDataFlange(ctx, data.FlangesData[0], data.Flanges.String(), data.Temp)
	if err != nil {
		return nil, err
	}
	var dataFlangeSecond models.InitialDataFlange
	if len(data.FlangesData) > 1 {
		dataFlangeSecond, err = s.getDataFlange(ctx, data.FlangesData[1], data.Flanges.String(), data.Temp)
		if err != nil {
			return nil, err
		}
	} else {
		dataFlangeSecond = dataFlangeFirst
	}

	//TODO добавить зависимость от типа фланца
	Tb := s.typeFlangesTD[data.Flanges.String()] * float64(data.Temp)
	// if ($TipF1 == 2) {

	//TODO
	boltMat, err := s.materials.GetMatFotCalculate(ctx, data.Bolts.MarkId, Tb)
	if err != nil {
		return nil, err
	}
	logger.Debug(boltMat)

	gasket, err := s.gasket.Get(ctx, models.GetGasket{TypeGasket: "", Env: "", Thickness: 3.2})
	if err != nil {
		return nil, err
	}
	var Lb0 float64
	// TODO здесь условие надо if ($TipF1 == 2)
	Lb0 = float64(gasket.Thickness)
	Lb0 += dataFlangeFirst.H + dataFlangeSecond.H
	//TODO доп действия при выполнении условия

	logger.Debug(gasket)

	bp := float64(data.Gasket.DOut-data.Gasket.DIn) / 2

	var b0, Dcp float64

	if gasket.IsOval {
		b0 = bp / 4
		Dcp = float64(data.Gasket.DOut) - bp/2
	} else {
		if bp <= constants.Bp {
			b0 = bp
		} else {
			b0 = constants.B0 * math.Sqrt(bp)
		}
		Dcp = float64(data.Gasket.DOut) - b0
	}

	// yp = 0.0
	// if ($TipP == 0) {
	// 	$yp = $hp * $Kp / ($Ep * pi() * $Dcp * $bp)
	// }
	var yp float64 = 0
	//? иногда она не считается
	yp = float64(gasket.Thickness*gasket.Compression) / (float64(gasket.Epsilon) * math.Pi * Dcp * bp)

	Lb := Lb0 + s.typeBolt[data.Type.String()]*float64(dataFlangeFirst.Diameter)

	yb := Lb / (boltMat.EpsilonAt20 * float64(dataFlangeFirst.Area) * float64(dataFlangeFirst.Count))
	Ab := float64(dataFlangeFirst.Count) * float64(dataFlangeFirst.Area)

	logger.Debug(yp, yb, Ab, b0)

	resFirst, err := s.getCalculatedData(ctx, data.FlangesData[0], dataFlangeFirst, Dcp)
	if err != nil {
		return nil, err
	}

	var resSecond models.CalculatedData
	if len(data.FlangesData) > 1 {
		resSecond, err = s.getCalculatedData(ctx, data.FlangesData[1], dataFlangeSecond, Dcp)
		if err != nil {
			return nil, err
		}
	} else {
		resSecond = resFirst
	}

	logger.Debug(resFirst, resSecond)

	var alpha, dividend, divider float64

	divider = yp + yb*boltMat.EpsilonAt20/boltMat.Epsilon + (resFirst.Yf*dataFlangeFirst.EpsilonAt20/dataFlangeFirst.Epsilon)*math.Pow(resFirst.B, 2) + (resSecond.Yf*dataFlangeSecond.EpsilonAt20/dataFlangeSecond.Epsilon)*math.Pow(resSecond.B, 2)
	//TODO
	// if ($TipF1 == 2) {
	// 	$prom += ($yk1 * $Ek201 / $Ek1) * $a1 * $a1;
	// }
	// if ($TipF2 == 2) {
	// 	$prom += ($yk2 * $Ek202 / $Ek2) * $a2 * $a2;
	// }

	gamma := 1 / divider
	// if ($TipP == 2  or  $TipF1 == 2  or  $TipF2 == 2) {
	// 	$alfa = 1.0;
	// } else {
	alpha = 1 - (yp-(resFirst.Yf*resFirst.E*resFirst.B+resSecond.Yf*resSecond.E*resSecond.B))/(yp+yb+(resFirst.Yf*math.Pow(resFirst.B, 2)+resSecond.Yf*math.Pow(resSecond.B, 2)))
	// }

	dividend = yb + resFirst.Yfn*resFirst.B*(resFirst.B+resFirst.E-math.Pow(resFirst.E, 2)/Dcp) + resSecond.Yfn*resSecond.B*(resSecond.B+resSecond.E-math.Pow(resSecond.E, 2)/Dcp)
	divider = yb + yp*math.Pow(dataFlangeFirst.D6/Dcp, 2) + resFirst.Yfn*math.Pow(resFirst.B, 2) + resSecond.Yfn*math.Pow(resSecond.B, 2)
	/*
		if ($TipF1 == 2) {
			$dividend += $yfc1 * $a1 * $a1;
			$divider += $yfc1 * $a1 * $a1;
		}
		if ($TipF2 == 2) {
			$dividend += $yfc2 * $a2 * $a2;
			$divider += $yfc2 * $a2 * $a2;
		}
	*/
	alphaM := dividend / divider

	Pobg := 0.5 * math.Pi * Dcp * b0 * float64(gasket.SpecificPres)

	var Rp float64 = 0
	if data.Pressure >= 0 {
		Rp = math.Pi * Dcp * b0 * float64(gasket.M) * float64(data.Pressure)
	}

	Qd := 0.785 * math.Pow(Dcp, 2) * float64(data.Pressure)

	temp1 := float64(data.AxialForce) + 4*math.Abs(float64(gasket.M))/Dcp
	temp2 := float64(data.AxialForce) - 4*math.Abs(float64(gasket.M))/Dcp

	QFM := math.Max(temp1, temp2)
	logger.Debug(QFM)

	Pb2 := math.Max(Pobg, 4*Ab*boltMat.SigmaAt20)
	Pb1 := alpha*(Qd+float64(data.AxialForce)) + Rp + 4 + alphaM*math.Abs(float64(gasket.M))/Dcp

	PbmFirst := math.Max(Pb1, Pb2)
	Pbr1 := PbmFirst + (1-alpha)*(Qd+float64(data.AxialForce)) + 4*(1-alphaM)*math.Abs(float64(gasket.M))/Dcp

	sigmaB1 := PbmFirst / Ab
	sigmaB2 := Pbr1 / Ab

	var Kyp, Kyz float64
	if data.IsWork {
		Kyp = 1
	} else {
		Kyp = 1.35
	}
	if data.Condition == moment_proto.FlangeRequest_uncontrollable {
		Kyz = 1
	} else if data.Condition == moment_proto.FlangeRequest_controllable {
		Kyz = 1.1
	} else {
		Kyz = 1.3
	}

	temp1 = 1
	d_sigmaM := 1.2 * Kyp * Kyz * temp1 * boltMat.SigmaAt20
	d_sigmaR := Kyp * Kyz * temp1 * boltMat.Sigma
	// if ($TipP == 0) {
	// 	$qmax = max($Pbm_first, $Pbr1) / (pi() * $Dcp * $bp);
	// }

	logger.Debug(d_sigmaM, d_sigmaR)

	// if ($Moment != 1)

	temp1 = dataFlangeFirst.AlphaF*dataFlangeFirst.H*(dataFlangeFirst.Tf-20) + dataFlangeSecond.AlphaF*dataFlangeSecond.H*(dataFlangeSecond.Tf-20)
	temp2 = dataFlangeFirst.H + dataFlangeSecond.H
	/*
		if ($TipF1 == 2) {
			$prom3 += $alfak1 * $hk1 * ($Tk1 - 20.0);
			$prom4 += $hk1;
		}
		if ($TipF2 == 2) {
			$prom3 += $alfak2 * $hk2 * ($Tk2 - 20.0);
			$prom4 += $hk2;
		}
		if ($ZakDet == 1) {
			$prom3 += $alfaz * $hz * ($T - 20.0);
			$prom4 += $hz;
		}
	*/

	//TODO здесь должна быть новая формула (Qt)
	Qt := gamma * (temp1 - boltMat.AlphaF*temp2*(Tb-20))

	Pb1_117 := math.Max(Pb1, Pb1-Qt)
	Pbm := math.Max(Pb1_117, Pb2)
	Pbr := Pbm + (1-alpha)*(Qd+float64(data.AxialForce)) + Qt + 4*(1-alphaM)*math.Abs(float64(gasket.M))/Dcp

	sigmaB1 = Pbm / Ab
	sigmaB2 = Pbr / Ab

	Kyt := 1.3
	DsigmaM := 1.2 * Kyp * Kyz * Kyt * boltMat.SigmaAt20
	DsigmaR := Kyp * Kyz * Kyt * boltMat.Sigma
	// if ($TipP == 0) {
	// 	$qmax = max($Pbm, $Pbr) / (pi() * $Dcp * $bp);
	// }
	var qmax float64

	var v_sigmab1, v_sigmab2 bool
	if sigmaB1 <= DsigmaM {
		v_sigmab1 = true
	}
	if sigmaB2 <= DsigmaR {
		v_sigmab2 = true
	}

	// if ($Moment == 1) {

	// if (($v_sigmab1 == 0 and $TipP != 0 and $v_sigmab2 == 0)  or  ($v_sigmab1 == 0 and $TipP == 0 and $v_qmax == 0 and $v_sigmab2 == 0)) {
	if (v_sigmab1 && v_sigmab2) || (v_sigmab1 && v_sigmab2 && qmax <= float64(gasket.PermissiblePres)) {
		var Mkp float64
		if sigmaB1 > 120.0 && dataFlangeFirst.Diameter >= 20 && dataFlangeFirst.Diameter <= 52 {
			Mkp = s.graphic.CalculateMkp(dataFlangeFirst.Diameter, sigmaB1)
		} else {
			//TODO вроде как формула изменилась
			// зачем-то делится на 1000
			Mkp = (0.3 * Pbm * float64(dataFlangeFirst.Diameter) / float64(dataFlangeFirst.Count)) / 1000
			Mkp1 := 0.75 * Mkp

			Prek := 0.8 * Ab * boltMat.SigmaAt20
			Qrek := Prek / (math.Pi * Dcp * bp)
			Mrek := (0.3 * Prek * float64(dataFlangeFirst.Diameter) / float64(dataFlangeFirst.Count)) / 1000

			Pmax := DsigmaM * Ab
			Qmax := Pmax / (math.Pi * Dcp * bp)
			// if ($TipP == 0) {
			// 	if (Qmax > float64(gasket.PermissiblePres)) {
			// 		Pmax = float64(gasket.PermissiblePres) * (math.Pi * Dcp * bp);
			// 		Qmax = float64(gasket.PermissiblePres);
			// 	}
			// }

			Mmax := (0.3 * Pmax * float64(dataFlangeFirst.Diameter) / float64(dataFlangeFirst.Count)) / 1000

			logger.Debug(Mmax, Qmax, Mrek, Qrek, Mkp1)
		}
	}

	// else {

	return &moment_proto.FlangeResponse{}, nil
}

// Функция для получения данных необходимых для расчетов
func (s *FlangeService) getDataFlange(
	ctx context.Context,
	flange *moment_proto.FlangeData,
	typeFlange string,
	temp float32,
) (models.InitialDataFlange, error) {
	size, err := s.repo.GetSize(ctx, float64(flange.Dy), float64(flange.Py))
	if err != nil {
		return models.InitialDataFlange{}, fmt.Errorf("failed to get size. error: %w", err)
	}

	var dataFlange models.InitialDataFlange

	//TODO добавить зависимость от типа фланца (можно добавить мапу из которой в зависимости от фланца будет выбираться нужная константа)
	dataFlange.Tf = s.typeFlangesTF[typeFlange] * float64(temp)
	// TODO тут еще Tk1 считается если фланцы свободные

	mat, err := s.materials.GetMatFotCalculate(ctx, flange.MarkId, dataFlange.Tf)
	if err != nil {
		return models.InitialDataFlange{}, err
	}

	dataFlange = models.InitialDataFlange{
		DOut:        size.D1,
		D:           size.D,
		H:           size.B,
		S0:          size.S0,
		S1:          size.S1,
		L:           size.Lenght,
		D6:          size.D2,
		Diameter:    size.Diameter,
		Count:       size.Count,
		Area:        size.Area,
		AlphaF:      mat.AlphaF,
		EpsilonAt20: mat.EpsilonAt20,
		Epsilon:     mat.Epsilon,
		SigmaAt20:   mat.SigmaAt20,
		Sigma:       mat.Sigma,
	}

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
		calculated.Ds = 0.5 * (data.D + data.Dk + 2*data.H)
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

	calculated.E = 0.5 * (Dcp - data.DOut - calculated.Se)
	calculated.L0 = math.Sqrt(data.DOut * data.S0)
	calculated.K = data.D / data.DOut

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
		calculated.BettaF = 0.908920
		calculated.BettaV = 0.550103
		calculated.F = 1.0
	}

	calculated.Lymda = (calculated.BettaF+data.H+calculated.L0)/calculated.BettaT*data.L + calculated.BettaV*math.Pow(data.H, 3)/(calculated.BettaU*calculated.L0*math.Pow(data.S0, 2))
	calculated.Yf = (0.91 * calculated.BettaV) / (data.EpsilonAt20 * calculated.Lymda * math.Pow(data.S0, 2) * calculated.L0)

	if flange.Type == moment_proto.FlangeData_free {
		calculated.Psik = 1.28 * (math.Log(data.Dnk/data.Dk) / math.Log(10))
		calculated.Yk = 1 / (data.EpsilonKAt20 * math.Pow(data.Hk, 3) * calculated.Psik)
	}

	if flange.Type != moment_proto.FlangeData_free {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.D * math.Pow(data.H, 3)))
	} else {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.Ds / (data.EpsilonAt20 * data.D * math.Pow(data.H, 3)))
		calculated.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonKAt20 * data.Dnk * math.Pow(data.Hk, 3)))
	}

	return calculated, nil
}
