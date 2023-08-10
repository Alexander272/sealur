package models

type RingType struct {
	Id            string `db:"id"`
	Code          string `db:"code"`
	Title         string `db:"title"`
	Description   string `db:"description"`
	HasRotaryPlug bool   `db:"has_rotary_plug"`
	HasDensity    bool   `db:"has_density"`
	HasThickness  bool   `db:"has_thickness"`
	MaterialType  string `db:"material_type"`
}
