package data

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *DataService) getBoltData(ctx context.Context, bolt *dev_cooling_model.BoltData, temp float64,
) (boltData *dev_cooling_model.BoltResult, err error) {
	boltData = &dev_cooling_model.BoltResult{
		Distance: bolt.Distance,
		Count:    bolt.Count,
		Lenght:   bolt.Lenght,
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
		mat, err := s.materials.GetMatFotCalculate(ctx, bolt.MarkId, temp)
		if err != nil {
			return nil, err
		}
		boltData.Material = mat.Title
		boltData.Epsilon = mat.Epsilon
		boltData.SigmaAt20 = mat.SigmaAt20
		boltData.Sigma = mat.Sigma
	} else {
		boltData.Material = bolt.Material.Title
		boltData.Epsilon = bolt.Material.Epsilon
		boltData.SigmaAt20 = bolt.Material.SigmaAt20
		boltData.Sigma = bolt.Material.Sigma
	}

	return boltData, nil
}
