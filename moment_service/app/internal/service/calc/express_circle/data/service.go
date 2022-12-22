package data

import (
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
)

type DataService struct {
	flange    *flange.FlangeService
	materials *materials.MaterialsService
	gasket    *gasket.GasketService
	graphic   *graphic.GraphicService
}

func NewDataService(flange *flange.FlangeService, materials *materials.MaterialsService,
	gasket *gasket.GasketService, graphic *graphic.GraphicService) *DataService {

	return &DataService{
		flange:    flange,
		materials: materials,
		gasket:    gasket,
		graphic:   graphic,
	}
}
