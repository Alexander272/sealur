package gasket

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/gasket_model"
)

type GasketService struct {
	repo repository.Gasket
}

func NewGasketService(repo repository.Gasket) *GasketService {
	return &GasketService{
		repo: repo,
	}
}

func (s *GasketService) GetFullData(ctx context.Context, gasket models.GetGasket) (models.FullDataGasket, error) {
	g, err := s.repo.GetFullData(ctx, gasket)
	if err != nil {
		return models.FullDataGasket{}, fmt.Errorf("failed to get gasket. error: %w", err)
	}
	return g, nil
}

func (s *GasketService) GetData(ctx context.Context, gasket *gasket_api.GetFullDataRequest) (*gasket_api.FullDataResponse, error) {
	gasketData, err := s.repo.GetGasketData(ctx, gasket.GasketId)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasket data. error: %w", err)
	}
	GasketData := []*gasket_model.Full_GasketData{}
	for _, gdd := range gasketData {
		GasketData = append(GasketData, &gasket_model.Full_GasketData{
			Id:              gdd.Id,
			GasketId:        gdd.GasketId,
			PermissiblePres: gdd.PermissiblePres,
			Compression:     gdd.Compression,
			Epsilon:         gdd.Epsilon,
			Thickness:       gdd.Thickness,
			TypeId:          gdd.TypeId,
		})
	}

	gasketType, err := s.repo.GetTypeGasket(ctx, &gasket_api.GetGasketTypeRequest{TypeGasket: []gasket_api.TypeGasket{gasket_api.TypeGasket_All}})
	if err != nil {
		return nil, fmt.Errorf("failed to get type gasket. error: %w", err)
	}
	GasketType := []*gasket_model.GasketType{}
	for _, tgd := range gasketType {
		GasketType = append(GasketType, &gasket_model.GasketType{
			Id:    tgd.Id,
			Title: tgd.Title,
			Label: tgd.Label,
		})
	}

	envData, err := s.repo.GetEnvData(ctx, gasket.GasketId)
	if err != nil {
		return nil, fmt.Errorf("failed to get env data. error: %w", err)
	}
	EnvData := []*gasket_model.Full_EnvData{}
	for _, edd := range envData {
		EnvData = append(EnvData, &gasket_model.Full_EnvData{
			Id:           edd.Id,
			GasketId:     edd.GasketId,
			EnvId:        edd.EnvId,
			M:            edd.M,
			SpecificPres: edd.SpecificPres,
		})
	}

	envType, err := s.repo.GetEnv(ctx, &gasket_api.GetEnvRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type env. error: %w", err)
	}
	EnvType := []*gasket_model.Env{}
	for _, tgd := range envType {
		EnvType = append(EnvType, &gasket_model.Env{
			Id:    tgd.Id,
			Title: tgd.Title,
		})
	}

	data := gasket_api.FullDataResponse{
		GasketData: GasketData,
		GasketType: GasketType,
		EnvData:    EnvData,
		EnvType:    EnvType,
	}

	return &data, nil
}

func (s *GasketService) GetGasket(ctx context.Context, req *gasket_api.GetGasketRequest) (gasket []*gasket_model.Gasket, err error) {
	data, err := s.repo.GetGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasket. error: %w", err)
	}

	for _, item := range data {
		gasket = append(gasket, &gasket_model.Gasket{
			Id:    item.Id,
			Title: item.Title,
		})
	}

	return gasket, nil
}

func (s *GasketService) GetGasketWithThick(ctx context.Context, req *gasket_api.GetGasketRequest) (gasket []*gasket_model.GasketWithThick, err error) {
	data, err := s.repo.GetGasketWithThick(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasket. error: %w", err)
	}

	curId := ""
	for _, item := range data {
		item.Thickness = math.Round(item.Thickness*1000) / 1000
		if item.Id != curId {
			curId = item.Id
			gasket = append(gasket, &gasket_model.GasketWithThick{
				Id:    item.Id,
				Title: item.Title,
			})
			gasket[len(gasket)-1].Thickness = append(gasket[len(gasket)-1].Thickness, item.Thickness)
		} else {
			gasket[len(gasket)-1].Thickness = append(gasket[len(gasket)-1].Thickness, item.Thickness)
		}
	}

	return gasket, nil
}

func (s *GasketService) CreateGasket(ctx context.Context, gasket *gasket_api.CreateGasketRequest) (id string, err error) {
	id, err = s.repo.CreateGasket(ctx, gasket)
	if err != nil {
		return "", fmt.Errorf("failed to create gasket. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateGasket(ctx context.Context, gasket *gasket_api.UpdateGasketRequest) error {
	if err := s.repo.UpdateGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to update gasket. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasket(ctx context.Context, gasket *gasket_api.DeleteGasketRequest) error {
	if err := s.repo.DeleteGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete gasket. error: %w", err)
	}
	return nil
}

//---
func (s *GasketService) CreateManyGasketData(ctx context.Context, data *gasket_api.CreateManyGasketDataRequest) error {
	if err := s.repo.CreateManyGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to create many gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) CreateGasketData(ctx context.Context, data *gasket_api.CreateGasketDataRequest) error {
	if err := s.repo.CreateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to create gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) UpdateGasketData(ctx context.Context, data *gasket_api.UpdateGasketDataRequest) error {
	if err := s.repo.UpdateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to update gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) UpdateGasketTypeId(ctx context.Context, data *gasket_api.UpdateGasketTypeIdRequest) error {
	if err := s.repo.UpdateGasketTypeId(ctx, data); err != nil {
		return fmt.Errorf("failed to update gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasketData(ctx context.Context, data *gasket_api.DeleteGasketDataRequest) error {
	if err := s.repo.DeleteGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to delete gasket data. error: %w", err)
	}
	return nil
}
