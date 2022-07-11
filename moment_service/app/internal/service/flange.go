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
	typeFlangesTF map[string]float64
	typeFlangesTD map[string]float64
	typeBolt      map[string]float64
}

func NewFlangeService(repo repository.Flange, materials *MaterialsService, gasket *GasketService) *FlangeService {
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
		typeFlangesTF: flangesTF,
		typeFlangesTD: flangesTD,
		typeBolt:      bolt,
	}
}

// TODO в зависимости от госта можно будет вызывать отдельные функции
func (s *FlangeService) Calculation(ctx context.Context, data *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error) {
	dataFlangeFirst, err := s.getDataFlange(ctx, data.FlangesData[0], data.Flanges, data.Temp)
	if err != nil {
		return nil, err
	}
	var dataFlangeSecond models.InitialDataFlange
	if len(data.FlangesData) > 1 {
		dataFlangeSecond, err = s.getDataFlange(ctx, data.FlangesData[1], data.Flanges, data.Temp)
		if err != nil {
			return nil, err
		}
	} else {
		dataFlangeSecond = dataFlangeFirst
	}

	//TODO добавить зависимость от типа фланца
	Tb := s.typeFlangesTD[data.Flanges] * float64(data.Temp)
	// if ($TipF1 == 2) {

	//TODO
	boltMat, err := s.materials.GetMatFotCalculate(ctx, data.Bolt.MarkId, Tb)
	if err != nil {
		return nil, err
	}
	logger.Debug(boltMat)

	gasket, err := s.gasket.Get(ctx, models.GetGasket{TypeGasket: "", Env: "", Thickness: 3.2})
	if err != nil {
		return nil, err
	}
	var Lb0 float64
	// TODO здесь какое-то условие надо
	Lb0 = float64(gasket.Thickness)
	Lb0 += dataFlangeFirst.H + dataFlangeSecond.H

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

	// yp = 0.0;
	// if ($TipP == 0) {
	// 	$yp = $hp * $Kp / ($Ep * pi() * $Dcp * $bp);
	// }
	var yp float64 = 0
	//? иногда она не считается
	yp = float64(gasket.Thickness*gasket.Compression) / (float64(gasket.Epsilon) * math.Pi * Dcp * bp)

	Lb := Lb0 + s.typeBolt[data.Type]*float64(dataFlangeFirst.Diameter)

	yb := Lb / (dataFlangeFirst.EpsilonAt20 * float64(dataFlangeFirst.Area) * float64(dataFlangeFirst.Count))
	Ab := float64(dataFlangeFirst.Count) * float64(dataFlangeFirst.Area)

	logger.Debug(yp, yb, Ab, b0)

	return &moment_proto.FlangeResponse{}, nil
}

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

func (s *FlangeService) getCalculatedDate(
	ctx context.Context,
	flange *moment_proto.FlangeData,
	data models.InitialDataFlange,
	Dcp float64,
) (models.CalculatedData, error) {
	var calculated models.CalculatedData
	if flange.Type != "free" {
		calculated.B = 0.5 * (data.D6 - Dcp)
	} else {
		calculated.Ds = 0.5 * (data.D + data.Dk + 2*data.H)
		calculated.A = 0.5 * (data.D6 - data.Ds)
		calculated.B = 0.5 * (data.Ds - Dcp)
	}

	if flange.Type != "welded" {
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

	if flange.Type == "welded" && data.S0 != data.S1 {
		//TODO
	} else {
		calculated.BettaF = 0.91
		calculated.BettaV = 0.55
		calculated.F = 1.0
	}

	return calculated, nil
}
