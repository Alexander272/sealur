package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type DataService struct {
	flange        *FlangeService
	materials     *MaterialsService
	gasket        *GasketService
	graphic       *GraphicService
	typeFlangesTF map[string]float64
	typeFlangesTB map[string]float64
	typeFlangesTK map[string]float64
}

func NewDataService(flange *FlangeService, materials *MaterialsService, gasket *GasketService, graphic *GraphicService) *DataService {
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

	return &DataService{
		flange:        flange,
		materials:     materials,
		gasket:        gasket,
		typeFlangesTF: flangesTF,
		typeFlangesTB: flangesTB,
		typeFlangesTK: flangeTK,
	}
}

func (s *DataService) GetDataForFlange(ctx context.Context, data *moment_api.CalcFlangeRequest) (result models.DataFlange, err error) {
	//* формула из Таблицы В.1
	Tb := s.typeFlangesTB[data.Flanges.String()] * data.Temp
	if data.FlangesData[0].Type == moment_api.FlangeData_free {
		Tb = s.typeFlangesTB[data.Flanges.String()+"-free"] * data.Temp
	}

	flange1, boltSize, err := s.getDataFlange(ctx, data.FlangesData[0], data.Bolts, data.Flanges.String(), data.Temp)
	if err != nil {
		return result, err
	}

	result.Type1, result.Type2 = data.FlangesData[0].Type, data.FlangesData[0].Type
	flange2 := flange1

	if len(data.FlangesData) > 1 {
		flange2, _, err = s.getDataFlange(ctx, data.FlangesData[1], data.Bolts, data.Flanges.String(), data.Temp)
		if err != nil {
			return result, err
		}
		result.Type2 = data.FlangesData[1].Type
	}

	result.Bolt, err = s.getBoltData(ctx, data.Bolts, boltSize, flange1.L, Tb)
	if err != nil {
		return result, err
	}

	//? я использую температуру фланца. хз верно илил нет.
	if data.IsUseWasher {
		result.Washer1, err = s.getWasherData(ctx, data.Washer[0], flange1.Tf)
		if err != nil {
			return result, err
		}
		if !data.IsSameFlange {
			result.Washer2, err = s.getWasherData(ctx, data.Washer[1], flange2.Tf)
			if err != nil {
				return result, err
			}
		} else {
			result.Washer2 = result.Washer1
		}
	}
	if data.IsEmbedded {
		result.Embed, err = s.getEmbedData(ctx, data.Embed, data.Temp)
		if err != nil {
			return result, err
		}
	}
	bp := (data.Gasket.DOut - data.Gasket.DIn) / 2
	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket, bp)
	if err != nil {
		return result, err
	}

	if result.TypeGasket == "Oval" {
		// фомула 4
		result.B0 = bp / 4
		// фомула ?
		result.Dcp = data.Gasket.DOut - bp/2

	} else {
		if bp <= constants.Bp {
			// фомула 2
			result.B0 = bp
		} else {
			// фомула 3
			result.B0 = constants.B0 * math.Sqrt(bp)
		}
		// фомула 5
		result.Dcp = data.Gasket.DOut - result.B0
	}

	flange1, err = s.getCalculatedDataFlange(ctx, data.FlangesData[0].Type, flange1, result.Dcp)
	if err != nil {
		return result, err
	}

	if len(data.FlangesData) > 1 {
		flange2, err = s.getCalculatedDataFlange(ctx, data.FlangesData[1].Type, flange2, result.Dcp)
		if err != nil {
			return result, err
		}
	} else {
		flange2 = flange1
	}

	result.Flange1 = flange1
	result.Flange2 = flange2

	return result, nil
}

