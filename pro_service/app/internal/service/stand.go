package service

import (
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
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
		logger.Error("failed to get stands. error: %w", err)
		return stands, err
	}
	return stands, nil
}

func (s *StandService) Create(stand *proto.CreateStandRequest) (st *proto.Id, err error) {
	id, err := s.repo.Create(stand)
	if err != nil {
		logger.Error("failed to create stand. error: %w", err)
		return nil, err
	}
	return &proto.Id{Id: id}, nil
}

func (s *StandService) Update(stand *proto.UpdateStandRequest) error {
	if err := s.repo.Update(stand); err != nil {
		logger.Error("failed to update stand. error: %w", err)
		return err
	}
	return nil
}

func (s *StandService) Delete(stand *proto.DeleteStandRequest) error {
	if err := s.repo.Delete(stand); err != nil {
		logger.Error("failed to delete stand. error: %w", err)
		return err
	}
	return nil
}
