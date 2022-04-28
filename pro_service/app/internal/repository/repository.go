package repository

import (
	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Stand interface {
	GetAll(stand *proto.GetStandsRequest) ([]*proto.Stand, error)
	GetByTitle(title string) ([]*proto.Stand, error)
	Create(stand *proto.CreateStandRequest) (id string, err error)
	Update(stand *proto.UpdateStandRequest) error
	Delete(stand *proto.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*proto.Flange, error)
	GetByTitle(title, short string) ([]*proto.Flange, error)
	Create(*proto.CreateFlangeRequest) (id string, err error)
	Update(*proto.UpdateFlangeRequest) error
	Delete(*proto.DeleteFlangeRequest) error
}

type StFl interface {
	Get() ([]*proto.StFl, error)
	Create(*proto.CreateStFlRequest) (string, error)
	Update(*proto.UpdateStFlRequest) error
	Delete(*proto.DeleteStFlRequest) error
}

type TypeFl interface {
	Get() ([]*proto.TypeFl, error)
	GetAll() ([]*proto.TypeFl, error)
	Create(*proto.CreateTypeFlRequest) (string, error)
	Update(*proto.UpdateTypeFlRequest) error
	Delete(*proto.DeleteTypeFlRequest) error
}

type Addit interface {
	GetAll() ([]models.Addit, error)
	Create(*proto.CreateAddRequest) error
	UpdateMat(models.UpdateMat) error
	UpdateMod(models.UpdateMod) error
	UpdateTemp(models.UpdateTemp) error
	UpdateMoun(models.UpdateMoun) error
	UpdateGrap(models.UpdateGrap) error
	UpdateFillers(models.UpdateFill) error
	UpdateCoating(models.UpdateCoating) error
	UpdateConstruction(models.UpdateConstr) error
	UpdateObturator(models.UpdateObturator) error
}

type Size interface {
	Get(req *proto.GetSizesRequest) ([]*proto.Size, error)
	Create(size *proto.CreateSizeRequest) (id string, err error)
	Update(size *proto.UpdateSizeRequest) error
	Delete(size *proto.DeleteSizeRequest) error
	DeleteAll(size *proto.DeleteAllSizeRequest) error
}

type SNP interface {
	Get(req *proto.GetSNPRequest) ([]models.SNP, error)
	Create(snp models.SnpDTO) (id string, err error)
	Update(snp models.SnpDTO) error
	Delete(snp *proto.DeleteSNPRequest) error

	GetByCondition(cond string) ([]models.SNP, error)
	UpdateAddit(snp models.UpdateAdditDTO) error
}

type PutgImage interface {
	Get(req *proto.GetPutgImageRequest) ([]*proto.PutgImage, error)
	Create(image *proto.CreatePutgImageRequest) (id string, err error)
	Update(image *proto.UpdatePutgImageRequest) error
	Delete(image *proto.DeletePutgImageRequest) error
}

type Putg interface {
	Get(req *proto.GetPutgRequest) ([]models.Putg, error)
	Create(putg models.PutgDTO) (id string, err error)
	Update(putg models.PutgDTO) error
	Delete(putg *proto.DeletePutgRequest) error
}

type Repositories struct {
	Stand
	Flange
	StFl
	TypeFl
	Addit
	Size
	SNP
	PutgImage
	Putg
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Stand:     NewStandRepo(db),
		Flange:    NewFlangeRepo(db),
		StFl:      NewStFlRepo(db),
		TypeFl:    NewTypeFlRepo(db),
		Addit:     NewAdditRepo(db),
		Size:      NewSizesRepo(db),
		SNP:       NewSNPRepo(db),
		PutgImage: NewPutgImageRepo(db),
		Putg:      NewPutgRepo(db),
	}
}
