package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Flange interface {
	GetSize(ctx context.Context, dn, pn float64, standId string) (models.FlangeSize, error)
}

type Materials interface {
	GetMaterials(context.Context, *moment_proto.GetMaterialsRequest) (materials []models.MaterialsDTO, err error)
	GetAllData(context.Context, string) (models.MaterialsAll, error)

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
	Get(ctx context.Context, gasket models.GetGasket) (models.Gasket, error)
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
