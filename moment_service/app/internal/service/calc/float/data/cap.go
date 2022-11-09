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
		mat, err = s.materials.GetMatForCalculate(ctx, cap.MarkId, capData.T)
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
	flange *float_model.FlangeResult,
	cap *float_model.CapResult,
) (*float_model.CapResult, error) {
	data := cap

	data.Lambda = (cap.H / flange.D) * math.Sqrt(data.Radius/cap.S)
	data.Omega = 1 / (1 + 1.285*data.Lambda + 1.63*data.Lambda*math.Pow((cap.H/cap.S), 2)*math.Log10(flange.DOut/flange.D))
	// Угловая податливость фланца со сферической неотбортованной крышкой
	data.Y = ((1 - data.Omega*(1+1.285*data.Lambda)) / (data.EpsilonAt20 * math.Pow(cap.H, 3))) *
		((flange.DOut + flange.D) / (flange.DOut - flange.D))

	return data, nil
}
