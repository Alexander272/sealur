package data

import (
	"context"

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
		mat, err = s.materials.GetMatForCalculate(ctx, cap.MarkId, capData.T)
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
