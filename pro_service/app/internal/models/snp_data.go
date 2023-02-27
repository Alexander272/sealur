package models

type SnpData struct {
	Id           string `db:"id"`
	TypeId       string `db:"type_id"`
	HasInnerRing bool   `db:"has_inner_ring"`
	HasFrame     bool   `db:"has_frame"`
	HasOuterRing bool   `db:"has_outer_ring"`
	HasHole      bool   `db:"has_hole"`
	HasJumper    bool   `db:"has_jumper"`
	HasMounting  bool   `db:"has_mounting"`
}
