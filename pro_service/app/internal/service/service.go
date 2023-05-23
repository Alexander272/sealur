package service

import (
	"bytes"
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/file_api"
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
	"github.com/Alexander272/sealur_proto/api/pro/putg_api"
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
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/xuri/excelize/v2"
)

type PutgImage interface {
	Get(req *pro_api.GetPutgImageRequest) ([]*pro_api.PutgImage, error)
	Create(image *pro_api.CreatePutgImageRequest) (*pro_api.IdResponse, error)
	Update(image *pro_api.UpdatePutgImageRequest) error
	Delete(image *pro_api.DeletePutgImageRequest) error
}

// type Putg interface {
// 	Get(*pro_api.GetPutgRequest) ([]*pro_api.Putg, error)
// 	Create(*pro_api.CreatePutgRequest) (*pro_api.IdResponse, error)
// 	Update(*pro_api.UpdatePutgRequest) error
// 	Delete(*pro_api.DeletePutgRequest) error

// 	DeleteGrap(id string) error
// 	DeleteTemp(id string) error
// 	DeleteMod(id string) error
// 	DeleteMat(id string, materials []*pro_api.AddMaterials) error
// 	DeleteCon(id string) error
// 	DeleteObt(id string) error
// 	DeleteMoun(id string) error
// 	DeleteCoating(id string) error
// }

type PutgmImage interface {
	Get(req *pro_api.GetPutgmImageRequest) ([]*pro_api.PutgmImage, error)
	Create(image *pro_api.CreatePutgmImageRequest) (*pro_api.IdResponse, error)
	Update(image *pro_api.UpdatePutgmImageRequest) error
	Delete(image *pro_api.DeletePutgmImageRequest) error
}

