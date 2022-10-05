package data

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (s *DataService) getEmbedData(ctx context.Context, data *moment_api.EmbedData, temp float64) (*moment_api.EmbedResult, error) {
	if data.MarkId != "another" {
		washer, err := s.materials.GetMatFotCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &moment_api.EmbedResult{
			Material:  washer.Title,
			Thickness: data.Thickness,
			Alpha:     washer.AlphaF,
			Temp:      temp,
		}
		return res, nil
	}

	res := &moment_api.EmbedResult{
		Material:  data.Material.Title,
		Thickness: data.Thickness,
		Alpha:     data.Material.AlphaF,
		Temp:      temp,
	}
	return res, nil
}
