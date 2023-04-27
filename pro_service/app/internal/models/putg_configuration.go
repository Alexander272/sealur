package models

type PutgConfiguration struct {
	Id          string `db:"id"`
	Title       string `db:"title"`
	Code        string `db:"code"`
	HasStandard bool   `db:"has_standard"`
	HasDrawing  bool   `db:"has_drawing"`
}
