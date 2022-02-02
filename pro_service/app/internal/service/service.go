package service

import (
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type Stand interface {
	GetAll(req *proto.GetStandsRequest) (stands []*proto.Stand, err error)
	Create(stand *proto.CreateStandRequest) (st *proto.Id, err error)
	Update(stand *proto.UpdateStandRequest) error
	Delete(stand *proto.DeleteStandRequest) error
}

type Services struct {
	Stand
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Stand: NewStandService(repos.Stand),
	}
}
