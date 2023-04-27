package models

type PutgData struct {
	Id           string `db:"id"`
	FillerId     string `db:"filler_id"`
	HasJumper    bool   `db:"has_jumper"`
	HasHole      bool   `db:"has_hole"`
	HasRemovable bool   `db:"has_removable"`
	HasMounting  bool   `db:"has_mounting"`
	HasCoating   bool   `db:"has_coating"`
}
