package models

import "github.com/lib/pq"

type RingsKit struct {
	Id               string         `db:"id"`
	PositionId       string         `db:"position_id"`
	TypeId           string         `db:"type_id"`
	TypeCode         string         `db:"type_code"`
	Count            string         `db:"count"`
	Size             string         `db:"size"`
	Thickness        string         `db:"thickness"`
	Material         string         `db:"materials"`
	Modifying        string         `db:"modifying"`
	Drawing          string         `db:"drawing"`
	ConstructionId   string         `db:"construction_id"`
	ConstructionCode string         `db:"construction_code"`
	Title            string         `db:"title"`
	Image            string         `db:"image"`
	SameRings        bool           `db:"same_rings"`
	MaterialTypes    string         `db:"material_types"`
	HasThickness     bool           `db:"has_thickness"`
	EnabledMaterials pq.StringArray `db:"enabled_materials"`
}
