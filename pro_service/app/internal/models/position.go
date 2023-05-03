package models

type PositionAnalytics struct {
	OrderCount       int64 `db:"order_count"`
	UserCount        int64 `db:"user_count"`
	PositionCount    int64 `db:"position_count"`
	PositionSnpCount int64 `db:"position_snp_count"`
}

type PositionNew struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	Amount string `db:"amount"`
	Type   string `db:"type"`
	Count  int64  `db:"count"`
	Info   string `db:"info"`
}

type FullPosition struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	Amount string `db:"amount"`
	Type   string `db:"type"`
	Count  int64  `db:"count"`
	Info   string `db:"info"`

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
