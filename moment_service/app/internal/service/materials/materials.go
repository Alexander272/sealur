package materials

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/material_model"
)

type MaterialsService struct {
	repo repository.Materials
}

func NewMaterialsService(repo repository.Materials) *MaterialsService {
	return &MaterialsService{repo: repo}
}

// Получение данных о материале
func (s *MaterialsService) GetMatForCalculate(ctx context.Context, markId string, temp float64) (models.MaterialsResult, error) {
	mats, err := s.repo.GetAllData(ctx, &material_api.GetMaterialsDataRequest{MarkId: markId})
	if err != nil {
		return models.MaterialsResult{}, fmt.Errorf("failed to get materials. error: %w", err)
	}

	var alphaF, epsilonAt20, epsilon, sigmaAt20, sigma float64 = 0, 0, 0, 0, 0

	epsilonAt20 = math.Round(mats.Elasticity[0].Elasticity*1000) / 1000 * math.Pow10(5)
	sigmaAt20 = math.Round(mats.Voltage[0].Voltage*1000) / 1000

	if temp <= mats.Alpha[0].Temperature {
		alphaF = math.Round(mats.Alpha[0].Alpha*1000) / 1000 * math.Pow10(-6)
	} else if temp >= mats.Alpha[len(mats.Alpha)-1].Temperature {
		alphaF = math.Round(mats.Alpha[len(mats.Alpha)-1].Alpha*1000) / 1000 * math.Pow10(-6)
	} else {
		for i, m := range mats.Alpha {
			if i == 0 {
				continue
			}

			if temp >= mats.Alpha[i-1].Temperature && temp < m.Temperature {
				temps := (temp - mats.Alpha[i-1].Temperature) / (m.Temperature - mats.Alpha[i-1].Temperature)
				al1 := math.Round(m.Alpha*1000) / 1000
				al2 := math.Round(mats.Alpha[i-1].Alpha*1000) / 1000
				alphaF = (temps*(al1-al2) + al2) * math.Pow10(-6)
				break
			}
		}
	}

	if temp <= mats.Elasticity[0].Temperature {
		epsilon = math.Round(mats.Elasticity[0].Elasticity*1000) / 1000 * math.Pow10(5)
	} else if temp >= mats.Elasticity[len(mats.Elasticity)-1].Temperature {
		epsilon = math.Round(mats.Elasticity[len(mats.Elasticity)-1].Elasticity*1000) / 1000 * math.Pow10(5)
	} else {
		for i, m := range mats.Elasticity {
			if i == 0 {
				continue
			}

			if temp >= mats.Elasticity[i-1].Temperature && temp < m.Temperature {
				temps := (temp - mats.Elasticity[i-1].Temperature) / (m.Temperature - mats.Elasticity[i-1].Temperature)
				el1 := math.Round(m.Elasticity*1000) / 1000
				el2 := math.Round(mats.Elasticity[i-1].Elasticity*1000) / 1000
				epsilon = (temps*(el1-el2) + el2) * math.Pow10(5)
				break
			}
		}
	}

	if temp <= mats.Voltage[0].Temperature {
		sigma = math.Round(mats.Voltage[0].Voltage*1000) / 1000
	} else if temp >= mats.Voltage[len(mats.Voltage)-1].Temperature {
		sigma = math.Round(mats.Voltage[len(mats.Voltage)-1].Voltage*1000) / 1000
	} else {
		for i, m := range mats.Voltage {
			if i == 0 {
				continue
			}

			if temp >= mats.Voltage[i-1].Temperature && temp < m.Temperature {
				temps := (temp - mats.Voltage[i-1].Temperature) / (m.Temperature - mats.Voltage[i-1].Temperature)
				v1 := math.Round(m.Voltage*1000) / 1000
				v2 := math.Round(mats.Voltage[i-1].Voltage*1000) / 1000
				sigma = (temps*(v1-v2) + v2)
				break
			}
		}
	}

	res := models.MaterialsResult{
		Title:       mats.Title,
		AlphaF:      alphaF,
		EpsilonAt20: epsilonAt20,
		Epsilon:     epsilon,
		SigmaAt20:   sigmaAt20,
		Sigma:       sigma,
	}

	return res, nil
}