type Putgm interface {
	Get(*pro_api.GetPutgmRequest) ([]*pro_api.Putgm, error)
	Create(*pro_api.CreatePutgmRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdatePutgmRequest) error
	Delete(*pro_api.DeletePutgmRequest) error

	DeleteGrap(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
	DeleteMat(id string, materials []*pro_api.AddMaterials) error
	DeleteCon(id string) error
	DeleteObt(id string) error
	DeleteSeal(id string) error
	DeleteMoun(id string) error
	DeleteCoating(id string) error
}

type BoltMaterials interface {
	GetAll(*pro_api.GetBoltMaterialsRequest) ([]*pro_api.BoltMaterials, error)
	Create(*pro_api.CreateBoltMaterialsRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateBoltMaterialsRequest) error
	Delete(*pro_api.DeleteBoltMaterialsRequest) error
}

type SizeInt interface {
	Get(*pro_api.GetSizesIntRequest) ([]*pro_api.SizeInt, []*pro_api.Dn, error)
	GetAll(*pro_api.GetAllSizeIntRequest) ([]*pro_api.SizeInt, []*pro_api.Dn, error)
	Create(*pro_api.CreateSizeIntRequest) (*pro_api.IdResponse, error)
	CreateMany(*pro_api.CreateSizesIntRequest) error
	Update(*pro_api.UpdateSizeIntRequest) error
	Delete(*pro_api.DeleteSizeIntRequest) error
	DeleteAll(*pro_api.DeleteAllSizeIntRequest) error
}

type Interview interface {
	SendInterview(context.Context, *pro_api.SendInterviewRequest) error
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

type Snp interface {
	Get(context.Context, *snp_api.GetSnp) (*snp_api.Snp, error)
	GetData(context.Context, *snp_api.GetSnpData) (*snp_model.SnpData, error)
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
	GetNew(context.Context, *putg_construction_api.GetPutgConstruction_New) ([]*putg_construction_type_model.PutgConstruction, error)
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
	GetNew(context.Context, *putg_filler_api.GetPutgFiller_New) ([]*putg_filler_model.PutgFiller, error)
	Create(context.Context, *putg_filler_api.CreatePutgFiller) error
	Update(context.Context, *putg_filler_api.UpdatePutgFiller) error
	Delete(context.Context, *putg_filler_api.DeletePutgFiller) error
}

type PutgFlangeType interface {
	Get(context.Context, *putg_flange_type_api.GetPutgFlangeType) ([]*putg_flange_type_model.PutgFlangeType, error)
	Create(context.Context, *putg_flange_type_api.CreatePutgFlangeType) error
	Update(context.Context, *putg_flange_type_api.UpdatePutgFlangeType) error
	Delete(context.Context, *putg_flange_type_api.DeletePutgFlangeType) error
}

type PutgStandard interface {
	Get(context.Context, *putg_standard_api.GetPutgStandard) ([]*putg_standard_model.PutgStandard, error)
	Create(context.Context, *putg_standard_api.CreatePutgStandard) error
	Update(context.Context, *putg_standard_api.UpdatePutgStandard) error
	Delete(context.Context, *putg_standard_api.DeletePutgStandard) error
}

type PutgMaterial interface {
	Get(context.Context, *putg_material_api.GetPutgMaterial) (*putg_material_model.PutgMaterials, error)
	Create(context.Context, *putg_material_api.CreatePutgMaterial) error
	Update(context.Context, *putg_material_api.UpdatePutgMaterial) error
	Delete(context.Context, *putg_material_api.DeletePutgMaterial) error
}

type PutgData interface {
	Get(context.Context, *putg_data_api.GetPutgData) (*putg_data_model.PutgData, error)
	Create(context.Context, *putg_data_api.CreatePutgData) error
	Update(context.Context, *putg_data_api.UpdatePutgData) error
	Delete(context.Context, *putg_data_api.DeletePutgData) error
}

type PutgSize interface {
	Get(context.Context, *putg_size_api.GetPutgSize) ([]*putg_size_model.PutgSize, error)
	GetNew(context.Context, *putg_size_api.GetPutgSize_New) ([]*putg_size_model.PutgSize, error)
	Create(context.Context, *putg_size_api.CreatePutgSize) error
	CreateSeveral(context.Context, *putg_size_api.CreateSeveralPutgSize) error
	Update(context.Context, *putg_size_api.UpdatePutgSize) error
	Delete(context.Context, *putg_size_api.DeletePutgSize) error
}

type PutgType interface {
	Get(context.Context, *putg_type_api.GetPutgType) ([]*putg_type_model.PutgType, error)
	Create(context.Context, *putg_type_api.CreatePutgType) error
	Update(context.Context, *putg_type_api.UpdatePutgType) error
	Delete(context.Context, *putg_type_api.DeletePutgType) error
}

type Putg interface {
	GetBase(context.Context, *putg_api.GetPutgBase) (*putg_api.PutgBase, error)
	GetData(context.Context, *putg_api.GetPutgData) (*putg_api.PutgData, error)
	Get(context.Context, *putg_api.GetPutg) (*putg_api.Putg, error)
}

type OrderNew interface {
	Get(context.Context, *order_api.GetOrder) (*order_model.CurrentOrder, error)
	GetCurrent(context.Context, *order_api.GetCurrentOrder) (*order_model.CurrentOrder, error)
	GetAll(context.Context, *order_api.GetAllOrders) ([]*order_model.Order, error)
	GetFile(context.Context, *order_api.GetOrder) (*bytes.Buffer, string, error)
	GetOpen(context.Context, *order_api.GetManagerOrders) ([]*order_model.ManagerOrder, error)
	GetAnalytics(context.Context, *order_api.GetOrderAnalytics) (*order_api.Analytics, error)
	GetFullAnalytics(context.Context, *order_api.GetFullOrderAnalytics) ([]*analytic_model.FullOrder, error)
	Save(context.Context, *order_api.CreateOrder) (*order_api.OrderNumber, error)
	Copy(context.Context, *order_api.CopyOrder) error
	Create(context.Context, *order_api.CreateOrder) (string, error)
	SetInfo(context.Context, *order_api.Info) error
	SetStatus(context.Context, *order_api.Status) error
	SetManager(context.Context, *order_api.Manager) error
}

type Position interface {
	Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error)
	GetAll(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error)
	GetFull(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error)
	GetAnalytics(context.Context, *order_api.GetOrderAnalytics) (*order_api.Analytics, error)
	Create(context.Context, *position_model.FullPosition) (string, error)
	CreateSeveral(ctx context.Context, positions []*position_model.FullPosition, orderId string) error
	Update(context.Context, *position_model.FullPosition) error
	Copy(context.Context, *position_api.CopyPosition) (string, error)
	Delete(ctx context.Context, positionId string) error
}

type PositionSnp interface {
	Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error)
	GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionSnp, error)
	Create(context.Context, *position_model.FullPosition) error
	CreateSeveral(context.Context, []*position_model.FullPosition) error
	Update(context.Context, *position_model.FullPosition) error
	Copy(ctx context.Context, targetId string, position *position_api.CopyPosition) (string, error)
	Delete(ctx context.Context, positionId string) error
}

type Zip interface {
	Create(fileName string, excel *excelize.File) (*bytes.Buffer, error)
	InsertDrawings(file bytes.Buffer, drawings []string, buffer *bytes.Buffer) (*bytes.Buffer, error)
	CreateWithDrawings(excelName string, excel *excelize.File, file bytes.Buffer, drawings []string) (*bytes.Buffer, error)
}

