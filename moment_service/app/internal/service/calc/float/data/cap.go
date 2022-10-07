package data

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *DataService) getDataCap(ctx context.Context, cap *float_model.CapData) (capData *float_model.CapResult, err error) {
	capData = &float_model.CapResult{
		H:      cap.H,
		Radius: cap.Radius,
		S:      cap.S,
		T:      cap.T,
	}

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
			EpsilonAt20: cap.Material.EpsilonAt20,
			Epsilon:     cap.Material.Epsilon,
		}
	}

	capData.Material = mat.Title
	capData.EpsilonAt20 = mat.EpsilonAt20
	capData.Epsilon = mat.Epsilon

	return capData, nil
}

func (s *DataService) getCalculatedDataCap(
	ctx context.Context,
	calc *float_model.CapResult,
	h, D, S0, DOut, Dcp float64,
) (*float_model.CapResult, error) {
	data := calc

	data.Lambda = (h / D) * math.Sqrt(data.Radius/S0)
	data.Omega = 1 / (1 + 1.285*data.Lambda + 1.63*data.Lambda*math.Pow((h/S0), 2)*math.Log10(DOut/D))
	data.Y = ((1 - data.Omega*(1+1.285*data.Lambda)) / (data.EpsilonAt20 * math.Pow(h, 3))) * ((DOut + D) / (DOut - D))

	return data, nil
}
