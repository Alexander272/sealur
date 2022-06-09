package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type SizeIntService struct {
	repo repository.SizeInt
}

func NewSizeIntService(repo repository.SizeInt) *SizeIntService {
	return &SizeIntService{repo: repo}
}

func (s *SizeIntService) Get(req *proto.GetSizesIntRequest) (sizes []*proto.SizeInt, err error) {
	var data []models.SizeInterview
	data, err = s.repo.Get(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes interview. error: %w", err)
	}

	for _, d := range data {
		s := proto.SizeInt{
			Id:        d.Id,
			Dy:        d.Dy,
			Py:        d.Py,
			D1:        d.D1,
			D2:        d.D2,
			DUp:       d.DUp,
			D:         d.D,
			H1:        d.H1,
			H2:        d.H2,
			Bolt:      d.Bolt,
			CountBolt: d.Count,
		}
		sizes = append(sizes, &s)
	}

	return sizes, nil
}

func (s *SizeIntService) Create(size *proto.CreateSizeIntRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create size interview. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *SizeIntService) Update(size *proto.UpdateSizeIntRequest) error {
	if err := s.repo.Update(size); err != nil {
		return fmt.Errorf("failed to update size interview. error: %w", err)
	}
	return nil
}

func (s *SizeIntService) Delete(size *proto.DeleteSizeIntRequest) error {
	if err := s.repo.Delete(size); err != nil {
		return fmt.Errorf("failed to delete size interview. error: %w", err)
	}
	return nil
}
