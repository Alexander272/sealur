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

type Gasket struct {
	TypeGasket      string  `db:"type_gasket"`
	Env             string  `db:"env"`
	M               float64 `db:"m"`
	SpecificPres    float64 `db:"q_s_pres"`
	PermissiblePres float64 `db:"q_p_pres"`
	Compression     float64 `db:"k"`
	Epsilon         float64 `db:"e"`
	Thickness       float64 `db:"thickness"`
	Type            string  `db:"gasket_type"`
}
