package data

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
)

func (s *DataService) getDataFlange(
	ctx context.Context,
	flange *cap_model.FlangeData,
	bolt *cap_model.BoltData,
	typeFlange string,
	temp float64,
) (flangeData *cap_model.FlangeResult, boltSize *cap_model.BoltResult, err error) {
	if flange.StandartId == "another" {
		flangeData = &cap_model.FlangeResult{
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

		if bolt.BoltId != "another" {
			b, err := s.flange.GetBolt(ctx, bolt.BoltId)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get bolt size. error: %w", err)
			}
			boltSize = &cap_model.BoltResult{
				Diameter: b.Diameter,
				Count:    bolt.Count,
				Area:     b.Area,
			}
		} else {
			boltSize = &cap_model.BoltResult{
				Diameter: bolt.Diameter,
				Count:    bolt.Count,
				Area:     bolt.Area,
			}
		}
	} else {
		size, err := s.flange.GetFlangeSize(ctx, &flange_api.GetFlangeSizeRequest{
			// D:       flange.Dy,
			Pn:      flange.Py,
			StandId: flange.StandartId,
			Dn:      flange.Dn,
			Row:     flange.Row,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		if size.D == 0 {
			size.D = flange.B
			size.S0 = (size.A - flange.B) / 2
			size.S1 = (size.X - flange.B) / 2
		}
		flangeData = &cap_model.FlangeResult{
			DOut: size.DOut,
			D:    size.D,
			H:    size.H,
			S0:   size.S0,
			S1:   size.S1,
			L:    size.Length,
			D6:   size.D6,
			C:    flange.Corrosion,
		}

		boltSize = &cap_model.BoltResult{
			Diameter: size.Diameter,
			Count:    size.Count,
			Area:     size.Area,
		}
	}

	flangeData.Tf = s.typeFlangesTF[typeFlange] * temp

	if flange.Type == cap_model.FlangeData_free {
		flangeData.Tk = s.typeFlangesTK[typeFlange] * temp

		//? при свободных фланцах добавляется еще один материал
		var mat models.MaterialsResult
		if flange.RingMarkId != "another" {
			var err error
			mat, err = s.materials.GetMatForCalculate(ctx, flange.RingMarkId, flangeData.Tk)
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
		mat, err = s.materials.GetMatForCalculate(ctx, flange.MarkId, flangeData.Tf)
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

func (s *DataService) getCalculatedDataFlange(
	ctx context.Context,
	flangeType cap_model.FlangeData_Type,
	data *cap_model.FlangeResult,
	Dcp float64,
) (*cap_model.FlangeResult, error) {
	calculated := data
	if flangeType != cap_model.FlangeData_free {
		calculated.B = 0.5 * (data.D6 - Dcp)
	} else {
		calculated.Ds = 0.5 * (data.DOut + data.Dk + 2*data.H0)
		calculated.A = 0.5 * (data.D6 - data.Ds)
		calculated.B = 0.5 * (data.Ds - Dcp)
	}

	if flangeType != cap_model.FlangeData_welded {
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

	if flangeType == cap_model.FlangeData_welded && data.S0 != data.S1 {
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

	if flangeType == cap_model.FlangeData_free {
		calculated.Psik = 1.28 * (math.Log(data.Dnk/data.Dk) / math.Log(10))
		calculated.Yk = 1 / (data.EpsilonKAt20 * math.Pow(data.Hk, 3) * calculated.Psik)
	}

	if flangeType != cap_model.FlangeData_free {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
	} else {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.Ds / (data.EpsilonAt20 * data.DOut * math.Pow(data.H, 3)))
		calculated.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonKAt20 * data.Dnk * math.Pow(data.Hk, 3)))
	}

	return calculated, nil
}
