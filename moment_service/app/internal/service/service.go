package service

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur/moment_service/internal/service/read"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type CalcFlange interface {
	Calculation(ctx context.Context, data *moment_api.CalcFlangeRequest) (*moment_api.FlangeResponse, error)
}

type CalcCap interface {
	Calculation(ctx context.Context, data *moment_api.CalcCapRequest) (*moment_api.CapResponse, error)
}

type Flange interface {
	GetFlangeSize(context.Context, *moment_api.GetFlangeSizeRequest) (models.FlangeSize, error)
	GetFullFlangeSize(context.Context, *moment_api.GetFullFlangeSizeRequest) (*moment_api.FullFlangeSizeResponse, error)
	GetBasisFlangeSize(context.Context, *moment_api.GetBasisFlangeSizeRequest) (*moment_api.BasisFlangeSizeResponse, error)
	CreateFlangeSize(context.Context, *moment_api.CreateFlangeSizeRequest) error
	CreateFlangeSizes(context.Context, *moment_api.CreateFlangeSizesRequest) error
	UpdateFlangeSize(context.Context, *moment_api.UpdateFlangeSizeRequest) error
	DeleteFlangeSize(context.Context, *moment_api.DeleteFlangeSizeRequest) error

	GetBolts(context.Context, *moment_api.GetBoltsRequest) ([]*moment_api.Bolt, error)
	GetAllBolts(context.Context, *moment_api.GetBoltsRequest) ([]*moment_api.Bolt, error)
	CreateBolt(context.Context, *moment_api.CreateBoltRequest) error
	CreateBolts(context.Context, *moment_api.CreateBoltsRequest) error
	UpdateBolt(context.Context, *moment_api.UpdateBoltRequest) error
	DeleteBolt(context.Context, *moment_api.DeleteBoltRequest) error

	GetTypeFlange(context.Context, *moment_api.GetTypeFlangeRequest) ([]*moment_api.TypeFlange, error)
	CreateTypeFlange(context.Context, *moment_api.CreateTypeFlangeRequest) (id string, err error)
	UpdateTypeFlange(context.Context, *moment_api.UpdateTypeFlangeRequest) error
	DeleteTypeFlange(context.Context, *moment_api.DeleteTypeFlangeRequest) error

	GetStandarts(context.Context, *moment_api.GetStandartsRequest) ([]*moment_api.Standart, error)
	GetStandartsWithSize(context.Context, *moment_api.GetStandartsRequest) ([]*moment_api.StandartWithSize, error)
	CreateStandart(context.Context, *moment_api.CreateStandartRequest) (id string, err error)
	UpdateStandart(context.Context, *moment_api.UpdateStandartRequest) error
	DeleteStandart(context.Context, *moment_api.DeleteStandartRequest) error
}

type Materials interface {
	GetMatFotCalculate(ctx context.Context, markId string, temp float64) (models.MaterialsResult, error)

	GetMaterials(context.Context, *moment_api.GetMaterialsRequest) ([]*moment_api.Material, error)
	GetMaterialsWithIsEmpty(context.Context, *moment_api.GetMaterialsRequest) ([]*moment_api.MaterialWithIsEmpty, error)
	GetMaterialsData(context.Context, *moment_api.GetMaterialsDataRequest) (*moment_api.MaterialsDataResponse, error)
	CreateMaterial(context.Context, *moment_api.CreateMaterialRequest) (id string, err error)
	UpdateMaterial(context.Context, *moment_api.UpdateMaterialRequest) error
	DeleteMaterial(context.Context, *moment_api.DeleteMaterialRequest) error

	CreateVoltage(context.Context, *moment_api.CreateVoltageRequest) error
	UpdateVoltage(context.Context, *moment_api.UpdateVoltageRequest) error
	DeleteVoltage(context.Context, *moment_api.DeleteVoltageRequest) error

	CreateElasticity(context.Context, *moment_api.CreateElasticityRequest) error
	UpdateElasticity(context.Context, *moment_api.UpdateElasticityRequest) error
	DeleteElasticity(context.Context, *moment_api.DeleteElasticityRequest) error

	CreateAlpha(context.Context, *moment_api.CreateAlphaRequest) error
	UpateAlpha(context.Context, *moment_api.UpdateAlphaRequest) error
	DeleteAlpha(context.Context, *moment_api.DeleteAlphaRequest) error
}

