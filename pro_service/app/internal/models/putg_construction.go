package models

type PutgConstruction struct {
	Id            string `db:"id"`
	Title         string `db:"title"`
	Code          string `db:"code"`
	HasD4         bool   `db:"has_d4"`
	HasD3         bool   `db:"has_d3"`
	HasD2         bool   `db:"has_d2"`
	HasD1         bool   `db:"has_d1"`
	HasRotaryPlug bool   `db:"has_rotary_plug"`
	HasInnerRing  bool   `db:"has_inner_ring"`
	HasOuterRing  bool   `db:"has_outer_ring"`
}