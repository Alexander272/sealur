package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type FlangeService struct {
	repo repository.Flange
}

func NewFlangeService(repo repository.Flange) *FlangeService {
	return &FlangeService{repo: repo}
}

func (s *FlangeService) GetAll() (flanges []*pro_api.Flange, err error) {
	flanges, err = s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get flanges. error: %w", err)
	}
	return flanges, nil
}

func (s *FlangeService) Create(flange *pro_api.CreateFlangeRequest) (fl *pro_api.IdResponse, err error) {
	candidate, err := s.repo.GetByTitle(flange.Title, flange.Short)
	if err != nil {
		return nil, fmt.Errorf("failed to get flange by title. error: %w", err)
	}
	if candidate != nil {
		return nil, models.ErrFlangeAlreadyExists
	}

	id, err := s.repo.Create(flange)
	if err != nil {
		return nil, fmt.Errorf("failed to create flange. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *FlangeService) Update(fl *pro_api.UpdateFlangeRequest) error {
	if err := s.repo.Update(fl); err != nil {
		return fmt.Errorf("failed to update flange. error: %w", err)
	}
	return nil
}

func (s *FlangeService) Delete(fl *pro_api.DeleteFlangeRequest) error {
	if err := s.repo.Delete(fl); err != nil {
		return fmt.Errorf("failed to delete flange. error: %w", err)
	}
	return nil
}
