package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
)

type FlangeService struct {
	repo      repository.Flange
	materials *MaterialsService
	gasket    *GasketService
}

func NewFlangeService(repo repository.Flange, materials *MaterialsService, gasket *GasketService) *FlangeService {
	return &FlangeService{
		repo:      repo,
		materials: materials,
		gasket:    gasket,
	}
}

func (s *FlangeService) Calculation(ctx context.Context, data *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error) {
	size, err := s.repo.GetSize(ctx, 400, 1.0)
	if err != nil {
		return nil, fmt.Errorf("failed to get size. error: %w", err)
	}

	var dataFlangeFirst models.InitialDataFlange

	//TODO добавить зависимость от типа фланца
	dataFlangeFirst.Tf = constants.IsolatedFlatTf * float64(data.Temp)

	mat, err := s.materials.GetMatFotCalculate(ctx, data.FlangesData[0].MarkId, dataFlangeFirst.Tf)
	if err != nil {
		return nil, err
	}

	dataFlangeFirst = models.InitialDataFlange{
		DOut:        size.D1,
		D:           size.D,
		H:           size.B,
		S0:          size.S0,
		S1:          size.S1,
		L:           size.Lenght,
		D6:          size.D2,
		AlphaF:      mat.AlphaF,
		EpsilonAt20: mat.EpsilonAt20,
		Epsilon:     mat.Epsilon,
		SigmaAt20:   mat.SigmaAt20,
		Sigma:       mat.Sigma,
	}

	dataFlangeFirst.SigmaM = 1.5 * mat.Sigma
	dataFlangeFirst.SigmaMAt20 = 1.5 * mat.SigmaAt20
	dataFlangeFirst.SigmaR = 3 * mat.Sigma
	dataFlangeFirst.SigmaRAt20 = 3 * mat.SigmaAt20

	logger.Debug(dataFlangeFirst)

	//TODO добавить зависимость от типа фланца
	Tb := constants.IsolatedFlatTb * float64(data.Temp)

	//TODO
	boltMat, err := s.materials.GetMatFotCalculate(ctx, "1", Tb)
	if err != nil {
		return nil, err
	}
	logger.Debug(boltMat)

	gasket, err := s.gasket.Get(ctx, models.GetGasket{TypeGasket: "", Env: "", Thickness: 3.2})
	if err != nil {
		return nil, err
	}
	logger.Debug(gasket)

	return &moment_proto.FlangeResponse{}, nil
}
