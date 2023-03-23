package service

import (
	"bytes"
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_snp_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/mounting_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
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
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
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
)

type Stand interface {
	GetAll(*pro_api.GetStandsRequest) (stands []*pro_api.Stand, err error)
	Create(*pro_api.CreateStandRequest) (st *pro_api.IdResponse, err error)
	Update(*pro_api.UpdateStandRequest) error
	Delete(*pro_api.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*pro_api.Flange, error)
	Create(*pro_api.CreateFlangeRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateFlangeRequest) error
	Delete(*pro_api.DeleteFlangeRequest) error
}

type StFl interface {
	Get() ([]*pro_api.StFl, error)
	Create(*pro_api.CreateStFlRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateStFlRequest) error
	Delete(*pro_api.DeleteStFlRequest) error
}

type TypeFl interface {
	Get() ([]*pro_api.TypeFl, error)
	GetAll() ([]*pro_api.TypeFl, error)
	Create(*pro_api.CreateTypeFlRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateTypeFlRequest) error
	Delete(*pro_api.DeleteTypeFlRequest) error
}

type Addit interface {
	GetAll() ([]*pro_api.Additional, error)
	Create(*pro_api.CreateAddRequest) (*pro_api.SuccessResponse, error)
	UpdateMat(*pro_api.UpdateAddMatRequest) (*pro_api.SuccessResponse, error)
	UpdateMod(*pro_api.UpdateAddModRequest) (*pro_api.SuccessResponse, error)
	UpdateTemp(*pro_api.UpdateAddTemRequest) (*pro_api.SuccessResponse, error)
	UpdateMoun(*pro_api.UpdateAddMounRequest) (*pro_api.SuccessResponse, error)
	UpdateGrap(*pro_api.UpdateAddGrapRequest) (*pro_api.SuccessResponse, error)
	UpdateFillers(*pro_api.UpdateAddFillersRequest) (*pro_api.SuccessResponse, error)
	UpdateCoating(*pro_api.UpdateAddCoatingRequest) (*pro_api.SuccessResponse, error)
	UpdateConstruction(*pro_api.UpdateAddConstructionRequest) (*pro_api.SuccessResponse, error)
	UpdateObturator(*pro_api.UpdateAddObturatorRequest) (*pro_api.SuccessResponse, error)
	UpdateBasis(*pro_api.UpdateAddBasisRequest) (*pro_api.SuccessResponse, error)
	UpdatePObturator(*pro_api.UpdateAddPObturatorRequest) (*pro_api.SuccessResponse, error)
	UpdateSealant(*pro_api.UpdateAddSealantRequest) (*pro_api.SuccessResponse, error)
}

type Size interface {
	Get(*pro_api.GetSizesRequest) ([]*pro_api.Size, []*pro_api.Dn, error)
	GetAll(*pro_api.GetSizesRequest) ([]*pro_api.Size, []*pro_api.Dn, error)
	Create(*pro_api.CreateSizeRequest) (*pro_api.IdResponse, error)
	CreateMany(*pro_api.CreateSizesRequest) error
	Update(*pro_api.UpdateSizeRequest) error
	Delete(*pro_api.DeleteSizeRequest) error
	DeleteAll(*pro_api.DeleteAllSizeRequest) error
}

