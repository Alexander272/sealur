package service

import "github.com/Alexander272/sealur/pro_service/internal/repository"

type Services struct {
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{}
}
