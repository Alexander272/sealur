package data

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
)

func (s *DataService) getDataFlange(
	ctx context.Context,
	fl *flange_model.FlangeData,
	bolt *flange_model.BoltData,
	typeFlange string,
	temp float64,
) (flangeData *flange_model.FlangeResult, boltSize *flange_model.BoltResult, err error) {
	if fl.StandartId == "another" {
		flangeData = &flange_model.FlangeResult{
			DOut: fl.Size.DOut,
			D:    fl.Size.D,
			Dk:   fl.Size.Dk,
			Dnk:  fl.Size.Dnk,
			Ds:   fl.Size.Ds,
			H:    fl.Size.H,
			H0:   fl.Size.H0,
			Hk:   fl.Size.Hk,
			S0:   fl.Size.S0,
			S1:   fl.Size.S1,
			L:    fl.Size.L,
			D6:   fl.Size.D6,
			C:    fl.Corrosion,
		}
		boltSize = &flange_model.BoltResult{
			Diameter: bolt.Diameter,
			Count:    bolt.Count,
			Area:     bolt.Area,
		}
	} else {
		size, err := s.flange.GetFlangeSize(ctx, &flange_api.GetFlangeSizeRequest{
			// D:       flange.Dy,
			Pn:      fl.Py,
			StandId: fl.StandartId,
			Dn:      fl.Dn,
			Row:     fl.Row,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		flangeData = &flange_model.FlangeResult{
			DOut: size.DOut,
			D:    size.D,
			H:    size.H,
			S0:   size.S0,
			S1:   size.S1,
			L:    size.Length,
			D6:   size.D6,
			C:    fl.Corrosion,
		}
		if size.D == 0 {
			size.D = fl.B
			size.S0 = (size.A - fl.B) / 2
			size.S1 = (size.X - fl.B) / 2
		}

		boltSize = &flange_model.BoltResult{
			Diameter: size.Diameter,
			Count:    size.Count,
			Area:     size.Area,
		}
	}

	flangeData.Tf = s.typeFlangesTF[typeFlange] * temp

	if fl.Type == flange_model.FlangeData_free {
		flangeData.Tk = s.typeFlangesTK[typeFlange] * temp

		//? при свободных фланцах добавляется еще один материал
		var mat models.MaterialsResult
		if fl.RingMarkId != "another" {
			var err error
			mat, err = s.materials.GetMatFotCalculate(ctx, fl.RingMarkId, flangeData.Tk)
			if err != nil {
				return nil, nil, err
			}
		} else {
			mat = models.MaterialsResult{
				Title:       fl.RingMaterial.Title,
				AlphaF:      fl.RingMaterial.AlphaF,
				EpsilonAt20: fl.RingMaterial.EpsilonAt20,
				Epsilon:     fl.RingMaterial.Epsilon,
				SigmaAt20:   fl.RingMaterial.SigmaAt20,
				Sigma:       fl.RingMaterial.Sigma,
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
	if fl.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatFotCalculate(ctx, fl.MarkId, flangeData.Tf)
		if err != nil {
			return nil, nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:       fl.Material.Title,
			AlphaF:      fl.Material.AlphaF,
			EpsilonAt20: fl.Material.EpsilonAt20,
			Epsilon:     fl.Material.Epsilon,
			SigmaAt20:   fl.Material.SigmaAt20,
			Sigma:       fl.Material.Sigma,
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
	flangeData.Type = fl.Type.String()

	return flangeData, boltSize, nil
}

func (s *DataService) getCalculatedDataFlange(
	ctx context.Context,
	flangeType flange_model.FlangeData_Type,
	data *flange_model.FlangeResult,
	Dcp float64,
) (*flange_model.FlangeResult, error) {
	calculated := data
	if flangeType != flange_model.FlangeData_free {
		calculated.B = 0.5 * (data.D6 - Dcp)
	} else {
		calculated.Ds = 0.5 * (data.DOut + data.Dk + 2*data.H0)
		calculated.A = 0.5 * (data.D6 - data.Ds)
		calculated.B = 0.5 * (data.Ds - Dcp)
	}

	if flangeType != flange_model.FlangeData_welded {
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

	if flangeType == flange_model.FlangeData_welded && data.S0 != data.S1 {
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

	if flangeType == flange_model.FlangeData_free {
		calculated.Psik = 1.28 * (math.Log(data.Dnk/data.Dk) / math.Log(10))
		calculated.Yk = 1 / (data.EpsilonKAt20 * math.Pow(data.Hk, 3) * calculated.Psik)
	}

	if flangeType != flange_model.FlangeData_free {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
	} else {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.Ds / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
		calculated.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonKAt20 * data.Dnk * math.Pow(data.Hk, 3)))
	}

	return calculated, nil
}
