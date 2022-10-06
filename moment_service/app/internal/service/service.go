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
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/flange_model"
	"github.com/Alexander272/sealur_proto/api/moment/models/gasket_model"
	"github.com/Alexander272/sealur_proto/api/moment/models/material_model"
)

type Calc interface {
	CalculationFlange(ctx context.Context, data *calc_api.FlangeRequest) (*calc_api.FlangeResponse, error)
	CalculationCap(ctx context.Context, data *calc_api.CapRequest) (*calc_api.CapResponse, error)
}

type Flange interface {
	GetFlangeSize(context.Context, *flange_api.GetFlangeSizeRequest) (models.FlangeSize, error)
	GetFullFlangeSize(context.Context, *flange_api.GetFullFlangeSizeRequest) (*flange_api.FullFlangeSizeResponse, error)
	GetBasisFlangeSize(context.Context, *flange_api.GetBasisFlangeSizeRequest) (*flange_model.BasisFlangeSizeResponse, error)
	CreateFlangeSize(context.Context, *flange_api.CreateFlangeSizeRequest) error
	CreateFlangeSizes(context.Context, *flange_api.CreateFlangeSizesRequest) error
	UpdateFlangeSize(context.Context, *flange_api.UpdateFlangeSizeRequest) error
	DeleteFlangeSize(context.Context, *flange_api.DeleteFlangeSizeRequest) error

	GetBolts(context.Context, *flange_api.GetBoltsRequest) ([]*flange_model.Bolt, error)
	GetAllBolts(context.Context, *flange_api.GetBoltsRequest) ([]*flange_model.Bolt, error)
	CreateBolt(context.Context, *flange_api.CreateBoltRequest) error
	CreateBolts(context.Context, *flange_api.CreateBoltsRequest) error
	UpdateBolt(context.Context, *flange_api.UpdateBoltRequest) error
	DeleteBolt(context.Context, *flange_api.DeleteBoltRequest) error

	GetTypeFlange(context.Context, *flange_api.GetTypeFlangeRequest) ([]*flange_model.TypeFlange, error)
	CreateTypeFlange(context.Context, *flange_api.CreateTypeFlangeRequest) (id string, err error)
	UpdateTypeFlange(context.Context, *flange_api.UpdateTypeFlangeRequest) error
	DeleteTypeFlange(context.Context, *flange_api.DeleteTypeFlangeRequest) error

	GetStandarts(context.Context, *flange_api.GetStandartsRequest) ([]*flange_model.Standart, error)
	GetStandartsWithSize(context.Context, *flange_api.GetStandartsRequest) ([]*flange_model.StandartWithSize, error)
	CreateStandart(context.Context, *flange_api.CreateStandartRequest) (id string, err error)
	UpdateStandart(context.Context, *flange_api.UpdateStandartRequest) error
	DeleteStandart(context.Context, *flange_api.DeleteStandartRequest) error
}

type Materials interface {
	GetMatFotCalculate(ctx context.Context, markId string, temp float64) (models.MaterialsResult, error)

	GetMaterials(context.Context, *material_api.GetMaterialsRequest) ([]*material_model.Material, error)
	GetMaterialsWithIsEmpty(context.Context, *material_api.GetMaterialsRequest) ([]*material_model.MaterialWithIsEmpty, error)
	GetMaterialsData(context.Context, *material_api.GetMaterialsDataRequest) (*material_api.MaterialsDataResponse, error)
	CreateMaterial(context.Context, *material_api.CreateMaterialRequest) (id string, err error)
	UpdateMaterial(context.Context, *material_api.UpdateMaterialRequest) error
	DeleteMaterial(context.Context, *material_api.DeleteMaterialRequest) error

	CreateVoltage(context.Context, *material_api.CreateVoltageRequest) error
	UpdateVoltage(context.Context, *material_api.UpdateVoltageRequest) error
	DeleteVoltage(context.Context, *material_api.DeleteVoltageRequest) error

	CreateElasticity(context.Context, *material_api.CreateElasticityRequest) error
	UpdateElasticity(context.Context, *material_api.UpdateElasticityRequest) error
	DeleteElasticity(context.Context, *material_api.DeleteElasticityRequest) error

	CreateAlpha(context.Context, *material_api.CreateAlphaRequest) error
	UpateAlpha(context.Context, *material_api.UpdateAlphaRequest) error
	DeleteAlpha(context.Context, *material_api.DeleteAlphaRequest) error
}

type Gasket interface {
	GetFullData(context.Context, models.GetGasket) (models.FullDataGasket, error)

	GetData(ctx context.Context, gasket *gasket_api.GetFullDataRequest) (*gasket_api.FullDataResponse, error)
	GetGasket(context.Context, *gasket_api.GetGasketRequest) ([]*gasket_model.Gasket, error)
	GetGasketWithThick(ctx context.Context, req *gasket_api.GetGasketRequest) (gasket []*gasket_model.GasketWithThick, err error)
	CreateGasket(context.Context, *gasket_api.CreateGasketRequest) (id string, err error)
	UpdateGasket(context.Context, *gasket_api.UpdateGasketRequest) error
	DeleteGasket(context.Context, *gasket_api.DeleteGasketRequest) error

	GetTypeGasket(context.Context, *gasket_api.GetGasketTypeRequest) ([]*gasket_model.GasketType, error)
	CreateTypeGasket(context.Context, *gasket_api.CreateGasketTypeRequest) (id string, err error)
	UpdateTypeGasket(context.Context, *gasket_api.UpdateGasketTypeRequest) error
	DeleteTypeGasket(context.Context, *gasket_api.DeleteGasketTypeRequest) error

	GetEnv(context.Context, *gasket_api.GetEnvRequest) ([]*gasket_model.Env, error)
	CreateEnv(context.Context, *gasket_api.CreateEnvRequest) (id string, err error)
	UpdateEnv(context.Context, *gasket_api.UpdateEnvRequest) error
	DeleteEnv(context.Context, *gasket_api.DeleteEnvRequest) error

	CreateManyEnvData(context.Context, *gasket_api.CreateManyEnvDataRequest) error
	CreateEnvData(context.Context, *gasket_api.CreateEnvDataRequest) error
	UpdateEnvData(context.Context, *gasket_api.UpdateEnvDataRequest) error
	DeleteEnvData(context.Context, *gasket_api.DeleteEnvDataRequest) error

	CreateManyGasketData(context.Context, *gasket_api.CreateManyGasketDataRequest) error
	CreateGasketData(context.Context, *gasket_api.CreateGasketDataRequest) error
	UpdateGasketTypeId(context.Context, *gasket_api.UpdateGasketTypeIdRequest) error
	UpdateGasketData(context.Context, *gasket_api.UpdateGasketDataRequest) error
	DeleteGasketData(context.Context, *gasket_api.DeleteGasketDataRequest) error
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
	Calc
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
		Calc:      calc,
		Flange:    flange,
		Materials: materials,
		Gasket:    gasket,
		Graphic:   graphic,
		Read:      read.NewReadService(flange, materials, gasket),
	}
}
