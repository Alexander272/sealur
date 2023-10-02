package pro_model

import "github.com/Alexander272/sealur_proto/api/pro/models/position_model"

type KitData struct {
	TypeId       string          `json:"typeId"`
	Type         string          `json:"type"`
	Construction KitConstruction `json:"construction"`
	Count        string          `json:"count"`
	Size         string          `json:"sizes"`
	Thickness    string          `json:"thickness"`
	Material     string          `json:"materials"`
	Modifying    string          `json:"modifying"`
	Drawing      string          `json:"drawing"`
}

func (d *KitData) Parse() *position_model.PositionRingsKit {
	return &position_model.PositionRingsKit{
		TypeId:           d.TypeId,
		Type:             d.Type,
		ConstructionId:   d.Construction.Id,
		ConstructionCode: d.Construction.Code,
		Count:            d.Count,
		Size:             d.Size,
		Thickness:        d.Thickness,
		Material:         d.Material,
		Modifying:        d.Modifying,
		Drawing:          d.Drawing,
	}
}

type KitConstruction struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	// Title string `json:"title"`
}
