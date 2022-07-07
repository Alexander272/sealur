package service

import "github.com/Alexander272/sealur/moment_service/internal/repository"

type MaterialsService struct {
	repo repository.Materials
}

func NewMaterialsService(repo repository.Materials) *MaterialsService {
	return &MaterialsService{repo: repo}
}
