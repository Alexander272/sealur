package data

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *DataService) getBoltData(ctx context.Context, data *float_model.BoltData, bolt *float_model.BoltResult, L, temp float64,
) (*float_model.BoltResult, error) {
	if data.MarkId != "another" {
		mat, err := s.materials.GetMatForCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &float_model.BoltResult{
			Diameter:    bolt.Diameter,
			Area:        bolt.Area,
			Count:       bolt.Count,
			Lenght:      L,
			Temp:        temp,
			EpsilonAt20: mat.EpsilonAt20,
			Epsilon:     mat.Epsilon,
			SigmaAt20:   mat.SigmaAt20,
			Sigma:       mat.Sigma,
			Material:    mat.Title,
		}
		return res, nil
	}

	res := &float_model.BoltResult{
		Diameter:    bolt.Diameter,
		Area:        bolt.Area,
		Count:       bolt.Count,
		Lenght:      L,
		Temp:        temp,
		EpsilonAt20: data.Material.EpsilonAt20,
		Epsilon:     data.Material.Epsilon,
		SigmaAt20:   data.Material.SigmaAt20,
		Sigma:       data.Material.Sigma,
		Material:    data.Material.Title,
	}
	return res, nil
}
