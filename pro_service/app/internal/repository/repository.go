package repository

import (
	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type Stand interface {
	GetAll(stand *pro_api.GetStandsRequest) ([]*pro_api.Stand, error)
	GetByTitle(title string) ([]*pro_api.Stand, error)
	Create(stand *pro_api.CreateStandRequest) (id string, err error)
	Update(stand *pro_api.UpdateStandRequest) error
	Delete(stand *pro_api.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*pro_api.Flange, error)
	GetByTitle(title, short string) ([]*pro_api.Flange, error)
	Create(*pro_api.CreateFlangeRequest) (id string, err error)
	Update(*pro_api.UpdateFlangeRequest) error
	Delete(*pro_api.DeleteFlangeRequest) error
}

type StFl interface {
	Get() ([]*pro_api.StFl, error)
	Create(*pro_api.CreateStFlRequest) (string, error)
	Update(*pro_api.UpdateStFlRequest) error
	Delete(*pro_api.DeleteStFlRequest) error
}

type TypeFl interface {
	Get() ([]*pro_api.TypeFl, error)
	GetAll() ([]*pro_api.TypeFl, error)
	Create(*pro_api.CreateTypeFlRequest) (string, error)
	Update(*pro_api.UpdateTypeFlRequest) error
	Delete(*pro_api.DeleteTypeFlRequest) error
}

type Addit interface {
	GetAll() ([]models.Addit, error)
	Create(*pro_api.CreateAddRequest) error
	UpdateMat(models.UpdateMat) error
	UpdateMod(models.UpdateMod) error
	UpdateTemp(models.UpdateTemp) error
	UpdateMoun(models.UpdateMoun) error
	UpdateGrap(models.UpdateGrap) error
	UpdateFillers(models.UpdateFill) error
	UpdateCoating(models.UpdateCoating) error
	UpdateConstruction(models.UpdateConstr) error
	UpdateObturator(models.UpdateObturator) error
	UpdateBasis(models.UpdateBasis) error
	UpdatePObturator(models.UpdatePObturator) error
	UpdateSealant(models.UpdateSealant) error
}

type Size interface {
	Get(req *pro_api.GetSizesRequest) ([]*pro_api.Size, error)
	GetAll(req *pro_api.GetSizesRequest) ([]*pro_api.Size, error)
	Create(size *pro_api.CreateSizeRequest) (id string, err error)
	Update(size *pro_api.UpdateSizeRequest) error
	Delete(size *pro_api.DeleteSizeRequest) error
	DeleteAll(size *pro_api.DeleteAllSizeRequest) error
}

type SNP interface {
	Get(req *pro_api.GetSNPRequest) ([]models.SNP, error)
	Create(snp models.SnpDTO) (id string, err error)
	Update(snp models.SnpDTO) error
	Delete(snp *pro_api.DeleteSNPRequest) error

	GetByCondition(cond string) ([]models.SNP, error)
	UpdateAddit(snp models.UpdateAdditDTO) error
}

type PutgImage interface {
	Get(req *pro_api.GetPutgImageRequest) ([]*pro_api.PutgImage, error)
	Create(image *pro_api.CreatePutgImageRequest) (id string, err error)
	Update(image *pro_api.UpdatePutgImageRequest) error
	Delete(image *pro_api.DeletePutgImageRequest) error
}

type Putg interface {
	Get(req *pro_api.GetPutgRequest) ([]models.Putg, error)
	Create(putg models.PutgDTO) (id string, err error)
	Update(putg models.PutgDTO) error
	Delete(putg *pro_api.DeletePutgRequest) error

	GetByCondition(cond string) ([]models.Putg, error)
	UpdateAddit(putg models.UpdateAdditDTO) error
}

type PutgmImage interface {
	Get(req *pro_api.GetPutgmImageRequest) ([]*pro_api.PutgmImage, error)
	Create(image *pro_api.CreatePutgmImageRequest) (id string, err error)
	Update(image *pro_api.UpdatePutgmImageRequest) error
	Delete(image *pro_api.DeletePutgmImageRequest) error
}

type Putgm interface {
	Get(req *pro_api.GetPutgmRequest) ([]models.Putgm, error)
	Create(putg models.PutgmDTO) (id string, err error)
	Update(putg models.PutgmDTO) error
	Delete(putg *pro_api.DeletePutgmRequest) error

	GetByCondition(cond string) ([]models.Putgm, error)
	UpdateAddit(putgm models.UpdateAdditDTO) error
}

type Materials interface {
	GetAll(*pro_api.GetMaterialsRequest) ([]models.Materials, error)
	Create(*pro_api.CreateMaterialsRequest) (string, error)
	Update(*pro_api.UpdateMaterialsRequest) error
	Delete(*pro_api.DeleteMaterialsRequest) error
}

type BoltMaterials interface {
	GetAll(*pro_api.GetBoltMaterialsRequest) ([]models.BoltMaterials, error)
	Create(*pro_api.CreateBoltMaterialsRequest) (string, error)
	Update(*pro_api.UpdateBoltMaterialsRequest) error
	Delete(*pro_api.DeleteBoltMaterialsRequest) error
}

type SizeInt interface {
	Get(*pro_api.GetSizesIntRequest) ([]models.SizeInterview, error)
	GetAll(*pro_api.GetAllSizeIntRequest) ([]models.SizeInterview, error)
	Create(*pro_api.CreateSizeIntRequest) (id string, err error)
	Update(*pro_api.UpdateSizeIntRequest) error
	Delete(*pro_api.DeleteSizeIntRequest) error
	DeleteAll(*pro_api.DeleteAllSizeIntRequest) error
}

type Order interface {
	GetAll(*pro_api.GetAllOrdersRequest) ([]models.Order, error)
	GetCur(req *pro_api.GetCurOrderRequest) (order models.Order, err error)
	Create(*pro_api.CreateOrderRequest) error
	Copy(*pro_api.CopyOrderRequest) error
	Delete(*pro_api.DeleteOrderRequest) error
	Save(*pro_api.SaveOrderRequest) error
	GetPositions(*pro_api.GetPositionsRequest) ([]models.Position, error)
}

type OrderPosition interface {
	Get(*pro_api.GetPositionsRequest) ([]models.Position, error)
	GetCur(*pro_api.GetCurPositionsRequest) ([]models.Position, error)
	Add(*pro_api.AddPositionRequest) (id string, err error)
	Update(*pro_api.UpdatePositionRequest) error
	Remove(*pro_api.RemovePositionRequest) (string, error)
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
	PutgmImage
	Putgm
	Materials
	BoltMaterials
	SizeInt
	Order
	OrderPosition
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Stand:         NewStandRepo(db),
		Flange:        NewFlangeRepo(db),
		StFl:          NewStFlRepo(db),
		TypeFl:        NewTypeFlRepo(db),
		Addit:         NewAdditRepo(db),
		Size:          NewSizesRepo(db),
		SNP:           NewSNPRepo(db),
		PutgImage:     NewPutgImageRepo(db),
		Putg:          NewPutgRepo(db),
		PutgmImage:    NewPutgmImageRepo(db),
		Putgm:         NewPutgmRepo(db),
		Materials:     NewMatRepo(db),
		BoltMaterials: NewBoltMatRepo(db),
		SizeInt:       NewSizeIntRepo(db),
		Order:         NewOrderRepo(db),
		OrderPosition: NewPositionRepo(db),
	}
}
