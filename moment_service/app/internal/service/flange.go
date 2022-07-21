package service

import (
	"github.com/Alexander272/sealur/moment_service/internal/repository"
)

type FlangeService struct {
	repo repository.Flange
}

func NewFlangeService(repo repository.Flange) *FlangeService {
	return &FlangeService{repo: repo}

}
