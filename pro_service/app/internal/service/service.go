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

type StFl interface {
	Get() ([]*proto.StFl, error)
	Create(*proto.CreateStFlRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateStFlRequest) error
	Delete(*proto.DeleteStFlRequest) error
}

type TypeFl interface {
	Get() ([]*proto.TypeFl, error)
	GetAll() ([]*proto.TypeFl, error)
	Create(*proto.CreateTypeFlRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateTypeFlRequest) error
	Delete(*proto.DeleteTypeFlRequest) error
}

type Addit interface {
	GetAll() ([]*proto.Additional, error)
	Create(*proto.CreateAddRequest) (*proto.SuccessResponse, error)
	UpdateMat(*proto.UpdateAddMatRequest) (*proto.SuccessResponse, error)
	UpdateMod(*proto.UpdateAddModRequest) (*proto.SuccessResponse, error)
	UpdateTemp(*proto.UpdateAddTemRequest) (*proto.SuccessResponse, error)
	UpdateMoun(*proto.UpdateAddMounRequest) (*proto.SuccessResponse, error)
	UpdateGrap(*proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error)
	UpdateFillers(*proto.UpdateAddFillersRequest) (*proto.SuccessResponse, error)
	UpdateCoating(*proto.UpdateAddCoatingRequest) (*proto.SuccessResponse, error)
	UpdateConstruction(*proto.UpdateAddConstructionRequest) (*proto.SuccessResponse, error)
	UpdateObturator(*proto.UpdateAddObturatorRequest) (*proto.SuccessResponse, error)
}

type Size interface {
	Get(*proto.GetSizesRequest) ([]*proto.Size, []*proto.Dn, error)
	Create(*proto.CreateSizeRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateSizeRequest) error
	Delete(*proto.DeleteSizeRequest) error
	DeleteAll(*proto.DeleteAllSizeRequest) error
}

type SNP interface {
	Get(*proto.GetSNPRequest) ([]*proto.SNP, error)
	Create(*proto.CreateSNPRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateSNPRequest) error
	Delete(*proto.DeleteSNPRequest) error

	AddMat(id string) error
	DeleteMat(id string, materials []*proto.AddMaterials) error
	AddMoun(id string) error
	DeleteMoun(id string) error
	AddGrap(id string) error
	DeleteGrap(id string) error
	DeleteFiller(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
}

type PutgImage interface {
	Get(req *proto.GetPutgImageRequest) ([]*proto.PutgImage, error)
	Create(image *proto.CreatePutgImageRequest) (*proto.IdResponse, error)
	Update(image *proto.UpdatePutgImageRequest) error
	Delete(image *proto.DeletePutgImageRequest) error
}

type Services struct {
	Stand
	Flange
	StFl
	TypeFl
	Addit
	Size
	SNP
	PutgImage
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Stand:     NewStandService(repos.Stand),
		Flange:    NewFlangeService(repos.Flange),
		StFl:      NewStFlService(repos.StFl),
		TypeFl:    NewTypeFlService(repos.TypeFl),
		Addit:     NewAdditService(repos.Addit),
		Size:      NewSizeService(repos.Size),
		SNP:       NewSNPService(repos.SNP),
		PutgImage: NewPutgImageService(repos.PutgImage),
	}
}
