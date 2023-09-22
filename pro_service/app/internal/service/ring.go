package service

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro/models/ring_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_construction_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_density_api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_type_api"
)

type RingService struct {
	constructions RingConstruction
	density       RingDensity
	ringTypes     RingType
}

func NewRingService(constructions RingConstruction, density RingDensity, ringTypes RingType) *RingService {
	return &RingService{
		constructions: constructions,
		density:       density,
		ringTypes:     ringTypes,
	}
}

type Ring interface {
	Get(ctx context.Context, req *ring_api.GetRings) (*ring_model.Ring, error)
}

func (s *RingService) Get(ctx context.Context, req *ring_api.GetRings) (*ring_model.Ring, error) {
	ringTypes, err := s.ringTypes.GetAll(ctx, &ring_type_api.GetRingTypes{})
	if err != nil {
		return nil, err
	}

	constructions, err := s.constructions.GetAll(ctx, &ring_construction_api.GetRingConstructions{})
	if err != nil {
		return nil, err
	}

	density, err := s.density.GetAll(ctx, &ring_density_api.GetRingDensity{})
	if err != nil {
		return nil, err
	}

	rings := &ring_model.Ring{
		RingTypes:        ringTypes,
		ConstructionsMap: constructions.Constructions,
		DensityMap:       density.Density,
	}

	return rings, nil
}
