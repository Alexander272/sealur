package models

type Material struct {
	Id       string `db:"id"`
	Title    string `db:"title"`
	Code     string `db:"code"`
	ShortEn  string `db:"short_en"`
	ShortRus string `db:"short_rus"`
}

type SNPMaterial struct {
	Id string `db:"id"`
	// MaterialId pq.StringArray `db:"material_id"`
	// Default    string         `db:"default_id"`
	// Type       string         `db:"type"`
	Title      string `db:"title"`
	Code       string `db:"code"`
	ShortEn    string `db:"short_en"`
	ShortRus   string `db:"short_rus"`
	MaterialId string `db:"material_id"`
	Default    string `db:"default_id"`
	Type       string `db:"type"`
}

type SnpMaterial struct {
	Id         string `db:"id"`
	MaterialId string `db:"material_id"`
	Type       string `db:"type"`
	IsDefault  bool   `db:"is_default"`
	Code       string `db:"code"`
	IsStandard bool   `db:"is_standard"`
	BaseCode   string `db:"base_code"`
	Title      string `db:"title"`
}
