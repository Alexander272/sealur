package models

import "github.com/lib/pq"

type SnpStandard struct {
	Id             string         `db:"id"`
	DnTitle        string         `db:"dn_title"`
	PnTitle        string         `db:"pn_title"`
	StandardTitle  string         `db:"standard_title"`
	StandardFormat pq.StringArray `db:"standard_format"`
	FlangeTitle    string         `db:"flange_title"`
	FlangeCode     string         `db:"flange_code"`
}
