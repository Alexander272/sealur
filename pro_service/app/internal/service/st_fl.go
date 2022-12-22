package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type StFlService struct {
	repo repository.StFl
}

func NewStFlService(repo repository.StFl) *StFlService {
	return &StFlService{repo: repo}
}

func (s *StFlService) Get() (st []*pro_api.StFl, err error) {
	st, err = s.repo.Get()
	if err != nil {
		return st, fmt.Errorf("failed to get st/fl. error: %w", err)
	}
	return st, nil
}

func (s *StFlService) Create(st *pro_api.CreateStFlRequest) (*pro_api.IdResponse, error) {
	id, err := s.repo.Create(st)
	if err != nil {
		return nil, fmt.Errorf("failed to create st/fl. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *StFlService) Update(st *pro_api.UpdateStFlRequest) error {
	if err := s.repo.Update(st); err != nil {
		return fmt.Errorf("failed to update st/fl. error: %w", err)
	}
	return nil
}

func (s *StFlService) Delete(st *pro_api.DeleteStFlRequest) error {
	if err := s.repo.Delete(st); err != nil {
		return fmt.Errorf("failed to delete st/fl. error: %w", err)
	}
	return nil
}
