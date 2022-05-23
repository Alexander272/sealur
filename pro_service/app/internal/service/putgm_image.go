package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type PutgmImageService struct {
	repo repository.PutgmImage
}

func NewPutgmImageService(repo repository.PutgmImage) *PutgmImageService {
	return &PutgmImageService{repo: repo}
}

func (s *PutgmImageService) Get(req *proto.GetPutgmImageRequest) (images []*proto.PutgmImage, err error) {
	images, err = s.repo.Get(req)
	if err != nil {
		return images, fmt.Errorf("failed to get putg image. error: %w", err)
	}
	return images, nil
}

func (s *PutgmImageService) Create(image *proto.CreatePutgmImageRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(image)
	if err != nil {
		return nil, fmt.Errorf("failed to create putg image. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *PutgmImageService) Update(image *proto.UpdatePutgmImageRequest) error {
	if err := s.repo.Update(image); err != nil {
		return fmt.Errorf("failed to update putg image. error: %w", err)
	}
	return nil
}

func (s *PutgmImageService) Delete(image *proto.DeletePutgmImageRequest) error {
	if err := s.repo.Delete(image); err != nil {
		return fmt.Errorf("failed to delete putg image. error: %w", err)
	}
	return nil
}
