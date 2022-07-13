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

type Graphic interface {
	CalculateBettaF(betta, x float64) float64
	CalculateBettaV(betta, x float64) float64
	CalculateF(betta, x float64) float64
	CalculateMkp(diameter int32, sigma float64) float64
}

type Services struct {
	Flange
	Materials
	Gasket
	Graphic
}

func NewServices(repos *repository.Repositories) *Services {
	Materials := NewMaterialsService(repos.Materials)
	Gasket := NewGasketService(repos.Gasket)
	Graphic := NewGraphicService()

	return &Services{
		Flange:    NewFlangeService(repos.Flange, Materials, Gasket, Graphic),
		Materials: Materials,
		Gasket:    Gasket,
		Graphic:   Graphic,
	}
}
