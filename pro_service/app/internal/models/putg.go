package models

//* old
// type Putg struct {
// 	Id           string `db:"id"`
// 	TypeFlId     string `db:"type_fl_id"`
// 	TypePr       string `db:"type_pr"`
// 	Form         string `db:"form"`
// 	Construction string `db:"construction"`
// 	Temperatures string `db:"temperatures"`
// 	Reinforce    string `db:"reinforce"`
// 	Obturator    string `db:"obturator"`
// 	ILimiter     string `db:"i_limiter"`
// 	OLimiter     string `db:"o_limiter"`
// 	Coating      string `db:"coating"`
// 	Mounting     string `db:"mounting"`
// 	Graphite     string `db:"graphite"`
// }

// type PutgDTO struct {
// 	Id           string `db:"id"`
// 	FlangeId     string `db:"flange_id"`
// 	TypeFlId     string `db:"type_fl_id"`
// 	TypePr       string `db:"type_pr"`
// 	Form         string `db:"form"`
// 	Construction string `db:"construction"`
// 	Temperatures string `db:"temperatures"`
// 	Reinforce    string `db:"reinforce"`
// 	Obturator    string `db:"obturator"`
// 	ILimiter     string `db:"i_limiter"`
// 	OLimiter     string `db:"o_limiter"`
// 	Coating      string `db:"coating"`
// 	Mounting     string `db:"mounting"`
// 	Graphite     string `db:"graphite"`
// }

type PutgMainBlock struct {
	Id         string `db:"id"`
	PositionId string `db:"position_id"`

	// standard
	PutgStandardId string `db:"putg_standard_id"`
	PutgStandardDn string `db:"dn_title"`
	PutgStandardPn string `db:"pn_title"`
	StandardId     string `db:"standard_id"`
	StandardTitle  string `db:"standard_title"`
	FlangeId       string `db:"flange_standard_id"`
	FlangeTitle    string `db:"flange_title"`
	FlangeCode     string `db:"flange_code"`

	// flangeType
	FlangeTypeId    string `db:"flange_type_id"`
	FlangeTypeCode  string `db:"flange_type_code"`
	FlangeTypeTitle string `db:"flange_type_title"`

	// configuration
	ConfId          string `db:"configuration_id"`
	ConfTitle       string `db:"conf_title"`
	ConfCode        string `db:"conf_code"`
	ConfHasStandard bool   `db:"conf_has_standard"`
	ConfHasDrawing  bool   `db:"conf_has_drawing"`
}

type PutgMaterialBlock struct {
	Id         string `db:"id"`
	PositionId string `db:"position_id"`

	// filler
	FillerId          string `db:"filler_id"`
	FillerBaseId      string `db:"f_base_id"`
	FillerCode        string `db:"f_code"`
	FillerTitle       string `db:"f_title"`
	FillerDescription string `db:"f_description"`
	FillerDesignation string `db:"f_designation"`

	// type
	TypeId           string  `db:"type_id"`
	TypeTitle        string  `db:"t_title"`
	TypeCode         string  `db:"t_code"`
	TypeMinThickness float64 `db:"t_min"`
	TypeMaxThickness float64 `db:"t_max"`
	TypeDescription  string  `db:"t_description"`
	TypeBaseCode     string  `db:"t_type_code"`
	TypeHasReinforce bool    `db:"t_has_reinforce"`

	// construction
	ConstructionId            string `db:"construction_id"`
	ConstructionCode          string `db:"construction_code"`
	ConstructionTitle         string `db:"c_title"`
	ConstructionBaseId        string `db:"c_base_id"`
	ConstructionDescription   string `db:"c_description"`
	ConstructionHasD4         bool   `db:"c_has_d4"`
	ConstructionHasD3         bool   `db:"c_has_d3"`
	ConstructionHasD2         bool   `db:"c_has_d2"`
	ConstructionHasD1         bool   `db:"c_has_d1"`
	ConstructionHasRotaryPlug bool   `db:"c_has_rotary_plug"`
	ConstructionHasInnerRing  bool   `db:"c_has_inner_ring"`
	ConstructionHasOuterRing  bool   `db:"c_has_outer_ring"`

	// reinforce
	ReinforceId         string  `db:"reinforce_id"`
	ReinforceBaseCode   string  `db:"reinforce_code"`
	ReinforceTitle      string  `db:"reinforce_title"`
	ReinforceCode       *string `db:"r_code"`
	ReinforceMaterialId *string `db:"r_material_id"`
	ReinforceType       *string `db:"r_type"`
	ReinforceIsDefault  *bool   `db:"r_is_default"`

	// rotaryPlug
	RotaryPlugId         string  `db:"rotary_plug_id"`
	RotaryPlugBaseCode   string  `db:"rotary_plug_code"`
	RotaryPlugTitle      string  `db:"rotary_plug_title"`
	RotaryPlugCode       *string `db:"rp_code"`
	RotaryPlugMaterialId *string `db:"rp_material_id"`
	RotaryPlugType       *string `db:"rp_type"`
	RotaryPlugIsDefault  *bool   `db:"rp_is_default"`

	// innerRing
	InnerRingId         string  `db:"inner_ring_id"`
	InnerRingBaseCode   string  `db:"inner_ring_code"`
	InnerRingTitle      string  `db:"inner_ring_title"`
	InnerRingCode       *string `db:"ir_code"`
	InnerRingMaterialId *string `db:"ir_material_id"`
	InnerRingType       *string `db:"ir_type"`
	InnerRingIsDefault  *bool   `db:"ir_is_default"`

	// outerRing
	OuterRingId         string  `db:"outer_ring_id"`
	OuterRingBaseCode   string  `db:"outer_ring_code"`
	OuterRingTitle      string  `db:"outer_ring_title"`
	OuterRingCode       *string `db:"or_code"`
	OuterRingMaterialId *string `db:"or_material_id"`
	OuterRingType       *string `db:"or_type"`
	OuterRingIsDefault  *bool   `db:"or_is_default"`
}

type PutgSizeBlock struct {
	Id            string `db:"id"`
	PositionId    string `db:"position_id"`
	Dn            string `db:"dn"`
	DnMm          string `db:"dn_mm"`
	PnMpa         string `db:"pn_mpa"`
	PnKg          string `db:"pn_kg"`
	D4            string `db:"d4"`
	D3            string `db:"d3"`
	D2            string `db:"d2"`
	D1            string `db:"d1"`
	H             string `db:"h"`
	Another       string `db:"another"`
	UseDimensions bool   `db:"use_dimensions"`
}

type PutgDesignBlock struct {
	Id           string `db:"id"`
	PositionId   string `db:"position_id"`
	HasJumper    bool   `db:"has_jumper"`
	JumperCode   string `db:"jumper_code"`
	JumperWidth  string `db:"jumper_width"`
	HasHole      bool   `db:"has_hole"`
	HasCoating   bool   `db:"has_coating"`
	HasRemovable bool   `db:"has_removable"`
	HasMounting  bool   `db:"has_mounting"`
	MountingCode string `db:"mounting_code"`
	Drawing      string `db:"drawing"`
}
