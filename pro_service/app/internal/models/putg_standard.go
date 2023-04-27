package models

type PutgStandard struct {
	Id      string `db:"id"`
	Title   string `db:"title"`
	Code    string `db:"code"`
	DnTitle string `db:"dn_title"`
	PnTitle string `db:"pn_title"`
}
