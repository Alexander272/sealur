package data

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
)

func (s *DataService) getBoltData(ctx context.Context, bolt *express_circle_model.BoltData, temp float64,
) (boltData *express_circle_model.BoltResult, err error) {
	boltData = &express_circle_model.BoltResult{
		Count: bolt.Count,
	}

	if bolt.BoltId != "another" {
		b, err := s.flange.GetBolt(ctx, bolt.BoltId)
		if err != nil {
			return nil, fmt.Errorf("failed to get bolt size. error: %w", err)
		}
		boltData.Diameter = b.Diameter
		boltData.Area = b.Area
	} else {
		boltData.Diameter = bolt.Diameter
		boltData.Area = bolt.Area
	}

	if bolt.MarkId != "another" {
		mat, err := s.materials.GetMatForCalculate(ctx, bolt.MarkId, temp)
		if err != nil {
			return nil, err
		}
		boltData.Material = mat.Title
		boltData.EpsilonAt20 = mat.EpsilonAt20
		boltData.SigmaAt20 = mat.SigmaAt20
	} else {
		boltData.Material = bolt.Material.Title
		boltData.EpsilonAt20 = bolt.Material.EpsilonAt20
		boltData.SigmaAt20 = bolt.Material.SigmaAt20
	}

	return boltData, nil
}
