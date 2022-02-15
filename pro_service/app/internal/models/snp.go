package models

type SNP struct {
	Id          string `db:"id"`
	TypeFlId    string `db:"type_fl_id"`
	TypePr      string `db:"type_pr"`
	Fillers     string `db:"filler"`
	Materials   string `db:"materials"`
	Mod         string `db:"mod"`
	Temperature string `db:"temperature"`
	Mounting    string `db:"mounting"`
	Graphite    string `db:"graphite"`
}
