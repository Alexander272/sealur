package service

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type CalcFlange interface {
	Calculation(ctx context.Context, data *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error)
}

type Flange interface {
	GetFlangeSize(context.Context, *moment_proto.GetFlangeSizeRequest) (models.FlangeSize, error)
	CreateFlangeSize(context.Context, *moment_proto.CreateFlangeSizeRequest) error
	UpdateFlangeSize(context.Context, *moment_proto.UpdateFlangeSizeRequest) error
	DeleteFlangeSize(context.Context, *moment_proto.DeleteFlangeSizeRequest) error

	GetBolts(context.Context, *moment_proto.GetBoltsRequest) ([]*moment_proto.Bolt, error)
	CreateBolt(context.Context, *moment_proto.CreateBoltRequest) error
	UpdateBolt(context.Context, *moment_proto.UpdateBoltRequest) error
	DeleteBolt(context.Context, *moment_proto.DeleteBoltRequest) error

	GetTypeFlange(context.Context, *moment_proto.GetTypeFlangeRequest) ([]*moment_proto.TypeFlange, error)
	CreateTypeFlange(context.Context, *moment_proto.CreateTypeFlangeRequest) (id string, err error)
	UpdateTypeFlange(context.Context, *moment_proto.UpdateTypeFlangeRequest) error
	DeleteTypeFlange(context.Context, *moment_proto.DeleteTypeFlangeRequest) error

	GetStandarts(context.Context, *moment_proto.GetStandartsRequest) ([]*moment_proto.Standart, error)
	CreateStandart(context.Context, *moment_proto.CreateStandartRequest) (id string, err error)
	UpdateStandart(context.Context, *moment_proto.UpdateStandartRequest) error
	DeleteStandart(context.Context, *moment_proto.DeleteStandartRequest) error
}

type Materials interface {
	GetMatFotCalculate(ctx context.Context, markId string, temp float64) (models.MaterialsResult, error)

	GetMaterials(context.Context, *moment_proto.GetMaterialsRequest) ([]*moment_proto.Material, error)
	GetMaterialsWithIsEmpty(context.Context, *moment_proto.GetMaterialsRequest) ([]*moment_proto.MaterialWithIsEmpty, error)
	GetMaterialsData(context.Context, *moment_proto.GetMaterialsDataRequest) (*moment_proto.MaterialsDataResponse, error)
	CreateMaterial(context.Context, *moment_proto.CreateMaterialRequest) (id string, err error)
	UpdateMaterial(context.Context, *moment_proto.UpdateMaterialRequest) error
	DeleteMaterial(context.Context, *moment_proto.DeleteMaterialRequest) error

	CreateVoltage(context.Context, *moment_proto.CreateVoltageRequest) error
	UpdateVoltage(context.Context, *moment_proto.UpdateVoltageRequest) error
	DeleteVoltage(context.Context, *moment_proto.DeleteVoltageRequest) error

	CreateElasticity(context.Context, *moment_proto.CreateElasticityRequest) error
	UpdateElasticity(context.Context, *moment_proto.UpdateElasticityRequest) error
	DeleteElasticity(context.Context, *moment_proto.DeleteElasticityRequest) error

	CreateAlpha(context.Context, *moment_proto.CreateAlphaRequest) error
	UpateAlpha(context.Context, *moment_proto.UpdateAlphaRequest) error
	DeleteAlpha(context.Context, *moment_proto.DeleteAlphaRequest) error
}

type Gasket interface {
	GetFullData(context.Context, models.GetGasket) (models.FullDataGasket, error)

	GetGasket(context.Context, *moment_proto.GetGasketRequest) ([]*moment_proto.Gasket, error)
	CreateGasket(context.Context, *moment_proto.CreateGasketRequest) (id string, err error)
	UpdateGasket(context.Context, *moment_proto.UpdateGasketRequest) error
	DeleteGasket(context.Context, *moment_proto.DeleteGasketRequest) error

	GetTypeGasket(context.Context, *moment_proto.GetGasketTypeRequest) ([]*moment_proto.GasketType, error)
	CreateTypeGasket(context.Context, *moment_proto.CreateGasketTypeRequest) (id string, err error)
	UpdateTypeGasket(context.Context, *moment_proto.UpdateGasketTypeRequest) error
	DeleteTypeGasket(context.Context, *moment_proto.DeleteGasketTypeRequest) error

	GetEnv(context.Context, *moment_proto.GetEnvRequest) ([]*moment_proto.Env, error)
	CreateEnv(context.Context, *moment_proto.CreateEnvRequest) (id string, err error)
	UpdateEnv(context.Context, *moment_proto.UpdateEnvRequest) error
	DeleteEnv(context.Context, *moment_proto.DeleteEnvRequest) error

	CreateEnvData(context.Context, *moment_proto.CreateEnvDataRequest) error
	UpdateEnvData(context.Context, *moment_proto.UpdateEnvDataRequest) error
	DeleteEnvData(context.Context, *moment_proto.DeleteEnvDataRequest) error

	CreateGasketData(context.Context, *moment_proto.CreateGasketDataRequest) error
	UpdateGasketData(context.Context, *moment_proto.UpdateGasketDataRequest) error
	DeleteGasketData(context.Context, *moment_proto.DeleteGasketDataRequest) error
}

type Graphic interface {
	CalculateBettaF(betta, x float64) float64
	CalculateBettaV(betta, x float64) float64
	CalculateF(betta, x float64) float64
	CalculateMkp(diameter int32, sigma float64) float64
}

type Services struct {
	CalcFlange
	Flange
	Materials
	Gasket
	Graphic
}

func NewServices(repos *repository.Repositories) *Services {
	flange := NewFlangeService(repos.Flange)
	materials := NewMaterialsService(repos.Materials)
	gasket := NewGasketService(repos.Gasket)
	graphic := NewGraphicService()

	return &Services{
		CalcFlange: NewCalcFlangeService(flange, materials, gasket, graphic),
		Flange:     flange,
		Materials:  materials,
		Gasket:     gasket,
		Graphic:    graphic,
	}
}
