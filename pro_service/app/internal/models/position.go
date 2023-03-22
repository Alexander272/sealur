package models

type Position struct {
	Id          string `db:"id"`
	Designation string `db:"designation"`
	Descriprion string `db:"description"`
	Count       int32  `db:"count"`
	Sizes       string `db:"sizes"`
	Drawing     string `db:"drawing"`
	OrderId     string `db:"order_id"`
}

type PositionNew struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	Amount string `db:"amount"`
	Type   string `db:"type"`
	Count  int64  `db:"count"`
}

type FullPosition struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	Amount string `db:"amount"`
	Type   string `db:"type"`
	Count  int64  `db:"count"`

	MainSnpStandardId   string `db:"snp_standard_id"`
	MainTypeId          string `db:"snp_type_id"`
	MainFlangeTypeCode  string `db:"flange_type_code"`
	MainFlangeTypeTitle string `db:"flange_type_title"`

	FillerCode    string `db:"filler_code"`
	FrameCode     string `db:"frame_code"`
	InnerRingCode string `db:"inner_ring_code"`
	OuterRingCode string `db:"outer_ring_code"`

	D4      string `db:"d4"`
	D3      string `db:"d3"`
	D2      string `db:"d2"`
	D1      string `db:"d1"`
	H       string `db:"h"`
	Another string `db:"another"`

	HasJumper    bool   `db:"has_jumper"`
	JumperCode   string `db:"jumper_code"`
	JumperWidth  string `db:"jumper_width"`
	HasHole      bool   `db:"has_hole"`
	HasMounting  bool   `db:"has_mounting"`
	MountingCode string `db:"mounting_code"`
	Drawing      string `db:"drawing"`
}
