package service

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type Flange interface {
	Calculation(ctx context.Context, data *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error)
}

type Materials interface {
	GetMatFotCalculate(ctx context.Context, markId string, temp float64) (models.MaterialsResult, error)
}

type Gasket interface {
	Get(ctx context.Context, gasket models.GetGasket) (models.Gasket, error)
}

type Services struct {
	Flange
	Materials
	Gasket
}

func NewServices(repos *repository.Repositories) *Services {
	Materials := NewMaterialsService(repos.Materials)
	Gasket := NewGasketService(repos.Gasket)

	return &Services{
		Flange:    NewFlangeService(repos.Flange, Materials, Gasket),
		Materials: Materials,
		Gasket:    Gasket,
	}
}
