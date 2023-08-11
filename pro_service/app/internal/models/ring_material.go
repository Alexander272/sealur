package models

type RingMaterial struct {
	Id        string `db:"id"`
	Type      string `db:"type"`
	Title     string `db:"title"`
	IsDefault bool   `db:"is_default"`
}
