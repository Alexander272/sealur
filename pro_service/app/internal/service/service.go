package service

import (
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type Stand interface {
	GetAll(*proto.GetStandsRequest) (stands []*proto.Stand, err error)
	Create(*proto.CreateStandRequest) (st *proto.IdResponse, err error)
	Update(*proto.UpdateStandRequest) error
	Delete(*proto.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*proto.Flange, error)
	Create(*proto.CreateFlangeRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateFlangeRequest) error
	Delete(*proto.DeleteFlangeRequest) error
}

type Addit interface {
	GetAll() ([]*proto.Additional, error)
	Create(*proto.CreateAddRequest) (*proto.SuccessResponse, error)
	UpdateMat(*proto.UpdateAddMatRequest) (*proto.SuccessResponse, error)
	UpdateMod(*proto.UpdateAddModRequest) (*proto.SuccessResponse, error)
	UpdateTemp(*proto.UpdateAddTemRequest) (*proto.SuccessResponse, error)
	UpdateMoun(*proto.UpdateAddMounRequest) (*proto.SuccessResponse, error)
	UpdateGrap(*proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error)
}

type Size interface {
	Get(*proto.GetSizesRequest) ([]*proto.Size, error)
	Create(*proto.CreateSizeRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateSizeRequest) error
	Delete(*proto.DeleteSizeRequest) error
}

type SNP interface {
	Get(*proto.GetSNPRequest) ([]*proto.SNP, error)
	Create(*proto.CreateSNPRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateSNPRequest) error
	Delete(*proto.DeleteSNPRequest) error
}

type Services struct {
	Stand
	Flange
	Addit
	Size
	SNP
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Stand:  NewStandService(repos.Stand),
		Flange: NewFlangeService(repos.Flange),
		Addit:  NewAdditService(repos.Addit),
		Size:   NewSizeService(repos.Size),
		SNP:    NewSNPService(repos.SNP),
	}
}
