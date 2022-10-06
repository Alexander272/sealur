package data

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

func (s *DataService) getEmbedData(ctx context.Context, data *flange_model.EmbedData, temp float64) (*flange_model.EmbedResult, error) {
	if data.MarkId != "another" {
		washer, err := s.materials.GetMatFotCalculate(ctx, data.MarkId, temp)
		if err != nil {
			return nil, err
		}
		res := &flange_model.EmbedResult{
			Material:  washer.Title,
			Thickness: data.Thickness,
			Alpha:     washer.AlphaF,
			Temp:      temp,
		}
		return res, nil
	}

	res := &flange_model.EmbedResult{
		Material:  data.Material.Title,
		Thickness: data.Thickness,
		Alpha:     data.Material.AlphaF,
		Temp:      temp,
	}
	return res, nil
}
