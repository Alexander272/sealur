package models

import "github.com/lib/pq"

type SNP struct {
	Id       string `db:"id"`
	TypeFlId string `db:"type_fl_id"`
	TypePr   string `db:"type_pr"`
	Fillers  string `db:"filler"`
	Frame    string `db:"frame"`
	Ir       string `db:"in_ring"`
	Or       string `db:"ou_ring"`
	Mounting string `db:"mounting"`
	Graphite string `db:"graphite"`
}

type SnpDTO struct {
	Id       string `db:"id"`
	StandId  string `db:"stand_id"`
	FlangeId string `db:"flange_id"`
	TypeFlId string `db:"type_fl_id"`
	TypePr   string `db:"type_pr"`
	Fillers  string `db:"filler"`
	Frame    string `db:"frame"`
	Ir       string `db:"in_ring"`
	Or       string `db:"ou_ring"`
	Mounting string `db:"mounting"`
	Graphite string `db:"graphite"`
}

type SnpMainBlock struct {
	Id               string         `db:"id"`
	PositionId       string         `db:"position_id"`
	FlangeTypeCode   string         `db:"flange_type_code"`
	FlangeTypeTitle  string         `db:"flange_type_title"`
	SnpTypeId        string         `db:"snp_type_id"`
	SnpTypeCode      string         `db:"code"`
	SnpTypeTitle     string         `db:"title"`
	SnpStandardId    string         `db:"snp_standard_id"`
	SnpStandardDn    string         `db:"dn_title"`
	SnpStandardPn    string         `db:"pn_title"`
	SnpStandardHasD2 bool           `db:"has_d2"`
	StandardId       string         `db:"standard_id"`
	StandardTitle    string         `db:"standard_title"`
	StandardFormat   pq.StringArray `db:"standard_format"`
	FlangeId         string         `db:"flange_standard_id"`
	FlangeTitle      string         `db:"flange_title"`
	FlangeCode       string         `db:"flange_code"`
}

type SnpMaterialBlock struct {
	Id                 string  `db:"id"`
	PositionId         string  `db:"position_id"`
	FillerId           string  `db:"filler_id"`
	FillerCode         string  `db:"code"`
	FillerTitle        string  `db:"title"`
	FillerAnotherTitle string  `db:"another_title"`
	FillerDescription  string  `db:"description"`
	FillerDesignation  string  `db:"designation"`
	FrameId            string  `db:"m1_id"`
	FrameCode          *string `db:"m1_code"`
	FrameTitle         *string `db:"m1_title"`
	FrameShortEn       *string `db:"m1_short_en"`
	FrameShortRus      *string `db:"m1_short_rus"`
	InnerRingId        string  `db:"m2_id"`
	InnerRingCode      *string `db:"m2_code"`
	InnerRingTitle     *string `db:"m2_title"`
	InnerRingShortEn   *string `db:"m2_short_en"`
	InnerRingShortRus  *string `db:"m2_short_rus"`
	OuterRingId        string  `db:"m3_id"`
	OuterRingCode      *string `db:"m3_code"`
	OuterRingTitle     *string `db:"m3_title"`
	OuterRingShortEn   *string `db:"m3_short_en"`
	OuterRingShortRus  *string `db:"m3_short_rus"`
}

type SnpSizeBlock struct {
	Id         string `db:"id"`
	PositionId string `db:"position_id"`
	Dn         string `db:"dn"`
	PnMpa      string `db:"pn_mpa"`
	PnKg       string `db:"pn_kg"`
	D4         string `db:"d4"`
	D3         string `db:"d3"`
	D2         string `db:"d2"`
	D1         string `db:"d1"`
	H          string `db:"h"`
	S2         string `db:"s2"`
	S3         string `db:"s3"`
	Another    string `db:"another"`
}

type SnpDesignBlock struct {
	Id           string `db:"id"`
	PositionId   string `db:"position_id"`
	HasJumper    bool   `db:"has_jumper"`
	JumperCode   string `db:"jumper_code"`
	JumperWidth  string `db:"jumper_width"`
	HasHole      bool   `db:"has_hole"`
	HasMounting  bool   `db:"has_mounting"`
	MountingCode string `db:"mounting_code"`
	Drawing      string `db:"drawing"`
}
