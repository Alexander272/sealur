package repository

import (
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Stand interface {
	GetAll(stand *proto.GetStandsRequest) ([]*proto.Stand, error)
	GetByTitle(title string) (*proto.Stand, error)
	Create(stand *proto.CreateStandRequest) (id string, err error)
	Update(stand *proto.UpdateStandRequest) error
	Delete(stand *proto.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*proto.Flange, error)
	GetByTitle(title, short string) (*proto.Flange, error)
	Create(*proto.CreateFlangeRequest) (id string, err error)
	Update(*proto.UpdateFlangeRequest) error
	Delete(*proto.DeleteFlangeRequest) error
}

type Addit interface {
	GetAll() ([]*proto.Additional, error)
	Create(*proto.CreateAddRequest) error
	UpdateMat(*proto.UpdateAddMatRequest) error
	UpdateMod(*proto.UpdateAddModRequest) error
	UpdateTemp(*proto.UpdateAddTemRequest) error
	UpdateMoun(*proto.UpdateAddMounRequest) error
	UpdateGrap(*proto.UpdateAddGrapRequest) error
	UpdateTypeFl(*proto.UpdateAddTypeFlRequest) error
}

type Size interface {
	Get(req *proto.GetSizesRequest) ([]*proto.Size, error)
	Create(size *proto.CreateSizeRequest) (id string, err error)
	Update(size *proto.UpdateSizeRequest) error
	Delete(size *proto.DeleteSizeRequest) error
}

type SNP interface {
	Get(req *proto.GetSNPRequest) ([]*proto.SNP, error)
	Create(snp *proto.CreateSNPRequest) (id string, err error)
	Update(snp *proto.UpdateSNPRequest) error
	Delete(snp *proto.DeleteSNPRequest) error
}

type Repositories struct {
	Stand
	Flange
	Addit
	Size
	SNP
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Stand:  NewStandRepo(db),
		Flange: NewFlangeRepo(db),
		Addit:  NewAdditRepo(db),
		Size:   NewSizesRepo(db),
		SNP:    NewSNPRepo(db),
	}
}
