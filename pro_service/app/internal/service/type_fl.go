package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type TypeFLService struct {
	repo repository.TypeFl
}

func NewTypeFlService(repo repository.TypeFl) *TypeFLService {
	return &TypeFLService{repo: repo}
}

func (s *TypeFLService) Get() (fl []*pro_api.TypeFl, err error) {
	fl, err = s.repo.Get()
	if err != nil {
		return fl, fmt.Errorf("failed to get type flange. error: %w", err)
	}
	return fl, nil
}

func (s *TypeFLService) GetAll() (fl []*pro_api.TypeFl, err error) {
	fl, err = s.repo.GetAll()
	if err != nil {
		return fl, fmt.Errorf("failed to get all types flange. error: %w", err)
	}
	return fl, nil
}

func (s *TypeFLService) Create(fl *pro_api.CreateTypeFlRequest) (*pro_api.IdResponse, error) {
	id, err := s.repo.Create(fl)
	if err != nil {
		return nil, fmt.Errorf("failed to create type flange. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *TypeFLService) Update(fl *pro_api.UpdateTypeFlRequest) error {
	if err := s.repo.Update(fl); err != nil {
		return fmt.Errorf("failed to update type flange. error: %w", err)
	}
	return nil
}

func (s *TypeFLService) Delete(fl *pro_api.DeleteTypeFlRequest) error {
	if err := s.repo.Delete(fl); err != nil {
		return fmt.Errorf("failed to delete type flange. error: %w", err)
	}
	return nil
}
