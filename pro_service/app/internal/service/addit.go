package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type AdditService struct {
	repo repository.Addit
}

func NewAdditService(repo repository.Addit) *AdditService {
	return &AdditService{repo: repo}
}

func (s *AdditService) GetAll() (addit []*proto.Additional, err error) {
	addit, err = s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get additional. error: %w", err)
	}
	return addit, nil
}

func (s *AdditService) Create(addit *proto.CreateAddRequest) (*proto.SuccessResponse, error) {
	err := s.repo.Create(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to create additional. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMat(addit *proto.UpdateAddMatRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateMat(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update materials. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMod(addit *proto.UpdateAddModRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateMod(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update mod. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateTemp(addit *proto.UpdateAddTemRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateTemp(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update temperature. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMoun(addit *proto.UpdateAddMounRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateMoun(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update mounting. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateGrap(addit *proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error) {
	if err := s.repo.UpdateGrap(addit); err != nil {
		return nil, fmt.Errorf("failed to update graphite. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateTypeFl(addit *proto.UpdateAddTypeFlRequest) (*proto.SuccessResponse, error) {
	if err := s.repo.UpdateTypeFl(addit); err != nil {
		return nil, fmt.Errorf("failed to update type_fl. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}