type SNP interface {
	Get(*pro_api.GetSNPRequest) ([]*pro_api.SNP, error)
	Create(*pro_api.CreateSNPRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateSNPRequest) error
	Delete(*pro_api.DeleteSNPRequest) error

	AddMat(id string) error
	DeleteMat(id string, materials []*pro_api.AddMaterials) error
	AddMoun(id string) error
	DeleteMoun(id string) error
	AddGrap(id string) error
	DeleteGrap(id string) error
	DeleteFiller(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
}

type PutgImage interface {
	Get(req *pro_api.GetPutgImageRequest) ([]*pro_api.PutgImage, error)
	Create(image *pro_api.CreatePutgImageRequest) (*pro_api.IdResponse, error)
	Update(image *pro_api.UpdatePutgImageRequest) error
	Delete(image *pro_api.DeletePutgImageRequest) error
}

type Putg interface {
	Get(*pro_api.GetPutgRequest) ([]*pro_api.Putg, error)
	Create(*pro_api.CreatePutgRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdatePutgRequest) error
	Delete(*pro_api.DeletePutgRequest) error

	DeleteGrap(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
	DeleteMat(id string, materials []*pro_api.AddMaterials) error
	DeleteCon(id string) error
	DeleteObt(id string) error
	DeleteMoun(id string) error
	DeleteCoating(id string) error
}

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

type Materials interface {
	GetAll(*pro_api.GetMaterialsRequest) ([]*pro_api.Materials, error)
	Create(*pro_api.CreateMaterialsRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateMaterialsRequest) error
	Delete(*pro_api.DeleteMaterialsRequest) error
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

type Order interface {
	GetAll(*pro_api.GetAllOrdersRequest) ([]*pro_api.Order, error)
	Create(*pro_api.CreateOrderRequest) (*pro_api.IdResponse, error)
	Delete(*pro_api.DeleteOrderRequest) (*pro_api.IdResponse, error)
	Save(context.Context, *pro_api.SaveOrderRequest) (*bytes.Buffer, error)
	Send(context.Context, *pro_api.SaveOrderRequest) error
	Copy(*pro_api.CopyOrderRequest) error
}

type OrderPosition interface {
	Get(*pro_api.GetPositionsRequest) ([]*pro_api.OrderPosition, error)
	GetCur(*pro_api.GetCurPositionsRequest) ([]*pro_api.OrderPosition, error)
	Add(*pro_api.AddPositionRequest) (*pro_api.IdResponse, error)
	Copy(*pro_api.CopyPositionRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdatePositionRequest) (*pro_api.IdResponse, error)
	Remove(*pro_api.RemovePositionRequest) (*pro_api.IdResponse, error)
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

type Snp interface {
	Get(context.Context, *snp_api.GetSnp) (*snp_api.Snp, error)
	GetData(context.Context, *snp_api.GetSnpData) (*snp_model.SnpData, error)
}

type OrderNew interface {
	Get(context.Context, *order_api.GetOrder) (*order_model.CurrentOrder, error)
	GetCurrent(context.Context, *order_api.GetCurrentOrder) (*order_model.CurrentOrder, error)
	GetAll(context.Context, *order_api.GetAllOrders) ([]*order_model.Order, error)
	GetFile(context.Context, *order_api.GetOrder) (*bytes.Buffer, string, error)
	GetOpen(context.Context, *order_api.GetManagerOrders) ([]*order_model.ManagerOrder, error)
	Save(context.Context, *order_api.CreateOrder) (*order_api.OrderNumber, error)
	Copy(context.Context, *order_api.CopyOrder) error
	Create(context.Context, *order_api.CreateOrder) (string, error)
	SetStatus(ctx context.Context, status *order_api.Status) error
}

type Position interface {
	Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error)
	GetAll(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error)
	GetFull(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error)
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

type Services struct {
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
	Interview
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
	Snp
	OrderNew
	Position
	PositionSnp
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

	positionSnp := NewPositionSnpService(repos.PositionSnp)
	position := NewPositionService_New(repos.Position, positionSnp)
	order := NewOrderService_New(repos.OrderNew, position)

	return &Services{
		Stand:         NewStandService(repos.Stand),
		Flange:        NewFlangeService(repos.Flange),
		StFl:          NewStFlService(repos.StFl),
		TypeFl:        NewTypeFlService(repos.TypeFl),
		Addit:         NewAdditService(repos.Addit),
		Size:          NewSizeService(repos.Size),
		SNP:           NewSNPService(repos.SNP),
		PutgImage:     NewPutgImageService(repos.PutgImage),
		Putg:          NewPutgService(repos.Putg),
		PutgmImage:    NewPutgmImageService(repos.PutgmImage),
		Putgm:         NewPutgmService(repos.Putgm),
		Materials:     NewMatService(repos.Materials),
		BoltMaterials: NewBoltMatRepo(repos.BoltMaterials),
		SizeInt:       NewSizeIntService(repos.SizeInt),
		Interview:     NewInterviewService(email, file),
		Order:         NewOrderService(repos.Order, email, file, user),
		OrderPosition: NewPositionService(repos.OrderPosition, file),

		FlangeStandard: NewFlangeStandardService(repos.FlangeStandard),
		Material:       NewMaterialService(repos.Material),
		Mounting:       mounting,
		Standard:       standard,
		Temperature:    NewTemperatureService(repos.Temperature),
		FlangeTypeSnp:  NewFlangeTypeSnpService(repos.FlangeTypeSnp),
		SnpFiller:      filler,
		SnpStandard:    snpStandard,
		SnpType:        snpType,
		SnpMaterial:    snpMaterial,
		SnpData:        snpData,
		SnpSize:        snpSize,
		Snp:            NewSnpService(filler, snpMaterial, snpType, mounting, standard, snpData, snpStandard, snpSize),

		PositionSnp: positionSnp,
		Position:    position,
		OrderNew:    order,
	}
}
