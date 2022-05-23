package models

type Putgm struct {
	Id           string `db:"id"`
	TypeFlId     string `db:"type_fl_id"`
	TypePr       string `db:"type_pr"`
	Form         string `db:"form"`
	Construction string `db:"construction"`
	Temperatures string `db:"temperatures"`
	Basis        string `db:"basis"`
	Obturator    string `db:"obturator"`
	Coating      string `db:"coating"`
	Mounting     string `db:"mounting"`
	Graphite     string `db:"graphite"`
}

type PutgmDTO struct {
	Id           string `db:"id"`
	FlangeId     string `db:"flange_id"`
	TypeFlId     string `db:"type_fl_id"`
	TypePr       string `db:"type_pr"`
	Form         string `db:"form"`
	Construction string `db:"construction"`
	Temperatures string `db:"temperatures"`
	Basis        string `db:"basis"`
	Obturator    string `db:"obturator"`
	Coating      string `db:"coating"`
	Mounting     string `db:"mounting"`
	Graphite     string `db:"graphite"`
}
