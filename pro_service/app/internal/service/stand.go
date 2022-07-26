package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type StandService struct {
	repo repository.Stand
}

func NewStandService(repo repository.Stand) *StandService {
	return &StandService{repo: repo}
}

func (s *StandService) GetAll(req *pro_api.GetStandsRequest) (stands []*pro_api.Stand, err error) {
	stands, err = s.repo.GetAll(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get stands. error: %w", err)
	}
	return stands, nil
}

func (s *StandService) Create(stand *pro_api.CreateStandRequest) (st *pro_api.IdResponse, err error) {
	candidate, err := s.repo.GetByTitle(stand.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to get stand by title. error: %w", err)
	}
	if candidate != nil {
		return nil, models.ErrStandAlreadyExists
	}

	id, err := s.repo.Create(stand)
	if err != nil {
		return nil, fmt.Errorf("failed to create stand. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *StandService) Update(stand *pro_api.UpdateStandRequest) error {
	if err := s.repo.Update(stand); err != nil {
		return fmt.Errorf("failed to update stand. error: %w", err)
	}
	return nil
}

func (s *StandService) Delete(stand *pro_api.DeleteStandRequest) error {
	if err := s.repo.Delete(stand); err != nil {
		return fmt.Errorf("failed to delete stand. error: %w", err)
	}
	return nil
}
