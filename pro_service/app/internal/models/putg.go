package models

type Putg struct {
	Id           string `db:"id"`
	TypeFlId     string `db:"type_fl_id"`
	TypePr       string `db:"type_pr"`
	Form         string `db:"form"`
	Construction string `db:"construction"`
	Temperatures string `db:"temperatures"`
	Reinforce    string `db:"reinforce"`
	Obturator    string `db:"obturator"`
	ILimiter     string `db:"i_limiter"`
	OLimiter     string `db:"o_limiter"`
	Coating      string `db:"coating"`
	Mounting     string `db:"mounting"`
	Graphite     string `db:"graphite"`
}

type PutgDTO struct {
	Id           string `db:"id"`
	FlangeId     string `db:"flange_id"`
	TypeFlId     string `db:"type_fl_id"`
	TypePr       string `db:"type_pr"`
	Form         string `db:"form"`
	Construction string `db:"construction"`
	Temperatures string `db:"temperatures"`
	Reinforce    string `db:"reinforce"`
	Obturator    string `db:"obturator"`
	ILimiter     string `db:"i_limiter"`
	OLimiter     string `db:"o_limiter"`
	Coating      string `db:"coating"`
	Mounting     string `db:"mounting"`
	Graphite     string `db:"graphite"`
}
