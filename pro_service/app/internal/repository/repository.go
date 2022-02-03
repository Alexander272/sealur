package repository

import (
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Stand interface {
	GetAll(stand *proto.GetStandsRequest) ([]*proto.Stand, error)
	Create(stand *proto.CreateStandRequest) (id string, err error)
	Update(stand *proto.UpdateStandRequest) error
	Delete(stand *proto.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*proto.Flange, error)
	Create(fl *proto.CreateFlangeRequest) (id string, err error)
	Update(fl *proto.UpdateFlangeRequest) error
	Delete(fl *proto.DeleteFlangeRequest) error
}

type Addit interface {
	GetAll() ([]*proto.Additional, error)
	Create(add *proto.CreateAddRequest) error
	UpdateMat(mat *proto.UpdateAddMatRequest) error
	UpdateMod(mod *proto.UpdateAddModRequest) error
	UpdateTemp(temp *proto.UpdateAddTemRequest) error
	UpdateMoun(moun *proto.UpdateAddMounRequest) error
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
