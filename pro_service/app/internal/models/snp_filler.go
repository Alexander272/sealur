package models

import "github.com/lib/pq"

type SNPFiller struct {
	Id           string `db:"id"`
	Title        string `db:"title"`
	AnotherTitle string `db:"another_title"`
	Code         string `db:"code"`
	Description  string `db:"description"`
	Designation  string `db:"designation"`
	Temperature  string `db:"temperature"`
}

type SnpFillerNew struct {
	Id            string         `db:"id"`
	Temperature   string         `db:"temperature"`
	BaseCode      string         `db:"base_code"`
	Code          string         `db:"code"`
	Title         string         `db:"title"`
	Description   string         `db:"description"`
	Designation   string         `db:"designation"`
	DisabledTypes pq.StringArray `db:"disabled_types"`
}
