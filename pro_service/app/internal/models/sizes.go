package models

type Size struct {
	Id       string  `db:"id"`
	Dn       string  `db:"dn"`
	Pn       string  `db:"pn"`
	D4       float32 `db:"d4"`
	D3       float32 `db:"d3"`
	D2       float32 `db:"d2"`
	D1       float32 `db:"d1"`
	H        string  `db:"h"`
	S2       string  `db:"s2"`
	S3       string  `db:"s3"`
	TypePr   string  `db:"type_pr"`
	TypeFlId string  `db:"type_fl_id"`
	Adn      int32   `db:"adn"`
}
