package models

type RingMaterial struct {
	Id          string `db:"id"`
	Type        string `db:"type"`
	Title       string `db:"title"`
	Description string `db:"description"`
	IsDefault   bool   `db:"is_default"`
	Designation string `db:"designation"`
}
