package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type MaterialsService struct {
	repo repository.Materials
}

func NewMaterialsService(repo repository.Materials) *MaterialsService {
	return &MaterialsService{repo: repo}
}

func (s *MaterialsService) GetMatFotCalculate(ctx context.Context, markId string, temp float64) (models.MaterialsResult, error) {
	mats, err := s.repo.GetAllData(ctx, markId)
	if err != nil {
		return models.MaterialsResult{}, fmt.Errorf("failed to get materials. error: %w", err)
	}

	var alphaF, epsilonAt20, epsilon, sigmaAt20, sigma float64 = 0, 0, 0, 0, 0

	epsilonAt20 = mats.Elasticity[0].Elasticity * math.Pow10(5)
	sigmaAt20 = mats.Voltage[0].Voltage

	if temp < mats.Alpha[0].Temperature {
		alphaF = mats.Alpha[0].Alpha * math.Pow10(-6)
	} else if temp > mats.Alpha[len(mats.Alpha)-1].Temperature {
		alphaF = mats.Alpha[len(mats.Alpha)-1].Alpha * math.Pow10(-6)
	} else {
		for i, m := range mats.Alpha {
			if i == 0 {
				continue
			}

			if temp >= mats.Alpha[i-1].Temperature && temp < m.Temperature {
				temps := (temp - mats.Alpha[i-1].Temperature) / (m.Temperature - mats.Alpha[i-1].Temperature)
				alphaF = (temps*(m.Alpha-mats.Alpha[i-1].Alpha) + mats.Alpha[i-1].Alpha) * math.Pow10(-6)
				break
			}
		}
	}

	if temp < mats.Elasticity[0].Temperature {
		epsilon = mats.Elasticity[0].Elasticity * math.Pow10(5)
	} else if temp > mats.Elasticity[len(mats.Elasticity)-1].Temperature {
		epsilon = mats.Elasticity[len(mats.Elasticity)-1].Elasticity * math.Pow10(5)
	} else {
		for i, m := range mats.Elasticity {
			if i == 0 {
				continue
			}

			if temp >= mats.Elasticity[i-1].Temperature && temp < m.Temperature {
				temps := (temp - mats.Elasticity[i-1].Temperature) / (m.Temperature - mats.Elasticity[i-1].Temperature)
				epsilon = (temps*(m.Elasticity-mats.Elasticity[i-1].Elasticity) + mats.Elasticity[i-1].Elasticity) * math.Pow10(5)
				break
			}
		}
	}

	if temp < mats.Voltage[0].Temperature {
		sigma = mats.Voltage[0].Voltage
	} else if temp > mats.Voltage[len(mats.Voltage)-1].Temperature {
		sigma = mats.Voltage[len(mats.Voltage)-1].Voltage
	} else {
		for i, m := range mats.Voltage {
			if i == 0 {
				continue
			}

			if temp >= mats.Voltage[i-1].Temperature && temp < m.Temperature {
				temps := (temp - mats.Voltage[i-1].Temperature) / (m.Temperature - mats.Voltage[i-1].Temperature)
				sigma = (temps*(m.Voltage-mats.Voltage[i-1].Voltage) + mats.Voltage[i-1].Voltage)
				break
			}
		}
	}

	res := models.MaterialsResult{
		AlphaF:      alphaF,
		EpsilonAt20: epsilonAt20,
		Epsilon:     epsilon,
		SigmaAt20:   sigmaAt20,
		Sigma:       sigma,
	}

	return res, nil
}

func (s *MaterialsService) GetMaterials(ctx context.Context, req *moment_proto.GetMaterialsRequest) (materials []*moment_proto.Material, err error) {
	mats, err := s.repo.GetMaterials(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get materials. error: %w", err)
	}

	for _, item := range mats {
		m := moment_proto.Material(item)
		materials = append(materials, &m)
	}

	return materials, nil
}

func (s *MaterialsService) CreateMaterial(ctx context.Context, material *moment_proto.CreateMaterialRequest) (id string, err error) {
	id, err = s.repo.CreateMaterial(ctx, material)
	if err != nil {
		return "", fmt.Errorf("failed to create material. error: %w", err)
	}
	return id, nil
}

func (s *MaterialsService) UpdateMaterial(ctx context.Context, material *moment_proto.UpdateMaterialRequest) error {
	if err := s.repo.UpdateMaterial(ctx, material); err != nil {
		return fmt.Errorf("failed to update material. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteMaterial(ctx context.Context, material *moment_proto.DeleteMaterialRequest) error {
	if err := s.repo.DeleteMaterial(ctx, material); err != nil {
		return fmt.Errorf("failed to delete material. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) CreateVoltage(ctx context.Context, voltage *moment_proto.CreateVoltageRequest) error {
	if err := s.repo.CreateVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to create voltage. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpdateVoltage(ctx context.Context, voltage *moment_proto.UpdateVoltageRequest) error {
	if err := s.repo.UpdateVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to update voltage. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteVoltage(ctx context.Context, voltage *moment_proto.DeleteVoltageRequest) error {
	if err := s.repo.DeleteVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to delete voltage. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) CreateElasticity(ctx context.Context, elasticity *moment_proto.CreateElasticityRequest) error {
	err := s.repo.CreateElasticity(ctx, elasticity)
	if err != nil {
		return fmt.Errorf("failed to create elasticity. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpdateElasticity(ctx context.Context, elasticity *moment_proto.UpdateElasticityRequest) error {
	if err := s.repo.UpdateElasticity(ctx, elasticity); err != nil {
		return fmt.Errorf("failed to update elasticity. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteElasticity(ctx context.Context, elasticity *moment_proto.DeleteElasticityRequest) error {
	if err := s.repo.DeleteElasticity(ctx, elasticity); err != nil {
		return fmt.Errorf("failed to delete elasticity. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) CreateAlpha(ctx context.Context, alpha *moment_proto.CreateAlphaRequest) error {
	err := s.repo.CreateAlpha(ctx, alpha)
	if err != nil {
		return fmt.Errorf("failed to create alpha. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpateAlpha(ctx context.Context, alpha *moment_proto.UpdateAlphaRequest) error {
	if err := s.repo.UpateAlpha(ctx, alpha); err != nil {
		return fmt.Errorf("failed to update alpha. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteAlpha(ctx context.Context, alpha *moment_proto.DeleteAlphaRequest) error {
	if err := s.repo.DeleteAlpha(ctx, alpha); err != nil {
		return fmt.Errorf("failed to delete alpha. error: %w", err)
	}
	return nil
}
