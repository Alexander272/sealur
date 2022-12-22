package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type MatSerive struct {
	repo repository.Materials
}

func NewMatService(repo repository.Materials) *MatSerive {
	return &MatSerive{repo: repo}
}

func (s *MatSerive) GetAll(req *pro_api.GetMaterialsRequest) (mats []*pro_api.Materials, err error) {
	var data []models.Materials
	data, err = s.repo.GetAll(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get materials. error: %w", err)
	}

	for _, d := range data {
		mats = append(mats, &pro_api.Materials{
			Id:      d.Id,
			Title:   d.Title,
			TypeMat: d.TypeMat,
		})
	}

	return mats, nil
}

func (s *MatSerive) Create(mat *pro_api.CreateMaterialsRequest) (*pro_api.IdResponse, error) {
	id, err := s.repo.Create(mat)
	if err != nil {
		return nil, fmt.Errorf("failed to create material. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *MatSerive) Update(mat *pro_api.UpdateMaterialsRequest) error {
	if err := s.repo.Update(mat); err != nil {
		return fmt.Errorf("failed to update material. error: %w", err)
	}
	return nil
}

func (s *MatSerive) Delete(mat *pro_api.DeleteMaterialsRequest) error {
	if err := s.repo.Delete(mat); err != nil {
		return fmt.Errorf("failed to delete material. error: %w", err)
	}
	return nil
}
