package models

type Addit struct {
	Id          string `db:"id"`
	Materials   string `db:"materials"`
	Mod         string `db:"mod"`
	Temperature string `db:"temperature"`
	Mounting    string `db:"mounting"`
	Graphite    string `db:"graphite"`
	Fillers     string `db:"fillers"`
}

type UpdateGrap struct {
	Id       string
	Graphite string
}

type UpdateMat struct {
	Id        string
	Materials string
}

type UpdateTemp struct {
	Id          string
	Temperature string
}

type UpdateMod struct {
	Id  string
	Mod string
}

type UpdateMoun struct {
	Id       string
	Mounting string
}

type UpdateFill struct {
	Id      string
	Fillers string
}

type UpdateAdditDTO struct {
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
