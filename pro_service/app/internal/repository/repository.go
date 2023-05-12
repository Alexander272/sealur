package repository

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository/postgres"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/analytic_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_snp_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/mounting_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_configuration_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_construction_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_flange_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_type_model"
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
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_base_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_conf_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_base_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_type_api"
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
	// UpdateAddit(putg models.UpdateAdditDTO) error
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
	// UpdateAddit(putgm models.UpdateAdditDTO) error
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

type FlangeType interface {
	Get(context.Context, *flange_type_api.GetFlangeType) ([]*flange_type_model.FlangeType, error)
	Create(context.Context, *flange_type_api.CreateFlangeType) error
	Update(context.Context, *flange_type_api.UpdateFlangeType) error
	Delete(context.Context, *flange_type_api.DeleteFlangeType) error
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
	Get(context.Context, *snp_material_api.GetSnpMaterial) (*snp_material_model.SnpMaterials, error)
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
	CreateSeveral(context.Context, *snp_size_api.CreateSeveralSnpSize) error
	Update(context.Context, *snp_size_api.UpdateSnpSize) error
	Delete(context.Context, *snp_size_api.DeleteSnpSize) error
}

type PutgConfiguration interface {
	Get(context.Context, *putg_conf_api.GetPutgConfiguration) ([]*putg_configuration_model.PutgConfiguration, error)
	Create(context.Context, *putg_conf_api.CreatePutgConfiguration) error
	Update(context.Context, *putg_conf_api.UpdatePutgConfiguration) error
	Delete(context.Context, *putg_conf_api.DeletePutgConfiguration) error
}

type PutgBaseConstruction interface {
	Get(context.Context, *putg_base_construction_api.GetPutgBaseConstruction) ([]*putg_construction_type_model.PutgConstruction, error)
	Create(context.Context, *putg_base_construction_api.CreatePutgBaseConstruction) error
	Update(context.Context, *putg_base_construction_api.UpdatePutgBaseConstruction) error
	Delete(context.Context, *putg_base_construction_api.DeletePutgBaseConstruction) error
}

type PutgConstruction interface {
	Get(context.Context, *putg_construction_api.GetPutgConstruction) ([]*putg_construction_type_model.PutgConstruction, error)
	Create(context.Context, *putg_construction_api.CreatePutgConstruction) error
	Update(context.Context, *putg_construction_api.UpdatePutgConstruction) error
	Delete(context.Context, *putg_construction_api.DeletePutgConstruction) error
}

type PutgBaseFiller interface {
	Get(context.Context, *putg_filler_base_api.GetPutgBaseFiller) ([]*putg_filler_model.PutgFiller, error)
	Create(context.Context, *putg_filler_base_api.CreatePutgBaseFiller) error
	Update(context.Context, *putg_filler_base_api.UpdatePutgBaseFiller) error
	Delete(context.Context, *putg_filler_base_api.DeletePutgBaseFiller) error
}

type PutgFiller interface {
	Get(context.Context, *putg_filler_api.GetPutgFiller) ([]*putg_filler_model.PutgFiller, error)
	Create(context.Context, *putg_filler_api.CreatePutgFiller) error
	Update(context.Context, *putg_filler_api.UpdatePutgFiller) error
	Delete(context.Context, *putg_filler_api.DeletePutgFiller) error
}

type PutgFlangeType interface {
	Get(context.Context, *putg_flange_type_api.GetPutgFlangeType) ([]*putg_flange_type_model.PutgFlangeType, error)
	Create(ctx context.Context, fl *putg_flange_type_api.CreatePutgFlangeType) error
	Update(ctx context.Context, fl *putg_flange_type_api.UpdatePutgFlangeType) error
	Delete(ctx context.Context, fl *putg_flange_type_api.DeletePutgFlangeType) error
}

type PutgStandard interface {
	Get(context.Context, *putg_standard_api.GetPutgStandard) ([]*putg_standard_model.PutgStandard, error)
	Create(context.Context, *putg_standard_api.CreatePutgStandard) error
	Update(context.Context, *putg_standard_api.UpdatePutgStandard) error
	Delete(context.Context, *putg_standard_api.DeletePutgStandard) error
}

type PutgData interface {
	Get(context.Context, *putg_data_api.GetPutgData) (*putg_data_model.PutgData, error)
	GetByConstruction(context.Context, *putg_data_api.GetPutgData) ([]*putg_data_model.PutgData, error)
	Create(context.Context, *putg_data_api.CreatePutgData) error
	Update(context.Context, *putg_data_api.UpdatePutgData) error
	Delete(context.Context, *putg_data_api.DeletePutgData) error
}

type PutgSize interface {
	Get(context.Context, *putg_size_api.GetPutgSize) ([]*putg_size_model.PutgSize, error)
	Create(context.Context, *putg_size_api.CreatePutgSize) error
	CreateSeveral(context.Context, *putg_size_api.CreateSeveralPutgSize) error
	Update(context.Context, *putg_size_api.UpdatePutgSize) error
	Delete(context.Context, *putg_size_api.DeletePutgSize) error
}

type PutgMaterial interface {
	Get(context.Context, *putg_material_api.GetPutgMaterial) (*putg_material_model.PutgMaterials, error)
	Create(context.Context, *putg_material_api.CreatePutgMaterial) error
	Update(context.Context, *putg_material_api.UpdatePutgMaterial) error
	Delete(context.Context, *putg_material_api.DeletePutgMaterial) error
}