func (s *DataService) GetDataForCap(ctx context.Context, data *moment_api.CalcCapRequest) (result models.DataCap, err error) {
	//* формула из Таблицы В.1
	Tb := s.typeFlangesTB[data.Flanges.String()] * data.Temp
	if data.FlangeData.Type == moment_api.FlangeData_free {
		Tb = s.typeFlangesTB[data.Flanges.String()+"-free"] * data.Temp
	}

	flange, boltSize, err := s.getDataFlange(ctx, data.FlangeData, data.Bolts, data.Flanges.String(), data.Temp)
	if err != nil {
		return result, err
	}
	cap, err := s.getDataCap(ctx, data.CapData, data.Flanges.String(), data.Temp)
	if err != nil {
		return result, err
	}

	result.FType = data.FlangeData.Type
	result.CType = data.CapData.Type

	result.Bolt, err = s.getBoltData(ctx, data.Bolts, boltSize, flange.L, Tb)
	if err != nil {
		return result, err
	}

	//? я использую температуру фланца. хз верно илил нет.
	if data.IsUseWasher {
		result.Washer1, err = s.getWasherData(ctx, data.Washer[0], flange.Tf)
		if err != nil {
			return result, err
		}

		result.Washer2, err = s.getWasherData(ctx, data.Washer[1], cap.T)
		if err != nil {
			return result, err
		}
	}
	if data.IsEmbedded {
		result.Embed, err = s.getEmbedData(ctx, data.Embed, data.Temp)
		if err != nil {
			return result, err
		}
	}
	bp := (data.Gasket.DOut - data.Gasket.DIn) / 2
	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket, bp)
	if err != nil {
		return result, err
	}

	if result.TypeGasket == "Oval" {
		// фомула 4
		result.B0 = bp / 4
		// фомула ?
		result.Dcp = data.Gasket.DOut - bp/2

	} else {
		if bp <= constants.Bp {
			// фомула 2
			result.B0 = bp
		} else {
			// фомула 3
			result.B0 = constants.B0 * math.Sqrt(bp)
		}
		// фомула 5
		result.Dcp = data.Gasket.DOut - result.B0
	}

	flange, err = s.getCalculatedDataFlange(ctx, data.FlangeData.Type, flange, result.Dcp)
	if err != nil {
		return result, err
	}
	cap, err = s.getCalculatedDataCap(ctx, data.CapData.Type, cap, flange.H, flange.D, flange.S0, flange.DOut, result.Dcp)
	if err != nil {
		return result, err
	}

	result.Flange = flange
	result.Cap = cap

	return result, nil
}

func (s *DataService) getWasherData(ctx context.Context, data *moment_api.WasherData, temp float64) (*moment_api.WasherResult, error) {
	if data.MarkId != "another" {
		washer, err := s.materials.GetMatFotCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &moment_api.WasherResult{
			Material:  washer.Title,
			Thickness: data.Thickness,
			Alpha:     washer.AlphaF,
			Temp:      temp,
		}
		return res, nil
	}

	res := &moment_api.WasherResult{
		Material:  data.Material.Title,
		Thickness: data.Thickness,
		Alpha:     data.Material.AlphaF,
		Temp:      temp,
	}
	return res, nil
}

func (s *DataService) getEmbedData(ctx context.Context, data *moment_api.EmbedData, temp float64) (*moment_api.EmbedResult, error) {
	if data.MarkId != "another" {
		washer, err := s.materials.GetMatFotCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &moment_api.EmbedResult{
			Material:  washer.Title,
			Thickness: data.Thickness,
			Alpha:     washer.AlphaF,
			Temp:      temp,
		}
		return res, nil
	}

	res := &moment_api.EmbedResult{
		Material:  data.Material.Title,
		Thickness: data.Thickness,
		Alpha:     data.Material.AlphaF,
		Temp:      temp,
	}
	return res, nil
}

func (s *DataService) getGasketData(ctx context.Context, data *moment_api.GasketData, bp float64) (*moment_api.GasketResult, string, error) {

	if data.GasketId != "another" {
		g := models.GetGasket{GasketId: data.GasketId, EnvId: data.EnvId, Thickness: data.Thickness}
		gasket, err := s.gasket.GetFullData(ctx, g)
		if err != nil {
			return nil, "", err
		}

		res := &moment_api.GasketResult{
			Gasket:          gasket.Gasket,
			Env:             gasket.Env,
			Type:            gasket.TypeTitle,
			Thickness:       data.Thickness,
			DOut:            data.DOut,
			Width:           bp,
			M:               gasket.M,
			Pres:            gasket.SpecificPres,
			PermissiblePres: gasket.PermissiblePres,
			Compression:     gasket.Compression,
			Epsilon:         gasket.Epsilon,
		}
		return res, gasket.TypeTitle, nil
	}
	//? наверное это не лучшее решение
	titles := map[string]string{
		"Soft":  "Мягкая",
		"Oval":  "Восьмигранная",
		"Metal": "Металлическая",
	}

	res := &moment_api.GasketResult{
		Gasket:          data.Data.Title,
		Type:            data.Data.Type.String(),
		Thickness:       data.Thickness,
		DOut:            data.DOut,
		Width:           bp,
		M:               data.Data.M,
		Pres:            data.Data.Qo,
		PermissiblePres: data.Data.PermissiblePres,
		Compression:     data.Data.Compression,
		Epsilon:         data.Data.Epsilon,
	}
	return res, titles[data.Data.Type.String()], nil
}

