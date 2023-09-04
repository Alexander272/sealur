package models

type Ring struct {
	Id                   string `db:"id"`
	PositionId           string `db:"position_id"`
	TypeId               string `db:"type_id"`
	TypeCode             string `db:"type_code"`
	TypeImage            string `db:"t_image"`
	TypeMT               string `db:"t_mt"`
	TypeHRP              bool   `db:"t_hrp"`
	TypeHD               bool   `db:"t_hd"`
	TypeHT               bool   `db:"t_ht"`
	TypeDesignation      string `db:"t_designation"`
	DensityId            string `db:"density_id"`
	DensityCode          string `db:"d_code"`
	DensityHRP           bool   `db:"d_hrp"`
	ConstructionCode     string `db:"construction_code"`
	ConstructionWRP      bool   `db:"construction_wrp"`
	ConstructionBaseCode string `db:"construction_bc"`
	Size                 string `db:"size"`
	Thickness            string `db:"thickness"`
	Material             string `db:"material"`
	Modifying            string `db:"modifying"`
	Drawing              string `db:"drawing"`
}
