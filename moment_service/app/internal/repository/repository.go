package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/jmoiron/sqlx"
)

type Flange interface {
	GetFlangeSize(context.Context, *flange_api.GetFlangeSizeRequest) (models.FlangeSize, error)
	GetBasisFlangeSizes(context.Context, models.GetBasisSize) ([]models.FlangeSize, error)
	GetFullFlangeSize(context.Context, *flange_api.GetFullFlangeSizeRequest, int32) ([]models.FlangeSizeDTO, error)
	CreateFlangeSize(context.Context, *flange_api.CreateFlangeSizeRequest) error
	CreateFlangeSizes(context.Context, *flange_api.CreateFlangeSizesRequest) error
	UpdateFlangeSize(context.Context, *flange_api.UpdateFlangeSizeRequest) error
	DeleteFlangeSize(context.Context, *flange_api.DeleteFlangeSizeRequest) error

	GetBolt(ctx context.Context, boltId string) (bolt models.BoltSize, err error)
	GetBolts(context.Context, *flange_api.GetBoltsRequest) ([]models.BoltsDTO, error)
	GetAllBolts(context.Context, *flange_api.GetBoltsRequest) ([]models.BoltsDTO, error)
	CreateBolt(context.Context, *flange_api.CreateBoltRequest) error
	CreateBolts(context.Context, *flange_api.CreateBoltsRequest) error
	UpdateBolt(context.Context, *flange_api.UpdateBoltRequest) error
	DeleteBolt(context.Context, *flange_api.DeleteBoltRequest) error

	GetTypeFlange(context.Context, *flange_api.GetTypeFlangeRequest) ([]models.TypeFlangeDTO, error)
	CreateTypeFlange(context.Context, *flange_api.CreateTypeFlangeRequest) (id string, err error)
	UpdateTypeFlange(context.Context, *flange_api.UpdateTypeFlangeRequest) error
	DeleteTypeFlange(context.Context, *flange_api.DeleteTypeFlangeRequest) error

	GetStandarts(context.Context, *flange_api.GetStandartsRequest) ([]models.StandartDTO, error)
	CreateStandart(context.Context, *flange_api.CreateStandartRequest) (id string, err error)
	UpdateStandart(context.Context, *flange_api.UpdateStandartRequest) error
	DeleteStandart(context.Context, *flange_api.DeleteStandartRequest) error
}

type Materials interface {
	GetMaterials(context.Context, *material_api.GetMaterialsRequest) ([]models.MaterialsDTO, error)
	GetMaterialsWithIsEmpty(context.Context, *material_api.GetMaterialsRequest) ([]models.MaterialsWithIsEmpty, error)
	GetAllData(context.Context, *material_api.GetMaterialsDataRequest) (models.MaterialsAll, error)

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

	GetGasket(context.Context, *gasket_api.GetGasketRequest) ([]models.GasketDTO, error)
	GetGasketWithThick(ctx context.Context, req *gasket_api.GetGasketRequest) (gasket []models.GasketWithThick, err error)
	CreateGasket(context.Context, *gasket_api.CreateGasketRequest) (id string, err error)
	UpdateGasket(context.Context, *gasket_api.UpdateGasketRequest) error
	DeleteGasket(context.Context, *gasket_api.DeleteGasketRequest) error

	GetTypeGasket(context.Context, *gasket_api.GetGasketTypeRequest) ([]models.TypeGasketDTO, error)
	CreateTypeGasket(context.Context, *gasket_api.CreateGasketTypeRequest) (id string, err error)
	UpdateTypeGasket(context.Context, *gasket_api.UpdateGasketTypeRequest) error
	DeleteTypeGasket(context.Context, *gasket_api.DeleteGasketTypeRequest) error

	GetEnv(context.Context, *gasket_api.GetEnvRequest) ([]models.EnvDTO, error)
	CreateEnv(context.Context, *gasket_api.CreateEnvRequest) (id string, err error)
	UpdateEnv(context.Context, *gasket_api.UpdateEnvRequest) error
	DeleteEnv(context.Context, *gasket_api.DeleteEnvRequest) error

	CreateManyEnvData(context.Context, *gasket_api.CreateManyEnvDataRequest) error
	GetEnvData(context.Context, string) ([]models.EnvDataDTO, error)
	CreateEnvData(context.Context, *gasket_api.CreateEnvDataRequest) error
	UpdateEnvData(context.Context, *gasket_api.UpdateEnvDataRequest) error
	DeleteEnvData(context.Context, *gasket_api.DeleteEnvDataRequest) error

	CreateManyGasketData(context.Context, *gasket_api.CreateManyGasketDataRequest) error
	GetGasketData(context.Context, string) ([]models.GasketDataDTO, error)
	CreateGasketData(context.Context, *gasket_api.CreateGasketDataRequest) error
	UpdateGasketTypeId(context.Context, *gasket_api.UpdateGasketTypeIdRequest) error
	UpdateGasketData(context.Context, *gasket_api.UpdateGasketDataRequest) error
	DeleteGasketData(context.Context, *gasket_api.DeleteGasketDataRequest) error
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
