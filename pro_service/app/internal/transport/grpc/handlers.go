package grpc

import (
	"github.com/Alexander272/sealur/pro_service/internal/config"
	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_base_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_conf_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_base_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_density_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_modifying_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
)

type FlangeStandard interface {
	flange_standard_api.FlangeStandardServiceServer
}
type FlangeType interface {
	flange_type_api.FlangeTypeServiceServer
}
type FlangeTypeSnp interface {
	flange_type_snp_api.FlangeTypeSnpServiceServer
}
type Material interface {
	material_api.MaterialServiceServer
}
type Mounting interface {
	mounting_api.MountingServiceServer
}

type SnpData interface {
	snp_data_api.SnpDataServiceServer
}
type SnpFiller interface {
	snp_filler_api.SnpFillerServiceServer
}
type SnpMaterial interface {
	snp_material_api.SnpMaterialServiceServer
}
type SnpStandard interface {
	snp_standard_api.SnpStandardServiceServer
}
type SnpType interface {
	snp_type_api.SnpTypeServiceServer
}
type Standard interface {
	standard_api.StandardServiceServer
}
type Temperature interface {
	temperature_api.TemperatureServiceServer
}
type SnpSize interface {
	snp_size_api.SnpSizeServiceServer
}
type Snp interface {
	snp_api.SnpDataServiceServer
}

type Putg interface {
	putg_api.PutgDataServiceServer
}
type PutgConfiguration interface {
	putg_conf_api.PutgConfigurationServiceServer
}
type PutgConstruction interface {
	putg_construction_api.PutgConstructionServiceServer
}
type PutgBaseConstruction interface {
	putg_base_construction_api.PutgBaseConstructionServiceServer
}
type PutgData interface {
	putg_data_api.PutgDataServiceServer
}
type PutgBaseFiller interface {
	putg_filler_base_api.PutgBaseFillerServiceServer
}
type PutgFiller interface {
	putg_filler_api.PutgFillerServiceServer
}
type PutgFlangeType interface {
	putg_flange_type_api.PutgFlangeTypeServiceServer
}
type PutgMaterial interface {
	putg_material_api.PutgMaterialServiceServer
}
type PutgStandard interface {
	putg_standard_api.PutgStandardServiceServer
}
type PutgSize interface {
	putg_size_api.PutgSizeServiceServer
}
type PutgType interface {
	putg_type_api.PutgTypeServiceServer
}

type RingConstruction interface {
	ring_construction_api.RingConstructionServiceServer
}
type RingDensity interface {
	ring_density_api.RingDensityServiceServer
}
type RingMaterial interface {
	ring_material_api.RingMaterialServiceServer
}
type RingModifying interface {
	ring_modifying_api.RingModifyingServiceServer
}
type RingType interface {
	ring_type_api.RingTypeServiceServer
}
type RingSize interface {
	ring_size_api.RingSizeServiceServer
}
type Ring interface {
	ring_api.RingServiceServer
}

type Order interface {
	order_api.OrderServiceServer
}
type Position interface {
	position_api.PositionServiceServer
}

type Handler struct {
	service *service.Services
	conf    config.ApiConfig

	FlangeStandard
	FlangeType
	FlangeTypeSnp
	Material
	Mounting
	SnpData
	SnpFiller
	SnpMaterial
	SnpStandard
	SnpType
	Standard
	Temperature
	SnpSize
	Snp

	Putg
	PutgConfiguration
	PutgBaseConstruction
	PutgConstruction
	PutgData
	PutgBaseFiller
	PutgFiller
	PutgFlangeType
	PutgMaterial
	PutgStandard
	PutgSize
	PutgType

	RingConstruction
	RingDensity
	RingType
	RingMaterial
	RingModifying
	RingSize
	Ring

	Order
	Position
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,

		FlangeStandard: NewFlangeStandardHandlers(service.FlangeStandard),
		FlangeType:     NewFlangeTypeHandlers(service.FlangeType),
		FlangeTypeSnp:  NewFlangeTypeSnpHandlers(service.FlangeTypeSnp),
		Material:       NewMaterialHandlers(service.Material),
		Mounting:       NewMountingHandlers(service.Mounting),
		SnpData:        NewSnpDataHandlers(service.SnpData),
		SnpFiller:      NewSnpFillerHandlers(service.SnpFiller),
		SnpMaterial:    NewSnpMaterialHandlers(service.SnpMaterial),
		SnpStandard:    NewSnpStandardHandlers(service.SnpStandard),
		SnpType:        NewSnpTypeHandlers(service.SnpType),
		Standard:       NewStandardHandlers(service.Standard),
		Temperature:    NewTemperatureHandlers(service.Temperature),
		SnpSize:        NewSnpSizeHandlers(service.SnpSize),
		Snp:            NewSnpHandlers(service.Snp),

		Putg:                 NewPutgHandlers(service.Putg),
		PutgConfiguration:    NewPutgConfigurationHandlers(service.PutgConfiguration),
		PutgBaseConstruction: NewPutgBaseConstructionHandlers(service.PutgBaseConstruction),
		PutgConstruction:     NewPutgConstructionHandlers(service.PutgConstruction),
		PutgData:             NewPutgDataHandlers(service.PutgData),
		PutgBaseFiller:       NewPutgBaseFillerHandlers(service.PutgBaseFiller),
		PutgFiller:           NewPutgFillerHandlers(service.PutgFiller),
		PutgFlangeType:       NewPutgFlangeTypeHandlers(service.PutgFlangeType),
		PutgMaterial:         NewPutgMaterialHandlers(service.PutgMaterial),
		PutgStandard:         NewPutgStandardHandlers(service.PutgStandard),
		PutgSize:             NewPutgSizeHandlers(service.PutgSize),
		PutgType:             NewPutgTypeHandlers(service.PutgType),

		RingConstruction: NewRingConstructionHandlers(service.RingConstruction),
		RingDensity:      NewRingDensityHandlers(service.RingDensity),
		RingType:         NewRingTypeHandlers(service.RingType),
		RingMaterial:     NewRingMaterialHandlers(service.RingMaterial),
		RingModifying:    NewRingModifyingHandlers(service.RingModifying),
		RingSize:         NewRingSizeHandlers(service.RingSize),
		Ring:             NewRingHandlers(service.Ring),

		Order:    NewOrderHandlers(service.OrderNew),
		Position: NewPositionHandlers(service.Position),
	}
}

// func (h *Handler) Ping(ctx context.Context, req *pro_api.PingRequest) (*pro_api.PingResponse, error) {
// 	return &pro_api.PingResponse{Ping: "pong"}, nil
// }
