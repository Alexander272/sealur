package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type StandService struct {
	repo repository.Stand
}

func NewStandService(repo repository.Stand) *StandService {
	return &StandService{repo: repo}
}

func (s *StandService) GetAll(req *proto.GetStandsRequest) (stands []*proto.Stand, err error) {
	stands, err = s.repo.GetAll(req)
	if err != nil {
		return stands, fmt.Errorf("failed to get stands. error: %w", err)
	}
	return stands, nil
}

func (s *StandService) Create(stand *proto.CreateStandRequest) (st *proto.IdResponse, err error) {
	id, err := s.repo.Create(stand)
	if err != nil {
		return nil, fmt.Errorf("failed to create stand. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *StandService) Update(stand *proto.UpdateStandRequest) error {
	if err := s.repo.Update(stand); err != nil {
		return fmt.Errorf("failed to update stand. error: %w", err)
	}
	return nil
}

func (s *StandService) Delete(stand *proto.DeleteStandRequest) error {
	if err := s.repo.Delete(stand); err != nil {
		return fmt.Errorf("failed to delete stand. error: %w", err)
	}
	return nil
}
