package models

type GetGasket struct {
	TypeGasket string  `db:"type_gasket"`
	Env        string  `db:"env"`
	Thickness  float64 `db:"thickness"`
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
	IsOval          bool    `db:"is_oval"`
	IsMetal         bool    `db:"is_metal"`
}
