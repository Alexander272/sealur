package data

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *DataService) getDataFlange(
	ctx context.Context,
	fl *float_model.FlangeData,
	bolt *float_model.BoltData,
) (flangeData *float_model.FlangeResult, boltSize *float_model.BoltResult, err error) {

	flangeData = &float_model.FlangeResult{
		DOut: fl.DOut,
		D:    fl.D,
		D6:   fl.D6,
		Tf:   fl.T,
	}

	if bolt.BoltId != "another" {
		b, err := s.flange.GetBolt(ctx, bolt.BoltId)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get bolt size. error: %w", err)
		}
		boltSize = &float_model.BoltResult{
			Diameter: b.Diameter,
			Count:    bolt.Count,
			Area:     b.Area,
		}
	} else {
		boltSize = &float_model.BoltResult{
			Diameter: bolt.Diameter,
			Count:    bolt.Count,
			Area:     bolt.Area,
		}
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
			EpsilonAt20: fl.Material.EpsilonAt20,
			Epsilon:     fl.Material.Epsilon,
			SigmaAt20:   fl.Material.SigmaAt20,
			Sigma:       fl.Material.Sigma,
		}
	}

	flangeData.Material = mat.Title
	flangeData.EpsilonAt20 = mat.EpsilonAt20
	flangeData.Epsilon = mat.Epsilon
	flangeData.SigmaAt20 = mat.SigmaAt20
	flangeData.Sigma = mat.Sigma

	return flangeData, boltSize, nil
}
