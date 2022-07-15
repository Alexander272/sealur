package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
)

type MaterialsService struct {
	repo repository.Materials
}

func NewMaterialsService(repo repository.Materials) *MaterialsService {
	return &MaterialsService{repo: repo}
}

func (s *MaterialsService) GetMatFotCalculate(ctx context.Context, markId string, temp float64) (models.MaterialsResult, error) {
	mats, err := s.repo.GetAll(ctx, markId)
	if err != nil {
		return models.MaterialsResult{}, fmt.Errorf("failed to get materials. error: %w", err)
	}

	var alphaF, epsilonAt20, epsilon, sigmaAt20, sigma float64 = 0, 0, 0, 0, 0

	epsilonAt20 = mats[0].Elasticity * math.Pow10(5)
	sigmaAt20 = mats[0].Voltage

	//TODO разделить на три части (будет 3 таблицы)
	if temp < mats[0].Temp {
		alphaF = mats[0].Alpha * math.Pow10(-6)
		epsilon = mats[0].Elasticity * math.Pow10(5)
		sigma = mats[0].Voltage
	} else {
		for i, m := range mats {
			if i == 0 {
				continue
			}
			temps := (temp - mats[i-1].Temp) / (m.Temp - mats[i-1].Temp)

			if temp >= mats[i-1].Temp && temp < m.Temp {
				alphaF = (temps*(m.Alpha-mats[i-1].Alpha) + mats[i-1].Alpha) * math.Pow10(-6)
			}

			if temp >= mats[i-1].Temp && temp < m.Temp {
				epsilon = (temps*(m.Elasticity-mats[i-1].Elasticity) + mats[i-1].Elasticity) * math.Pow10(5)
			}

			if temp >= mats[i-1].Temp && temp < m.Temp {
				sigma = temps*(m.Voltage-mats[i-1].Voltage) + mats[i-1].Voltage
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
