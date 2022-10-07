package data

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *DataService) getDataCap(
	ctx context.Context,
	cap *cap_model.CapData,
	typeFlange string,
	temp float64,
) (capData *cap_model.CapResult, err error) {
	capData = &cap_model.CapResult{
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
	capData.Type = cap.Type.String()

	return capData, nil
}

func (s *DataService) getCalculatedDataCap(
	ctx context.Context,
	capType cap_model.CapData_Type,
	calc *cap_model.CapResult,
	h, D, S0, DOut, Dcp float64,
) (*cap_model.CapResult, error) {
	data := calc

	if capType == cap_model.CapData_flat {
		data.K = DOut / Dcp
		data.X = (0.67*math.Pow(data.K, 2)*(1+8.55*math.Log10(data.K)) - 1) / ((data.K - 1) *
			(math.Pow(data.K, 2) - 1 + (1.857*math.Pow(data.K, 2)+1)*(math.Pow(data.H, 3)/math.Pow(data.Delta, 3))))
		data.Y = data.X / (math.Pow(data.Delta, 3) * data.EpsilonAt20)
	} else {
		data.Lambda = (h / D) * math.Sqrt(data.Radius/S0)
		data.Omega = 1 / (1 + 1.285*data.Lambda + 1.63*data.Lambda*math.Pow((h/S0), 2)*math.Log10(DOut/D))
		data.Y = ((1 - data.Omega*(1+1.285*data.Lambda)) / (data.EpsilonAt20 * math.Pow(h, 3))) * ((DOut + D) / (DOut - D))
	}

	return data, nil
}
