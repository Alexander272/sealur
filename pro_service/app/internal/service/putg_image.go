package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type PutgImageService struct {
	repo repository.PutgImage
}

func NewPutgImageService(repo repository.PutgImage) *PutgImageService {
	return &PutgImageService{repo: repo}
}

func (s *PutgImageService) Get(req *pro_api.GetPutgImageRequest) (images []*pro_api.PutgImage, err error) {
	images, err = s.repo.Get(req)
	if err != nil {
		return images, fmt.Errorf("failed to get putg image. error: %w", err)
	}
	return images, nil
}

func (s *PutgImageService) Create(image *pro_api.CreatePutgImageRequest) (*pro_api.IdResponse, error) {
	id, err := s.repo.Create(image)
	if err != nil {
		return nil, fmt.Errorf("failed to create putg image. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *PutgImageService) Update(image *pro_api.UpdatePutgImageRequest) error {
	if err := s.repo.Update(image); err != nil {
		return fmt.Errorf("failed to update putg image. error: %w", err)
	}
	return nil
}

func (s *PutgImageService) Delete(image *pro_api.DeletePutgImageRequest) error {
	if err := s.repo.Delete(image); err != nil {
		return fmt.Errorf("failed to delete putg image. error: %w", err)
	}
	return nil
}
