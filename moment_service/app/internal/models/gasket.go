package models

type GetGasket struct {
	TypeGasket string  `db:"type_gasket"`
	Env        string  `db:"env"`
	Thickness  float32 `db:"thickness"`
}

type Gasket struct {
	TypeGasket      string  `db:"type_gasket"`
	Env             string  `db:"env"`
	M               float32 `db:"m"`
	SpecificPres    float32 `db:"q_s_pres"`
	PermissiblePres float32 `db:"q_p_pres"`
	Compression     float32 `db:"k"`
	Epsilon         float32 `db:"e"`
	Thickness       float32 `db:"thickness"`
	IsOval          bool    `db:"is_oval"`
}
