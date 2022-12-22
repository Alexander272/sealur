package models

type EnvDTO struct {
	Id    string `db:"id"`
	Title string `db:"title"`
}

type EnvDataDTO struct {
	Id           string  `db:"id"`
	EnvId        string  `db:"env_id"`
	GasketId     string  `db:"gasket_id"`
	M            float64 `db:"m"`
	SpecificPres float64 `db:"specific_pres"`
}