func (s *DataService) getBoltData(ctx context.Context, data *moment_api.BoltData, bolt *moment_api.BoltResult, L, temp float64,
) (*moment_api.BoltResult, error) {
	if data.MarkId != "another" {
		mat, err := s.materials.GetMatFotCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &moment_api.BoltResult{
			Diameter:    bolt.Diameter,
			Area:        bolt.Area,
			Count:       bolt.Count,
			Lenght:      L,
			Temp:        temp,
			Alpha:       mat.AlphaF,
			EpsilonAt20: mat.EpsilonAt20,
			Epsilon:     mat.Epsilon,
			SigmaAt20:   mat.SigmaAt20,
			Sigma:       mat.Sigma,
			Material:    mat.Title,
		}
		return res, nil
	}

	res := &moment_api.BoltResult{
		Diameter:    bolt.Diameter,
		Area:        bolt.Area,
		Count:       bolt.Count,
		Lenght:      L,
		Temp:        temp,
		Alpha:       data.Material.AlphaF,
		EpsilonAt20: data.Material.EpsilonAt20,
		Epsilon:     data.Material.Epsilon,
		SigmaAt20:   data.Material.SigmaAt20,
		Sigma:       data.Material.Sigma,
		Material:    data.Material.Title,
	}
	return res, nil
}

