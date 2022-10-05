package read

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type Flange interface {
	Get(ctx context.Context, req *moment_api.GetFlangeRequest) (*moment_api.GetFlangeResponse, error)
}

type ReadService struct {
	Flange
}

func NewReadService(flange *flange.FlangeService, materials *materials.MaterialsService, gasket *gasket.GasketService) *ReadService {
	return &ReadService{
		Flange: NewFlangeService(flange, materials, gasket),
	}
}