type PutgType interface {
	Get(context.Context, *putg_type_api.GetPutgType) ([]*putg_type_model.PutgType, error)
	Create(context.Context, *putg_type_api.CreatePutgType) error
	Update(context.Context, *putg_type_api.UpdatePutgType) error
	Delete(context.Context, *putg_type_api.DeletePutgType) error
}

type OrderNew interface {
	Get(context.Context, *order_api.GetOrder) (*order_model.FullOrder, error)
	GetCurrent(context.Context, *order_api.GetCurrentOrder) (*order_model.CurrentOrder, error)
	GetAll(context.Context, *order_api.GetAllOrders) ([]*order_model.Order, error)
	GetNumber(ctx context.Context, order *order_api.CreateOrder, date string) (int64, error)
	GetOpen(ctx context.Context, managerId string) ([]*order_model.ManagerOrder, error)
	GetAnalytics(context.Context, *order_api.GetOrderAnalytics) ([]*analytic_model.Order, error)
	GetFullAnalytics(context.Context, *order_api.GetFullOrderAnalytics) ([]*analytic_model.FullOrder, error)
	Create(ctx context.Context, order *order_api.CreateOrder, date string) error
	SetInfo(context.Context, *order_api.Info) error
	SetStatus(context.Context, *order_api.Status) error
	SetManager(context.Context, *order_api.Manager) error
}

type Position interface {
	GetById(ctx context.Context, positionId string) (position *position_model.FullPosition, err error)
	Get(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error)
	GetByTitle(ctx context.Context, title, orderId string) (string, error)
	GetAnalytics(context.Context, *order_api.GetOrderAnalytics) (*order_api.Analytics, error)
	Create(context.Context, *position_model.FullPosition) (string, error)
	CreateSeveral(context.Context, []*position_model.FullPosition) error
	Update(context.Context, *position_model.FullPosition) error
	Copy(context.Context, *position_api.CopyPosition) (string, error)
	Delete(ctx context.Context, positionId string) error
}

type PositionSnp interface {
	Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error)
	GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionSnp, error)
	Create(ctx context.Context, position *position_model.FullPosition) error
	CreateSeveral(context.Context, []*position_model.FullPosition) error
	Update(ctx context.Context, position *position_model.FullPosition) error
	Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error)
	Delete(ctx context.Context, positionId string) error
}

type Repositories struct {
	PutgImage
	Putg
	PutgmImage
	Putgm
	BoltMaterials
	SizeInt

	//* new
	FlangeStandard
	Material
	Mounting
	Standard
	Temperature

	FlangeType
	FlangeTypeSnp
	SnpFiller
	SnpStandard
	SnpType
	SnpMaterial
	SnpData
	SnpSize

	PutgConfiguration
	PutgBaseConstruction
	PutgConstruction
	PutgFiller
	PutgFlangeType
	PutgStandard
	PutgMaterial
	PutgData
	PutgSize
	PutgType

	OrderNew
	Position
	PositionSnp
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		// PutgImage:     NewPutgImageRepo(db),
		// Putg:          NewPutgRepo(db),
		// PutgmImage:    NewPutgmImageRepo(db),
		// Putgm:         NewPutgmRepo(db),
		// BoltMaterials: NewBoltMatRepo(db),
		// SizeInt:       NewSizeIntRepo(db),

		FlangeStandard: postgres.NewFlangeStandardRepo(db),
		Material:       postgres.NewMaterialRepo(db),
		Mounting:       postgres.NewMountingRepo(db),
		Standard:       postgres.NewStandardRepo(db),
		Temperature:    postgres.NewTemperatureRepo(db),

		FlangeType:    postgres.NewFlangeTypeRepo(db),
		FlangeTypeSnp: postgres.NewFlangeTypeSnpRepo(db),
		SnpFiller:     postgres.NewSNPFillerRepo(db),
		SnpStandard:   postgres.NewSnpStandardRepo(db),
		SnpType:       postgres.NewSNPTypeRepo(db),
		SnpMaterial:   postgres.NewSnpMaterialRepo(db),
		SnpData:       postgres.NewSnpDataRepo(db),
		SnpSize:       postgres.NewSnpSizeRepo(db),

		PutgConfiguration:    postgres.NewPutgConfigurationRepo(db),
		PutgBaseConstruction: postgres.NewPutgConstructionBaseRepo(db),
		PutgConstruction:     postgres.NewPutgConstructionRepo(db),
		PutgFiller:           postgres.NewPutgFillerRepo(db),
		PutgFlangeType:       postgres.NewPutgFlangeTypeRepo(db),
		PutgStandard:         postgres.NewPutgStandardRepo(db),
		PutgMaterial:         postgres.NewPutgMaterialRepo(db),
		PutgData:             postgres.NewPutgDataRepo(db),
		PutgSize:             postgres.NewPutgSizeRepo(db),
		PutgType:             postgres.NewPutgTypeRepo(db),

		OrderNew:    postgres.NewOrderRepo(db),
		Position:    postgres.NewPositionRepo(db),
		PositionSnp: postgres.NewPositionSnpRepo(db),
	}
}
