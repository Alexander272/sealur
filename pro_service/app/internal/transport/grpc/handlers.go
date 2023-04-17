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
		Order:          NewOrderHandlers(service.OrderNew),
		Position:       NewPositionHandlers(service.Position),
	}
}

// func (h *Handler) Ping(ctx context.Context, req *pro_api.PingRequest) (*pro_api.PingResponse, error) {
// 	return &pro_api.PingResponse{Ping: "pong"}, nil
// }
