package service

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_api"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_type_api"
)

type RingsKitService struct {
	types         RingsKitType
	constructions RingsKitConstruction
}

func NewRingsKitService(types RingsKitType, constructions RingsKitConstruction) *RingsKitService {
	return &RingsKitService{
		types:         types,
		constructions: constructions,
	}
}

type RingsKit interface {
	Get(context.Context, *rings_kit_api.GetRingsKit) (*rings_kit_model.RingsKit, error)
}

func (s *RingsKitService) Get(ctx context.Context, req *rings_kit_api.GetRingsKit) (*rings_kit_model.RingsKit, error) {
	kitTypes, err := s.types.GetAll(ctx, &rings_kit_type_api.GetRingsKitTypes{})
	if err != nil {
		return nil, err
	}

	kitConstructions, err := s.constructions.GetAll(ctx, &rings_kit_construction_api.GetRingsKitConstructions{})
	if err != nil {
		return nil, err
	}

	ringsKit := &rings_kit_model.RingsKit{
		RingsKitTypes:   kitTypes,
		ConstructionMap: kitConstructions.Constructions,
	}

	return ringsKit, nil
}
