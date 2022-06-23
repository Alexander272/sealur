package models

type Size struct {
	Id       string `db:"id"`
	Dn       string `db:"dn"`
	Pn       string `db:"pn"`
	D4       string `db:"d4"`
	D3       string `db:"d3"`
	D2       string `db:"d2"`
	D1       string `db:"d1"`
	H        string `db:"h"`
	S2       string `db:"s2"`
	S3       string `db:"s3"`
	TypePr   string `db:"type_pr"`
	TypeFlId string `db:"type_fl_id"`
}

type SizeInterview struct {
	Id    string `db:"id"`
	Dy    string `db:"dy"`
	Py    string `db:"py"`
	D1    string `db:"d1"`
	D2    string `db:"d2"`
	DUp   string `db:"d_up"`
	D     string `db:"d"`
	H1    string `db:"h1"`
	H2    string `db:"h2"`
	Bolt  string `db:"bolt"`
	Count int32  `db:"count_bolt"`
	Row   int32  `db:"row_count"`
}
