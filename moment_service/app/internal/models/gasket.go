package models

type GasketType int

const (
	Soft GasketType = iota
	Oval
	Metal
)

type GetGasket struct {
	GasketId  string  `db:"gasket_id"`
	EnvId     string  `db:"env_id"`
	Thickness float64 `db:"thickness"`
}

type GasketDTO struct {
	Id    string `db:"id"`
	Title string `db:"title"`
}

type GasketDataDTO struct {
	Id              string  `db:"id"`
	GasketId        string  `db:"gasket_id"`
	PermissiblePres float64 `db:"permissible_pres"`
	Compression     float64 `db:"compression"`
	Epsilon         float64 `db:"epsilon"`
	Thickness       float64 `db:"thickness"`
	TypeId          string  `db:"type_id"`
}

type TypeGasketDTO struct {
	Id    string `db:"id"`
	Title string `db:"title"`
}

type Gasket struct {
	Id              string  `db:"id"`
	GasketId        string  `db:"gasket_id"`
	EnvId           string  `db:"env_id"`
	M               float64 `db:"m"`
	SpecificPres    float64 `db:"specific_pres"`
	PermissiblePres float64 `db:"permissible_pres"`
	Compression     float64 `db:"compression"`
	Epsilon         float64 `db:"epsilon"`
	Thickness       float64 `db:"thickness"`
	Type            string  `db:"type"`
}
