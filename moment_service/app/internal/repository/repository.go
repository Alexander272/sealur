package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Flange interface {
	GetFlangeSize(context.Context, *moment_proto.GetFlangeSizeRequest) (models.FlangeSize, error)
	CreateFlangeSize(context.Context, *moment_proto.CreateFlangeSizeRequest) error
	UpdateFlangeSize(context.Context, *moment_proto.UpdateFlangeSizeRequest) error
	DeleteFlangeSize(context.Context, *moment_proto.DeleteFlangeSizeRequest) error

	GetBolts(context.Context, *moment_proto.GetBoltsRequest) ([]models.BoltsDTO, error)
	CreateBolt(context.Context, *moment_proto.CreateBoltRequest) error
	UpdateBolt(context.Context, *moment_proto.UpdateBoltRequest) error
	DeleteBolt(context.Context, *moment_proto.DeleteBoltRequest) error

	GetTypeFlange(context.Context, *moment_proto.GetTypeFlangeRequest) ([]models.TypeFlangeDTO, error)
	CreateTypeFlange(context.Context, *moment_proto.CreateTypeFlangeRequest) (id string, err error)
	UpdateTypeFlange(context.Context, *moment_proto.UpdateTypeFlangeRequest) error
	DeleteTypeFlange(context.Context, *moment_proto.DeleteTypeFlangeRequest) error

	GetStandarts(context.Context, *moment_proto.GetStandartsRequest) ([]models.StandartDTO, error)
	CreateStandart(context.Context, *moment_proto.CreateStandartRequest) (id string, err error)
	UpdateStandart(context.Context, *moment_proto.UpdateStandartRequest) error
	DeleteStandart(context.Context, *moment_proto.DeleteStandartRequest) error
}

type Materials interface {
	GetMaterials(context.Context, *moment_proto.GetMaterialsRequest) ([]models.MaterialsDTO, error)
	GetMaterialsWithIsEmpty(context.Context, *moment_proto.GetMaterialsRequest) ([]models.MaterialsWithIsEmpty, error)
	GetAllData(context.Context, *moment_proto.GetMaterialsDataRequest) (models.MaterialsAll, error)

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

	GetGasket(context.Context, *moment_proto.GetGasketRequest) ([]models.GasketDTO, error)
	CreateGasket(context.Context, *moment_proto.CreateGasketRequest) (id string, err error)
	UpdateGasket(context.Context, *moment_proto.UpdateGasketRequest) error
	DeleteGasket(context.Context, *moment_proto.DeleteGasketRequest) error

	GetTypeGasket(context.Context, *moment_proto.GetGasketTypeRequest) ([]models.TypeGasketDTO, error)
	CreateTypeGasket(context.Context, *moment_proto.CreateGasketTypeRequest) (id string, err error)
	UpdateTypeGasket(context.Context, *moment_proto.UpdateGasketTypeRequest) error
	DeleteTypeGasket(context.Context, *moment_proto.DeleteGasketTypeRequest) error

	GetEnv(context.Context, *moment_proto.GetEnvRequest) ([]models.EnvDTO, error)
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

type Repositories struct {
	Flange
	Materials
	Gasket
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Flange:    NewFlangeRepo(db),
		Materials: NewMaterialsRepo(db),
		Gasket:    NewGasketRepo(db),
	}
}
