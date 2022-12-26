package data

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
)

func (s *DataService) getDataFlange(
	ctx context.Context,
	flange *cap_model.FlangeData,
	bolt *cap_model.BoltData,
	typeFlange string,
	temp float64,
) (flangeData *cap_model.FlangeResult, boltSize *cap_model.BoltResult, err error) {
	if flange.StandartId == "another" {
		flangeData = &cap_model.FlangeResult{
			DOut: flange.Size.DOut,
			D:    flange.Size.D,
			H:    flange.Size.H,
			H0:   flange.Size.H0,
			S0:   flange.Size.S0,
			S1:   flange.Size.S1,
			L:    flange.Size.L,
			D6:   flange.Size.D6,
			C:    flange.Corrosion,
			Ring: &cap_model.FlangeResult_Ring{
				Dk:  flange.Size.Dk,
				Dnk: flange.Size.Dnk,
				Ds:  flange.Size.Ds,
				Hk:  flange.Size.Hk,
			},
		}

		if bolt.BoltId != "another" {
			b, err := s.flange.GetBolt(ctx, bolt.BoltId)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get bolt size. error: %w", err)
			}
			boltSize = &cap_model.BoltResult{
				Diameter: b.Diameter,
				Count:    bolt.Count,
				Area:     b.Area,
			}
		} else {
			boltSize = &cap_model.BoltResult{
				Diameter: bolt.Diameter,
				Count:    bolt.Count,
				Area:     bolt.Area,
			}
		}
	} else {
		size, err := s.flange.GetFlangeSize(ctx, &flange_api.GetFlangeSizeRequest{
			// D:       flange.Dy,
			Pn:      flange.Py,
			StandId: flange.StandartId,
			Dn:      flange.Dn,
			Row:     flange.Row,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		if size.D == 0 {
			size.D = flange.B
			size.S0 = (size.A - flange.B) / 2
			size.S1 = (size.X - flange.B) / 2
		}
		flangeData = &cap_model.FlangeResult{
			DOut: size.DOut,
			D:    size.D,
			H:    size.H,
			S0:   size.S0,
			S1:   size.S1,
			L:    size.Length,
			D6:   size.D6,
			C:    flange.Corrosion,
		}

		boltSize = &cap_model.BoltResult{
			Diameter: size.Diameter,
			Count:    size.Count,
			Area:     size.Area,
		}
	}

	flangeData.T = s.typeFlangesTF[typeFlange] * temp

	if flange.Type == cap_model.FlangeData_free {
		flangeData.Ring.T = s.typeFlangesTK[typeFlange] * temp

		//? при свободных фланцах добавляется еще один материал
		var mat models.MaterialsResult
		if flange.RingMarkId != "another" {
			var err error
			mat, err = s.materials.GetMatForCalculate(ctx, flange.RingMarkId, flangeData.T)
			if err != nil {
				return nil, nil, err
			}
		} else {
			mat = models.MaterialsResult{
				Title:       flange.RingMaterial.Title,
				AlphaF:      flange.RingMaterial.AlphaF,
				EpsilonAt20: flange.RingMaterial.EpsilonAt20,
				Epsilon:     flange.RingMaterial.Epsilon,
				SigmaAt20:   flange.RingMaterial.SigmaAt20,
				Sigma:       flange.RingMaterial.Sigma,
			}
		}
		flangeData.Ring.Material = mat.Title
		flangeData.Ring.Alpha = mat.AlphaF
		flangeData.Ring.EpsilonAt20 = mat.EpsilonAt20
		flangeData.Ring.Epsilon = mat.Epsilon
		flangeData.Ring.SigmaAt20 = mat.SigmaAt20
		flangeData.Ring.Sigma = mat.Sigma
	}

	var mat models.MaterialsResult
	if flange.MarkId != "another" {
		var err error
		mat, err = s.materials.GetMatForCalculate(ctx, flange.MarkId, flangeData.T)
		if err != nil {
			return nil, nil, err
		}
	} else {
		mat = models.MaterialsResult{
			Title:       flange.Material.Title,
			AlphaF:      flange.Material.AlphaF,
			EpsilonAt20: flange.Material.EpsilonAt20,
			Epsilon:     flange.Material.Epsilon,
			SigmaAt20:   flange.Material.SigmaAt20,
			Sigma:       flange.Material.Sigma,
		}
	}

	flangeData.Material = mat.Title
	flangeData.Alpha = mat.AlphaF
	flangeData.EpsilonAt20 = mat.EpsilonAt20
	flangeData.Epsilon = mat.Epsilon
	flangeData.SigmaAt20 = mat.SigmaAt20
	flangeData.Sigma = mat.Sigma

	flangeData.SigmaM = constants.SigmaM * mat.Sigma
	flangeData.SigmaMAt20 = constants.SigmaM * mat.SigmaAt20
	flangeData.SigmaR = constants.SigmaR * mat.Sigma
	flangeData.SigmaRAt20 = constants.SigmaR * mat.SigmaAt20
	flangeData.Type = flange.Type.String()

	return flangeData, boltSize, nil
}