type Gasket interface {
	GetFullData(context.Context, models.GetGasket) (models.FullDataGasket, error)

	GetData(ctx context.Context, gasket *moment_api.GetFullDataRequest) (*moment_api.FullDataResponse, error)
	GetGasket(context.Context, *moment_api.GetGasketRequest) ([]*moment_api.Gasket, error)
	GetGasketWithThick(ctx context.Context, req *moment_api.GetGasketRequest) (gasket []*moment_api.GasketWithThick, err error)
	CreateGasket(context.Context, *moment_api.CreateGasketRequest) (id string, err error)
	UpdateGasket(context.Context, *moment_api.UpdateGasketRequest) error
	DeleteGasket(context.Context, *moment_api.DeleteGasketRequest) error

	GetTypeGasket(context.Context, *moment_api.GetGasketTypeRequest) ([]*moment_api.GasketType, error)
	CreateTypeGasket(context.Context, *moment_api.CreateGasketTypeRequest) (id string, err error)
	UpdateTypeGasket(context.Context, *moment_api.UpdateGasketTypeRequest) error
	DeleteTypeGasket(context.Context, *moment_api.DeleteGasketTypeRequest) error

	GetEnv(context.Context, *moment_api.GetEnvRequest) ([]*moment_api.Env, error)
	CreateEnv(context.Context, *moment_api.CreateEnvRequest) (id string, err error)
	UpdateEnv(context.Context, *moment_api.UpdateEnvRequest) error
	DeleteEnv(context.Context, *moment_api.DeleteEnvRequest) error

	CreateManyEnvData(context.Context, *moment_api.CreateManyEnvDataRequest) error
	CreateEnvData(context.Context, *moment_api.CreateEnvDataRequest) error
	UpdateEnvData(context.Context, *moment_api.UpdateEnvDataRequest) error
	DeleteEnvData(context.Context, *moment_api.DeleteEnvDataRequest) error

	CreateManyGasketData(context.Context, *moment_api.CreateManyGasketDataRequest) error
	CreateGasketData(context.Context, *moment_api.CreateGasketDataRequest) error
	UpdateGasketTypeId(context.Context, *moment_api.UpdateGasketTypeIdRequest) error
	UpdateGasketData(context.Context, *moment_api.UpdateGasketDataRequest) error
	DeleteGasketData(context.Context, *moment_api.DeleteGasketDataRequest) error
}

type Graphic interface {
	CalculateBetaF(betta, x float64) float64
	CalculateBetaV(betta, x float64) float64
	CalculateF(betta, x float64) float64
	CalculateMkp(diameter float64, sigma float64) float64
}

type Read interface {
	read.Flange
}

type Services struct {
	CalcFlange
	CalcCap
	Flange
	Materials
	Gasket
	Graphic
	Read
}

func NewServices(repos *repository.Repositories) *Services {
	flange := flange.NewFlangeService(repos.Flange)
	materials := materials.NewMaterialsService(repos.Materials)
	gasket := gasket.NewGasketService(repos.Gasket)
	graphic := graphic.NewGraphicService()
	calc := calc.NewCalcServices(flange, gasket, materials)

	return &Services{
		CalcFlange: calc.Flange,
		CalcCap:    calc.Cap,
		Flange:     flange,
		Materials:  materials,
		Gasket:     gasket,
		Graphic:    graphic,
		Read:       read.NewReadService(flange, materials, gasket),
	}
}
