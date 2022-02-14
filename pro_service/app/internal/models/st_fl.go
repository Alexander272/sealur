package models

type StFl struct {
	Id       string `db:"id"`
	StandId  string `db:"stand_id"`
	Stand    string `db:"stand"`
	FlangeId string `db:"fl_id"`
	Flange   string `db:"flange"`
}
