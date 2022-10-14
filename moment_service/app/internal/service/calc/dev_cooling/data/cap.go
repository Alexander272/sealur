package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *DataService) getDataCap(ctx context.Context, cap *dev_cooling_model.CapData, temp float64) (capData *dev_cooling_model.CapResult, err error) {
	capData = &dev_cooling_model.CapResult{
		BottomThick:   cap.BottomThick,
		WallThick:     cap.WallThick,
		FlangeThick:   cap.FlangeThick,
		SideWallThick: cap.SideWallThick,
		InnerSize:     cap.InnerSize,
		OuterSize:     cap.OuterSize,
		Depth:         cap.Depth,
		L:             cap.L,
		Strength:      cap.Strength,
		Corrosion:     cap.Corrosion,
		Radius:        cap.Radius,
	}

	var mat models.MaterialsResult
	if cap.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatFotCalculate(ctx, cap.MarkId, temp)
		if err != nil {
			return nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:       cap.Material.Title,
			EpsilonAt20: cap.Material.Epsilon,
			SigmaAt20:   cap.Material.SigmaAt20,
			Sigma:       cap.Material.Sigma,
		}
	}

	capData.Material = mat.Title
	capData.Epsilon = mat.EpsilonAt20
	capData.SigmaAt20 = mat.SigmaAt20
	capData.Sigma = mat.Sigma

	return capData, nil
}
