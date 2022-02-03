package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type SizeService struct {
	repo repository.Size
}

func NewSizeService(repo repository.Size) *SizeService {
	return &SizeService{repo: repo}
}

func (s *SizeService) Get(req *proto.GetSizesRequest) (sizes []*proto.Size, err error) {
	sizes, err = s.repo.Get(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes. error: %w", err)
	}
	return sizes, nil
}

func (s *SizeService) Create(size *proto.CreateSizeRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create size. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *SizeService) Update(size *proto.UpdateSizeRequest) error {
	if err := s.repo.Update(size); err != nil {
		return fmt.Errorf("failed to update query. error: %w", err)
	}
	return nil
}

func (s *SizeService) Delete(size *proto.DeleteSizeRequest) error {
	if err := s.repo.Delete(size); err != nil {
		return fmt.Errorf("failed to delete query. error: %w", err)
	}
	return nil
}
