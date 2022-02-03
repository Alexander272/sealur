package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type FlangeService struct {
	repo repository.Flange
}

func NewFlangeService(repo repository.Flange) *FlangeService {
	return &FlangeService{repo: repo}
}

func (s *FlangeService) GetAll() (flanges []*proto.Flange, err error) {
	flanges, err = s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get flanges. error: %w", err)
	}
	return flanges, nil
}

func (s *FlangeService) Create(flange *proto.CreateFlangeRequest) (fl *proto.IdResponse, err error) {
	id, err := s.repo.Create(flange)
	if err != nil {
		return nil, fmt.Errorf("failed to create flange. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *FlangeService) Update(fl *proto.UpdateFlangeRequest) error {
	if err := s.repo.Update(fl); err != nil {
		return fmt.Errorf("failed to update flange. error: %w", err)
	}
	return nil
}

func (s *FlangeService) Delete(fl *proto.DeleteFlangeRequest) error {
	if err := s.repo.Delete(fl); err != nil {
		return fmt.Errorf("failed to delete flange. error: %w", err)
	}
	return nil
}
