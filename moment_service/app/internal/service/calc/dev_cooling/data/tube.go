package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *DataService) getTubeData(ctx context.Context, tube *dev_cooling_model.TubeData, temp float64,
) (tubeData *dev_cooling_model.TubeResult, err error) {
	tubeData = &dev_cooling_model.TubeResult{
		Length:        tube.Length,
		ReducedLength: tube.ReducedLength,
		Diameter:      tube.Diameter,
		Thickness:     tube.Thickness,
		Corrosion:     tube.Corrosion,
		Depth:         tube.Depth,
	}

	var mat models.MaterialsResult
	if tube.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatForCalculate(ctx, tube.MarkId, temp)
		if err != nil {
			return nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:       tube.Material.Title,
			EpsilonAt20: tube.Material.Epsilon,
			SigmaAt20:   tube.Material.SigmaAt20,
			Sigma:       tube.Material.Sigma,
		}
	}

	tubeData.Material = mat.Title
	tubeData.Epsilon = mat.EpsilonAt20
	tubeData.SigmaAt20 = mat.SigmaAt20
	tubeData.Sigma = mat.Sigma

	return tubeData, nil
}
