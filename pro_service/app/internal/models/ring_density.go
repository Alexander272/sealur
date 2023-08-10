package models

type RingDensity struct {
	Id            string `db:"id"`
	TypeId        string `db:"type_id"`
	Code          string `db:"code"`
	Title         string `db:"title"`
	Description   string `db:"description"`
	HasRotaryPlug bool   `db:"has_rotary_plug"`
}
