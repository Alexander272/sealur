package service

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro/models/mounting_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_configuration_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_conf_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_type_api"
)

type PutgService struct {
	configuration PutgConfiguration
	construction  PutgConstruction
	filler        PutgFiller
	flangeType    PutgFlangeType
	standard      PutgStandard
	materials     PutgMaterial
	data          PutgData
	sizes         PutgSize
	putgType      PutgType
	mounting      Mounting
}

type PutgDeps struct {
	Configuration PutgConfiguration
	Construction  PutgConstruction
	Filler        PutgFiller
	FlangeType    PutgFlangeType
	Standard      PutgStandard
	Materials     PutgMaterial
	Data          PutgData
	Sizes         PutgSize
	PutgType      PutgType
	Mounting      Mounting
}

func NewPutgService(deps PutgDeps) *PutgService {
	return &PutgService{
		configuration: deps.Configuration,
		construction:  deps.Construction,
		filler:        deps.Filler,
		flangeType:    deps.FlangeType,
		standard:      deps.Standard,
		materials:     deps.Materials,
		data:          deps.Data,
		sizes:         deps.Sizes,
		putgType:      deps.PutgType,
		mounting:      deps.Mounting,
	}
}

func (s *PutgService) GetBase(ctx context.Context, req *putg_api.GetPutgBase) (*putg_api.PutgBase, error) {
	var standards []*putg_standard_model.PutgStandard
	var mounting []*mounting_model.Mounting
	var configurations []*putg_configuration_model.PutgConfiguration

	if req.StandardId == "" {
		var err error
		standards, err = s.standard.Get(ctx, &putg_standard_api.GetPutgStandard{})
		if err != nil {
			return nil, err
		}
		mounting, err = s.mounting.GetAll(ctx, &mounting_api.GetAllMountings{})
		if err != nil {
			return nil, err
		}
		configurations, err = s.configuration.Get(ctx, &putg_conf_api.GetPutgConfiguration{})
		if err != nil {
			return nil, err
		}

		req.StandardId = standards[0].Id
	}

	// var flangeTypes []*putg_flange_type_model.PutgFlangeType
	// var materials *putg_material_model.PutgMaterials

	// if req.FlangeTypeId == "" {
	// 	var err error
	flangeTypes, err := s.flangeType.Get(ctx, &putg_flange_type_api.GetPutgFlangeType{StandardId: req.StandardId})
	if err != nil {
		return nil, err
	}
	materials, err := s.materials.Get(ctx, &putg_material_api.GetPutgMaterial{StandardId: req.StandardId})
	if err != nil {
		return nil, err
	}
	fillers, err := s.filler.GetNew(ctx, &putg_filler_api.GetPutgFiller_New{StandardId: req.StandardId})
	if err != nil {
		return nil, err
	}

	// req.FlangeTypeId = flangeTypes[0].Id
	// }

	// constructions, err := s.construction.Get(ctx, &putg_construction_api.GetPutgConstruction{StandardId: req.StandardId, FlangeTypeId: req.FlangeTypeId})
	// if err != nil {
	// 	return nil, err
	// }

	putgBase := &putg_api.PutgBase{
		Configurations: configurations,
		Standards:      standards,
		Mounting:       mounting,
		// Constructions:  constructions,
		FlangeTypes: flangeTypes,
		Materials:   materials,
		Fillers:     fillers,
	}

	return putgBase, nil
}

func (s *PutgService) GetData(ctx context.Context, req *putg_api.GetPutgData) (*putg_api.PutgData, error) {
	// flangeTypes, err := s.flangeType.Get(ctx, &putg_flange_type_api.GetPutgFlangeType{StandardId: req.StandardId})
	// if err != nil {
	// 	return nil, err
	// }
	// materials, err := s.materials.Get(ctx, &putg_material_api.GetPutgMaterial{StandardId: req.StandardId})
	// if err != nil {
	// 	return nil, err
	// }
	// TODO нужно свои размеры отдавать овальным и прямоугольным прокладкам
	sizes, err := s.sizes.Get(ctx, &putg_size_api.GetPutgSize{PutgStandardId: req.StandardId, ConstructionId: req.ConstructionId})
	if err != nil {
		return nil, err
	}

	//TODO привязать к стандарту
	// fillers, err := s.filler.Get(ctx, &putg_filler_api.GetPutgFiller{ConstructionId: req.BaseConstructionId})
	// if err != nil {
	// 	return nil, err
	// }

	// data, err := s.data.GetByConstruction(ctx, &putg_data_api.GetPutgData{ConstructionId: req.ConstructionId})
	// if err != nil {
	// 	return nil, err
	// }

	putgData := &putg_api.PutgData{
		// FlangeTypes: flangeTypes,
		// Materials:   materials,
		Sizes: sizes,
		// Fillers: fillers,
		// Data:        data,
	}

	return putgData, nil
}

func (s *PutgService) Get(ctx context.Context, req *putg_api.GetPutg) (*putg_api.Putg, error) {
	types, err := s.putgType.Get(ctx, &putg_type_api.GetPutgType{BaseId: req.BaseId})
	if err != nil {
		return nil, err
	}
	data, err := s.data.Get(ctx, &putg_data_api.GetPutgData{FillerId: req.FillerId})
	if err != nil {
		return nil, err
	}
	constructions, err := s.construction.GetNew(ctx, &putg_construction_api.GetPutgConstruction_New{FillerId: req.BaseId, FlangeTypeId: req.FlangeTypeId})
	if err != nil {
		return nil, err
	}

	putg := &putg_api.Putg{
		PutgTypes:     types,
		Data:          data,
		Constructions: constructions,
	}

	return putg, nil
}
