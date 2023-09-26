package models

import "github.com/lib/pq"

type RingsKitConstruction struct {
	Id               string         `db:"id"`
	TypeId           string         `db:"type_id"`
	Code             string         `db:"code"`
	Title            string         `db:"title"`
	Image            string         `db:"image"`
	SameRings        bool           `db:"same_rings"`
	MaterialTypes    string         `db:"material_types"`
	HasThickness     bool           `db:"has_thickness"`
	DefaultCount     string         `db:"default_count"`
	DefaultMaterials string         `db:"default_materials"`
	EnabledMaterials pq.StringArray `db:"enabled_materials"`
}
