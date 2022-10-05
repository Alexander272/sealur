package data

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (s *DataService) getWasherData(ctx context.Context, data *moment_api.WasherData, temp float64) (*moment_api.WasherResult, error) {
	if data.MarkId != "another" {
		washer, err := s.materials.GetMatFotCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &moment_api.WasherResult{
			Material:  washer.Title,
			Thickness: data.Thickness,
			Alpha:     washer.AlphaF,
			Temp:      temp,
		}
		return res, nil
	}

	res := &moment_api.WasherResult{
		Material:  data.Material.Title,
		Thickness: data.Thickness,
		Alpha:     data.Material.AlphaF,
		Temp:      temp,
	}
	return res, nil
}
