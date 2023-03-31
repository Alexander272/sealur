package pro_model

import (
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_size_model"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type SNPDTO struct {
	StandId  string                `json:"standId" binding:"required"`
	FlangeId string                `json:"flangeId" binding:"required"`
	TypeFlId string                `json:"typeFlId" binding:"required"`
	TypePr   string                `json:"typePr" binding:"required"`
	Fillers  []*pro_api.Filler     `json:"fillers" binding:"required"`
	Frame    *pro_api.SnpMaterials `json:"frame"`
	Ir       *pro_api.SnpMaterials `json:"ir"`
	Or       *pro_api.SnpMaterials `json:"or"`
	Mounting []string              `json:"mounting" binding:"required"`
	Graphite []string              `json:"graphite" binding:"required"`
}

type DefResponse struct {
	TypeFl []*pro_api.TypeFl     `json:"typeFl"`
	Snp    []*pro_api.SNP        `json:"snp"`
	Sizes  *pro_api.SizeResponse `json:"sizes"`
}

type SnpData struct {
	Main     MainSnp     `json:"main"`
	Size     SizeSnp     `json:"size"`
	Material MaterialSnp `json:"material"`
	Design   DesignSnp   `json:"design"`
}

func (s *SnpData) Parse() *position_model.PositionSnp {
	return &position_model.PositionSnp{
		Main:     s.Main.Parse(),
		Size:     s.Size.Parse(),
		Material: s.Material.Parse(),
		Design:   s.Design.Parse(),
	}
}

type MainSnp struct {
	FlangeTypeCode  string `json:"flangeTypeCode"`
	FlangeTypeTitle string `json:"flangeTypeTitle"`
	SnpStandardId   string `json:"snpStandardId"`
	SnpTypeId       string `json:"snpTypeId"`
}

func (s *MainSnp) Parse() *position_model.PositionSnp_Main {
	return &position_model.PositionSnp_Main{
		SnpStandardId:   s.SnpStandardId,
		SnpTypeId:       s.SnpTypeId,
		FlangeTypeCode:  s.FlangeTypeCode,
		FlangeTypeTitle: s.FlangeTypeTitle,
	}
}

type Pn struct {
	Mpa string `json:"mpa"`
	Kg  string `json:"kg"`
}

func (pn *Pn) Parse() *snp_size_model.Pn {
	return &snp_size_model.Pn{
		Mpa: pn.Mpa,
		Kg:  pn.Kg,
	}
}

type SizeSnp struct {
	Dn      string `json:"dn"`
	Pn      Pn     `json:"pn"`
	D4      string `json:"d4"`
	D3      string `json:"d3"`
	D2      string `json:"d2"`
	D1      string `json:"d1"`
	H       string `json:"h"`
	Another string `json:"another"`
	S2      string `json:"s2"`
	S3      string `json:"s3"`
}

func (s *SizeSnp) Parse() *position_model.PositionSnp_Size {
	return &position_model.PositionSnp_Size{
		Dn:      s.Dn,
		Pn:      s.Pn.Parse(),
		D4:      s.D4,
		D3:      s.D3,
		D2:      s.D2,
		D1:      s.D1,
		H:       s.H,
		Another: s.Another,
		S2:      s.S2,
		S3:      s.S3,
	}
}

type MaterialSnp struct {
	Ir     Material  `json:"ir"`
	Or     Material  `json:"or"`
	Fr     Material  `json:"fr"`
	Filler SnpFiller `json:"filler"`
}

func (s *MaterialSnp) Parse() *position_model.PositionSnp_Material {
	return &position_model.PositionSnp_Material{
		FillerId:      s.Filler.Id,
		FillerCode:    s.Filler.BaseCode,
		FrameId:       s.Fr.Id,
		FrameCode:     s.Fr.Code,
		InnerRingId:   s.Ir.Id,
		InnerRingCode: s.Ir.Code,
		OuterRingId:   s.Or.Id,
		OuterRingCode: s.Or.Code,
	}
}

type DesignSnp struct {
	HasHole  bool     `json:"hasHole"`
	Jumper   Jumper   `json:"jumper"`
	Mounting Mounting `json:"mounting"`
}

func (s *DesignSnp) Parse() *position_model.PositionSnp_Design {
	return &position_model.PositionSnp_Design{
		HasJumper:    s.Jumper.HasJumper,
		JumperCode:   s.Jumper.Code,
		JumperWidth:  s.Jumper.Width,
		HasHole:      s.HasHole,
		HasMounting:  s.Mounting.HasMounting,
		MountingCode: s.Mounting.Code,
	}
}
