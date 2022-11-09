package data

import (
	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
)

type DataService struct {
	flange        *flange.FlangeService
	materials     *materials.MaterialsService
	gasket        *gasket.GasketService
	graphic       *graphic.GraphicService
	typeFlangesTF map[string]float64
	typeFlangesTB map[string]float64
	typeFlangesTK map[string]float64
}

func NewDataService(flange *flange.FlangeService, materials *materials.MaterialsService,
	gasket *gasket.GasketService, graphic *graphic.GraphicService) *DataService {
	flangesTF := map[string]float64{
		"isolated":    constants.IsolatedFlatTf,
		"nonIsolated": constants.NonIsolatedFlatTf,
	}

	flangesTB := map[string]float64{
		"isolated":         constants.IsolatedFlatTb,
		"nonIsolated":      constants.NonIsolatedFlatTb,
		"isolated-free":    constants.IsolatedFreeTb,
		"nonIsolated-free": constants.NonIsolatedFlatTb,
	}

	flangeTK := map[string]float64{
		"isolated":    constants.IsolatedFreeTk,
		"nonIsolated": constants.NonIsolatedFreeTk,
	}

	return &DataService{
		flange:        flange,
		materials:     materials,
		gasket:        gasket,
		graphic:       graphic,
		typeFlangesTF: flangesTF,
		typeFlangesTB: flangesTB,
		typeFlangesTK: flangeTK,
	}
}
