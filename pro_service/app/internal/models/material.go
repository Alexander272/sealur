package models

type Material struct {
	Id       string `db:"id"`
	Title    string `db:"title"`
	Code     string `db:"code"`
	ShortEn  string `db:"short_en"`
	ShortRus string `db:"short_rus"`
}
