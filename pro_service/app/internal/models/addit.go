package models

type Addit struct {
	Id           string `db:"id"`
	Materials    string `db:"materials"`
	Mod          string `db:"mod"`
	Temperature  string `db:"temperature"`
	Mounting     string `db:"mounting"`
	Graphite     string `db:"graphite"`
	Fillers      string `db:"fillers"`
	Coating      string `db:"coating"`
	Construction string `db:"construction"`
	Obturator    string `db:"obturator"`
	Basis        string `db:"basis"`
	PObturator   string `db:"p_obturator"`
	Sealant      string `db:"sealant"`
}

type UpdateGrap struct {
	Id       string
	Graphite string
}

type UpdateMat struct {
	Id        string
	Materials string
}

type UpdateTemp struct {
	Id          string
	Temperature string
}

type UpdateMod struct {
	Id  string
	Mod string
}

type UpdateMoun struct {
	Id       string
	Mounting string
}

type UpdateFill struct {
	Id      string
	Fillers string
}

type UpdateCoating struct {
	Id      string
	Coating string
}

type UpdateConstr struct {
	Id           string
	Construction string
}

type UpdateObturator struct {
	Id        string
	Obturator string
}

type UpdateBasis struct {
	Id    string
	Basis string
}

type UpdatePObturator struct {
	Id         string
	PObturator string
}

type UpdateSealant struct {
	Id      string
	Sealant string
}

type UpdateAdditDTO struct {
	Id       string `db:"id"`
	TypeFlId string `db:"type_fl_id"`
	TypePr   string `db:"type_pr"`
	Fillers  string `db:"filler"`
	Frame    string `db:"frame"`
	Ir       string `db:"in_ring"`
	Or       string `db:"ou_ring"`
	Mounting string `db:"mounting"`
	Graphite string `db:"graphite"`

	Construction string `db:"construction"`
	Temperature  string `db:"temperature"`
	Obturator    string `db:"obturator"`
	ILimiter     string `db:"i_limiter"`
	OLimiter     string `db:"o_limiter"`
	Coating      string `db:"coating"`

	Basis      string `db:"basis"`
	PObturator string `db:"p_obturator"`
	Sealant    string `db:"sealant"`
}
