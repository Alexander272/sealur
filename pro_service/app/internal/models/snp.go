package models

type SNP struct {
	Id       string `db:"id"`
	TypeFlId string `db:"type_fl_id"`
	TypePr   string `db:"type_pr"`
	Fillers  string `db:"filler"`
	Frame    string `db:"frame"`
	Ir       string `db:"in_ring"`
	Or       string `db:"ou_ring"`
	Mounting string `db:"mounting"`
	Graphite string `db:"graphite"`
}

type SnpDTO struct {
	Id       string `db:"id"`
	StandId  string `db:"stand_id"`
	FlangeId string `db:"flange_id"`
	TypeFlId string `db:"type_fl_id"`
	TypePr   string `db:"type_pr"`
	Fillers  string `db:"filler"`
	Frame    string `db:"frame"`
	Ir       string `db:"in_ring"`
	Or       string `db:"ou_ring"`
	Mounting string `db:"mounting"`
	Graphite string `db:"graphite"`
}