func (s *MaterialsService) GetMaterials(ctx context.Context, req *material_api.GetMaterialsRequest) (materials []*material_model.Material, err error) {
	mats, err := s.repo.GetMaterials(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get materials. error: %w", err)
	}

	for _, item := range mats {
		materials = append(materials, &material_model.Material{
			Id:    item.Id,
			Title: item.Title,
		})
	}

	return materials, nil
}

func (s *MaterialsService) GetMaterialsWithIsEmpty(ctx context.Context, req *material_api.GetMaterialsRequest,
) (materials []*material_model.MaterialWithIsEmpty, err error) {
	mats, err := s.repo.GetMaterialsWithIsEmpty(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get materials. error: %w", err)
	}

	for _, item := range mats {
		materials = append(materials, &material_model.MaterialWithIsEmpty{
			Id:                item.Id,
			Title:             item.Title,
			IsEmptyAlpha:      item.IsEmptyAlpha,
			IsEmptyElasticity: item.IsEmptyElasticity,
			IsEmptyVoltage:    item.IsEmptyVoltage,
		})
	}

	return materials, nil
}

func (s *MaterialsService) GetMaterialsData(ctx context.Context, req *material_api.GetMaterialsDataRequest,
) (materials *material_api.MaterialsDataResponse, err error) {
	mats, err := s.repo.GetAllData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get materials. error: %w", err)
	}

	voltage := make([]*material_api.MaterialsDataResponse_Voltage, 0, len(mats.Voltage))
	for _, item := range mats.Voltage {
		item.Voltage = math.Round(item.Voltage*1000) / 1000
		voltage = append(voltage, &material_api.MaterialsDataResponse_Voltage{
			Id:          item.Id,
			Temperature: item.Temperature,
			Voltage:     item.Voltage,
		})
	}

	elasticity := make([]*material_api.MaterialsDataResponse_Elasticity, 0, len(mats.Elasticity))
	for _, item := range mats.Elasticity {
		item.Elasticity = math.Round(item.Elasticity*1000) / 1000

		elasticity = append(elasticity, &material_api.MaterialsDataResponse_Elasticity{
			Id:          item.Id,
			Temperature: item.Temperature,
			Elasticity:  item.Elasticity,
		})
	}

	alpha := make([]*material_api.MaterialsDataResponse_Alpha, 0, len(mats.Alpha))
	for _, item := range mats.Alpha {
		item.Alpha = math.Round(item.Alpha*1000) / 1000

		alpha = append(alpha, &material_api.MaterialsDataResponse_Alpha{
			Id:          item.Id,
			Temperature: item.Temperature,
			Alpha:       item.Alpha,
		})
	}

	materials = &material_api.MaterialsDataResponse{
		Voltage:    voltage,
		Elasticity: elasticity,
		Alpha:      alpha,
	}

	return materials, nil
}

func (s *MaterialsService) CreateMaterial(ctx context.Context, material *material_api.CreateMaterialRequest) (id string, err error) {
	id, err = s.repo.CreateMaterial(ctx, material)
	if err != nil {
		return "", fmt.Errorf("failed to create material. error: %w", err)
	}
	return id, nil
}

func (s *MaterialsService) UpdateMaterial(ctx context.Context, material *material_api.UpdateMaterialRequest) error {
	if err := s.repo.UpdateMaterial(ctx, material); err != nil {
		return fmt.Errorf("failed to update material. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteMaterial(ctx context.Context, material *material_api.DeleteMaterialRequest) error {
	if err := s.repo.DeleteMaterial(ctx, material); err != nil {
		return fmt.Errorf("failed to delete material. error: %w", err)
	}
	return nil
}