type Services struct {
	PutgImage
	PutgmImage
	Putgm
	BoltMaterials
	SizeInt
	Interview

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
	Snp

	PutgConfiguration
	PutgBaseConstruction
	PutgConstruction
	PutgBaseFiller
	PutgFiller
	PutgFlangeType
	PutgStandard
	PutgMaterial
	PutgData
	PutgSize
	PutgType
	Putg

	OrderNew
	Position
	PositionSnp
	Zip
}

func NewServices(repos *repository.Repositories, email email_api.EmailServiceClient,
	file file_api.FileServiceClient, user user_api.UserServiceClient,
) *Services {
	mounting := NewMountingService(repos.Mounting)
	standard := NewStandardService(repos.Standard)
	filler := NewSnpFillerService(repos.SnpFiller)
	snpStandard := NewSnpStandardService(repos.SnpStandard)
	snpType := NewSnpTypeService(repos.SnpType)
	snpMaterial := NewSnpMaterialService(repos.SnpMaterial)
	snpData := NewSnpDataService(repos.SnpData)
	snpSize := NewSnpSizeService(repos.SnpSize)

	putgConfiguration := NewPutgConfigurationService(repos.PutgConfiguration)
	putgBaseConstruction := NewPutgBaseConstructionService(repos.PutgBaseConstruction)
	putgConstruction := NewPutgConstructionService(repos.PutgConstruction)
	putgBaseFiller := NewPutgBaseFillerService(repos.PutgBaseFiller)
	putgFiller := NewPutgFillerService(repos.PutgFiller)
	putgFlangeType := NewPutgFlangeTypeService(repos.PutgFlangeType)
	putgStandard := NewPutgStandardService(repos.PutgStandard)
	putgMaterial := NewPutgMaterialService(repos.PutgMaterial)
	putgData := NewPutgDataService(repos.PutgData)
	putgSize := NewPutgSizeService(repos.PutgSize)
	putgType := NewPutgTypeService(repos.PutgType)

	putg := NewPutgService(PutgDeps{
		Configuration: putgConfiguration,
		Construction:  putgConstruction,
		Filler:        putgFiller,
		FlangeType:    putgFlangeType,
		Standard:      putgStandard,
		Materials:     putgMaterial,
		Data:          putgData,
		Sizes:         putgSize,
		PutgType:      putgType,
		Mounting:      mounting,
	})

	zip := NewZipService()

	positionSnp := NewPositionSnpService(repos.PositionSnp)
	position := NewPositionService_New(repos.Position, positionSnp, file)
	order := NewOrderService_New(repos.OrderNew, position, zip, user, file)

	return &Services{
		// PutgImage:     NewPutgImageService(repos.PutgImage),
		// Putg:          NewPutgService(repos.Putg),
		// PutgmImage:    NewPutgmImageService(repos.PutgmImage),
		// Putgm:         NewPutgmService(repos.Putgm),
		// BoltMaterials: NewBoltMatRepo(repos.BoltMaterials),
		// SizeInt:       NewSizeIntService(repos.SizeInt),
		// Interview:     NewInterviewService(email, file),

		FlangeStandard: NewFlangeStandardService(repos.FlangeStandard),
		Material:       NewMaterialService(repos.Material),
		Mounting:       mounting,
		Standard:       standard,

		Temperature:   NewTemperatureService(repos.Temperature),
		FlangeType:    NewFlangeTypeService(repos.FlangeType),
		FlangeTypeSnp: NewFlangeTypeSnpService(repos.FlangeTypeSnp),
		SnpFiller:     filler,
		SnpStandard:   snpStandard,
		SnpType:       snpType,
		SnpMaterial:   snpMaterial,
		SnpData:       snpData,
		SnpSize:       snpSize,
		Snp:           NewSnpService(filler, snpMaterial, snpType, mounting, standard, snpData, snpStandard, snpSize),

		PutgConfiguration:    putgConfiguration,
		PutgBaseConstruction: putgBaseConstruction,
		PutgConstruction:     putgConstruction,
		PutgBaseFiller:       putgBaseFiller,
		PutgFiller:           putgFiller,
		PutgFlangeType:       putgFlangeType,
		PutgStandard:         putgStandard,
		PutgMaterial:         putgMaterial,
		PutgData:             putgData,
		PutgSize:             putgSize,
		PutgType:             putgType,
		Putg:                 putg,

		PositionSnp: positionSnp,
		Position:    position,
		OrderNew:    order,
	}
}
