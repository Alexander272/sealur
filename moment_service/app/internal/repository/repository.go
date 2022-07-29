package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/jmoiron/sqlx"
)

type Flange interface {
	GetFlangeSize(context.Context, *moment_api.GetFlangeSizeRequest) (models.FlangeSize, error)
	GetBasisFlangeSizes(context.Context, models.GetBasisSize) ([]models.FlangeSize, error)
	CreateFlangeSize(context.Context, *moment_api.CreateFlangeSizeRequest) error
	UpdateFlangeSize(context.Context, *moment_api.UpdateFlangeSizeRequest) error
	DeleteFlangeSize(context.Context, *moment_api.DeleteFlangeSizeRequest) error

	GetBolts(context.Context, *moment_api.GetBoltsRequest) ([]models.BoltsDTO, error)
	CreateBolt(context.Context, *moment_api.CreateBoltRequest) error
	UpdateBolt(context.Context, *moment_api.UpdateBoltRequest) error
	DeleteBolt(context.Context, *moment_api.DeleteBoltRequest) error

	GetTypeFlange(context.Context, *moment_api.GetTypeFlangeRequest) ([]models.TypeFlangeDTO, error)
	CreateTypeFlange(context.Context, *moment_api.CreateTypeFlangeRequest) (id string, err error)
	UpdateTypeFlange(context.Context, *moment_api.UpdateTypeFlangeRequest) error
	DeleteTypeFlange(context.Context, *moment_api.DeleteTypeFlangeRequest) error

	GetStandarts(context.Context, *moment_api.GetStandartsRequest) ([]models.StandartDTO, error)
	CreateStandart(context.Context, *moment_api.CreateStandartRequest) (id string, err error)
	UpdateStandart(context.Context, *moment_api.UpdateStandartRequest) error
	DeleteStandart(context.Context, *moment_api.DeleteStandartRequest) error
}

type Materials interface {
	GetMaterials(context.Context, *moment_api.GetMaterialsRequest) ([]models.MaterialsDTO, error)
	GetMaterialsWithIsEmpty(context.Context, *moment_api.GetMaterialsRequest) ([]models.MaterialsWithIsEmpty, error)
	GetAllData(context.Context, *moment_api.GetMaterialsDataRequest) (models.MaterialsAll, error)

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

	GetGasket(context.Context, *moment_api.GetGasketRequest) ([]models.GasketDTO, error)
	GetGasketWithThick(ctx context.Context, req *moment_api.GetGasketRequest) (gasket []models.GasketWithThick, err error)
	CreateGasket(context.Context, *moment_api.CreateGasketRequest) (id string, err error)
	UpdateGasket(context.Context, *moment_api.UpdateGasketRequest) error
	DeleteGasket(context.Context, *moment_api.DeleteGasketRequest) error

	GetTypeGasket(context.Context, *moment_api.GetGasketTypeRequest) ([]models.TypeGasketDTO, error)
	CreateTypeGasket(context.Context, *moment_api.CreateGasketTypeRequest) (id string, err error)
	UpdateTypeGasket(context.Context, *moment_api.UpdateGasketTypeRequest) error
	DeleteTypeGasket(context.Context, *moment_api.DeleteGasketTypeRequest) error

	GetEnv(context.Context, *moment_api.GetEnvRequest) ([]models.EnvDTO, error)
	CreateEnv(context.Context, *moment_api.CreateEnvRequest) (id string, err error)
	UpdateEnv(context.Context, *moment_api.UpdateEnvRequest) error
	DeleteEnv(context.Context, *moment_api.DeleteEnvRequest) error

	CreateEnvData(context.Context, *moment_api.CreateEnvDataRequest) error
	UpdateEnvData(context.Context, *moment_api.UpdateEnvDataRequest) error
	DeleteEnvData(context.Context, *moment_api.DeleteEnvDataRequest) error

	CreateGasketData(context.Context, *moment_api.CreateGasketDataRequest) error
	UpdateGasketData(context.Context, *moment_api.UpdateGasketDataRequest) error
	DeleteGasketData(context.Context, *moment_api.DeleteGasketDataRequest) error
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
