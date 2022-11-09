package data

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

func (s *DataService) getBoltData(ctx context.Context, data *flange_model.BoltData, bolt *flange_model.BoltResult, lenght, temp float64,
) (*flange_model.BoltResult, error) {
	res := &flange_model.BoltResult{
		Diameter: bolt.Diameter,
		Area:     bolt.Area,
		Count:    bolt.Count,
		Lenght:   lenght,
		Temp:     temp,
	}
	if data.MarkId != "another" {
		mat, err := s.materials.GetMatForCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}

		res.Alpha = mat.AlphaF
		res.EpsilonAt20 = mat.EpsilonAt20
		res.Epsilon = mat.Epsilon
		res.SigmaAt20 = mat.SigmaAt20
		res.Sigma = mat.Sigma
		res.Material = mat.Title
		return res, nil
	}

	res.Alpha = data.Material.AlphaF
	res.EpsilonAt20 = data.Material.EpsilonAt20
	res.Epsilon = data.Material.Epsilon
	res.SigmaAt20 = data.Material.SigmaAt20
	res.Sigma = data.Material.Sigma
	res.Material = data.Material.Title
	return res, nil
}
