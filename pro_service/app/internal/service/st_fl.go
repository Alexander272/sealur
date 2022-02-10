package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type StFlService struct {
	repo repository.StFl
}

func NewStFlService(repo repository.StFl) *StFlService {
	return &StFlService{repo: repo}
}

func (s *StFlService) Get() (st []*proto.StFl, err error) {
	st, err = s.repo.Get()
	if err != nil {
		return st, fmt.Errorf("failed to get st/fl. error: %w", err)
	}
	return st, nil
}

func (s *StFlService) Create(st *proto.CreateStFlRequest) error {
	if err := s.repo.Create(st); err != nil {
		return fmt.Errorf("failed to create st/fl. error: %w", err)
	}
	return nil
}

func (s *StFlService) Update(st *proto.UpdateStFlRequest) error {
	if err := s.repo.Update(st); err != nil {
		return fmt.Errorf("failed to update st/fl. error: %w", err)
	}
	return nil
}

func (s *StFlService) Delete(st *proto.DeleteStFlRequest) error {
	if err := s.repo.Delete(st); err != nil {
		return fmt.Errorf("failed to delete st/fl. error: %w", err)
	}
	return nil
}