func (s *DataService) getDataFlange(
	ctx context.Context,
	flange *moment_api.FlangeData,
	bolt *moment_api.BoltData,
	typeFlange string,
	temp float64,
) (flangeData *moment_api.FlangeResult, boltSize *moment_api.BoltResult, err error) {
	if flange.StandartId == "another" {
		flangeData = &moment_api.FlangeResult{
			DOut: flange.Size.DOut,
			D:    flange.Size.D,
			Dk:   flange.Size.Dk,
			Dnk:  flange.Size.Dnk,
			Ds:   flange.Size.Ds,
			H:    flange.Size.H,
			H0:   flange.Size.H0,
			Hk:   flange.Size.Hk,
			S0:   flange.Size.S0,
			S1:   flange.Size.S1,
			L:    flange.Size.L,
			D6:   flange.Size.D6,
			C:    flange.Corrosion,
		}
		boltSize = &moment_api.BoltResult{
			Diameter: bolt.Diameter,
			Count:    bolt.Count,
			Area:     bolt.Area,
		}
	} else {
		size, err := s.flange.GetFlangeSize(ctx, &moment_api.GetFlangeSizeRequest{
			D:       float64(flange.Dy),
			Pn:      flange.Py,
			StandId: flange.StandartId,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		flangeData = &moment_api.FlangeResult{
			DOut: size.DOut,
			D:    size.D,
			H:    size.H,
			S0:   size.S0,
			S1:   size.S1,
			L:    size.Length,
			D6:   size.D6,
			C:    flange.Corrosion,
		}
		boltSize = &moment_api.BoltResult{
			Diameter: size.Diameter,
			Count:    size.Count,
			Area:     size.Area,
		}
	}

	flangeData.Tf = s.typeFlangesTF[typeFlange] * temp

	if flange.Type == moment_api.FlangeData_free {
		flangeData.Tk = s.typeFlangesTK[typeFlange] * temp

		//? при свободных фланцах добавляется еще один материал
		var mat models.MaterialsResult
		if flange.RingMarkId != "another" {
			var err error
			mat, err = s.materials.GetMatFotCalculate(ctx, flange.RingMarkId, flangeData.Tk)
			if err != nil {
				return nil, nil, err
			}
		} else {
			mat = models.MaterialsResult{
				Title:       flange.RingMaterial.Title,
				AlphaF:      flange.RingMaterial.AlphaF,
				EpsilonAt20: flange.RingMaterial.EpsilonAt20,
				Epsilon:     flange.RingMaterial.Epsilon,
				SigmaAt20:   flange.RingMaterial.SigmaAt20,
				Sigma:       flange.RingMaterial.Sigma,
			}
		}
		flangeData.RingMaterial = mat.Title
		flangeData.AlphaK = mat.AlphaF
		flangeData.EpsilonKAt20 = mat.EpsilonAt20
		flangeData.EpsilonK = mat.Epsilon
		flangeData.SigmaKAt20 = mat.SigmaAt20
		flangeData.SigmaK = mat.Sigma
	}

	var mat models.MaterialsResult
	if flange.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatFotCalculate(ctx, flange.MarkId, flangeData.Tf)
		if err != nil {
			return nil, nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:       flange.Material.Title,
			AlphaF:      flange.Material.AlphaF,
			EpsilonAt20: flange.Material.EpsilonAt20,
			Epsilon:     flange.Material.Epsilon,
			SigmaAt20:   flange.Material.SigmaAt20,
			Sigma:       flange.Material.Sigma,
		}
	}

	flangeData.Material = mat.Title
	flangeData.AlphaF = mat.AlphaF
	flangeData.EpsilonAt20 = mat.EpsilonAt20
	flangeData.Epsilon = mat.Epsilon
	flangeData.SigmaAt20 = mat.SigmaAt20
	flangeData.Sigma = mat.Sigma

	flangeData.SigmaM = constants.SigmaM * mat.Sigma
	flangeData.SigmaMAt20 = constants.SigmaM * mat.SigmaAt20
	flangeData.SigmaR = constants.SigmaR * mat.Sigma
	flangeData.SigmaRAt20 = constants.SigmaR * mat.SigmaAt20
	flangeData.Type = flange.Type.String()

	return flangeData, boltSize, nil
}

func (s *DataService) getDataCap(
	ctx context.Context,
	cap *moment_api.CapData,
	typeFlange string,
	temp float64,
) (capData *moment_api.CapResult, err error) {
	capData = &moment_api.CapResult{
		H:      cap.H,
		Radius: cap.Radius,
		Delta:  cap.Delta,
	}
	capData.T = s.typeFlangesTF[typeFlange] * temp

	var mat models.MaterialsResult
	if cap.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatFotCalculate(ctx, cap.MarkId, capData.T)
		if err != nil {
			return nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:       cap.Material.Title,
			AlphaF:      cap.Material.AlphaF,
			EpsilonAt20: cap.Material.EpsilonAt20,
			Epsilon:     cap.Material.Epsilon,
		}
	}

	capData.Material = mat.Title
	capData.Alpha = mat.AlphaF
	capData.EpsilonAt20 = mat.EpsilonAt20
	capData.Epsilon = mat.Epsilon

	// flangeData.SigmaM = constants.SigmaM * mat.Sigma
	// flangeData.SigmaMAt20 = constants.SigmaM * mat.SigmaAt20
	// flangeData.SigmaR = constants.SigmaR * mat.Sigma
	// flangeData.SigmaRAt20 = constants.SigmaR * mat.SigmaAt20
	// capData.Type = cap.Type.String()

	return capData, nil
}

