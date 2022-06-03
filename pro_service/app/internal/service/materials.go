package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type MatSerive struct {
	repo repository.Materials
}

func NewMatService(repo repository.Materials) *MatSerive {
	return &MatSerive{repo: repo}
}

func (s *MatSerive) GetAll(req *proto.GetMaterialsRequest) (mats []*proto.Materials, err error) {
	var data []models.Materials
	data, err = s.repo.GetAll(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get materials. error: %w", err)
	}

	for _, d := range data {
		s := proto.Materials(d)
		mats = append(mats, &s)
	}

	return mats, nil
}

func (s *MatSerive) Create(mat *proto.CreateMaterialsRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(mat)
	if err != nil {
		return nil, fmt.Errorf("failed to create material. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *MatSerive) Update(mat *proto.UpdateMaterialsRequest) error {
	if err := s.repo.Update(mat); err != nil {
		return fmt.Errorf("failed to update material. error: %w", err)
	}
	return nil
}

func (s *MatSerive) Delete(mat *proto.DeleteMaterialsRequest) error {
	if err := s.repo.Delete(mat); err != nil {
		return fmt.Errorf("failed to delete material. error: %w", err)
	}
	return nil
}
