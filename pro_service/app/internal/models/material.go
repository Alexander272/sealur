package models

import "github.com/lib/pq"

type Material struct {
	Id       string `db:"id"`
	Title    string `db:"title"`
	Code     string `db:"code"`
	ShortEn  string `db:"short_en"`
	ShortRus string `db:"short_rus"`
}

type SNPMaterial struct {
	Id         string         `db:"id"`
	MaterialId pq.StringArray `db:"material_id"`
	Default    string         `db:"default"`
	Type       string         `db:"type"`
}
