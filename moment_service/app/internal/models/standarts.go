package models

type StandartDTO struct {
	Id             string `db:"id"`
	Title          string `db:"title"`
	TypeId         string `db:"type_id"`
	TitleDn        string `db:"title_dn"`
	TitlePn        string `db:"title_pn"`
	IsNeedRow      bool   `db:"is_need_row"`
	Rows           string `db:"rows"`
	IsInch         bool   `db:"is_inch"`
	HasDesignation bool   `db:"has_designation"`
}

type StandartWithSize struct {
	Id             string `db:"id"`
	Title          string `db:"title"`
	TypeId         string `db:"type_id"`
	TitleDn        string `db:"title_dn"`
	TitlePn        string `db:"title_pn"`
	IsNeedRow      bool   `db:"is_need_row"`
	Rows           string `db:"rows"`
	IsInch         bool   `db:"is_inch"`
	HasDesignation bool   `db:"has_designation"`
	Sizes          []FlangeSize
}
