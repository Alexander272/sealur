package models

type Voltage struct {
	Id          string  `db:"id"`
	Temperature float64 `db:"temperature"`
	Voltage     float64 `db:"voltage"`
}

type Elasticity struct {
	Id          string  `db:"id"`
	Temperature float64 `db:"temperature"`
	Elasticity  float64 `db:"elasticity"`
}

type Alpha struct {
	Id          string  `db:"id"`
	Temperature float64 `db:"temperature"`
	Alpha       float64 `db:"alpha"`
}

type MaterialsAll struct {
	Title      string `db:"title"`
	Alpha      []Alpha
	Elasticity []Elasticity
	Voltage    []Voltage
}

type MaterialsWithIsEmpty struct {
	Id                string `db:"id"`
	Title             string `db:"title"`
	IsEmptyAlpha      bool   `db:"is_empty_alpha"`
	IsEmptyElasticity bool   `db:"is_empty_elasticity"`
	IsEmptyVoltage    bool   `db:"is_empty_voltage"`
}

type MaterialsDTO struct {
	Id    string `db:"id"`
	Title string `db:"title"`
}

type VoltageDTO struct {
	Id          string  `db:"id"`
	MarkId      string  `db:"mark_id"`
	Temperature float64 `db:"temperature"`
	Voltage     float64 `db:"voltage"`
}

type ElasticityDTO struct {
	Id          string  `db:"id"`
	MarkId      string  `db:"mark_id"`
	Temperature float64 `db:"temperature"`
	Elasticity  float64 `db:"elasticity"`
}

type AlphaDTO struct {
	Id          string  `db:"id"`
	MarkId      string  `db:"mark_id"`
	Temperature float64 `db:"temperature"`
	Alpha       float64 `db:"alpha"`
}

type MaterialsResult struct {
	Title       string
	AlphaF      float64
	EpsilonAt20 float64
	Epsilon     float64
	SigmaAt20   float64
	Sigma       float64
}
