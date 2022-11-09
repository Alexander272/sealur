package data

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *DataService) getBoltData(ctx context.Context, data *cap_model.BoltData, bolt *cap_model.BoltResult, L, temp float64,
) (*cap_model.BoltResult, error) {
	if data.MarkId != "another" {
		mat, err := s.materials.GetMatForCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &cap_model.BoltResult{
			Diameter:    bolt.Diameter,
			Area:        bolt.Area,
			Count:       bolt.Count,
			Lenght:      L,
			Temp:        temp,
			Alpha:       mat.AlphaF,
			EpsilonAt20: mat.EpsilonAt20,
			Epsilon:     mat.Epsilon,
			SigmaAt20:   mat.SigmaAt20,
			Sigma:       mat.Sigma,
			Material:    mat.Title,
		}
		return res, nil
	}

	res := &cap_model.BoltResult{
		Diameter:    bolt.Diameter,
		Area:        bolt.Area,
		Count:       bolt.Count,
		Lenght:      L,
		Temp:        temp,
		Alpha:       data.Material.AlphaF,
		EpsilonAt20: data.Material.EpsilonAt20,
		Epsilon:     data.Material.Epsilon,
		SigmaAt20:   data.Material.SigmaAt20,
		Sigma:       data.Material.Sigma,
		Material:    data.Material.Title,
	}
	return res, nil
}