func (s *DataService) getCalculatedDataFlange(
	ctx context.Context,
	flangeType moment_api.FlangeData_Type,
	data *moment_api.FlangeResult,
	Dcp float64,
) (*moment_api.FlangeResult, error) {
	calculated := data
	if flangeType != moment_api.FlangeData_free {
		calculated.B = 0.5 * (data.D6 - Dcp)
	} else {
		calculated.Ds = 0.5 * (data.DOut + data.Dk + 2*data.H0)
		calculated.A = 0.5 * (data.D6 - data.Ds)
		calculated.B = 0.5 * (data.Ds - Dcp)
	}

	if flangeType != moment_api.FlangeData_welded {
		calculated.Se = data.S0
	} else {
		calculated.X = data.L / (math.Sqrt(data.D * data.S0))
		calculated.Beta = data.S1 / data.S0
		calculated.Xi = 1 + (calculated.Beta-1)*calculated.X/(calculated.X+(1+calculated.Beta)/4)
		calculated.Se = calculated.Xi * data.S0
	}

	calculated.E = 0.5 * (Dcp - data.D - calculated.Se)
	calculated.L0 = math.Sqrt(data.D * data.S0)
	calculated.K = data.DOut / data.D

	dividend := math.Pow(calculated.K, 2)*(1+8.55*(math.Log(calculated.K)/math.Log(10))) - 1
	divider := (1.05 + 1.945*math.Pow(calculated.K, 2)) * (calculated.K - 1)
	calculated.BetaT = dividend / divider

	divider = 1.36 * (math.Pow(calculated.K, 2) - 1) * (calculated.K - 1)
	calculated.BetaU = dividend / divider

	dividend = 1 / (calculated.K - 1)
	divider = 0.69 + 5.72*((math.Pow(calculated.K, 2)*(math.Log(calculated.K)/math.Log(10)))/(math.Pow(calculated.K, 2)-1))
	calculated.BetaY = dividend * divider

	dividend = math.Pow(calculated.K, 2) + 1
	divider = math.Pow(calculated.K, 2) - 1
	calculated.BetaZ = dividend / divider

	if flangeType == moment_api.FlangeData_welded && data.S0 != data.S1 {
		// Beta := data.S1 / data.S0
		// x := data.L / calculated.L0

		// calculated.BetaF = s.graphic.CalculateBetaF(Beta, x)
		// calculated.BetaV = s.graphic.CalculateBetaV(Beta, x)
		// calculated.F = s.graphic.CalculateF(Beta, x)
		calculated.BetaF = s.graphic.CalculateBetaF(calculated.Beta, calculated.X)
		calculated.BetaV = s.graphic.CalculateBetaV(calculated.Beta, calculated.X)
		calculated.F = s.graphic.CalculateF(calculated.Beta, calculated.X)
	} else {
		calculated.BetaF = constants.InitBetaF
		calculated.BetaV = constants.InitBetaV
		calculated.F = constants.InitF
	}

	calculated.Lymda = (calculated.BetaF*data.H+calculated.L0)/(calculated.BetaT*calculated.L0) +
		+(calculated.BetaV*math.Pow(data.H, 3))/(calculated.BetaU*calculated.L0*math.Pow(data.S0, 2))
	calculated.Yf = (0.91 * calculated.BetaV) / (data.EpsilonAt20 * calculated.Lymda * math.Pow(data.S0, 2) * calculated.L0)

	if flangeType == moment_api.FlangeData_free {
		calculated.Psik = 1.28 * (math.Log(data.Dnk/data.Dk) / math.Log(10))
		calculated.Yk = 1 / (data.EpsilonKAt20 * math.Pow(data.Hk, 3) * calculated.Psik)
	}

	if flangeType != moment_api.FlangeData_free {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
	} else {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.Ds / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
		calculated.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonKAt20 * data.Dnk * math.Pow(data.Hk, 3)))
	}

	return calculated, nil
}

func (s *DataService) getCalculatedDataCap(
	ctx context.Context,
	capType moment_api.CapData_Type,
	data *moment_api.CapResult,
	h, D, S0, DOut, Dcp float64,
) (*moment_api.CapResult, error) {
	calculated := data

	if capType == moment_api.CapData_flat {
		data.K = DOut / Dcp
		data.X = (0.67*math.Pow(data.K, 2)*(1+8.55*math.Log10(data.K)) - 1) / ((data.K - 1) *
			(math.Pow(data.K, 2) - 1 + (1.857*math.Pow(data.K, 2)+1)*(math.Pow(data.H, 3)/math.Pow(data.Delta, 3))))
		data.Y = data.X / (math.Pow(data.Delta, 3) * data.EpsilonAt20)
	} else {
		data.K = (h / D) * math.Sqrt(data.Radius/S0)
		data.X = 1 / (1 + 1.285*data.K + 1.63*data.K*math.Pow((h/S0), 2)*math.Log10(DOut/D))
		data.Y = ((1 - data.X*(1+1.285*data.K)) / (data.EpsilonAt20 * math.Pow(h, 3))) * ((DOut + D) / (DOut - D))
	}

	return calculated, nil
}
