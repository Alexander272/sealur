package models

type SNP struct {
	Id        string `db:"id"`
	TypeFlId  string `db:"type_fl_id"`
	TypePr    string `db:"type_pr"`
	Fillers   string `db:"filler"`
	Materials string `db:"materials"`
	DefMat    string `db:"def_mat"`
	Mounting  string `db:"mounting"`
	Graphite  string `db:"graphite"`
}
