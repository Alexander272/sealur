package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type BoltMatService struct {
	repo repository.BoltMaterials
}

func NewBoltMatRepo(repo repository.BoltMaterials) *BoltMatService {
	return &BoltMatService{repo: repo}
}

func (s *BoltMatService) GetAll(req *proto.GetBoltMaterialsRequest) (mats []*proto.BoltMaterials, err error) {
	var data []models.BoltMaterials
	data, err = s.repo.GetAll(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get bolt materials. error: %w", err)
	}

	for _, d := range data {
		s := proto.BoltMaterials(d)
		mats = append(mats, &s)
	}

	return mats, nil
}

func (s *BoltMatService) Create(mat *proto.CreateBoltMaterialsRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(mat)
	if err != nil {
		return nil, fmt.Errorf("failed to create bold material. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *BoltMatService) Update(mat *proto.UpdateBoltMaterialsRequest) error {
	if err := s.repo.Update(mat); err != nil {
		return fmt.Errorf("failed to update bolt material. error: %w", err)
	}
	return nil
}

func (s *BoltMatService) Delete(mat *proto.DeleteBoltMaterialsRequest) error {
	if err := s.repo.Delete(mat); err != nil {
		return fmt.Errorf("failed to delete bolt material. error: %w", err)
	}
	return nil
}
