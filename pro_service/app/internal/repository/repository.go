package repository

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository/postgres"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_snp_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/mounting_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/temperature_model"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
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

//* New

type FlangeStandard interface {
	GetAll(context.Context, *flange_standard_api.GetAllFlangeStandards) ([]*flange_standard_model.FlangeStandard, error)
	Create(context.Context, *flange_standard_api.CreateFlangeStandard) error
	CreateSeveral(context.Context, *flange_standard_api.CreateSeveralFlangeStandard) error
	Update(context.Context, *flange_standard_api.UpdateFlangeStandard) error
	Delete(context.Context, *flange_standard_api.DeleteFlangeStandard) error
}

type Material interface {
	GetAll(context.Context, *material_api.GetAllMaterials) ([]*material_model.Material, error)
	Create(context.Context, *material_api.CreateMaterial) error
	CreateSeveral(context.Context, *material_api.CreateSeveralMaterial) error
	Update(context.Context, *material_api.UpdateMaterial) error
	Delete(context.Context, *material_api.DeleteMaterial) error
}

type Mounting interface {
	GetAll(context.Context, *mounting_api.GetAllMountings) ([]*mounting_model.Mounting, error)
	Create(context.Context, *mounting_api.CreateMounting) error
	CreateSeveral(context.Context, *mounting_api.CreateSeveralMounting) error
	Update(context.Context, *mounting_api.UpdateMounting) error
	Delete(context.Context, *mounting_api.DeleteMounting) error
}

type Standard interface {
	GetAll(context.Context, *standard_api.GetAllStandards) ([]*standard_model.Standard, error)
	GetDefault(context.Context) (*standard_model.Standard, error)
	Create(context.Context, *standard_api.CreateStandard) error
	CreateSeveral(context.Context, *standard_api.CreateSeveralStandard) error
	Update(context.Context, *standard_api.UpdateStandard) error
	Delete(context.Context, *standard_api.DeleteStandard) error
}

type Temperature interface {
	GetAll(context.Context, *temperature_api.GetAllTemperatures) ([]*temperature_model.Temperature, error)
	Create(context.Context, *temperature_api.CreateTemperature) error
	CreateSeveral(context.Context, *temperature_api.CreateSeveralTemperature) error
	Update(context.Context, *temperature_api.UpdateTemperature) error
	Delete(context.Context, *temperature_api.DeleteTemperature) error
}

type FlangeTypeSnp interface {
	Get(context.Context, *flange_type_snp_api.GetFlangeTypeSnp) ([]*flange_type_snp_model.FlangeTypeSnp, error)
	Create(context.Context, *flange_type_snp_api.CreateFlangeTypeSnp) error
	CreateSeveral(context.Context, *flange_type_snp_api.CreateSeveralFlangeTypeSnp) error
	Update(context.Context, *flange_type_snp_api.UpdateFlangeTypeSnp) error
	Delete(context.Context, *flange_type_snp_api.DeleteFlangeTypeSnp) error
}

type SnpFiller interface {
	GetAll(context.Context, *snp_filler_api.GetSnpFillers) ([]*snp_filler_model.SnpFiller, error)
	Create(context.Context, *snp_filler_api.CreateSnpFiller) error
	CreateSeveral(context.Context, *snp_filler_api.CreateSeveralSnpFiller) error
	Update(context.Context, *snp_filler_api.UpdateSnpFiller) error
	Delete(context.Context, *snp_filler_api.DeleteSnpFiller) error
}

type SnpStandard interface {
	GetAll(context.Context, *snp_standard_api.GetAllSnpStandards) ([]*snp_standard_model.SnpStandard, error)
	GetDefault(context.Context) (*snp_standard_model.SnpStandard, error)
	Create(context.Context, *snp_standard_api.CreateSnpStandard) error
	CreateSeveral(context.Context, *snp_standard_api.CreateSeveralSnpStandard) error
	Update(context.Context, *snp_standard_api.UpdateSnpStandard) error
	Delete(context.Context, *snp_standard_api.DeleteSnpStandard) error
}

type SnpType interface {
	Get(context.Context, *snp_type_api.GetSnpTypes) ([]*snp_type_model.SnpType, error)
	GetWithFlange(context.Context, *snp_api.GetSnpData) ([]*snp_model.FlangeType, error)
	Create(context.Context, *snp_type_api.CreateSnpType) error
	CreateSeveral(context.Context, *snp_type_api.CreateSeveralSnpType) error
	Update(context.Context, *snp_type_api.UpdateSnpType) error
	Delete(context.Context, *snp_type_api.DeleteSnpType) error
}

type SnpMaterial interface {
	Get(context.Context, *snp_material_api.GetSnpMaterial) ([]*snp_material_model.SnpMaterial, error)
	Create(context.Context, *snp_material_api.CreateSnpMaterial) error
	Update(context.Context, *snp_material_api.UpdateSnpMaterial) error
	Delete(context.Context, *snp_material_api.DeleteSnpMaterial) error
}

type SnpData interface {
	Get(context.Context, *snp_data_api.GetSnpData) (*snp_data_model.SnpData, error)
	Create(context.Context, *snp_data_api.CreateSnpData) error
	Update(context.Context, *snp_data_api.UpdateSnpData) error
	Delete(context.Context, *snp_data_api.DeleteSnpData) error
}

type SnpSize interface {
	Get(context.Context, *snp_size_api.GetSnpSize) ([]*snp_size_model.SnpSize, error)
	Create(context.Context, *snp_size_api.CreateSnpSize) error
}

type OrderNew interface {
	GetNumber(ctx context.Context, orderId, date string) (int64, error)
	Create(ctx context.Context, order order_api.CreateOrder, date string) error
}

type Position interface {
	CreateSeveral(context.Context, []*position_model.Position) error
}

type PositionSnp interface {
	CreateSeveral(context.Context, []*position_model.Position) error
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

	//* new
	FlangeStandard
	Material
	Mounting
	Standard
	Temperature
	FlangeTypeSnp
	SnpFiller
	SnpStandard
	SnpType
	SnpMaterial
	SnpData
	SnpSize
	OrderNew
	Position
	PositionSnp
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

		FlangeStandard: postgres.NewFlangeStandardRepo(db),
		Material:       postgres.NewMaterialRepo(db),
		Mounting:       postgres.NewMountingRepo(db),
		Standard:       postgres.NewStandardRepo(db),
		Temperature:    postgres.NewTemperatureRepo(db),
		FlangeTypeSnp:  postgres.NewFlangeTypeRepo(db),
		SnpFiller:      postgres.NewSNPFillerRepo(db),
		SnpStandard:    postgres.NewSnpStandardRepo(db),
		SnpType:        postgres.NewSNPTypeRepo(db),
		SnpMaterial:    postgres.NewSnpMaterialRepo(db),
		SnpData:        postgres.NewSnpDataRepo(db),
		SnpSize:        postgres.NewSnpSizeRepo(db),
		OrderNew:       postgres.NewOrderRepo(db),
		Position:       postgres.NewPositionRepo(db),
		PositionSnp:    postgres.NewPositionSnpRepo(db),
	}
}
