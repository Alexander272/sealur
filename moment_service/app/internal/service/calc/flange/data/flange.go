package data

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
)

func (s *DataService) getFlangeData(ctx context.Context,
	fl *flange_model.FlangeData,
	bolt *flange_model.BoltData,
	typeFlange string,
	temp float64,
) (flangeData *flange_model.FlangeResult, boltSize *flange_model.BoltResult, err error) {
	if fl.StandartId == "another" {
		flangeData = &flange_model.FlangeResult{
			DOut: fl.Size.DOut,
			D:    fl.Size.D,
			Dk:   fl.Size.Dk,
			Dnk:  fl.Size.Dnk,
			Ds:   fl.Size.Ds,
			H:    fl.Size.H,
			H0:   fl.Size.H0,
			Hk:   fl.Size.Hk,
			S0:   fl.Size.S0,
			S1:   fl.Size.S1,
			L:    fl.Size.L,
			D6:   fl.Size.D6,
			C:    fl.Corrosion,
		}

		if bolt.BoltId != "another" {
			b, err := s.flange.GetBolt(ctx, bolt.BoltId)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get bolt size. error: %w", err)
			}
			boltSize = &flange_model.BoltResult{
				Diameter: b.Diameter,
				Count:    bolt.Count,
				Area:     b.Area,
			}
		} else {
			boltSize = &flange_model.BoltResult{
				Diameter: bolt.Diameter,
				Count:    bolt.Count,
				Area:     bolt.Area,
			}
		}

	} else {
		size, err := s.flange.GetFlangeSize(ctx, &flange_api.GetFlangeSizeRequest{
			// D:       fl.Dy,
			Pn:      fl.Py,
			StandId: fl.StandartId,
			Dn:      fl.Dn,
			Row:     fl.Row,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get size. error: %w", err)
		}

		flangeData = &flange_model.FlangeResult{
			DOut: size.DOut,
			D:    size.D,
			H:    size.H,
			S0:   size.S0,
			S1:   size.S1,
			L:    size.Length,
			D6:   size.D6,
			C:    fl.Corrosion,
		}
		if size.D == 0 {
			size.D = fl.B
			size.S0 = (size.A - fl.B) / 2
			size.S1 = (size.X - fl.B) / 2
		}

		boltSize = &flange_model.BoltResult{
			Diameter: size.Diameter,
			Count:    size.Count,
			Area:     size.Area,
		}
	}

	flangeData.Tf = s.typeFlangesTF[typeFlange] * temp

	if fl.Type == flange_model.FlangeData_free {
		flangeData.Ds = 0.5 * (flangeData.DOut + flangeData.Dk + 2*flangeData.H0)
		flangeData.Ring.Tk = s.typeFlangesTK[typeFlange] * temp

		//? при свободных фланцах добавляется еще один материал
		var mat models.MaterialsResult
		if fl.RingMarkId != "another" {
			var err error
			mat, err = s.materials.GetMatForCalculate(ctx, fl.RingMarkId, flangeData.Ring.Tk)
			if err != nil {
				return nil, nil, err
			}
		} else {
			mat = models.MaterialsResult{
				Title:       fl.RingMaterial.Title,
				AlphaF:      fl.RingMaterial.AlphaF,
				EpsilonAt20: fl.RingMaterial.EpsilonAt20,
				Epsilon:     fl.RingMaterial.Epsilon,
				SigmaAt20:   fl.RingMaterial.SigmaAt20,
				Sigma:       fl.RingMaterial.Sigma,
			}
		}
		flangeData.Ring.Material = mat.Title
		flangeData.Ring.AlphaK = mat.AlphaF
		flangeData.Ring.EpsilonKAt20 = mat.EpsilonAt20
		flangeData.Ring.EpsilonK = mat.Epsilon
		flangeData.Ring.SigmaKAt20 = mat.SigmaAt20
		flangeData.Ring.SigmaK = mat.Sigma
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
			AlphaF:      fl.Material.AlphaF,
			EpsilonAt20: fl.Material.EpsilonAt20,
			Epsilon:     fl.Material.Epsilon,
			SigmaAt20:   fl.Material.SigmaAt20,
			Sigma:       fl.Material.Sigma,
		}
	}

	flangeData.Material = mat.Title
	flangeData.AlphaF = mat.AlphaF
	flangeData.EpsilonAt20 = mat.EpsilonAt20
	flangeData.Epsilon = mat.Epsilon
	flangeData.SigmaAt20 = mat.SigmaAt20
	flangeData.Sigma = mat.Sigma

	flangeData.SigmaM = constants.SigmaM * mat.Sigma
	flangeData.SigmaMAt20 = constants.SigmaM * mat.SigmaAt20
	flangeData.SigmaR = constants.SigmaR * mat.Sigma
	flangeData.SigmaRAt20 = constants.SigmaR * mat.SigmaAt20
	flangeData.Type = fl.Type.String()

	return flangeData, boltSize, nil
}
