package data

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *DataService) getDataFlange(
	ctx context.Context,
	fl *float_model.FlangeData,
	bolt *float_model.BoltData,
) (flangeData *float_model.FlangeResult, boltSize *float_model.BoltResult, err error) {

	flangeData = &float_model.FlangeResult{
		DOut:  fl.DOut,
		D:     fl.D,
		D6:    fl.D6,
		Tf:    fl.T,
		Width: fl.Width,
		DIn:   fl.DIn,
	}

	if bolt.BoltId != "another" {
		b, err := s.flange.GetBolt(ctx, bolt.BoltId)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get bolt size. error: %w", err)
		}
		boltSize = &float_model.BoltResult{
			Diameter: b.Diameter,
			Count:    bolt.Count,
			Area:     b.Area,
		}
	} else {
		boltSize = &float_model.BoltResult{
			Diameter: bolt.Diameter,
			Count:    bolt.Count,
			Area:     bolt.Area,
		}
	}

	var mat models.MaterialsResult
	if fl.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatForCalculate(ctx, fl.MarkId, flangeData.Tf)
		if err != nil {
			return nil, nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:       fl.Material.Title,
			EpsilonAt20: fl.Material.EpsilonAt20,
			Epsilon:     fl.Material.Epsilon,
			SigmaAt20:   fl.Material.SigmaAt20,
			Sigma:       fl.Material.Sigma,
		}
	}

	flangeData.Material = mat.Title
	flangeData.EpsilonAt20 = mat.EpsilonAt20
	flangeData.Epsilon = mat.Epsilon
	flangeData.SigmaAt20 = mat.SigmaAt20
	flangeData.Sigma = mat.Sigma

	return flangeData, boltSize, nil
}

func (s *DataService) getCalculatedDataFlange(
	ctx context.Context,
	data *float_model.FlangeResult,
	cap *float_model.CapResult,
	Dcp float64,
) (*float_model.FlangeResult, error) {
	calculated := data

	// Плечи действия усилий в болтах/шпильках
	calculated.B = 0.5 * (data.D6 - Dcp)
	// Параметр длины обечайки
	calculated.L0 = math.Sqrt(data.D * cap.S)
	// Угловая податливость фланца при затяжке (она 0 т.к. S0 = 0)
	calculated.Y = 0

	return calculated, nil
}
