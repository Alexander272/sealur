package service

import "github.com/Alexander272/sealur/moment_service/internal/repository"

type Flange interface{}

type Materials interface{}

type Services struct {
	Flange
	Materials
}

func NewServices(repos *repository.Repositories) *Services {
	Materials := NewMaterialsService(repos.Materials)

	return &Services{
		Flange:    NewFlangeService(repos.Flange, Materials),
		Materials: Materials,
	}
}
