package service

import (
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type Stand interface {
	GetAll(req proto.GetStands) (stands []proto.Stand, err error)
	Create(stand proto.CreateStand) (st proto.Id, err error)
	Update(stand proto.UpdateStand) error
	Delete(stand proto.DeleteStand) error
}

type Services struct {
	Stand
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Stand: NewStandService(repos.Stand),
	}
}
