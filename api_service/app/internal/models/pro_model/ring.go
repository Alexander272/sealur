package pro_model

import "github.com/Alexander272/sealur_proto/api/pro/models/position_model"

type RingData struct {
	// TypeId               string `json:"typeId"`
	// TypeCode             string `json:"typeCode"`
	// DensityId            string `json:"densityId"`
	// DensityCode          string `json:"densityCode"`
	// ConstructionCode     string `json:"constructionCode"`
	// ConstructionWRP      bool   `json:"constructionWRP"`
	// ConstructionBaseCode string `json:"constructionBaseCode"`
	RingType     RingType     `json:"ringType"`
	Density      Density      `json:"density"`
	Construction Construction `json:"construction"`
	Size         string       `json:"size"`
	Thickness    string       `json:"thickness"`
	Material     string       `json:"material"`
	Modifying    string       `json:"modifying"`
	Drawing      string       `json:"drawing"`
}

func (s *RingData) Parse() *position_model.PositionRing {
	return &position_model.PositionRing{
		// TypeId:               s.TypeId,
		// TypeCode:             s.TypeCode,
		// DensityId:            s.DensityId,
		// DensityCode:          s.DensityCode,
		// ConstructionCode:     s.ConstructionCode,
		// ConstructionWRP:      s.ConstructionWRP,
		// ConstructionBaseCode: s.ConstructionBaseCode,
		TypeId:               s.RingType.Id,
		TypeCode:             s.RingType.Code,
		DensityId:            s.Density.Id,
		DensityCode:          s.Density.Code,
		ConstructionCode:     s.Construction.Code,
		ConstructionWRP:      s.Construction.WithoutRotaryPlug,
		ConstructionBaseCode: s.Construction.BaseCode,
		Size:                 s.Size,
		Thickness:            s.Thickness,
		Material:             s.Material,
		Modifying:            s.Modifying,
		Drawing:              s.Drawing,
	}
}

type RingType struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

type Density struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

type Construction struct {
	Code              string `json:"code"`
	WithoutRotaryPlug bool   `json:"withoutRotaryPlug"`
	BaseCode          string `json:"baseCode"`
}
