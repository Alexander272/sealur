package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

func (s *DataService) getTubeSheetData(ctx context.Context, tubeSheet *dev_cooling_model.TubeSheetData, temp float64,
) (tubeSheetData *dev_cooling_model.TubeSheetResult, err error) {
	tubeSheetData = &dev_cooling_model.TubeSheetResult{
		ZoneThick:    tubeSheet.ZoneThick,
		PlaceThick:   tubeSheet.PlaceThick,
		OutZoneThick: tubeSheet.OutZoneThick,
		Width:        tubeSheet.Width,
		StepLong:     tubeSheet.StepLong,
		StepTrans:    tubeSheet.StepTrans,
		Count:        tubeSheet.Count,
		Diameter:     tubeSheet.Diameter,
		Corrosion:    tubeSheet.Corrosion,
	}

	var mat models.MaterialsResult
	if tubeSheet.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatFotCalculate(ctx, tubeSheet.MarkId, temp)
		if err != nil {
			return nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:     tubeSheet.Material.Title,
			Epsilon:   tubeSheet.Material.Epsilon,
			SigmaAt20: tubeSheet.Material.SigmaAt20,
			Sigma:     tubeSheet.Material.Sigma,
		}
	}

	tubeSheetData.Material = mat.Title
	tubeSheetData.Epsilon = mat.Epsilon
	tubeSheetData.SigmaAt20 = mat.SigmaAt20
	tubeSheetData.Sigma = mat.Sigma

	return tubeSheetData, nil
}
