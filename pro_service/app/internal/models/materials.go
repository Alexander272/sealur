package models

type Materials struct {
	Id      string `db:"id"`
	Title   string `db:"title"`
	TypeMat string `db:"type_mat"`
}

type BoltMaterials struct {
	Id    string `db:"id"`
	Title string `db:"title"`
}
